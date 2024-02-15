package handlers

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// HandleChangeCalcTime страница изменения времени выполнения
func HandleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		timeValues := map[string]int{}
		timeFields := []string{"1", "2", "3", "4"}
		for _, field := range timeFields {
			timeStr := r.FormValue(field)
			timeValue, err := strconv.Atoi(timeStr)
			if err != nil || timeValue <= 0 {
				http.Error(w, fmt.Sprintf("Invalid input for %s", field), http.StatusBadRequest)
				return
			}
			timeValues[field] = timeValue
		}

		for id, value := range timeValues {
			_, err := databaseManager.DB.Exec("UPDATE operations SET time=? WHERE id=?", value, id)
			if err != nil {
				log.Println("Error updating database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/change-calc-time", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("static/assets/edit_time.html")
	if err != nil {
		log.Println("Error parsing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationTimes, err := databaseManager.GetTimes() // получение данных время операций
	if err != nil {
		log.Println("Error fetching times", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationsData := models.OperationTimes{
		Time1: operationTimes[0],
		Time2: operationTimes[1],
		Time3: operationTimes[2],
		Time4: operationTimes[3],
	}

	go cacheMaster.OperationCache.SetList(operationTimes)

	err = tmpl.Execute(w, operationsData)
	if err != nil {
		log.Println("Error executing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
