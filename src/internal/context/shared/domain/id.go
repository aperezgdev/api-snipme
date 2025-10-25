package domain

import "github.com/google/uuid"

type Id uuid.UUID

func NewID() (Id, error) {
	uuid, err := uuid.NewV7()
	return Id(uuid), err
}

func (id Id) String() string {
	return uuid.UUID(id).String()
}

func ParseID(idStr string) (Id, error) {
	uuid, err := uuid.Parse(idStr)
	return Id(uuid), err
}
