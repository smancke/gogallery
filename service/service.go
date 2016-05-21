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

	if Cfg("testOverwriteUsername") != "" {
		testOverwriteUsername = Cfg("testOverwriteUsername")
	}

	router.
		Middleware(web.LoggerMiddleware).
		//Middleware(web.ShowErrorsMiddleware).
		Middleware(web.StaticMiddleware(Cfg("galleryDir"), web.StaticOption{Prefix: "/gallery/image"})).
		Middleware(web.StaticMiddleware(Cfg("htmlDir"), web.StaticOption{Prefix: "/gallery/ui"})).
		Get("/", func(w web.ResponseWriter, r *web.Request) {
			w.Write([]byte(`<html>
  <body>
    <a href="/gallery/ui/pub/index.html">View gallery</a>
    <br><a href="/gallery/ui/user/index.html">Add images</a>
  </body>
</html>`))
		})

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
