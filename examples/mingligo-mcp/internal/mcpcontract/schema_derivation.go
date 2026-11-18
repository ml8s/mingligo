package mcpcontract

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

// inputSchemaFor derives the DTO-owned JSON Schema used for publication and
// official SDK validation.
func inputSchemaFor[T any](options *jsonschema.ForOptions) (json.RawMessage, error) {
	schema, err := jsonschema.For[T](options)
	if err != nil {
		return nil, err
	}
	raw, err := schema.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

// outputSchemaFor derives the DTO-owned JSON Schema for a typed handler result.
func outputSchemaFor[T any](options *jsonschema.ForOptions) (json.RawMessage, error) {
	schema, err := jsonschema.For[T](options)
	if err != nil {
		return nil, err
	}
	raw, err := schema.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}
