package wrappers

import "github.com/google/uuid"

type UUIDWrapper interface {
	New() uuid.UUID
}

type SimpleUUIDWrapper struct{}

func (self *SimpleUUIDWrapper) New() uuid.UUID {
	return uuid.New()
}
