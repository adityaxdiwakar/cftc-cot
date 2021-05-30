package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

type Product struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func createProduct(relation string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != conf.AccessToken {
		encode("You are not authorized to create products", w, 401)
		return
	}

	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		encode("Malformed body, refer to documentation to see how to add new product!", w, 400)
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (id, name) values ($1, $2)", relation), p.Id, p.Name)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			encode("Data already exists in database", w, 409)
		} else {
			encode("Could not insert data, server error", w, 500)
		}
		return
	}

	encode(p, w, 201)
	return

}

func disaggregatedProductsPost(w http.ResponseWriter, r *http.Request) {
	createProduct("disaggregatedprods", w, r)
}

func financialProductsPost(w http.ResponseWriter, r *http.Request) {
	createProduct("financialprods", w, r)
}
