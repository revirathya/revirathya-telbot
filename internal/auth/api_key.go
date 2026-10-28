package auth

import "context"

// Key
type ApiKey struct {
	token string
}

func NewApiKey(token string) ApiKey {
	return ApiKey{
		token: token,
	}
}

func (t ApiKey) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t ApiKey) RequireTransportSecurity() bool {
	return false
}
