package yamlconfig

import (
	"log"
)

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema
	pad := ""
	cv.recurValidate(pad, c.Data.(map[string]interface{}))
	return nil
}

func (cv ConfigValidator_t) recurValidate(pad string, data interface{}) {
	pad = pad + "  "
	switch val := data.(type) {
	case map[string]interface{}:
		for k, v := range val {
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
				cv.recurValidate(pad, v)
			case map[string]interface{}:
				log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
				cv.recurValidate(pad, v)
			default:
				log.Printf(pad+"Field - Key: %s, Value: %v\n", k, v)
				value := val[k]
				cv.checkValue(pad, &value)
			}
		}
	case []interface{}:
		for i, v := range val {
			switch v := v.(type) {
			case []interface{}:
				log.Printf(pad+"List - %d, Value: %v\n", i, v)
				cv.recurValidate(pad, v)
			case map[string]interface{}:
				log.Printf(pad+"Map - %d, Value: %v\n", i, v)
				cv.recurValidate(pad, v)
			default:
				log.Printf(pad+"Field - %d, Value: %v\n", i, v)
				cv.checkValue(pad, &val[i])
			}
		}
	}
}

func (cv ConfigValidator_t) checkValue(pad string, v *interface{}) {
	switch val := (*v).(type) {
	case string:
		log.Printf(pad+"String: %s\n", val)
		*v = "new value"
	case int:
		log.Printf(pad+"Int: %d\n", val)
		*v = val + 1
	case bool:
		log.Printf(pad+"Bool: %t\n", val)
		*v = !val
	default:
		log.Printf(pad+"Unknown: %v\n", val)
	}
}
