package processor

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type WorkoutClient interface {
	CreateClient() error
	GenerateWorkout(prompt string) (*openai.ChatCompletionResponse, error)
}

type Exercise struct {
	Name          string `json:"name"`
	Duration      string `json:"duration"`
	Movement_Cues string `json:"movement_cues"`
	Reps          int    `json:"reps"`
	Rest_Time     string `json:"rest_time"`
}

type WorkoutPlan struct {
	Workout struct {
		Warmup []struct {
			Summary  string
			Exercise []Exercise
		} `json:"warmup"`
		MainWorkout []struct {
			Summary  string
			Exercise []Exercise
		} `json:"main_workout"`
		Cooldown []struct {
			Summary  string
			Exercise []Exercise
		} `json:"cooldown"`
	} `json:"workout"`
}

type OpenAIChatClient struct {
	client *openai.Client
}

func (workout *OpenAIChatClient) CreateClient() error {
	workout.client = openai.NewClient(os.Getenv("API_KEY"))
	return nil
}

func generateSchema() (*jsonschema.Definition, error) {
	var exercise WorkoutPlan

	schema, err := jsonschema.GenerateSchemaForType(exercise)
	if err != nil {
		fmt.Println(err)
		return &jsonschema.Definition{}, err
	}

	return schema, nil
}

func (workout *OpenAIChatClient) GenerateWorkout(prompt string) (*openai.ChatCompletionResponse, error) {
	schema, err := generateSchema()
	if err != nil {
		fmt.Printf("Invalid Schema %v\n", err)
		return nil, err
	}

	resp, err := workout.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "For movement_cues, give physical/mental cues for each exercise only. Like what muscles to engage, posture, alignment etc. If there arent any, just keep it empty.",
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "For rest_time, use a time format such as 1s for \"One Second\" or 1m for \"One minute\"",
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type:       openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{Name: "Workout_Plan", Schema: schema, Strict: true},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}
	return &resp, nil
}
