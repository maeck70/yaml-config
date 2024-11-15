package yamlconfig

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// LoadConfig reads a YAML configuration file and loads it into a custom struct.
// It returns a map representation of the configuration.
//
// Parameters:
// - file: The path to the YAML configuration file.
// - customStruct: A pointer to the custom struct where the configuration will be loaded.
//
// Returns:
// - A map representation of the configuration.
func LoadConfig(file string, customStruct interface{}) any { // map[string]interface{} {
	// Read the file
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	// Build the full config struct
	dataStruct := reflect.ValueOf(customStruct)
	configStruct := &Config_t{}

	// Check if data is a pointer
	if dataStruct.Kind() == reflect.Ptr {
		configStruct.Data = dataStruct.Elem()
	} else {
		return nil
	}

	err = loadConfigFromString(data, configStruct)
	if err != nil {
		fmt.Println(err)
	}

	cv := &ConfigValidator_t{}
	err = cv.loadSchemaFromFile(configStruct.Metadata.Schema)
	if err != nil {
		fmt.Println(err)
	}

	cv.validateConfig(configStruct)

	// Convert the map structure back to the original struct
	res := reflect.New(reflect.TypeOf(customStruct).Elem()).Interface()
	mapstructure.Decode(configStruct.Data.(map[string]interface{}), &res)
	return res
	//return configStruct.Data.(map[string]interface{})
}

// loadConfigFromString unmarshals a YAML configuration string into a Config_t struct.
//
// Parameters:
// - configStr: The YAML configuration string.
// - c: A pointer to the Config_t struct where the configuration will be loaded.
//
// Returns:
// - An error if unmarshalling fails or if the schema is not defined in metadata.
func loadConfigFromString(configStr []byte, c *Config_t) error {
	err := yaml.Unmarshal(configStr, c)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	if c.Metadata.Schema == "" {
		return fmt.Errorf("schema not defined in metadata")
	}

	return nil
}

// loadSchemaFromFile reads a YAML schema file and loads it into the ConfigValidator_t struct.
//
// Parameters:
// - file: The path to the YAML schema file.
//
// Returns:
// - An error if reading the file or unmarshalling the schema fails.
func (cv *ConfigValidator_t) loadSchemaFromFile(file string) error {
	// Read the file
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = cv.loadSchemaFromString(data)
	if err != nil {
		return fmt.Errorf("error loading schema: %v", err)
	}

	return nil
}

// loadSchemaFromString unmarshals a YAML schema string into the ConfigValidator_t struct.
//
// Parameters:
// - schemaStr: The YAML schema string.
//
// Returns:
// - An error if unmarshalling the schema fails.
func (cv *ConfigValidator_t) loadSchemaFromString(schemaStr []byte) error {
	err := yaml.Unmarshal(schemaStr, cv)
	if err != nil {
		return fmt.Errorf("error unmarshalling Schema YAML: %v", err)
	}

	return nil
}

// validateConfig validates the configuration against the schema.
//
// Parameters:
// - c: A pointer to the Config_t struct containing the configuration.
//
// Returns:
// - An error if validation fails.
func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	errors := []error{}
	data := c.Data.(map[string]interface{})

	// Add missing attributes and default them
	cvo := cv.Schema
	cvo.recurValidateConfig(data, errors)

	// Check validity of the attributes
	cv.checkAttr(data, errors)

	// Check and report errors
	if len(errors) != 0 {
		fmt.Println("Validation failed:")
		for _, e := range errors {
			fmt.Printf("  %v\n", e)
		}
		fmt.Println("")
	}

	return nil
}

