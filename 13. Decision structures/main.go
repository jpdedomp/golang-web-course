package main

import "log"

func main() {

	//if statements
	cat := "dog"

	if cat == "cat" {
		log.Println("Cat is cat")
	} else {
		log.Println("Cat is not cat")
	}

	myNum := 100
	isTrue := true

	if myNum > 99 && isTrue {
		log.Println("myNum is greater than 99 and isTrue is true")
	}
	/*
		if (logical equation)
		&&: and
		!: not
		==: is equal
		||: or
		else
		else if (logical equation)
	*/

	//switch

	myVar := "fish"

	switch myVar {
	case "cat":
		log.Println("myVar is set to cat")
	case "fish":
		log.Println("myVar is set to fish")
	case "dog":
		log.Println("myVar is set to dog")
	default:
		log.Println("myVar is set to something else")
	}
}
