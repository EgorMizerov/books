package wrappers

import "github.com/google/uuid"

//go:generate mockery --name=UUIDWrapper
type UUIDWrapper interface {
	New() uuid.UUID
}

type SimpleUUIDWrapper struct{}

func (self *SimpleUUIDWrapper) New() uuid.UUID {
	return uuid.New()
}
