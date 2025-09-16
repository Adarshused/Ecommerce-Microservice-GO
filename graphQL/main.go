package main

import (
	"log"
	"net/http"
	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string 	`envconfig:"ACCOUNT_SERVICE_URL"`
	catalogURL string	`envconfig:"CATALOG_SERVICE_URL"`
	orderURL string 	`envconfig:"ORDER_SERVICE_URL"`
}


// func main() {

// 	var cfg AppConfig

// }