package main

import (
    "flag"
    "log"
    "net/http"
	"github.com/NYTimes/gziphandler"
)

var (
    listen = flag.String("listen", ":8181", "listen address")
    dir    = flag.String("dir", ".", "directory to serve")
)

func main() {

    flag.Parse()
    log.Printf("listening on %q...", *listen)
    log.Fatal(http.ListenAndServe(*listen, gziphandler.GzipHandler(http.FileServer(http.Dir(*dir)))))
}
