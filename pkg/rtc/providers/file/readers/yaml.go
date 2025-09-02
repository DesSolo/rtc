package readers

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/internal"
	"rtc/pkg/rtc/providers/file"
)

// SimpleYAML ...
func SimpleYAML() file.ReaderFunc {
	return func(data []byte) (map[rtc.Key]rtc.Value, error) {
		var items map[string]string
		if err := yaml.Unmarshal(data, &items); err != nil {
			return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
		}

		result := make(map[rtc.Key]rtc.Value, len(items))
		for k, v := range items {
			result[rtc.Key(k)] = internal.NewValue([]byte(v))
		}

		return result, nil
	}
}
