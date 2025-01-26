package processor

type Processor interface {
    GetWorkout() error
	GetRetry() error
}