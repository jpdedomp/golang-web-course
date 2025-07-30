package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tc = make(map[string]*template.Template) //map that will contain all our parsed templates... our cache

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	//check if we have the template already in our cache
	_, inMap := tc[t]
	if !inMap {
		//need to create the template
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		//we the template in the cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	//parse the templates
	tmpl, err := template.ParseFiles(templates...) //... adds all the strings from the slice consecutively
	if err != nil {
		return err
	}

	//add template to cache
	tc[t] = tmpl
	log.Println("added template to cache")

	return nil
}
