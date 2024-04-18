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

////////////////////////////////////////////////////////////////////

type PhonologicRhymePair struct {
	Question                 string                     `json:"question"`
	Answer                   []string                   `json:"answer"`
	PhonologicRhymePairItems []PhonologicRhymePairItems `json:"items"`
}

type PhonologicRhymePairItems struct {
	Words   []string `json:"words"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeMatch struct {
	Question                  string                      `json:"question"`
	PhonologicRhymeMatchItems []PhonologicRhymeMatchItems `json:"items"`
}

type PhonologicRhymeMatchItems struct {
	Words   []string `json:"words"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeMultipleMatch struct {
	Question string         `json:"question"`
	Column1  []string       `json:"column1"`
	Column2  []string       `json:"items"`
	Answers  map[string]int `json:"answers"`
}

type PhonologicRhymeSentence struct {
	Question  string                      `json:"question"`
	ItemsAnti []PhonologicRhymeMatchItems `json:"items"`
}

////////////////////////////////////////////////////////////////////

var allExercises []AllExercises
var exercises1 []Exercise1
var exercises2 []Exercise2
var exercises3 []Exercise3
var answers []Answer

////////////////////////////////////////////////////////////////

func loadExercises() {
	// Read the JSON file
	data, err := os.ReadFile("exercises.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal JSON data into the temporary slice
	err = json.Unmarshal(data, &allExercises)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Distribute the data into exercises1 and exercises2
	for _, ex := range allExercises {
		if len(ex.Items) > 0 {
			// Exercise1
			exercises1 = append(exercises1, Exercise1{
				Question: ex.Question,
				Items:    ex.Items,
				Answers:  ex.Answers,
			})
		} else if len(ex.Column3) > 0 {
			// Exercise3
			exercises3 = append(exercises3, Exercise3{
				Question: ex.Question,
				Column3:  ex.Column3,
				Column4:  ex.Column4,
			})
		} else {
			// Exercise2
			exercises2 = append(exercises2, Exercise2{
				Question: ex.Question,
				Column1:  ex.Column1,
				Column2:  ex.Column2,
			})
		}
	}

	log.Println("Successfully loaded exercises")

	// Print parsed data
	fmt.Println("Exercise 1:")
	for _, ex := range exercises1 {
		fmt.Printf("Question: %s\n", ex.Question)
		fmt.Printf("Items: %v\n", ex.Items)
		fmt.Printf("Answers: %v\n", ex.Answers)
	}

	fmt.Println("\nExercise 2:")
	for _, ex := range exercises2 {
		fmt.Printf("Question: %s\n", ex.Question)
		fmt.Printf("Column1: %v\n", ex.Column1)
		fmt.
			Printf("Column2: %v\n", ex.Column2)
	}

	fmt.Println("\nExercise 3:")
	for _, ex := range exercises3 {
		fmt.Printf("Question: %s\n", ex.Question)
		fmt.Printf("Column3: %v\n", ex.Column3)
		fmt.Printf("Column4: %v\n", ex.Column4)
	}

}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	// Create an instance of ExerciseData with both sets of exercises
	exerciseData := ExerciseData{
		Exercises1: exercises1,
		Exercises2: exercises2,
		Exercises3: exercises3,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("index.html")
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
