package yamlconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

var schemaMap = make(map[string]ConfigValidator_t)

// LoadConfig reads a YAML configuration file and loads it into a custom struct.
// It returns a map representation of the configuration.
//
// Parameters:
// - file: The path to the YAML configuration file.
// - customStruct: A pointer to the custom struct where the configuration will be loaded.
// - schemaPath: The path to the schema file. Multiple paths can be provided in an array.
//
// Returns:
// - A map representation of the configuration.
func LoadConfig(file string, customStruct interface{}, schemaPath ...string) any {

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

	schema := getSchema(datameta.Metadata.SchemaVersion, schemaPath)

	// Unmarshal the config schema
	err = yaml.Unmarshal(data, configStruct)
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

	// prettyPrint(configStruct, "ConfigStruct")
	res := reflect.New(reflect.TypeOf(customStruct).Elem()).Interface()
	err = mapstructure.Decode(configStruct.Data.(map[string]interface{}), &res)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	prettyPrint(res, "Result")
	return res
}

func prettyPrint(s interface{}, n string) {
	log.Print("=======================================================================================")
	log.Printf("===> %s out: %+v", n, s)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("%s out:", n)
	log.Print(string(data))
}

func getSchema(schemaFile string, schemaPath []string) ConfigValidator_t {
	var (
		err        error
		schemaData []byte
		file       string
	)

	// Check if version is in schemaMap, if so return it
	if sm, ok := schemaMap[schemaFile]; ok {
		return sm
	}

	// Scan the folders for the schema file
	for _, sp := range schemaPath {
		// If the schema definition file (schema_version) does not exist, we dont support it
		file = fmt.Sprintf("%s/%s.yaml", sp, schemaFile)
		schemaData, err = os.ReadFile(file)
		if err != nil {
			switch err {
			case os.ErrNotExist:
				continue
			default:
				log.Fatalf("Schema version %s is not supported in file %s", schemaFile, file)
			}
		}
	}

	if schemaData == nil {
		log.Fatalf("Schema file %s not found", file)
	}

	// Unmarshal the validation schema
	schema := ConfigValidator_t{}
	err = yaml.Unmarshal(schemaData, &schema)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Save in cache
	schemaMap[schemaFile] = schema

	return schema
}
