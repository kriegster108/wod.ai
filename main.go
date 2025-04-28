package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/kriegster108/wod.ai/processor"
)

type WorkoutHandler struct {
	client processor.WorkoutClient
}

func main() {
	var workoutClient processor.WorkoutClient = &processor.OpenAIChatClient{}

	mux := http.NewServeMux()

	err := workoutClient.CreateClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var workoutHandler = WorkoutHandler{client: workoutClient}
	mux.Handle("/workout", &workoutHandler)
	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", mux)
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	PostWorkout := regexp.MustCompile(`^/workout`)

	switch {
	case r.Method == http.MethodPost && PostWorkout.MatchString(r.URL.Path):
		err := h.generateWorkout(w, r)
		if err != nil {
			fmt.Println(err)
		}
		return
	default:
		return
	}
}

func (h *WorkoutHandler) generateWorkout(w http.ResponseWriter, r *http.Request) error {
	workout := processor.WorkoutPlan{}

	result, err := h.client.GenerateWorkout("give me a crossfit workout")
	if err != nil {
		http.Error(w, "Failed to get workout", http.StatusInternalServerError)
		return err
	}

	err = json.Unmarshal([]byte(result.Choices[0].Message.Content), &workout)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}

	workoutJSON, err := json.Marshal(workout)
	if err != nil {
		http.Error(w, "Failed to encode workout to JSON", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(workoutJSON)

	return nil
}
