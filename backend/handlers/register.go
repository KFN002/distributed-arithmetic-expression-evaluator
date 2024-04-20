package handlers

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"html/template"
	"log"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/register.html")
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
	} else if r.Method == http.MethodPost {

		login := r.FormValue("username")
		password := r.FormValue("password")

		err := databaseManager.SignUpUser(login, password)

		message := models.CreateNewTemplateMessage()

		if err != nil {
			message.ChangeMessage(err.Error())

		} else {
			message.ChangeMessage("Sign up successful!")
		}

		err = tmpl.Execute(w, message)
		if err != nil {
			log.Println("Error executing create_expression.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
