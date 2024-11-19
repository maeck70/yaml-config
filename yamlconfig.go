package yamlconfig

import (
	"embed"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*.yaml
var schemaDir embed.FS

// LoadConfig reads a YAML configuration file and loads it into a custom struct.
// It returns a map representation of the configuration.
//
// Parameters:
// - file: The path to the YAML configuration file.
// - customStruct: A pointer to the custom struct where the configuration will be loaded.
//
// Returns:
// - A map representation of the configuration.
func LoadConfig(file string, customStruct interface{}) any {

	// Build the full config struct
	dataStruct := reflect.ValueOf(customStruct)
	configStruct := &Config_t{}

	// Check if data is a pointer
	if dataStruct.Kind() == reflect.Ptr {
		configStruct.Data = dataStruct.Elem()
	} else {
		return nil
	}

	// Load the workflow from the yaml file
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	datameta := Config_t{}
	err = yaml.Unmarshal(data, &datameta)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// If the schema definition file (schema_version) does not exist, we dont support it
	schemadata, err := schemaDir.ReadFile("schemas/" + datameta.Metadata.SchemaVersion + ".yaml")
	if err != nil {
		log.Fatalf("Schema version %s is not supported in file %s", datameta.Metadata.SchemaVersion, file)
	}

	// Unmarshal the config schema
	err = yaml.Unmarshal(data, customStruct)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the validation schema
	schema := ConfigValidator_t{}
	err = yaml.Unmarshal(schemadata, &schema)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Check if the schema id matches the workflow id
	if datameta.Metadata.Id != schema.Metadata.Id {
		log.Fatalf("Invalid workflow schema file, schema should contain id: %s", schema.Metadata.Id)
	}

	// Check if the schema version matches the internal schema validator
	if schema.Metadata.SchemaVersion != datameta.Metadata.SchemaVersion {
		log.Fatalf("Workflow schema version %s does not match internal schema validator %s",
			datameta.Metadata.SchemaVersion, schema.Metadata.SchemaVersion,
		)
	}

	schema.validateConfig(configStruct)

	// Convert the map structure back to the original struct
	res := reflect.New(reflect.TypeOf(customStruct).Elem()).Interface()
	mapstructure.Decode(configStruct.Data.(map[string]interface{}), &res)
	return res
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
