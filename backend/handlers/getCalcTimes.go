package handlers

import (
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func HandleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(float64)

	if r.Method == http.MethodPost {
		timeValues := map[int]int{}

		for i := 1; i <= 4; i++ {
			field := strconv.Itoa(i)
			timeStr := r.FormValue(field)
			timeValue, err := strconv.Atoi(timeStr)
			if err != nil || timeValue <= 0 {
				http.Error(w, fmt.Sprintf("Invalid input for %s", field), http.StatusBadRequest)
				return
			}
			timeValues[i] = timeValue
		}

		log.Println("Got from input form", timeValues)

		for id, value := range timeValues {
			err := databaseManager.DB.UpdateOperationTime(value, cacheMaster.OperatorByID[id-1], int(userID))
			if err != nil {
				log.Println("Error updating database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			cacheMaster.OperationCache.Set(int(userID), id-1, value)
		}

		log.Println("After post cache:", cacheMaster.OperationCache.GetList(int(userID)))

		http.Redirect(w, r, "/change-calc-time", http.StatusSeeOther)
		return

	} else if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("static/assets/edit_time.html")
		if err != nil {
			log.Println("Error parsing edit_time.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		operationTimes := cacheMaster.OperationCache.GetList(int(userID))

		if len(operationTimes) == 0 {
			operationTimes, err = databaseManager.DB.GetTimes(int(userID))
			if err != nil {
				log.Println("Error fetching times", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			cacheMaster.OperationCache.SetList(int(userID), operationTimes)
		}

		log.Println("Got from cache in get request:", operationTimes)

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
}
