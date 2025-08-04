package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	//create template cache
	templateCache, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	//get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println(err) //this tells me that the error comes from the value stored in the map
	}

	//render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
	cacheMap := map[string]*template.Template{}

	//we first parse the pages and then the layouts associated

	//we review in the file system the files
	//get all of the files named *.page.tmpl from ./templates

	pages_paths, err := filepath.Glob("./templates/*.page.tmpl") //creates a slice of strings with the file paths of the pages
	if err != nil {
		return cacheMap, err
	}

	//range through all files ending with *.page.tmpl
	for _, page_path := range pages_paths {
		//get the actual filename, without the path
		filename := filepath.Base(page_path)

		template_set, err := template.New(filename).ParseFiles(page_path)
		if err != nil {
			return cacheMap, err
		}

		//creates a slice of strings with the file paths of the layouts
		layout_pages, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return cacheMap, err
		}

		//if there are any layouts, parse them as well
		if len(layout_pages) > 0 {
			template_set, err = template_set.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return cacheMap, err
			}
		}

		cacheMap[filename] = template_set //we add the parsed template set to the cache map
	}
	return cacheMap, nil
}

//previous version
/*
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
*/
