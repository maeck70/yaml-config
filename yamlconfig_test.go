package yamlconfig

import (
	"testing"
)

type RabbitMQ_t struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

type Mysql_t struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Redis_t struct {
	Host     string `yaml:"host"`
	Database int    `yaml:"db"`
}

type myConfig_t struct {
	Name     string                `yaml:"Name"`
	City     string                `yaml:"City"`
	State    string                `yaml:"State"`
	Id       int                   `yaml:"Id"`
	Options  []string              `yaml:"Options"`
	Rabbitmq map[string]RabbitMQ_t `yaml:"Rabbitmq"`
	Mysql    Mysql_t               `yaml:"Mysql"`
	Redis    []Redis_t             `yaml:"Redis"`
}

func TestYamlConfig(t *testing.T) {
	myc := myConfig_t{}

	// simple test of LoadConfig()
	c := LoadConfig("testfiles/example.config.yaml", &myc)
	newConf := c.(*myConfig_t)
	if newConf.Name != "MyName" {
		t.Errorf("MyConf: %+v\n", newConf)
	}
}

/*
func equal(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
*/
