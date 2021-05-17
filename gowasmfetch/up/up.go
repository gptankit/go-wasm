package main

import (
	"flag"
	"fmt"
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
		fmt.Println(req.Header)
		fmt.Println(req.URL.Hostname)
		fmt.Println(req.Host)
		fmt.Println(req.Method)
		fmt.Println(req.RequestURI)
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(`{"TestServer": "OK!"}`))
	})

	// adding cors headers to the server responses
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8181"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
	})

	handler := crs.Handler(smux)
	log.Fatal(http.ListenAndServe(*listen, handler))
}
