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
	case GroupField_t:
		log.Printf(pad+"RVGroup - %s = %+v", key, schema)

		csa := s.Attributes
		log.Printf(pad+"CS - %+v", csa)
		for k, v := range s.Attributes {
			sf := *new(SchemaField_t)
			for rkey, rvalue := range v.(map[string]interface{}) {
				log.Printf(pad+"rkey %s rvalue - %+v", rkey, rvalue)

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
					sf.Items = rvalue.(sfitem_t)
				case "valid":
					sf.Valid = rvalue.([]string)
				case "group":
					sf.Group = rvalue.(GroupField_t)
				default:
					log.Printf(pad+"Unknown Field %s", rkey)
				}
			}
			d := data.(map[string]interface{})[key]
			log.Printf(pad+"  Validate - %s = %+v", k, d)
			validate(pad, d, sf, k)
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

	case "object":
		log.Printf(pad+"Object - Config %s = %+v", key, f)

		for k, v := range f.(map[string]interface{}) {
			log.Printf(pad+"  Group - %s = %+v", k, v)
			recurValidate(pad+"  ", f, schemaField.Group, k)
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

/*
func recurValidate(pad string, data any, key string, cv Schema_t) {
	pad = pad + "  "
	switch val := data.(type) {
	case map[string]interface{}:
		log.Printf("val: %+v\n", val)
		for k, v := range data.(map[string]interface{}) {
			ncv := cv
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - Key: %s Value: %v", k, v)
				ncv := cv
				// recurValidate(pad, v, k, ncv)
				log.Printf(pad+"NCV: %+v", ncv)
			case map[string]interface{}:
				log.Printf(pad+"Map - Key: %s Value: %v", k, v)
				recurValidate(pad, v, k, ncv)
				log.Printf(pad+"NCV: %+v", ncv)
				checkFields(pad, data, ncv[k])
			default:
				log.Printf(pad+"Field - Key: %s Value: %v", k, v)
			}
		}

			case []interface{}:
				for i, v := range val {
					ncv := cv[key]
					switch v := v.(type) {
					case []interface{}:
						log.Printf(pad+"List - %d Value: %v", i, v)
						//recurValidate(pad, v, key, ncv)
						log.Printf(pad+"NCV: %+v", ncv)
					case map[string]interface{}:
						log.Printf(pad+"Map - %d Value: %v", i, v)
						//recurValidate(pad, v, key, ncv)
						log.Printf(pad+"NCV: %+v", ncv)
					default:
						log.Printf(pad+"Field - %d Value: %v", i, v)
					}
				}
				checkFields(pad, data, cv)

	}
}
*/

/*
func checkFields(pad string, data any, cv Schema_t) {
	// Check if the schema fields are present in the data
	for k, v := range cv.Group {
		value := data.(map[string]interface{})[k]
		log.Printf(pad+"Validate %s, %+v = %v", k, v, value)

		// Check if required field is empty
		if v.Required {
			if value == nil {
				log.Fatalf(pad+"Required field %s is missing.\n", k)
			}
		}

		// Set default if non existent
		if data.(map[string]interface{})[k] == nil {
			log.Printf(pad+"Setting default value for %s to %v\n", k, v.Default)
			data.(map[string]interface{})[k] = v.Default
		}
		if _, ok := data.(map[string]interface{})[k]; !ok {
			log.Printf(pad+"Setting default value for %s to %v\n", k, v.Default)
			data.(map[string]interface{})[k] = v.Default
		}

		// Check value against fields
		cv.checkValue(pad, k, &value)
	}
}

func (cv Schema_t) checkValue(pad string, schemakey string, v *interface{}) {
	switch val := (*v).(type) {
	case string:
		// log.Printf(pad+"String: %s\n", val)
		// *v = "new value"

	case int:
		if cv[schemakey].Min > 0 || cv[schemakey].Max > 0 {
			if val > int(cv[schemakey].Max) {
				log.Fatalf(pad+"On field %s value %d is greater than max %d.\n",
					schemakey, val, cv[schemakey].Max)
			}
			if val < int(cv[schemakey].Min) {
				log.Fatalf(pad+"On field %s value %d is less than min %d.\n",
					schemakey, val, cv[schemakey].Min)
			}
		}
		// log.Printf(pad+"Int: %d\n", val)
		//*v = val + 1

	case float64:
		if cv[schemakey].Min > 0 || cv[schemakey].Max > 0 {
			if val > float64(cv[schemakey].Max) {
				log.Fatalf(pad+"On field %s value %f is greater than max %f.\n",
					schemakey, val, cv[schemakey].Max)
			}
			if val < float64(cv[schemakey].Min) {
				log.Fatalf(pad+"On field %s value %f is less than min %f.\n",
					schemakey, val, cv[schemakey].Min)
			}
		}
		// log.Printf(pad+"Int: %d\n", val)
		//*v = val + 1

	case bool:
		// log.Printf(pad+"Bool: %t\n", val)
		//*v = !val

	default:
		log.Printf(pad+"Unknown field type: %v\n", val)
	}

	// Get the details from cv for field key
	log.Printf(pad+"Schema field %s specs: %+v", schemakey, cv[schemakey])

}
*/
