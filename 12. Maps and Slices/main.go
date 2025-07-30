package main

import (
	"log"
	"sort"
)

type User struct {
	FirstName string
	LastName  string
}

func main() {
	//maps
	myMap := make(map[string]User)

	me := User{ //maps don't need pointers... they are immutable. You can't depend on them being in the same order
		FirstName: "Jean Paul",
		LastName:  "de Dompierre",
	}

	myMap["me"] = me

	log.Println("My first name is", myMap["me"].FirstName)
	log.Println("My last name is", myMap["me"].LastName)

	//slices (like an array)
	var myStringSlice []string
	myStringSlice = append(myStringSlice, "First appended string to my slice")
	myStringSlice = append(myStringSlice, "Second appended string to my slice")

	var myIntSlice []int
	myIntSlice = append(myIntSlice, 4)
	myIntSlice = append(myIntSlice, 2)
	myIntSlice = append(myIntSlice, 5)

	log.Println(myStringSlice)
	log.Println(myIntSlice)

	sort.Ints(myIntSlice) //you can sort slices
	log.Println(myIntSlice)

	otherNumbers := []int{1, 2, 3, 4, 5, 6}
	log.Println(otherNumbers[2:4]) //you can reference just a range of the slice

	names := []string{"Jean", "Peter", "Jonathan", "Trevor"}
	log.Println(names[2:2]) //beware of how the ranges work. 0:1 returns only the first element. If there are n elements 0:n then you get all the elements (not n-1, starts at 0 but ends at n)
}

/*
func main() {
	myMap := make(map[string]int) //the first value is the index type and the second is the stored data type

	myMap["First"] = 1
	myMap["Second"] = 2

	log.Println(myMap["First"])
	log.Println(myMap["Second"])
}


func main() {
	myMap := make(map[string]string) //the first value is the index type and the second is the stored data type

	myMap["dog"] = "Dipsy"

	log.Println(myMap["dog"])
}
*/
