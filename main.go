// main.go
package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Exercise struct {
	Question string   `json:"question"`
	Items    []string `json:"items"`
	Answers  []string `json:"answers"`
}

type Answer struct {
	Question string `json:"question"` 
	Answer   string `json:"answer"`
}

var exercises []Exercise
var answers []Answer

func loadExercises() {
	jsonFile, err := os.Open("exercises.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &exercises)
	if err != nil {
		log.Fatalf("Failed to decode the JSON file: %v", err)
	}

	log.Println("Successfully loaded exercises")
}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse the HTML template: %v", err)
	}

	err = tmpl.Execute(w, exercises)
	if err != nil {
		log.Fatalf("Failed to execute the HTML template: %v", err)
	}

	log.Println("Served exercises to a client")
}

func handleAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	question := r.FormValue("question")
	answer := r.FormValue("answer")

	answers = append(answers, Answer{Question: question, Answer: answer})

	jsonData, err := json.Marshal(answers)
	if err != nil {
		log.Fatalf("Failed to encode answers to JSON: %v", err)
	}

	err = os.WriteFile("answers.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write answers to JSON file: %v", err)
	}

	log.Println("Stored answer")
}

func main() {
	loadExercises()

	http.HandleFunc("/", serveExercises)
	http.HandleFunc("/answer", handleAnswer)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
