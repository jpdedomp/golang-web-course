package main

import (
	"log"
	"time"
)

type User struct {
	FirstName   string //things you declare with capital letters are visible outside your package
	LastName    string
	PhoneNumber string
	Age         int
	BirthDate   time.Time
}

func main() {
	user := User{
		FirstName:   "Trevor",
		LastName:    "Sawler",
		PhoneNumber: "1 555 555 1212",
	}

	log.Println(user.FirstName, user.LastName, user.BirthDate)
}

/*
func main() {
	var s2 = "six"

	s := "eight" //careful with this syntaxis. := creates a NEW variable within the scope of the main function (different from
	//the variable declared above as a package function)

	log.Println("s is", s)
	log.Println("s2 is", s2)

	saySomething("xxx")
}

func saySomething(s string) (string, string) {
	log.Println("s from the saySomething function is", s) // Variable shadowing. Be careful with the variable naming
	return s, "World"
}
*/
