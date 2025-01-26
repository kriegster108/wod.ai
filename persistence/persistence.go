package persistence

import "fmt"

type Persistence interface {
    GetRecord() error
    DeleteRecord() error
    UpdateRecord() error
}