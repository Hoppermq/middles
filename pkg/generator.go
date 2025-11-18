package pkg

import "github.com/google/uuid"

func Generator(newUUID func() (uuid.UUID, error)) string {
	generatedUUID, err := newUUID()
	if err != nil {
		return "invalid uuid"
	}

	return generatedUUID.String()
}
