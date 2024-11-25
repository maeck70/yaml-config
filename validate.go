package yamlconfig

import (
	"log"
)

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema
	pad := ""
	recurValidate(pad, c.Data, cv.Schema, "")
	return nil
}

func recurValidate(pad string, data any, schema interface{}, key string) {

	switch s := schema.(type) {

	case SchemaField_t:
		// Map of fields
		log.Printf(pad+"RV - %s = %+v", key, s)
		for k, v := range s.Attributes {
			log.Printf(pad+"RV - %s = %+v", k, v)
			validate(pad, data, v, k)
		}

	case GroupField_t:
		// Grouped map of fields
		log.Printf(pad+"RVGroup - %s = %+v", key, schema)

		csa := s.Attributes
		log.Printf(pad+"CS - %+v", csa)
		for k, v := range s.Attributes {

			sf2 := *new(SchemaField_t)
			for rkey, rvalue := range v.(map[string]interface{}) {
				log.Printf(pad+"rkey %s rvalue - %+v", rkey, rvalue)

				// I know this is stupid, but I can't figure out how to do this properly with reflection
				switch rkey {
				case "type":
					sf2.Type = rvalue.(string)
				case "description":
					sf2.Description = rvalue.(string)
				case "required":
					sf2.Required = rvalue.(bool)
				case "default":
					sf2.Default = rvalue
				case "options":
					sf2.Options = rvalue.([]any)
				case "optiontype":
					sf2.OptionType = rvalue.(string)
				case "min":
					sf2.Min = rvalue.(int)
				case "max":
					sf2.Max = rvalue.(int)
				case "attributes":
					sf2.Attributes = rvalue.(map[string]SchemaField_t)
				case "items":
					sf2.List = rvalue.(map[string]SchemaField_t)
				case "valid":
					sf2.Valid = rvalue.([]string)
				case "group":
					sf2.Group = rvalue.(GroupField_t)
				default:
					log.Printf(pad+"Unknown Field %s", rkey)
				}
			}

			d := data.(map[string]interface{})[key]
			log.Printf(pad+"  Validate - %s = %+v", k, d)

			log.Printf(pad+"  sf2 - %s = %+v", k, sf2)

			validate(pad, d, sf2, k)
		}

	case Schema_t:
		for k, v := range s {
			log.Printf(pad+"RV - %s = %+v", k, v)
			validate(pad, data, v, k)
		}
	default:
		log.Printf(pad+"Unknown Type %s", schema)
	}
}

func validate(pad string, data any, schemaField SchemaField_t, key string) {
	config := data.(map[string]interface{})
	f := getConfigField(config, key)

	switch schemaField.Type {
	case "string":
		log.Printf(pad+"String - Config %s = %+v", key, f)
		checkField(data, key, schemaField)
		log.Printf(pad+"  After %s = %+v", key, getConfigField(config, key))

	case "integer":
		log.Printf(pad+"Integer - Config %s = %+v", key, f)
		checkField(data, key, schemaField)
		log.Printf(pad+"  After %s = %+v", key, getConfigField(config, key))

	case "boolean":
		log.Printf(pad+"Boolean - Config %s = %+v", key, f)
		checkField(data, key, schemaField)
		log.Printf(pad+"  After %s = %+v", key, getConfigField(config, key))

	case "float":
		log.Printf(pad+"Float - Config %s = %+v", key, f)
		checkField(data, key, schemaField)
		log.Printf(pad+"  After %s = %+v", key, getConfigField(config, key))

	case "array":
		log.Printf(pad+"Array - Config %s = %+v", key, f)
		checkOptions(data, key, schemaField)
		log.Printf(pad+"  After %s = %+v", key, getConfigField(config, key))

	case "map":
		log.Printf(pad+"Map - Config %s = %+v", key, f)

		for k, v := range f.(map[string]interface{}) {
			log.Printf(pad+"  Map - %s = %+v", k, v)
			recurValidate(pad+"  ", f, schemaField.Group, k)
		}

	case "object":
		log.Printf(pad+"Object - Config %s = %+v", key, f)
		recurValidate(pad+"  ", f, schemaField, key)

	case "objectlist":
		log.Printf(pad+"objectlist - Config %s = %+v", key, f)
		for sk, sv := range schemaField.List {
			for _, cv := range f.([]interface{}) {
				checkField(cv, sk, sv)
			}
		}

	default:
		log.Printf(pad+"Unknown Type %s", schemaField.Type)
	}

	log.Print("")
}

func getConfigField(config interface{}, key string) interface{} {
	c := config.(map[string]interface{})
	cf := c[key]
	return cf
}

func checkOptions(data interface{}, schemaFieldKey string, schemaField SchemaField_t) {
	log.Printf("Check Options %s = %+v", schemaFieldKey, schemaField.Valid)
	log.Printf("         Data %s = %+v", schemaFieldKey, data.(map[string]interface{})[schemaFieldKey])

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
