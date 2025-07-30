package main

import "log"

type myStruct struct {
	FirstName string
}

func (m *myStruct) printFirstName() string { //this is a receiver, and it associates the function to the structure
	log.Println(m.FirstName) //this way you can access the information of the structure
	return m.FirstName
}

func main() {
	var myVar myStruct
	myVar.FirstName = "John"

	myVar2 := myStruct{
		FirstName: "Mary",
	}

	myVar.printFirstName()
	myVar2.printFirstName()
}
