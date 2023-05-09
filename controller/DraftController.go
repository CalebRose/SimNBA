package controller

import (
	"fmt"
	"net/http"

	"github.com/CalebRose/SimNBA/managers"
)

func GenerateDraftGrades(w http.ResponseWriter, r *http.Request) {
	managers.GenerateDraftLetterGrades()
	fmt.Println(w, "Congrats, you generated the Letter Grades!")
}
