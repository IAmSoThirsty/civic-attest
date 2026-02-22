package canonical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fxamacker/cbor/v2"
)

// Format represents a canonical encoding format
type Format string

const (
	// CBOR is the canonical CBOR format
	CBOR Format = "CBOR"
	// JSON is the deterministic JSON format
	JSON Format = "JSON"
)

// Encode encodes data in a canonical format
func Encode(data interface{}, format Format) ([]byte, error) {
	switch format {
	case CBOR:
		return encodeCBOR(data)
	case JSON:
		return encodeJSON(data)
	default:
		return nil, fmt.Errorf("unsupported canonical format: %s", format)
	}
}

// Decode decodes canonically encoded data
func Decode(data []byte, format Format, dest interface{}) error {
	switch format {
	case CBOR:
		return decodeCBOR(data, dest)
	case JSON:
		return decodeJSON(data, dest)
	default:
		return fmt.Errorf("unsupported canonical format: %s", format)
	}
}

func encodeCBOR(data interface{}) ([]byte, error) {
	// Use deterministic CBOR encoding mode
	encMode, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, fmt.Errorf("failed to create canonical CBOR encoder: %w", err)
	}

	return encMode.Marshal(data)
}

func decodeCBOR(data []byte, dest interface{}) error {
	decMode, err := cbor.DecOptions{}.DecMode()
	if err != nil {
		return fmt.Errorf("failed to create CBOR decoder: %w", err)
	}

	return decMode.Unmarshal(data, dest)
}

func encodeJSON(data interface{}) ([]byte, error) {
	// Ensure deterministic JSON by sorting keys
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")

	// Convert to map to sort keys if needed
	normalized, err := normalizeForJSON(data)
	if err != nil {
		return nil, err
	}

	if err := encoder.Encode(normalized); err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Remove trailing newline
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result, nil
}

func decodeJSON(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}

// normalizeForJSON ensures deterministic JSON encoding by sorting map keys
func normalizeForJSON(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return normalizemap(v)
	case []interface{}:
		normalized := make([]interface{}, len(v))
		for i, item := range v {
			n, err := normalizeForJSON(item)
			if err != nil {
				return nil, err
			}
			normalized[i] = n
		}
		return normalized, nil
	default:
		return data, nil
	}
}

func normalizemap(m map[string]interface{}) (map[string]interface{}, error) {
	normalized := make(map[string]interface{})
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		n, err := normalizeForJSON(m[k])
		if err != nil {
			return nil, err
		}
		normalized[k] = n
	}

	return normalized, nil
}
