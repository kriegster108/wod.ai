package processor

type Processor interface {
	CreateClient() error
	GetWorkout(prompt string) error
	// GetRetry() error
}
