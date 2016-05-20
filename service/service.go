package service

import (
	"github.com/gocraft/web"
	"github.com/smancke/gogallery/imglib"
	"net/http"
)

var globalLibObject *imglib.ImageLibrary

type Context struct {
}

func (cntx *Context) Lib() *imglib.ImageLibrary {
	return globalLibObject
}

func Handler(lib *imglib.ImageLibrary) http.Handler {
	globalLibObject = lib
	return createRouter()
}

func createRouter() *web.Router {
	router := web.New(Context{})

	router.
		Middleware(web.LoggerMiddleware).
		//Middleware(web.ShowErrorsMiddleware).
		Middleware(web.StaticMiddleware(Cfg("galleryDir"), web.StaticOption{Prefix: "/gallery/image"})).
		Middleware(web.StaticMiddleware(Cfg("htmlDir"), web.StaticOption{Prefix: "/gallery/ui"}))

	router.Subrouter(Api{}, "/gallery/api").
		Get("/images", (*Api).GetImages)

	router.Subrouter(UserApi{}, "/gallery/api").
		Middleware((*UserApi).UserRequired).
		Post("/upload", (*UserApi).UploadImage).
		Get("/myImages", (*UserApi).GetMyImages).
		Delete("/myImages/:id", (*UserApi).DeleteImage).
		Get("/myData", (*UserApi).GetMyData).
		Post("/myData", (*UserApi).UpdateMyData)

	return router
}
