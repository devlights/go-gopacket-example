package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var (
	appLog = log.New(os.Stderr, "", 0)
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "helloworld\n")
	})

	return http.ListenAndServe("localhost:12345", nil)
}