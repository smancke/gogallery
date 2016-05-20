package service

import (
	"bytes"
	"github.com/gocraft/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestReadCookie(t *testing.T) {

	c, r := givenACookieAndRequest()

	cookie, err := readCookie(r, c.Name, Cfg("session_secret"))

	assert.NoError(t, err)
	assert.Equal(t, "s.mancke@tarent.de", cookie.UserName)
	assert.Equal(t, "Sebastian Mancke", cookie.DisplayName)
	assert.Equal(t, int(1439039852142), int(cookie.LastSeen))
}

func TestVerifyCookie(t *testing.T) {

	_, r := givenACookieAndRequest()

	_, err := getVerifiedAuthCookie(r)

	assert.Error(t, err)
	assert.Equal(t, "session expired", err.Error())
}

func TestForbiddenOnApiUsage(t *testing.T) {

	router := createRouter()
	rr, req := newTestRequest("POST", "/gallery/api/upload", []byte("nothing"))
	router.ServeHTTP(rr, req)
	assertResponse(t, rr, "not logged in", 403)
}

func givenACookieAndRequest() (*http.Cookie, *web.Request) {
	c := http.Cookie{
		Name:  Cfg("cookieName"),
		Value: "7LlfIz+szcULggktMm5Py18sdIBmLYHj2DnztuIdAj2SZWfEva2peRsRGeE4HI3jCBCUKHVVyewlCoRE885ZqpYKiKX5SQn5cDOq+BRNdZmb1UXiO5iqQo4P61yAK5j7XKJbVUIlgjwj4Qil/gFCVVbK6/9LSUSogJGsA150rMgwCzQI+0b/s2Fdt5wUW3POdjvxqzrNmeePP1cfYe4gAg==",
	}
	httpR, _ := http.NewRequest("GET", "/", bytes.NewBuffer([]byte("")))
	httpR.AddCookie(&c)
	return &c, &web.Request{
		Request: httpR,
	}
}
