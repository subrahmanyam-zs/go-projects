package Employee

import (
	"github.com/google/uuid"
)

func generateUUID() uuid.UUID {
	id := uuid.New()
	return id
}
