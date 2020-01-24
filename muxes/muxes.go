package muxes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type (
	//Client - record for client
	Client struct {
		ID       int
		Name     string
		Lastname string
		Age      int
		Cell     string
		Email    string
	}
)

//Clients - slice of clients
type Clients = []Client

//SERVE - made to serve
func SERVE(db *sql.DB) *http.ServeMux {

	log.Print("Server started at http://127.0.0.1:3000 port.")

	mux := http.NewServeMux()

	mux.HandleFunc("/newClient", func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			http.NotFound(w, req)
			return
		}

		var newClient = convertRequestToClient(req)
		var lastInsertID int

		err := db.QueryRow("INSERT INTO clients(name, lastname, age, cell, email) VALUES($1, $2, $3, $4, $5) returning id;",
			newClient.Name, newClient.Lastname, newClient.Age, newClient.Cell, newClient.Email).Scan(&lastInsertID)

		checkErr(err)

		fmt.Println("last inserted id =", lastInsertID)

		okStatus(w)

		log.Printf("New Client %s %s added successfully.", newClient.Name, newClient.Lastname)

		json.NewEncoder(w).Encode(newClient)

		return
	})

	mux.HandleFunc("/getAll", func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			http.NotFound(w, req)
			return
		}

		okStatus(w)

		json.NewEncoder(w).Encode(getAllClients(db))

		log.Printf("All clients listed successfully.")

		return
	})

	mux.HandleFunc("/clients/", func(w http.ResponseWriter, req *http.Request) {

		var method = req.Method

		if method != "GET" && method != "DELETE" && method != "PUT" {
			http.NotFound(w, req)
			return
		}

		id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/clients/"))

		if err != nil {
			panic(err)
		}

		okStatus(w)

		if method == "GET" {

			rows, err := db.Query("SELECT * FROM clients where id=$1", id)

			checkErr(err)

			client := Client{}

			for rows.Next() {
				err = rows.Scan(&client.ID, &client.Name, &client.Lastname, &client.Age, &client.Cell, &client.Email)
				checkErr(err)
			}

			log.Printf("Client with id = %d listed successfully.", id)

			json.NewEncoder(w).Encode(client)

			return

		}

		if method == "DELETE" {

			stmt, err := db.Prepare("delete from clients where id=$1")

			checkErr(err)

			_, err = stmt.Exec(id)

			checkErr(err)

			log.Printf("Client with id = %d deleted successfully.", id)

			json.NewEncoder(w).Encode(nil)

			return
		}

		if method == "PUT" {

			client := convertRequestToClient(req)

			stmt, err := db.Prepare("update clients set name=$1, lastname=$2, age=$3, cell=$4, email=$5 where id=$6")

			checkErr(err)

			_, err = stmt.Exec(client.Name, client.Lastname, client.Age, client.Cell, client.Email, id)

			checkErr(err)

			log.Printf("Client with id = %d updated successfully.", id)

			json.NewEncoder(w).Encode(client)

			return
		}

		json.NewEncoder(w).Encode(nil)

		return

	})

	return mux
}

func getAllClients(db *sql.DB) Clients {

	rows, err := db.Query("SELECT * FROM clients")

	checkErr(err)

	client := Client{}
	allClients := Clients{}

	for rows.Next() {
		err = rows.Scan(&client.ID, &client.Name, &client.Lastname, &client.Age, &client.Cell, &client.Email)
		allClients = append(allClients, client)
		checkErr(err)
	}

	return allClients
}

func okStatus(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	return
}

func convertRequestToClient(req *http.Request) Client {

	body, err := ioutil.ReadAll(req.Body)

	checkErr(err)

	var newClient Client

	err = json.Unmarshal(body, &newClient)

	checkErr(err)

	return newClient
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
