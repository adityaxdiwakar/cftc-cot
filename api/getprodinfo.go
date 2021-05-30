package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type DisaggregatedProductEntry struct {
	Date       int `json:"date,omitempty"`
	ProdLong   int `json:"prod_long"`
	ProdShort  int `json:"prod_short"`
	SwapLong   int `json:"swap_long"`
	SwapShort  int `json:"swap_short"`
	MMLong     int `json:"mm_long"`
	MMShort    int `json:"mm_short"`
	OtherLong  int `json:"other_long"`
	OtherShort int `json:"other_short"`
}

type FinancialProductEntry struct {
	Date       int `json:"date,omitempty"`
	DealLong   int `json:"deal_long"`
	DealShort  int `json:"deal_short"`
	AssetLong  int `json:"asset_long"`
	AssetShort int `json:"asset_short"`
	LevLong    int `json:"lev_long"`
	LevShort   int `json:"lev_short"`
	OtherLong  int `json:"other_long"`
	OtherShort int `json:"other_short"`
}

type AllProductEntries struct {
	Name    string      `json:"name"`
	Entries interface{} `json:"entries"`
}

func checkForProductInfo(relation string, w http.ResponseWriter, r *http.Request) *string {
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM %sprods WHERE id=$1", relation),
		chi.URLParam(r, "productId"))
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			encode("Product not found in database", w, 404)
		} else {
			encode("Error when querying for product, server error", w, 500)
		}
		return nil
	}

	return &name
}

func getProductInfo(relation string, w http.ResponseWriter, r *http.Request) (*string, *sql.Rows) {
	name := checkForProductInfo(relation, w, r)
	if name == nil {
		return nil, nil
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE product=$1", relation),
		chi.URLParam(r, "productId"))

	if err != nil {
		encode("An issue occured when querying the products", w, 500)
		return nil, nil
	}

	return name, rows
}

func disaggregatedProductInfo(w http.ResponseWriter, r *http.Request) {
	name, rows := getProductInfo("disaggregated", w, r)
	if rows == nil {
		return
	}

	entries := map[int]DisaggregatedProductEntry{}
	var void string
	var date int
	for rows.Next() {
		entry := DisaggregatedProductEntry{}
		rows.Scan(&void, &void, &entry.Date, &entry.ProdLong, &entry.ProdShort,
			&entry.SwapLong, &entry.SwapShort, &entry.MMLong, &entry.MMShort,
			&entry.OtherLong, &entry.OtherShort)
		date = entry.Date
		entry.Date = 0
		entries[date] = entry
	}

	resp := AllProductEntries{
		Name:    *name,
		Entries: entries,
	}

	encode(resp, w, 200)
}

func financialProductInfo(w http.ResponseWriter, r *http.Request) {
	name, rows := getProductInfo("financial", w, r)
	if rows == nil {
		return
	}

	entries := map[int]FinancialProductEntry{}
	var void string
	var date int
	for rows.Next() {
		entry := FinancialProductEntry{}
		rows.Scan(&void, &void, &date, &entry.DealLong, &entry.DealShort,
			&entry.AssetLong, &entry.AssetShort, &entry.LevLong, &entry.LevShort,
			&entry.OtherLong, &entry.OtherShort)
		entries[date] = entry
	}

	resp := AllProductEntries{
		Name:    *name,
		Entries: entries,
	}

	encode(resp, w, 200)
}
