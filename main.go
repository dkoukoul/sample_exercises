// main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PhonologicRhymePairItems struct {
	Words   []string `json:"words"`
	Correct int      `json:"correct"`
}
type PhonologicRhymeMatchItems struct {
	Words   []string `json:"words"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
}

type PhonologicRhymePair struct {
	Question                 string                     `json:"question"`
	Answer                   []string                   `json:"answer"`
	PhonologicRhymePairItems []PhonologicRhymePairItems `json:"items"`
}

type PhonologicRhymeMatch struct {
	Question                  string                      `json:"question"`
	PhonologicRhymeMatchItems []PhonologicRhymeMatchItems `json:"items"`
}

type PhonologicRhymeMultipleMatch struct {
	Question string         `json:"question"`
	Column1  []string       `json:"column1"`
	Column2  []string       `json:"items"`
	Answers  map[string]int `json:"answers"`
}

type PhonologicRhymeSentence struct {
	Question string                      `json:"question"`
	Items    []PhonologicRhymeMatchItems `json:"items"`
}

type Answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

var answers []Answer
var phonologicExercises []interface{}

type ExerciseData struct {
	PhonologicExercises []interface{}
}

func readJSONFile(filePath string, v interface{}) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	err = json.Unmarshal(file, v)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
}

func loadExercises() {
	var exercise1 PhonologicRhymePair
	readJSONFile("data/PhonologicRhymePair.json", &exercise1)
	phonologicExercises = append(phonologicExercises, exercise1)

	var exercise2 PhonologicRhymeMatch
	readJSONFile("data/PhonologicRhymeMatch.json", &exercise2)
	phonologicExercises = append(phonologicExercises, exercise2)

	var exercise3 PhonologicRhymeMultipleMatch
	readJSONFile("data/PhonologicRhymeMultipleMatch.json", &exercise3)
	phonologicExercises = append(phonologicExercises, exercise3)

	var exercise4 PhonologicRhymeSentence
	readJSONFile("data/PhonologicRhymeSentence.json", &exercise4)
	phonologicExercises = append(phonologicExercises, exercise4)
}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	// Create an instance of ExerciseData with phonologicExercises
	exerciseData := ExerciseData{
		PhonologicExercises: phonologicExercises,
	}

	// Parse the HTML template
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"getType": func(i interface{}) string {
			switch i.(type) {
			case PhonologicRhymePair:
				return "PhonologicRhymePair"
			case PhonologicRhymeMatch:
				return "PhonologicRhymeMatch"
			case PhonologicRhymeMultipleMatch:
				return "PhonologicRhymeMultipleMatch"
			case PhonologicRhymeSentence:
				return "PhonologicRhymeSentence"
			default:
				return ""
			}
		},
	}).ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse the HTML template: %v", err)
	}
	// Execute the HTML template with the ExerciseData instance
	err = tmpl.Execute(w, exerciseData)
	if err != nil {
		log.Fatalf("Failed to execute the HTML template: %v", err)
	}
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
