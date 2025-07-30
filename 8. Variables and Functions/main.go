package main

import "fmt"

func main() {
	fmt.Println("Hello world.")

	var whatToSay string
	var i int

	whatToSay = "Goodbye, cruel world"
	fmt.Println(whatToSay)

	i = 7
	fmt.Println("i is set to", i)
	whatWasSaid, theOtherThingThatWasSaid := saySomeThing()

	fmt.Println("The function returned", whatWasSaid, "and", theOtherThingThatWasSaid)
}

func saySomeThing() (string, string) {
	return "something", "another thing"
}
