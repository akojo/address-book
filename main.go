package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bunrouter"
)

func createDb() *bun.DB {
	config, err := pgx.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	sqldb := stdlib.OpenDB(*config)
	return bun.NewDB(sqldb, pgdialect.New())
}

func main() {
	for _, e := range os.Environ() {
		log.Printf("%s", e)
	}
	db := createDb()

	router := bunrouter.New()
	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		fmt.Fprintln(w, "hello, cloud")
		return nil
	})
	router.GET("/health", func(w http.ResponseWriter, req bunrouter.Request) error {
		if err := db.PingContext(req.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return err
		}
		fmt.Fprintln(w, "ok")
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println(http.ListenAndServe(":"+port, router))
}
