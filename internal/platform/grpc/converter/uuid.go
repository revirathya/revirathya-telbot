package converter

import (
	"github.com/google/uuid"
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
)

// FromProtoUUID parses a required UUID field; a bad value is a ValidationError, so ToStatus gives INVALID_ARGUMENT
// with "field: must be a UUID".
func FromProtoUUID(field string, s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, &errDom.ValidationError{Field: field, Reason: "must be a UUID"}
	}
	return id, nil
}

// FromOptionalProtoUUID is FromProtoUUID for an optional field: nil means "not set"
func FromOptionalProtoUUID(field string, s *string) (*uuid.UUID, error) {
	if s == nil {
		return nil, nil
	}

	id, err := FromProtoUUID(field, *s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
