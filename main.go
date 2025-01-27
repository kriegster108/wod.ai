package main

import (
	"log"

	"github.com/kriegster108/wod.ai/processor"
)

func main() {

	var workout processor.Processor = &processor.OpenAIChatResult{}

	err := workout.CreateClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	workout.GenerateWorkout("give me a crossfit workout")
}
