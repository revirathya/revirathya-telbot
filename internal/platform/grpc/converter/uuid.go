package converter

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID parses a required UUID field; a bad value is a ValidationError, so ToStatus gives INVALID_ARGUMENT
// with "field: must be a UUID".
func ParseUUID(field string, s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("must be valid uuid")
	}
	return id, nil
}

// ParseOptionalUUID is ParseUUID for an optional field: nil means "not set"
func ParseOptionalUUID(field string, s *string) (*uuid.UUID, error) {
	if s == nil {
		return nil, nil
	}

	id, err := ParseUUID(field, *s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
