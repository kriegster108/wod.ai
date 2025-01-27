package processor

type Processor interface {
	CreateClient() error
	GenerateWorkout(prompt string) error
	// GetRetry() error
}
