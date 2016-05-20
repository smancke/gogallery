package main

import (
	"github.com/smancke/gogallery/imglib"
	"github.com/smancke/gogallery/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	lib := openLib()

	s := &http.Server{
		Addr:    service.Cfg("address"),
		Handler: service.Handler(lib)}

	go func() {
		log.Printf("starting up at %v", service.Cfg("address"))
		log.Fatal(s.ListenAndServe())
		os.Exit(1)
	}()

	waitForTermination(func() {

		lib.Close()

	})
}

func waitForTermination(callback func()) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Got singal '%v' .. exit more or less greacefully now", <-sigc)
	callback()
	log.Printf("Done.")
	os.Exit(0)
}

func openLib() *imglib.ImageLibrary {
	db := &imglib.ImageLibrary{}
	if err := db.Open(service.Cfg("galleryDir")); err != nil {
		log.Fatal(err.Error())
	}
	return db
}
