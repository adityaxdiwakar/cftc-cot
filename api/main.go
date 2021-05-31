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
	if _, err := toml.DecodeFile("config/config.toml", &conf); err != nil {
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
	r.NotFound(notFoundHandler)

	// disaggregated
	r.Get("/disaggregated", disaggregatedProducts)
	r.Post("/disaggregated", disaggregatedProductsPost)
	r.Get("/disaggregated/{productId}", disaggregatedProductInfo)
	r.Put("/disaggregated/{productId}", disaggregatedProductInfoPut)

	// financials
	r.Get("/financials", financialProducts)
	r.Post("/financials", financialProductsPost)
	r.Get("/financials/{productId}", financialProductInfo)
	r.Put("/financials/{productId}", financialProductInfoPut)

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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	encode("This is not a valid route, please refer to the documentation", w, 404)
}
