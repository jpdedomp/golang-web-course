package main

//channels are unique to go
//means of sending info to another very easily

import (
	"fmt"
	"log"

	"github.com/jp/myprogram/helpers"
)

func main() {
	log.Println("Welcome")

	var myVar helpers.SomeType
	myVar.TypeName = "Some name"

	log.Println(myVar.TypeName)

	results := make(chan string)
	defer close(results) //defers to when main is finished. Good practice to close channels

	for i := 0; i < 5; i++ {
		go func(i int) {
			results <- fmt.Sprintf("Done %d", i)
		}(i)
	} //we are launching several co routines... each race

	for i := 0; i < 5; i++ {
		fmt.Println(<-results) //the channel makes sure that info is collected when send
	}
}
