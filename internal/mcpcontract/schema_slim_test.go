package mcpcontract

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
)

func TestSlimSchemaRawCompactsNullableEnums(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"object",
		"additionalProperties":false,
		"properties":{
			"a":{"description":"甲。","anyOf":[{"type":"null"},{"type":"string","enum":["甲","乙"]}]},
			"b":{"description":"乙。","anyOf":[{"type":"null"},{"type":"string","enum":["甲","乙"]}]}
		}
	}`)
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties map[string]map[string]any `json:"properties"`
		Defs       map[string]map[string]any `json:"$defs"`
	}
	if err := json.Unmarshal(slimmed, &schema); err != nil {
		t.Fatal(err)
	}
	for name, property := range schema.Properties {
		ref, ok := property["$ref"].(string)
		if !ok || ref == "" {
			t.Fatalf("%s has no enum $ref: %#v", name, property)
		}
		if property["description"] == "" {
			t.Fatalf("%s lost its field description", name)
		}
	}
	if len(schema.Defs) != 1 {
		t.Fatalf("shared enum definitions = %d, want 1", len(schema.Defs))
	}
	for _, definition := range schema.Defs {
		types, ok := definition["type"].([]any)
		if !ok || len(types) != 2 || types[0] != "null" || types[1] != "string" {
			t.Fatalf("shared schema type = %#v, want [null string]", definition["type"])
		}
		values, ok := definition["enum"].([]any)
		if !ok || len(values) != 3 || values[0] != nil || values[1] != "甲" || values[2] != "乙" {
			t.Fatalf("shared enum = %#v, want [null 甲 乙]", definition["enum"])
		}
	}
}

func TestSlimSchemaRawKeepsUniqueEnumsInline(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"object",
		"properties":{
			"a":{"description":"甲。","type":"string","enum":["甲"]}
		}
	}`)
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(slimmed, &schema); err != nil {
		t.Fatal(err)
	}
	if _, exists := schema["$defs"]; exists {
		t.Fatalf("unique enum was hoisted: %s", slimmed)
	}
}

func TestSlimSchemaRawSharesRepeatedObjects(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"object",
		"properties":{
			"a":{
				"description":"甲列表。",
				"type":"array",
				"items":{
					"type":"object",
					"additionalProperties":false,
					"required":["name"],
					"properties":{"name":{"type":"string","description":"名称。"}}
				}
			},
			"b":{
				"description":"乙列表。",
				"type":"array",
				"items":{
					"type":"object",
					"additionalProperties":false,
					"required":["name"],
					"properties":{"name":{"type":"string","description":"名称。"}}
				}
			},
			"small":{
				"type":"object",
				"additionalProperties":false,
				"required":["name"],
				"properties":{"name":{"type":"string"}}
			},
			"other":{"type":"object","properties":{"unique":{"type":"string"}}}
		}
	}`)
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties map[string]map[string]any `json:"properties"`
		Defs       map[string]map[string]any `json:"$defs"`
	}
	if err := json.Unmarshal(slimmed, &schema); err != nil {
		t.Fatal(err)
	}
	aRef, ok := schema.Properties["a"]["items"].(map[string]any)["$ref"].(string)
	if !ok || !strings.HasPrefix(aRef, "#/$defs/Object_") {
		t.Fatalf("a items ref = %#v, want shared object ref", schema.Properties["a"]["items"])
	}
	bRef, ok := schema.Properties["b"]["items"].(map[string]any)["$ref"].(string)
	if !ok || bRef != aRef {
		t.Fatalf("b items ref = %#v, want %q", schema.Properties["b"]["items"], aRef)
	}
	if schema.Properties["a"]["description"] == "" || schema.Properties["b"]["description"] == "" {
		t.Fatal("array descriptions were lost")
	}
	if _, exists := schema.Properties["small"]["$ref"]; exists {
		t.Fatal("small object was hoisted")
	}
	if _, exists := schema.Properties["other"]["$ref"]; exists {
		t.Fatal("unique object was hoisted")
	}
	if len(schema.Defs) != 1 {
		t.Fatalf("definitions = %d, want 1", len(schema.Defs))
	}
	definition := schema.Defs["Object_"+aRef[len("#/$defs/Object_"):]]
	if definition["type"] != "object" || definition["additionalProperties"] != false {
		t.Fatalf("shared object definition = %#v", definition)
	}
}

func TestSlimSchemaRawSharesRepeatedObjectsInsideUniqueEnvelope(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"object",
		"properties":{
			"container":{
				"type":"object",
				"additionalProperties":false,
				"required":["a","b","unique"],
				"properties":{
					"a":{
						"type":"array",
						"items":{"type":"object","required":["name","kind","value"],"properties":{"name":{"type":"string","description":"名称。"},"kind":{"type":"string","description":"类型。"},"value":{"type":"string","description":"值。"}}}
					},
					"b":{
						"type":"array",
						"items":{"type":"object","required":["name","kind","value"],"properties":{"name":{"type":"string","description":"名称。"},"kind":{"type":"string","description":"类型。"},"value":{"type":"string","description":"值。"}}}
					},
					"unique":{"type":"string"}
				}
			}
		}
	}`)
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties struct {
			Container struct {
				Properties struct {
					A struct {
						Items map[string]any `json:"items"`
					} `json:"a"`
					B struct {
						Items map[string]any `json:"items"`
					} `json:"b"`
				} `json:"properties"`
			} `json:"container"`
		} `json:"properties"`
		Defs map[string]any `json:"$defs"`
	}
	if err := json.Unmarshal(slimmed, &schema); err != nil {
		t.Fatal(err)
	}
	aRef, ok := schema.Properties.Container.Properties.A.Items["$ref"].(string)
	if !ok || !strings.HasPrefix(aRef, "#/$defs/Object_") {
		t.Fatalf("nested a items ref = %#v, want shared object ref", schema.Properties.Container.Properties.A.Items)
	}
	if schema.Properties.Container.Properties.B.Items["$ref"] != aRef {
		t.Fatalf("nested b items ref = %#v, want %q", schema.Properties.Container.Properties.B.Items, aRef)
	}
	if len(schema.Defs) != 1 {
		t.Fatalf("definitions = %d, want 1", len(schema.Defs))
	}
}

