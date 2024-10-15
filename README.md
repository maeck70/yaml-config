## yamlconfig: Configuration Loader and Validator

This is the README file for the `yamlconfig` package, which provides functionality for loading and validating YAML configuration files.

**Features:**

- Loads YAML configuration files into custom structs.
- Validates the configuration against a YAML schema.
- Handles missing attributes and sets default values if provided in the schema.
- Provides detailed error messages for invalid configuration.

**Example**

A working example of using this module can be found here [https://github.com/maeck70/yaml-config-example].

**Installation:**

1. Clone this repository.
2. Run `go mod download` to download dependencies.

**Usage:**

1.  Import the package:

    ```go
    import "github.com/maeck70/yamlconfig"
    ```

2.  Define your custom struct to hold the configuration data.

3.  Use the `yamlconfig.LoadConfig` function to load and validate the configuration:

    ```go
    type MyConfig struct {
        // ... Your configuration fields here
    }

    configPath := "config.yaml"
    config := MyConfig{}
    configMap, err := yamlconfig.LoadConfig(configPath, &config)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Access the configuration data in the configMap
    fmt.Println(configMap["key1"])
    ```

**Schema Definition:**

The schema for your configuration should be a separate YAML file defining the structure and validation rules for each field. The schema itself is a map where keys are field names and values are objects describing the field properties.

Here's an example schema structure:

```yaml
# schema.yaml
metadata:
  version: 0.1alpha
schema:
  Name:
    type: string
    description: Name of the user
    required: true
  City:
    type: string
    description: City of the user
    required: true
  Id:
    type: integer
    description: Unique identifier
    required: false
    min: 100
    max: 200    
  Options:
    type: array
    description: List of options
    optiontype: string
    required: false
  Rabbitmq:
    type: object
    description: Rabbitmq configuration
    attributes:
      host:
        type: string
        description: Hostname
        required: true
      port:
        type: integer
        description: Port number
        default: 5672
      user:
        type: string
        description: Username
        required: true 
      password:
        type: string
        description: Password
        required: true
      vhost:
        type: string
        description: Virtual host
        default: /
```

**Validation Rules:**

The schema can define validation rules for each field:

- `type`: The expected data type of the field (string, integer, float, boolean, array, or object).
- `required`: A boolean indicating if the field is mandatory.
- `default`: A default value to be used if the field is missing in the configuration file.
- `min` and `max`: Integer values for minimum and maximum allowed values (applicable to integers and floats).
- `optionType`: Specifies the type of elements allowed in an array (applicable to the "array" type).
- `description`: Provides a description to the field.

The `yamlconfig` package performs validation based on these rules and will report any errors encountered during the loading process.


**The actuual Config Definition:**

The schema for your configuration is provided in the configuration yaml file under the metadata section.

Here's an example schema structure:

```yaml
# config file
# uses example.schema.yaml

metadata:
  schema: example.schema.yaml
data:
  Name: William
  City: New York
  Id: 12345
  Options:
    - Foo
    - "11"
    - Bar
  Rabbitmq:
    main:
      host: localhost
      user: guest
      password: guest
    baskets: 
      host: localhost
      port: 5672
      user: guest
      password: guest
      vhost: /external
```

**Additional Notes:**

- This package utilizes reflection for handling different custom struct types.
- Error handling is incorporated for file reading, unmarshalling, and validation.

**Contribution:**

Feel free to submit pull requests to improve this package or add additional features.
