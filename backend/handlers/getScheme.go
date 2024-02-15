package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func HandleGetScheme(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/scheme.html")
	if err != nil {
		log.Println("Error parsing expression_by_id.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodGet {
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Error executing expression_by_id.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
