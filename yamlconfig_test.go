package yamlconfig

import (
	"reflect"
	"testing"
)

/*
func TestYamlConfig1(t *testing.T) {

	myc := myConfig_t{}

	res := myConfig_t{
		Name:    "Marcel",
		City:    "Walnut Creek",
		State:   "CA",
		Id:      123,
		Options: []string{"foo", "baz"},
		Rabbitmq: map[string]RabbitMQ_t{
			"main":      {Host: "localhost", Port: 5672, User: "guest", Password: "guest", Vhost: "/"},
			"secondary": {Host: "localhost2", Port: 5672, User: "willem", Password: "waters", Vhost: "/dev"},
		},
		Mysql: Mysql_t{
			Host: "localhost",
			Port: 1234,
		},
		Redis: []Redis_t{
			{Host: "localhost", Db: 0},
			{Host: "localhost2", Db: 1},
		},
		MyArray: []Path_t{
			{Path: "./path1"},
			{Path: "./path2"},
			{Path: "./path3"},
		},
	}

	// simple test of LoadConfig()
	c := LoadConfig("./testfiles/example.config.yaml", &myc, "./schemas")
	newConf := c.(*myConfig_t)
	if !equalStructs(newConf, &res) {
		t.Errorf("MyConf: %+v\n", newConf)
		t.Errorf("Expected: %+v\n", res)
	}
}
*/

func TestYamlConfig2(t *testing.T) {

	myc := Workflow_t{}

	res := Workflow_t{
		Name:        "testwf02",
		Description: "This is a set of test workflows",
		Version:     "1.0",

		InitiatePaths: []WorkflowInitiate_t{
			{Path: "broadpath"},
		},

		/*
			Options: []string{"foo", "baz"},
			Rabbitmq: map[string]RabbitMQ_t{
				"main":      {Host: "localhost", Port: 5672, User: "guest", Password: "guest", Vhost: "/"},
				"secondary": {Host: "localhost2", Port: 5672, User: "willem", Password: "waters", Vhost: "/dev"},
			},
			Mysql: Mysql_t{
				Host: "localhost",
				Port: 1234,
			},
			Redis: []Redis_t{
				{Host: "localhost", Db: 0},
				{Host: "localhost2", Db: 1},
			},
			MyArray: []Path_t{
				{Path: "./path1"},
				{Path: "./path2"},
				{Path: "./path3"},
			},
		*/
	}

	// simple test of LoadConfig()
	c := LoadConfig("./testfiles/test_more.yaml", &myc, "./schemas")
	newConf := c.(*Workflow_t)
	if !equalStructs(newConf, &res) {
		t.Errorf("Workflow_t: %+v\n", newConf)
		t.Errorf("Expected: %+v\n", res)
	}
}

func equalStructs(a, b interface{}) bool {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.Kind() != bValue.Kind() {
		return false
	}

	if aValue.Kind() == reflect.Ptr {
		aValue = aValue.Elem()
		bValue = bValue.Elem()
	}

	if aValue.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < aValue.NumField(); i++ {
		if !reflect.DeepEqual(aValue.Field(i).Interface(), bValue.Field(i).Interface()) {
			return false
		}
	}

	return true
}
