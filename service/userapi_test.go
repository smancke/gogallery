package service

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/smancke/gogallery/imglib"
	"github.com/gocraft/web"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestUploadImage(t *testing.T) {
	setTestUsername()
	globalLibObject = openTempDB()
	router := createRouter()

	// upload image
	rr, req := newTestRequest("POST", "/gallery/api/upload", createDummyImage())
	router.ServeHTTP(rr, req)
	assertResponse(t, rr, ".jpeg", 201)

	// upload another image
	rr, req = newTestRequest("POST", "/gallery/api/upload", createDummyImage())
	router.ServeHTTP(rr, req)
	assertResponse(t, rr, ".jpeg", 201)

	// read the image list
	assertImageList(t, rr, req, router, "/gallery/api/myImages", 2)
	imageList := assertImageList(t, rr, req, router, "/gallery/api/images", 2)

	// delete one image
	rr, req = newTestRequest("DELETE", fmt.Sprintf("/gallery/api/myImages/%v", imageList[0].ID), []byte(""))
	router.ServeHTTP(rr, req)
	assertResponse(t, rr, "deleted", 200)

	assertImageList(t, rr, req, router, "/gallery/api/myImages", 1)
}

func TestMyData(t *testing.T) {
	setTestUsername()
	globalLibObject = openTempDB()
	router := createRouter()

	// changeMyData
	rr, req := newTestRequest("POST", "/gallery/api/myData",  []byte(`{"nickName": "BenCoolUtzer", "link": "http://nowhere"}`) )
	router.ServeHTTP(rr, req)
	assertResponse(t, rr, "updated", 200)

	// check the data
	rr, req = newTestRequest("GET", "/gallery/api/myData",  []byte("") )
	router.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	user := imglib.User{}
	marshallErr := json.Unmarshal(rr.Body.Bytes(), &user)
	assert.NoError(t, marshallErr)
	assert.Equal(t, "BenCoolUtzer", user.NickName)
	assert.Equal(t, "http://nowhere", user.Link)
}

func assertImageList(t *testing.T, rr *httptest.ResponseRecorder, req *http.Request, router *web.Router, url string, imgCount int) ([]*imglib.Image) {
	rr, req = newTestRequest("GET", url, []byte{})
	router.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	imageList := make([]*imglib.Image, 0, 0)
	err := json.Unmarshal(rr.Body.Bytes(), &imageList)
	assert.Nil(t, err)
	assert.Equal(t, imgCount, len(imageList))
	assert.Equal(t, testOverwriteUsername, imageList[0].User.UserName)
	return imageList
}

func setTestUsername() {
	rand.Seed(time.Now().Unix())
	log.Printf("testOverwriteUsername int: %v", rand.Int63())
	sum := md5.Sum([]byte(fmt.Sprintf("%v", rand.Int63())))
	hash := base64.URLEncoding.EncodeToString(sum[10:])
	testOverwriteUsername = "user-" + string(hash)
	log.Printf("testOverwriteUsername: %v", testOverwriteUsername)
}

func assertResponse(t *testing.T, rr *httptest.ResponseRecorder, bodySubstring string, code int) {
	if gotBody := string(rr.Body.Bytes()); !strings.Contains(gotBody, bodySubstring) {
		t.Errorf("assertResponse: expected body to contain '%s' but got '%s'. (caller: %s)", bodySubstring, gotBody, callerInfo())
	}
	if code != rr.Code {
		t.Errorf("assertResponse: expected code to be '%d' but got '%d'. (caller: %s)", code, rr.Code, callerInfo())
	}
}

func newTestRequest(method, path string, body []byte) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	return recorder, request
}

func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func openTempDB() *imglib.ImageLibrary {
	file, _ := ioutil.TempDir("", "galleryDir.")

	db := &imglib.ImageLibrary{}
	if err := db.Open(file); err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func createDummyImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	col := color.RGBA{255, 0, 0, 255} // Red
	for i := 100; i < 500; i++ {
		for j := 100; j < 500; j++ {
			img.Set(i, j, col)
		}
	}
	buff := new(bytes.Buffer)
	jpeg.Encode(buff, img, nil)
	return buff.Bytes()
}
