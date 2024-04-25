// main.go
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type PhonologicRhymePair struct {
	Question                 string                     `json:"question"`
	Answer                   []string                   `json:"answer"`
	Level                    int                        `json:"level"`
	ExerciseNumber           int                        `json:"exercise"`
	PhonologicRhymePairItems []PhonologicRhymePairItems `json:"items"`
}

type PhonologicRhymePairItems struct {
	Words   []string `json:"words"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeMatchItems struct {
	Word    string   `json:"word"`
	Answers []string `json:"answers"`
	Correct int      `json:"correct"`
}

type PhonologicRhymeSentenceItems struct {
	Sentence string   `json:"sentence"`
	Answers  []string `json:"answers"`
	Correct  int      `json:"correct"`
}

type PhonologicRhymeMatch struct {
	Question                  string                      `json:"question"`
	Level                     int                         `json:"level"`
	ExerciseNumber            int                         `json:"exercise"`
	PhonologicRhymeMatchItems []PhonologicRhymeMatchItems `json:"items"`
}

type PhonologicRhymeMultipleMatch struct {
	Question       string         `json:"question"`
	Level          int            `json:"level"`
	ExerciseNumber int            `json:"exercise"`
	Column1        []string       `json:"column1"`
	Column2        []string       `json:"column2"`
	Answers        map[string]int `json:"answers"`
}

type PhonologicRhymeSentence struct {
	Question                     string                         `json:"question"`
	Level                        int                            `json:"level"`
	ExerciseNumber               int                            `json:"exercise"`
	Sentence                     string                         `json:"sentence"`
	PhonologicRhymeSentenceItems []PhonologicRhymeSentenceItems `json:"items"`
}

type Answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type PhonologicExercises struct {
	PhonologicRhymePair          []PhonologicRhymePair
	PhonologicRhymeMatch         []PhonologicRhymeMatch
	PhonologicRhymeMultipleMatch []PhonologicRhymeMultipleMatch
	PhonologicRhymeSentence      []PhonologicRhymeSentence
}

var answers []Answer

var phonologicRhymePairExercise []PhonologicRhymePair
var phonologicRhymeMatchExercise []PhonologicRhymeMatch
var phonologicRhymeMultipleMatchExercise []PhonologicRhymeMultipleMatch
var phonologicRhymeSentenceExercise []PhonologicRhymeSentence

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

func loadExercises(exerciseNumber int) {
	// Filter exercises in phonologicRhymePairExercise with ExerciseNumber == 1
	readJSONFile("data/PhonologicRhymePair.json", &phonologicRhymePairExercise)
	filtered := make([]PhonologicRhymePair, 0)
	for _, exercise := range phonologicRhymePairExercise {
		if exercise.ExerciseNumber == exerciseNumber {
			filtered = append(filtered, exercise)
		}
	}
	phonologicRhymePairExercise = filtered

	readJSONFile("data/PhonologicRhymeMatch.json", &phonologicRhymeMatchExercise)

	readJSONFile("data/PhonologicRhymeMultipleMatch.json", &phonologicRhymeMultipleMatchExercise)

	readJSONFile("data/PhonologicRhymeSentence.json", &phonologicRhymeSentenceExercise)

}

func serveExercises(w http.ResponseWriter, r *http.Request) {
	// Create an instance of ExerciseData with phonologicExercises
	exerciseData := PhonologicExercises{
		PhonologicRhymePair:          phonologicRhymePairExercise,
		PhonologicRhymeMatch:         phonologicRhymeMatchExercise,
		PhonologicRhymeMultipleMatch: phonologicRhymeMultipleMatchExercise,
		PhonologicRhymeSentence:      phonologicRhymeSentenceExercise,
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
	exerciseNumber := 1
	loadExercises(exerciseNumber)

	http.HandleFunc("/", serveExercises)
	http.HandleFunc("/answer", handleAnswer)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
