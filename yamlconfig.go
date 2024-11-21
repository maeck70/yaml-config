package yamlconfig

import (
	"embed"
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
	err = yaml.Unmarshal(data, configStruct)
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
	if datameta.Metadata.Id != configStruct.Metadata.Id {
		log.Fatalf("Invalid workflow schema file, schema should contain id: %s", configStruct.Metadata.Id)
	}

	// Check if the schema version matches the internal schema validator
	if schema.Metadata.SchemaVersion != datameta.Metadata.SchemaVersion {
		log.Fatalf("Workflow schema version %s does not match internal schema validator %s",
			datameta.Metadata.SchemaVersion, configStruct.Metadata.SchemaVersion,
		)
	}

	schema.validateConfig(configStruct)

	// Convert the map structure back to the original struct
	res := reflect.New(reflect.TypeOf(customStruct).Elem()).Interface()
	mapstructure.Decode(configStruct.Data.(items_t), &res)
	return res
}

func (cv ConfigValidator_t) validateConfig(c *Config_t) error {
	// Validate the config against the schema

	pad := ""

	for k, v := range c.Data.(map[string]interface{}) {
		switch v.(type) {
		case []interface{}:
			log.Printf("List - Key: %s, Value: %v\n", k, v)
		case map[string]interface{}:
			log.Printf("Map - Key: %s, Value: %v\n", k, v)

			if _, ok := v.(map[string]interface{})["map"]; ok {
				mapRecursive(pad, v.(map[string]interface{})["map"].(map[string]interface{}))
			}

			if _, ok := v.(map[string]interface{})["array"]; ok {
				arrayRecursive(pad, v.(map[string]interface{})["array"].([]interface{}))
			}

		default:
			log.Printf("Default - Key: %s, Value: %v\n", k, v)
		}
	}

	return nil
}

func mapRecursive(pad string, data map[string]interface{}) {
	pad = pad + "  "
	for k, v := range data {
		switch v.(type) {
		case []interface{}:
			log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
		case map[string]interface{}:
			log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
			// mapRecursive(pad)
		default:
			log.Printf(pad+"Default - Key: %s, Value: %v\n", k, v)
		}
	}
}

func arrayRecursive(pad string, data []interface{}) {
	pad = pad + "  "
	for k, v := range data {
		switch v.(type) {
		case []interface{}:
			log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
		case map[string]interface{}:
			log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
			// mapRecursive(pad)
		default:
			log.Printf(pad+"Default - Key: %s, Value: %v\n", k, v)
		}
	}
}

/*
func mapRecursive(pad string, data map[string]interface{}) {
	pad = pad + "  "
	for k, v := range data {
		switch v.(type) {
		case []interface{}:
			log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
		case map[string]interface{}:
			log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
			// mapRecursive(pad)
		default:
			log.Printf(pad+"Default - Key: %s, Value: %v\n", k, v)
		}
	}
}

func arrayRecursive(pad string, data []interface{}) {
	pad = pad + "  "
	for k, v := range data {
		switch v.(type) {
		case []interface{}:
			log.Printf(pad+"List - Key: %s, Value: %v\n", k, v)
		case map[string]interface{}:
			log.Printf(pad+"Map - Key: %s, Value: %v\n", k, v)
			// mapRecursive(pad)
		default:
			log.Printf(pad+"Default - Key: %s, Value: %v\n", k, v)
		}
	}
}
*/
