package yamlconfig

import (
	"log"
)

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema
	pad := ""
	cv.recurValidate(pad, c.Data.(map[string]interface{}), "")
	return nil
}

func (cv ConfigValidator_t) recurValidate(pad string, data interface{}, key string) {
	pad = pad + "  "
	switch val := data.(type) {
	case map[string]interface{}:
		for k, v := range data.(map[string]interface{}) { //val {
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
				cv.recurValidate(pad, v, k)
			case map[string]interface{}:
				log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
				cv.recurValidate(pad, v, k)
			default:
				log.Printf(pad+"Field - Key: %s, Value: %v\n", k, v)
				value := val[k]
				cv.checkValue(pad, k, &value)
				data.(map[string]interface{})[k] = value
			}
		}
	case []interface{}:
		for i, v := range val {
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - %d, Value: %v\n", i, v)
				cv.recurValidate(pad, v, key)
			case map[string]interface{}:
				log.Printf(pad+"Map - %d, Value: %v\n", i, v)
				cv.recurValidate(pad, v, key)
			default:
				log.Printf(pad+"Field - %d, Value: %v\n", i, v)
				cv.checkValue(pad, key, &val[i])
			}
		}
	}
}

func (cv ConfigValidator_t) checkValue(pad string, schemakey string, v *interface{}) {
	switch val := (*v).(type) {
	case string:
		// log.Printf(pad+"String: %s\n", val)
		// *v = "new value"
	case int:
		if cv.Schema[schemakey].Min > 0 || cv.Schema[schemakey].Max > 0 {
			if val > int(cv.Schema[schemakey].Max) {
				log.Fatalf(pad+"On attribute %s value %d is greater than max %d.\n",
					schemakey, val, cv.Schema[schemakey].Max)
			}
			if val < int(cv.Schema[schemakey].Min) {
				log.Fatalf(pad+"On attribute %s value %d is less than min %d.\n",
					schemakey, val, cv.Schema[schemakey].Min)
			}
		}
		// log.Printf(pad+"Int: %d\n", val)
		//*v = val + 1
	case float64:
		if cv.Schema[schemakey].Min > 0 || cv.Schema[schemakey].Max > 0 {
			if val > float64(cv.Schema[schemakey].Max) {
				log.Fatalf(pad+"On attribute %s value %f is greater than max %f.\n",
					schemakey, val, cv.Schema[schemakey].Max)
			}
			if val < float64(cv.Schema[schemakey].Min) {
				log.Fatalf(pad+"On attribute %s value %f is less than min %f.\n",
					schemakey, val, cv.Schema[schemakey].Min)
			}
		}
		// log.Printf(pad+"Int: %d\n", val)
		//*v = val + 1
	case bool:
		// log.Printf(pad+"Bool: %t\n", val)
		//*v = !val
	default:
		// log.Printf(pad+"Unknown: %v\n", val)
	}

	// Get the details from cv for field key
	log.Printf(pad+"Schema field %s specs: %+v", schemakey, cv.Schema[schemakey])

}
