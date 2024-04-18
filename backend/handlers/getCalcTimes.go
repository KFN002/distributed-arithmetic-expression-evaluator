package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
)

func HandleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(float64)

	if r.Method == http.MethodPost {
		timeValues := map[int]int{}
		timeFields := []string{"1", "2", "3", "4"}
		for _, field := range timeFields {
			timeStr := r.FormValue(field)
			timeValue, err := strconv.Atoi(timeStr)
			if err != nil || timeValue <= 0 {
				http.Error(w, fmt.Sprintf("Invalid input for %s", field), http.StatusBadRequest)
				return
			}
			operationID, _ := strconv.Atoi(field)
			timeValues[operationID] = timeValue
		}

		for id, value := range timeValues {
			err := databaseManager.DB.UpdateOperationTime(value, cacheMaster.OperatorByID[id], int(userID))
			if err != nil {
				log.Println("Error updating database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			cacheMaster.OperationCache.Set(int(userID), id-1, value)
		}

		http.Redirect(w, r, "/change-calc-time", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("static/assets/edit_time.html")
	if err != nil {
		log.Println("Error parsing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationTimes := cacheMaster.OperationCache.GetList(int(userID))

	fmt.Println(operationTimes)

	if len(operationTimes) == 0 {
		operationTimes, err = databaseManager.DB.GetTimes(int(userID))
		if err != nil {
			log.Println("Error fetching times", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		go cacheMaster.OperationCache.SetList(int(userID), operationTimes)
	}

	operationsData := models.OperationTimes{
		Time1: operationTimes[0],
		Time2: operationTimes[1],
		Time3: operationTimes[2],
		Time4: operationTimes[3],
	}

	err = tmpl.Execute(w, operationsData)
	if err != nil {
		log.Println("Error executing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
