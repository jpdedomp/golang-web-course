package main

import "log"

func main() {
	var myString string
	myString = "Green"

	log.Println("myString is set to", myString)
	changeUsingPointer(&myString) //when you need the reference to a variable use &
	log.Println("after calling the function, myString is set to", myString)
}

func changeUsingPointer(s *string) {
	log.Println("s is set to", s)
	log.Println("s is referencing to", *s) //when you need the content of the reference use *
	newValue := "Red"
	*s = newValue
}
