package service

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"github.com/smancke/gogallery/imglib"
	"log"
	"net/http"
	"time"
)

type AuthCookie struct {
	DisplayName string   `json:"displayName"`
	UserName    string   `json:"userName"`
	Groups      []string `json:"groups"`
	LastSeen    int64    `json:"lastSeen"`
	UserId      string   `json:"userId"`
}

func getVerifiedAuthCookie(r *web.Request) (*AuthCookie, error) {
	if testOverwriteUsername != "" {
		return &AuthCookie{
			UserName:    testOverwriteUsername,
			DisplayName: testOverwriteUsername,
		}, nil
	}

	authCookie, err := readCookie(r, Cfg("cookieName"), Cfg("session_secret"))
	if err != nil {
		return nil, err
	}

	lastSeen := time.Unix(int64(authCookie.LastSeen/1000), 0)
	sessionDuration := time.Minute * time.Duration(CfgInt("sessionLifetimeMinutes"))
	validUntil := lastSeen.Add(sessionDuration)

	if validUntil.Before(time.Now()) {
		return authCookie, fmt.Errorf("session expired")
	}
	return authCookie, nil
}

func (userApi *UserApi) UserRequired(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	authCookie, err := getVerifiedAuthCookie(r)
	if err != nil {
		log.Printf("not logged in: %v, %#v", err, authCookie)
		sendError(rw, "not logged in", http.StatusForbidden)
		return
	}

	user, err := userApi.Lib().UserByUsername(authCookie.UserName)
	if err == nil && user != nil {
		userApi.user = user
	} else {
		userApi.user = &imglib.User{
			UserName: authCookie.UserName,
			NickName: authCookie.DisplayName,
		}
		err := userApi.Lib().CreateUser(userApi.user)
		if err != nil {
			panic(err)
		}
	}
	next(rw, r)
}

func readCookie(r *web.Request, cookieName string, secret string) (*AuthCookie, error) {
	rawCookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	jsonData, err := decrypt([]byte(rawCookie.Value), secret)
	if err != nil {
		return nil, err
	}

	jsonData, err = pkcs5UnPadding(jsonData)
	if err != nil {
		return nil, err
	}

	cookie := AuthCookie{}
	err = json.Unmarshal(jsonData, &cookie)
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}

func decrypt(data []byte, secret string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	aes, err := aes.NewCipher([]byte(secret))
	if err != nil {
		log.Printf("error on AES cipher creation: %v", err)
		return nil, err
	}

	bs := aes.BlockSize()
	if len(data)%bs != 0 {
		return nil, fmt.Errorf("AES need a multiple of the blocksize")
	}

	decrypted := make([]byte, len(data))
	decryptedComplete := decrypted
	for len(data) > 0 {
		aes.Decrypt(decrypted, data)
		decrypted = decrypted[bs:]
		data = data[bs:]
	}
	return decryptedComplete, nil
}

func pkcs5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	rindex := length - unpadding
	if rindex < 0 || rindex >= length {
		return nil, errors.New("wrong rindex for unpadding okcs5 cookie")
	}
	return src[:rindex], nil
}
