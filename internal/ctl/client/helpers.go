package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func encodePayload(payload any) (io.Reader, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return bytes.NewReader(data), nil
}
