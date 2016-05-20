package service

import (
	"github.com/gocraft/web"
	"log"
	"net/http"
	"runtime"
	"fmt"
)

func ok(rw web.ResponseWriter, message string) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, message)
}

func sendError(rw web.ResponseWriter, message string, code int) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	log.Printf("%v: %v %v", f.Name(), code, message)
	http.Error(rw, message, code)
}
