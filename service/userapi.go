package service

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"github.com/smancke/gogallery/imglib"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserApi struct {
	*Context
	user *imglib.User
}

var testOverwriteUsername = ""

func (userApi *UserApi) UploadImage(rw web.ResponseWriter, req *web.Request) {
	image, err := userApi.Lib().CreateImage(*userApi.user, req.Body)
	if err != nil {
		sendError(rw, "error on processing the image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendAsJson(rw, image, http.StatusCreated)
}

func (userApi *UserApi) GetMyImages(rw web.ResponseWriter, req *web.Request) {
	imageList, err := userApi.Lib().GetImagesByUsername(userApi.user.UserName)
	if err != nil {
		sendError(rw, "error on fetching the images: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendAsJson(rw, imageList, http.StatusOK)
}

func (userApi *UserApi) DeleteImage(rw web.ResponseWriter, req *web.Request) {
	id, err := strconv.Atoi(req.PathParams["id"])
	if err != nil {
		sendError(rw, fmt.Sprintf("invalid id '%v': %v", req.PathParams["id"], err),
			http.StatusBadRequest)
	}
	err = userApi.Lib().DeleteImage(userApi.user.ID, uint(id))
	if err != nil {
		sendError(rw, "error on deletion of image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ok(rw, "deleted")
}

func (userApi *UserApi) GetMyData(rw web.ResponseWriter, req *web.Request) {
	user, err := userApi.Lib().UserByUsername(userApi.user.UserName)
	if err != nil {
		sendError(rw, "error retrieving user data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendAsJson(rw, user, http.StatusOK)
}

func (userApi *UserApi) UpdateMyData(rw web.ResponseWriter, req *web.Request) {
	userFromRequest := imglib.User{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendError(rw, "error on reading request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &userFromRequest); err != nil {
		sendError(rw, "error parsing user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	userApi.user.NickName = userFromRequest.NickName
	userApi.user.Link = userFromRequest.Link
	if err := userApi.Lib().SaveUser(userApi.user); err != nil {
		sendError(rw, "error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ok(rw, "updated")
}

func sendAsJson(rw web.ResponseWriter, oject interface{}, code int) {
	result, err := json.Marshal(oject)
	if err != nil {
		sendError(rw, "error on json marshaling: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(code)
	fmt.Fprintf(rw, string(result))
}
