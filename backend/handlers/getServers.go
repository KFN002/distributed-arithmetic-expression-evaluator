package handlers

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"html/template"
	"log"
	"net/http"
	"sort"
)

// HandleCurrentServers получение данных о серверах
func HandleCurrentServers(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(float64)

	log.Println(userID)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("static/assets/server_data.html")
	if err != nil {
		log.Println("Error parsing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	models.Servers.Servers.Mu.Lock()

	var serverList []*models.Server

	for _, server := range models.Servers.Servers.Servers {
		serverList = append(serverList, server)
	}

	sort.Slice(serverList, func(i, j int) bool {
		return serverList[i].ID < serverList[j].ID
	})

	models.Servers.Servers.Mu.Unlock()

	err = tmpl.Execute(w, serverList)
	if err != nil {
		log.Println("Error executing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
