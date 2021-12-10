package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

const (
	width     = 64
	height    = 64
	tickspeed = time.Second
)

type data struct {
	W int
	H int
}

func main() {
	t := template.Must(template.ParseFiles("index.tmpl"))
	b := new(bytes.Buffer)

	d := data{width, height}

	err := t.Execute(b, d)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		fmt.Fprint(w, b)
	})

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("client connected: %v\n", r.RemoteAddr)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		f := w.(http.Flusher)

		for result := range newGame(r.Context()) {
			fmt.Fprintf(w, "data: %s\n\n", result)
			f.Flush()
		}

		fmt.Printf("client disconnected: %v\n", r.RemoteAddr)

	})

	fmt.Println("listen in port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
