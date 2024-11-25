package yamlconfig

import (
	"log"
)

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema
	recurValidate(c.Data, cv.Schema, "")
	return nil
}

func recurValidate(data any, schema interface{}, key string) {
	switch s := schema.(type) {

	case SchemaField_t:
		// Map of fields
		for k, v := range s.Attributes {
			validate(data, v, k)
		}

	case GroupField_t:
		// Grouped map of fields
		for k, v := range s.Attributes {
			sf := parseSchemaField(v.(map[string]interface{}))
			d := data.(map[string]interface{})[key]
			validate(d, sf, k)
		}

	case Schema_t:
		for k, v := range s {
			validate(data, v, k)
		}
	default:
		log.Printf("Unknown Type %T", schema)
	}
}

func parseSchemaField(v map[string]interface{}) SchemaField_t {
	sf := SchemaField_t{}
	for rkey, rvalue := range v {
		switch rkey {
		case "type":
			sf.Type = rvalue.(string)
		case "description":
			sf.Description = rvalue.(string)
		case "required":
			sf.Required = rvalue.(bool)
		case "default":
			sf.Default = rvalue
		case "options":
			sf.Options = rvalue.([]any)
		case "optiontype":
			sf.OptionType = rvalue.(string)
		case "min":
			sf.Min = rvalue.(int)
		case "max":
			sf.Max = rvalue.(int)
		case "attributes":
			sf.Attributes = rvalue.(map[string]SchemaField_t)
		case "items":
			sf.List = rvalue.(map[string]SchemaField_t)
		case "valid":
			sf.Valid = rvalue.([]string)
		case "group":
			sf.Group = rvalue.(GroupField_t)
		default:
			log.Printf("Unknown Field %s", rkey)
		}
	}
	return sf
}

func validate(data any, schemaField SchemaField_t, key string) {
	config := data.(map[string]interface{})
	f := getConfigField(config, key)

	switch schemaField.Type {
	case "string", "integer", "boolean", "float":
		checkField(data, key, schemaField)

	case "array":
		checkOptions(data, key, schemaField)

	case "map":
		for k := range f.(map[string]interface{}) {
			recurValidate(f, schemaField.Group, k)
		}

	case "object":
		recurValidate(f, schemaField, key)

	case "objectlist":
		for _, cv := range f.([]interface{}) {
			for sk, sv := range schemaField.List {
				checkField(cv, sk, sv)
			}
		}

	default:
		log.Printf("Unknown Type %s", schemaField.Type)
	}
}

func getConfigField(config interface{}, key string) interface{} {
	c := config.(map[string]interface{})
	cf := c[key]
	return cf
}

func checkOptions(data interface{}, schemaFieldKey string, schemaField SchemaField_t) {
	options := data.(map[string]interface{})[schemaFieldKey].([]interface{})
	validOptions := make(map[string]struct{}, len(schemaField.Valid))
	for _, v := range schemaField.Valid {
		validOptions[v] = struct{}{}
	}

	for _, o := range options {
		if _, ok := validOptions[o.(string)]; !ok {
			log.Fatalf("Option '%s' is not allowed. Valid options: %v", o, schemaField.Valid)
		}
	}
}

func checkField(data interface{}, schemaFieldKey string, schemaField SchemaField_t) {
	config := data.(map[string]interface{})
	value, exists := config[schemaFieldKey]
	log.Printf("Validate %s, %+v = %v", schemaFieldKey, schemaField, value)

	// Set default if non-existent
	if !exists || value == nil {
		log.Printf("Setting default value for %s to %v\n", schemaFieldKey, schemaField.Default)
		config[schemaFieldKey] = schemaField.Default
		value = schemaField.Default
	}

	// Check if required field is empty
	if schemaField.Required && (value == nil || value == "") {
		log.Fatalf("Required field %s is missing.\n", schemaFieldKey)
	}

	// Check value against field types attributes
	switch schemaField.Type {
	case "integer":
		val := value.(int)
		if (schemaField.Min > 0 && val < schemaField.Min) || (schemaField.Max > 0 && val > schemaField.Max) {
			log.Fatalf("Error on field %s: value %d is out of range [%d, %d].\n",
				schemaFieldKey, val, schemaField.Min, schemaField.Max)
		}
	}
}
