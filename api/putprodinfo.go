package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
)

func disaggregatedProductInfoPut(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != conf.AccessToken {
		encode("You are not authorized to create products", w, 401)
		return
	}

	name := checkForProductInfo("disaggregated", w, r)
	if name == nil {
		return
	}

	id := chi.URLParam(r, "productId")

	var p DisaggregatedProductEntry
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		encode("Malformed body, refer to documentation to see how to add a new entry!", w, 400)
		return
	}

	_, err = db.Exec("INSERT INTO disaggregated (id, product, date, prodlong, prodshort"+
		", swaplong, swapshort, mmlong, mmshort, otherlong, othershort) values "+
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		fmt.Sprintf("%s-%d", id, p.Date),
		id,
		p.Date, p.ProdLong, p.ProdShort,
		p.SwapLong, p.SwapShort,
		p.MMLong, p.MMShort,
		p.OtherLong, p.OtherShort)

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

func financialProductInfoPut(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != conf.AccessToken {
		encode("You are not authorized to create products", w, 401)
		return
	}

	name := checkForProductInfo("financial", w, r)
	if name == nil {
		return
	}

	id := chi.URLParam(r, "productId")

	var p FinancialProductEntry
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		encode("Malformed body, refer to documentation to see how to add a new entry!", w, 400)
		return
	}

	_, err = db.Exec("INSERT INTO financial (id, product, date, deallong, dealshort"+
		", assetlong, assetshort, levlong, levshort, otherlong, othershort)"+
		" values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		fmt.Sprintf("%s-%d", id, p.Date),
		id,
		p.Date, p.DealLong, p.DealShort,
		p.AssetLong, p.AssetShort,
		p.LevLong, p.LevShort,
		p.OtherLong, p.OtherShort)

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
