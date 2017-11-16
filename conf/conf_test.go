package conf

import (
	"testing"
	//"fmt"
)

var configFile = "./test/conf.json"
var errorFile = "./test/error.json"

func TestNewConfig(t *testing.T) {

	t.Run("No Errors", func(t *testing.T) {
		_, err := NewConfig(configFile)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Error file", func(t *testing.T) {
		_, err := NewConfig(errorFile)
		if err == nil {
			t.Error("no file errors")
		}
	})

	t.Run("Has Errors", func(t *testing.T) {
		_, err := NewConfig("./conf.json")
		if err == nil {
			t.Error("no errors")
		}
	})
}

func TestConfig_Load(t *testing.T) {
	conf := Config{}
	err := conf.Load()
	if err == nil {
		t.Error("no errors")
	}

	conf.FileName = "conf.json"
	err = conf.Load()
	if err == nil {
		t.Error("no errors")
	}

	conf.FileName = configFile
	err = conf.Load()
	if err != nil {
		t.Error(err)
	}
}

func TestConfig_Get(t *testing.T) {
	conf, _ := NewConfig(configFile)

	if res := conf.Get("name"); res.String() != "Name" || res.Int() != 0 || !res.Bool() || res.Float() != 0. {
		t.Error("Name Error")
	}

	if res := conf.Get("data").String(); res != "test" {
		t.Error("Data Error")
	}

	if res := conf.Get("number"); res.Float() != 10. || res.Int() != 10 || !res.Bool() || res.String() != "10" {
		t.Error("number error")
	}

	if res := conf.Get("bool"); res.Bool() != false || res.String() != "false" || res.Int() != 0 || res.Float() != 0. {
		t.Error("bool error")
	}
}
