package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/uptrace/bunrouter"
)

func main() {
	router := bunrouter.New()
	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		fmt.Fprintln(w, "hello, cloud")
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println(http.ListenAndServe(":"+port, router))
}
