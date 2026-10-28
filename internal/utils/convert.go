package utils

import (
	"github.com/google/uuid"
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
)

func ToStrPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func StringToUUID(field string, s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, &errDom.ValidationError{Field: field, Reason: "must be a UUID"}
	}
	return id, nil
}
