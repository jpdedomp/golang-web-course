package main

import (
	"errors"
	"fmt"
	"net/http"
)

const portNumber = ":8080"

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {
	sum := addValues(2, 3)
	_, _ = fmt.Fprintf(w, "This is the about page and 2 + 3 is %d", sum)
}

// lower case makes a function private. If first letter is UpperCase it will be public
// addValues adds two integers and returns the sum
func addValues(x int, y int) int {
	return x + y
}

func Divide(w http.ResponseWriter, r *http.Request) {
	x := 100.0
	y := 10.0
	f, err := divideValues(x, y)
	if err != nil {
		fmt.Fprintf(w, "Cannot divide by 0")
		return
	}
	fmt.Fprintf(w, "%f divided by %f is %f", x, y, f)
}

func divideValues(x float64, y float64) (float64, error) {
	if y == 0 {
		err := errors.New("cannot divide by 0")
		return 0, err
	}
	result := x / y
	return result, nil
}

// entry point
func main() {
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			n, err := fmt.Fprintf(w, "Hello, world!")
			fmt.Println("Bytes written:", n)
			fmt.Println("Error:", err)
		})
	*/

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/divide", Divide)

	fmt.Printf("Starting application on port %s", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
