package service

import (
	"github.com/gocraft/web"
	"net/http"
	"strconv"
)

type Api struct {
	*Context
}

func (api *Api) GetImages(rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 5
	}

	// in case of a conversion error, 0 is the default
	offset, _ := strconv.Atoi(req.FormValue("offset"))

	imageList, err := api.Lib().GetImages(offset, limit)
	if err != nil {
		sendError(rw, "error on fetching the images: "+err.Error(), http.StatusInternalServerError)
		return
	}
	sendAsJson(rw, imageList, http.StatusOK)
}
