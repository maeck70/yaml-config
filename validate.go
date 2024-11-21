package yamlconfig

import (
	"log"
)

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema
	pad := ""
	ncv := cv.Schema
	ncv.recurValidate(pad, c.Data.(map[string]interface{}), "")
	return nil
}

func (cv Schema_t) recurValidate(pad string, data interface{}, key string) {
	pad = pad + "  "
	switch val := data.(type) {
	case map[string]interface{}:

		for k, v := range data.(map[string]interface{}) {
			ncv := cv[k].Attributes
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
				ncv.recurValidate(pad, v, k)
			case map[string]interface{}:
				log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
				ncv.recurValidate(pad, v, k)
			default:
				log.Printf(pad+"Field - Key: %s, Value: %v\n", k, v)
			}
		}
		cv.checkFields(pad, data)

	case []interface{}:
		for i, v := range val {
			ncv := cv[key].Attributes
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - %d, Value: %v\n", i, v)
				ncv.recurValidate(pad, v, key)
			case map[string]interface{}:
				log.Printf(pad+"Map - %d, Value: %v\n", i, v)
				ncv.recurValidate(pad, v, key)
			default:
				log.Printf(pad+"Field - %d, Value: %v\n", i, v)
			}
		}
		cv.checkFields(pad, data)

	}
}

func (cv Schema_t) checkFields(pad string, data interface{}) {
	// Check if the schema fields are present in the data
	for k, v := range cv {
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
