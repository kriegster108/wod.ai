package persistence

type Persistence interface {
	GetRecord() error
	DeleteRecord() error
	UpdateRecord() error
}
