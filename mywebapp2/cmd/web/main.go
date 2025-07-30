package main

import (
	"fmt"
	"net/http"

	"github.com/jpdedomp/golang-web-course/mywebapp2/pkg/handlers"
)

const portNumber = ":8080"

// entry point
func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("Starting application on port %s", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
