package main

import (
	"fmt"
	"net/http"
)

func getProducts(relation string, w http.ResponseWriter) {
	rows, err := db.Query(fmt.Sprintf("SELECT * from %s", relation))
	defer rows.Close()
	fmt.Println(err)
	if err != nil {
		encode("An issue occured when querying the products", w, 500)
		return
	}

	products := make(map[string]string)

	if err != nil {
		encode("An issue occured when querying the products", w, 500)
		return
	}

	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)
		products[id] = name
	}

	encode(products, w, 200)
}

func disaggregatedProducts(w http.ResponseWriter, r *http.Request) {
	getProducts("disaggregatedprods", w)
}

func financialProducts(w http.ResponseWriter, r *http.Request) {
	getProducts("financialprods", w)
}
