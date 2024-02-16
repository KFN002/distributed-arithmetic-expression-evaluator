package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {
	validPattern := `^[0-9+\-*/()]+$`
	input := strings.ReplaceAll("(2)/*+2", " ", "")
	match, err := regexp.MatchString(validPattern, input)
	if err != nil {
		log.Println("Error in regular expression matching:", err)
	}
	fmt.Println(match)
}
