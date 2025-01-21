package uuid

import "github.com/google/uuid"

func Gen() string {
	id := uuid.New()
	return id.String()
}
