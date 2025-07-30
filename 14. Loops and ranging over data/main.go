package main

import (
	"fmt"
	"log"
)

func main() {
	//traditional for loop
	for i := 0; i <= 10; i++ {
		log.Println(i)
	}

	//for ranging through a slice
	animals := []string{"dog", "fish", "horse", "cow", "cat"}
	for i, animal := range animals { //the index of the loop count i, and the value of that range (in this case it makes sense to call it animal (singular) as it is a value of a slice animals (plural))
		log.Println(i, animal) //you can also ignore i with this syntax for _, animal := range animals
	}

	//we can also range through a map
	animalNames := make(map[string]string)

	animalNames["dog"] = "Dipsy"
	animalNames["cat"] = "Theo"

	for animalType, name := range animalNames {
		log.Println(animalType, name)
	}

	//we can range through a string
	firstLine := "Everybody on the lot somebody"
	firstLine = "xxxx" //be careful... strings are immutable, this means that when you do this you destroy the previous and create a new one (so pointers would not work here)

	for i, l := range firstLine {
		log.Println(i, l) //a string is a slice o bytes
	}

	//range through custom type

	type User struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
	}

	var users []User
	users = append(users, User{"John", "Smith", "john@smith.com", 30})
	users = append(users, User{"Mary", "Jones", "mary@jones.com", 20})
	users = append(users, User{"Sally", "Brown", "sally@brown.com", 45})
	users = append(users, User{"Alex", "Anderson", "alex@anderson.com", 17})

	for _, l := range users {
		fmt.Println(l.FirstName, l.LastName, l.Email, l.Age)
	}
}