// addMissingAttr adds missing attributes to the configuration data and sets default values.
//
// Parameters:
// - ks: The key of the attribute.
// - vs: The schema field definition of the attribute.
// - data: The configuration data map.
func addMissingAttr(ks string, vs SchemaField_t, data map[string]interface{}) {
	if _, ok := data[ks]; !ok {
		if vs.Default == nil {
			switch vs.Type {
			case "string":
				data[ks] = ""
			case "integer":
				if vs.Min != 0 {
					data[ks] = int(vs.Min)
				} else {
					data[ks] = 0
				}
			case "float":
				if vs.Min != 0 {
					data[ks] = float64(vs.Min)
				} else {
					data[ks] = 0.0
				}
			case "boolean":
				data[ks] = false
			case "array":
				data[ks] = []interface{}{}
			case "object":
				data[ks] = map[string]interface{}{}
			default:
				log.Fatalf("field %s has an unknown type", ks)
			}
		} else {
			data[ks] = vs.Default
		}
	}
}

// recurValidateConfig recursively validates the configuration data against the schema
// for nested objects.
//
// Parameters:
// - data: The configuration data map.
// - e: A slice of errors to collect validation errors.
func (cv sf_t) recurValidateConfig(data map[string]interface{}, e []error) {
	// Add any attributes that are not provided
	for ks, vs := range cv {
		switch vs.Type {
		case "object":
			// loop through the attributes in this object and add the missing attributes
			cvo := cv[ks].Attributes
			for _, datao := range data[ks].(map[string]interface{}) {
				cvo.recurValidateConfig(datao.(map[string]interface{}), e)
			}
		default:
			addMissingAttr(ks, vs, data)
		}
	}
}

// checkAttr checks the attribute values in the configuration data against the schema.
//
// Parameters:
// - data: The configuration data map.
// - e: A slice of errors to collect validation errors.
func (cv ConfigValidator_t) checkAttr(data map[string]interface{}, e []error) {
	// Check attribute values
	for k, v := range data {
		val := cv.Schema[k]

		// Check Required and use Default if not set
		if val.Required && val.Default == nil && v == nil {
			e = append(e, fmt.Errorf("required field %s is empty", k))
		}

		// Check if Default is set and value is empty
		if val.Default != nil && v == nil {
			data[k] = val.Default
		}

		// Check field value types
		switch val.Type {
		case "string":
			// check if value is a string
			if _, ok := v.(string); !ok {
				e = append(e, fmt.Errorf("field %s is not a string", k))
			}
		case "integer":
			// check if value is an integer
			if _, ok := v.(int); !ok {
				e = append(e, fmt.Errorf("field %s is not an integer", k))
			}
			// check if integer within min max range
			if val.Min != 0 && v.(int) < int(val.Min) {
				e = append(e, fmt.Errorf("field %s is less than min value", k))
			}
			if val.Max != 0 && v.(int) > int(val.Max) {
				e = append(e, fmt.Errorf("field %s is greater than max value", k))
			}
		case "float":
			// Check if value is a float
			if _, ok := v.(float64); !ok {
				e = append(e, fmt.Errorf("field %s is not a float", k))
			}
		case "boolean":
			// check if value is a boolean
			if _, ok := v.(bool); !ok {
				e = append(e, fmt.Errorf("field %s is not a boolean", k))
			}
		case "array":
			// Check if value is an array
			if _, ok := v.([]interface{}); !ok {
				e = append(e, fmt.Errorf("field %s is not an array", k))
			}
			for i, o := range v.([]interface{}) {
				// Check if the option is of the correct type
				switch val.OptionType {
				case "string":
					if _, ok := o.(string); !ok {
						e = append(e, fmt.Errorf("field %s option %d is not a string", k, i))
					}
				case "integer":
					if _, ok := o.(int); !ok {
						e = append(e, fmt.Errorf("field %s option %d is not an integer", k, i))
					}
				case "float":
					if _, ok := o.(float64); !ok {
						e = append(e, fmt.Errorf("field %s option %d is not a float", k, i))
					}
				default:
					e = append(e, fmt.Errorf("field %s has an unknown option type", k))
				}
			}
		case "object":
			// Do nothing for now
		default:
			e = append(e, fmt.Errorf("field %s has an unknown type", k))
		}
	}
}
