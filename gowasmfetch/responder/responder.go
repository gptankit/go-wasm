package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var (
	listen = flag.String("listen", ":9191", "listen address")
)

func main() {

	flag.Parse()

	log.Printf("listening on %q...", *listen)

	smux := http.NewServeMux()
	smux.HandleFunc("/fetchme", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(`{"Data": "OK!"}`))
	})

	// adding cors headers to the server responses
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8181"},
		AllowedMethods: []string{"GET", "POST"},
	})

	handler := crs.Handler(smux)
	log.Fatal(http.ListenAndServe(*listen, handler))
}
