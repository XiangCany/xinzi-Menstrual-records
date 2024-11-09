package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Period struct {
	ID        string
	StartDate time.Time
	EndDate   time.Time
	Symptoms  string
	Mood      string
	Weight    float64
	CreatedAt time.Time
}

var periods []Period

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/stats", statsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	startDateStr := r.FormValue("startDate")
	endDateStr := r.FormValue("endDate")
	symptoms := r.FormValue("symptoms")
	mood := r.FormValue("mood")
	weightStr := r.FormValue("weight")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		weight = 0
	}

	period := Period{
		ID:        uuid.New().String(),
		StartDate: startDate,
		EndDate:   endDate,
		Symptoms:  symptoms,
		Mood:      mood,
		Weight:    weight,
		CreatedAt: time.Now(),
	}

	periods = append(periods, period)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var cycleLengths []int
	for _, period := range periods {
		if len(periods) > 1 {
			prevPeriod := periods[len(periods)-2]
			cycleLength := int(period.StartDate.Sub(prevPeriod.EndDate).Seconds() / (24 * 60 * 60))
			cycleLengths = append(cycleLengths, cycleLength)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/stats.html"))
	tmpl.Execute(w, struct {
		CycleLengths []int
	}{
		CycleLengths: cycleLengths,
	})
}
