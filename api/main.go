package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var conf tomlConfig
var db *sql.DB

type tomlConfig struct {
	Database    postgresCreds
	AccessToken string `toml:"access_token"`
}

type postgresCreds struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatalf("error: could not parse config file: %v\n", err)
	}

	pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User,
		conf.Database.Password, conf.Database.DBName)

	var err error
	db, err = sql.Open("postgres", pSqlInfo)
	if err != nil {
		log.Fatalf("error: could not open database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error: could not ping database: %v\n", err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homePage)
	r.Get("/disaggregated", disaggregatedProducts)
	r.Post("/disaggregated", disaggregatedProductsPost)
	r.Get("/financials", financialProducts)
	r.Post("/financials", financialProductsPost)

	http.ListenAndServe(":3000", r)
}

type StringResponse struct {
	Payload string `json:"payload"`
	Code    int    `json:"status_code"`
}

type PayloadResponse struct {
	Payload interface{} `json:"payload"`
	Code    int         `json:"status_code"`
}

func encode(data interface{}, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := PayloadResponse{
		Payload: data,
		Code:    statusCode,
	}
	str, _ := json.MarshalIndent(resp, "", "    ")
	str = bytes.Replace(str, []byte("\\u003c"), []byte("<"), -1)
	str = bytes.Replace(str, []byte("\\u003e"), []byte(">"), -1)
	str = bytes.Replace(str, []byte("\\u0026"), []byte("&"), -1)
	w.Write(str)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	encode("Welcome to the Commitment of Traders API by Aditya Diwakar - Not affiliated with any government organization", w, 200)
}

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
