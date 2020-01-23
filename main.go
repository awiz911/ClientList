package main

import (
	"Learning/ClientList/dbconfig"
	"Learning/ClientList/muxes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", dbconfig.DbUser, dbconfig.DbPass, dbconfig.DbName, dbconfig.Port)

	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	log.Printf("Postgres started at %s PORT", dbconfig.Port)

	defer db.Close()

	s := &http.Server{
		Addr:           ":3000",
		Handler:        muxes.SERVE(db),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
