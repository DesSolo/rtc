package auth

// Token ...
type Token struct {
	items map[string]*Payload
}

// NewToken ...
func NewToken(items map[string]*Payload) *Token {
	return &Token{
		items: items,
	}
}

// Authenticate ...
func (t *Token) Authenticate(token string) (*Payload, error) {
	p, ok := t.items[token]
	if !ok {
		return nil, ErrAuthFailed
	}

	return p, nil
}
