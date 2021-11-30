package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bunrouter"
)

func createDb() *bun.DB {
	network := "postgres"
	if strings.HasSuffix(os.Getenv("PGHOST"), ".s.PGSQL.5432") {
		network = "unix"
	}
	insecure := false
	if os.Getenv("PGSSLMODE") == "disable" {
		insecure = true
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork(network),
		pgdriver.WithInsecure(insecure),
		pgdriver.WithAddr(os.Getenv("PGHOST")),
		pgdriver.WithDatabase(os.Getenv("PGDATABASE")),
		pgdriver.WithUser(os.Getenv("PGUSER")),
		pgdriver.WithPassword(os.Getenv("PGPASSWORD")),
	))
	return bun.NewDB(sqldb, pgdialect.New())
}

func main() {
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
