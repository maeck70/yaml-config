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

			sf := *new(SchemaField_t)
			for rkey, rvalue := range v.(map[string]interface{}) {
				// I know this is stupid, but I can't figure out how to do this properly with reflection
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
			d := data.(map[string]interface{})[key]
			validate(d, sf, k)
		}

	case Schema_t:
		for k, v := range s {
			validate(data, v, k)
		}
	default:
		log.Printf("Unknown Type %s", schema)
	}
}

func validate(data any, schemaField SchemaField_t, key string) {
	config := data.(map[string]interface{})
	f := getConfigField(config, key)

	switch schemaField.Type {
	case "string":
		checkField(data, key, schemaField)

	case "integer":
		checkField(data, key, schemaField)

	case "boolean":
		checkField(data, key, schemaField)

	case "float":
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
		for sk, sv := range schemaField.List {
			for _, cv := range f.([]interface{}) {
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
	for _, o := range data.(map[string]interface{})[schemaFieldKey].([]interface{}) {
		f := false
		for _, v := range schemaField.Valid {
			if v == o {
				f = true
				break
			}
		}
		if !f {
			log.Fatalf("Option '%s' is not allowed. Valid options: %v", o, schemaField.Valid)
		}
	}
}

func checkField(data interface{}, schemaFieldKey string, schemaField SchemaField_t) {
	// Check if the schema fields are present in the data
	value := data.(map[string]interface{})[schemaFieldKey]
	log.Printf("Validate %s, %+v = %v", schemaFieldKey, schemaField, value)

	// Set default if non existent
	if data.(map[string]interface{})[schemaFieldKey] == nil {
		log.Printf("Setting default value for %s to %v\n", schemaFieldKey, schemaField.Default)
		data.(map[string]interface{})[schemaFieldKey] = schemaField.Default
	}
	if _, ok := data.(map[string]interface{})[schemaFieldKey]; !ok {
		log.Printf("Setting default value for %s to %v\n", schemaFieldKey, schemaField.Default)
		data.(map[string]interface{})[schemaFieldKey] = schemaField.Default
	}

	// Check if required field is empty
	if schemaField.Required {
		if data.(map[string]interface{})[schemaFieldKey] == nil {
			log.Fatalf("Required field %s is missing.\n", schemaFieldKey)
		}
	}

	// Check value against field types attributes
	switch schemaField.Type {
	case "integer":
		val := data.(map[string]interface{})[schemaFieldKey].(int)
		if schemaField.Min > 0 || schemaField.Max > 0 {
			if val > schemaField.Max {
				log.Fatalf("Error on field %s value %d is greater than max %d.\n",
					schemaFieldKey, val, schemaField.Max)
			}
			if val < schemaField.Min {
				log.Fatalf("Error on field %s value %d is less than min %d.\n",
					schemaFieldKey, val, schemaField.Min)
			}
		}
	}
}
