package utils

func ToStrPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