func TestSlimSchemaRawSharesRepeatedNullableObjects(t *testing.T) {
	item := `{"type":["null","object"],"required":["west","east"],"properties":{"west":{"type":"number","description":"西界经度。"},"east":{"type":"number","description":"东界经度。"}}}`
	raw := json.RawMessage(`{
		"type":"object",
		"properties":{
			"a":{"type":["null","array"],"items":` + item + `},
			"b":{"type":["null","array"],"items":` + item + `}
		}
	}`)
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(slimmed, &schema); err != nil {
		t.Fatal(err)
	}
	properties := schema["properties"].(map[string]any)
	a := properties["a"].(map[string]any)
	b := properties["b"].(map[string]any)
	items := a["items"].(map[string]any)
	bItems := b["items"].(map[string]any)
	aRef, ok := items["$ref"].(string)
	if !ok || !strings.HasPrefix(aRef, "#/$defs/Object_") {
		t.Fatalf("nullable a items ref = %#v, want shared object ref", items)
	}
	if bItems["$ref"] != aRef {
		t.Fatalf("nullable b items ref = %#v, want %q", bItems, aRef)
	}
	definitions := schema["$defs"].(map[string]any)
	parts := strings.Split(aRef, "/")
	definition := definitions[parts[len(parts)-1]].(map[string]any)
	if !ok {
		t.Fatalf("definition %q missing or not an object; defs=%v, aRef=%s", parts[len(parts)-1], schema["$defs"], aRef)
	}
	types, ok := definition["type"].([]any)
	if !ok || len(types) != 2 || types[0] != "null" || types[1] != "object" {
		t.Fatalf("nullable definition type = %#v", definition["type"])
	}
}

func TestSlimBaziChartKeepsNullInHiddenStemEnums(t *testing.T) {
	options, err := cachedSchemaOptions()
	if err != nil {
		t.Fatal(err)
	}
	derived, err := jsonschema.For[BaziChartOutput](options)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := derived.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	slimmed, err := slimSchemaRaw(raw)
	if err != nil {
		t.Fatal(err)
	}

	var root map[string]any
	if err := json.Unmarshal(slimmed, &root); err != nil {
		t.Fatal(err)
	}
	found := false
	var walk func(node any)
	walk = func(node any) {
		object, ok := node.(map[string]any)
		if !ok {
			if items, ok := node.([]any); ok {
				for _, item := range items {
					walk(item)
				}
			}
			return
		}
		types, _ := object["type"].([]any)
		values, _ := object["enum"].([]any)
		hasNullableString := slices.Contains(types, "null") && slices.Contains(types, "string")
		hasStemAndNull := slices.Contains(values, nil) && slices.Contains(values, "癸")
		if hasNullableString && hasStemAndNull {
			found = true
		}
		for _, field := range object {
			walk(field)
		}
	}
	walk(root)
	if !found {
		t.Fatalf("bazi chart schema has no nullable hidden-stem enum: %s", slimmed)
	}
}
