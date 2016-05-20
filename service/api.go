package service

import (
	"github.com/gocraft/web"
	"net/http"
)

type Api struct {
	*Context
}

func (api *Api) GetImages(rw web.ResponseWriter, req *web.Request) {
	imageList, err := api.Lib().GetImages()
	if err != nil {
		sendError(rw, "error on fetching the images: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendAsJson(rw, imageList, http.StatusOK)
}
