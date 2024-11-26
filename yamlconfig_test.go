package yamlconfig

import (
	"testing"
)

func TestYamlConfig(t *testing.T) {
	myc := myConfig_t{}

	// simple test of LoadConfig()
	c := LoadConfig("./testfiles/example.config.yaml", &myc, "./schemas")
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
