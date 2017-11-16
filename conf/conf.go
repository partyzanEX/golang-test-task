package conf

import (
	"os"
	"encoding/json"
	"reflect"
	"strconv"
)

type Interface interface {
	Load() error
	Get(key string) (interface{}, error)
}

type Config struct {
	// storage
	data     map[string]interface{}
	FileName string
}

// load json file
func (c *Config) Load() error {
	file, err := os.Open(c.FileName)
	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	return decoder.Decode(&c.data)
}

// getting value by key
func (c Config) Get(key string) TypeConvert {
	if value, success := c.data[key]; success {
		v := Value{data: value}
		return v
	}

	return nil
}

// constructor
func NewConfig(fileName string) (*Config, error) {
	conf := Config{FileName: fileName}
	err := conf.Load()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

type TypeConvert interface {
	Int() int
	String() string
	Float() float64
	Bool() bool
}

type Value struct {
	data interface{}
}

// convert data to int
func (v Value) Int() int {
	var res = 0

	switch reflect.TypeOf(v.data).String() {
	case "float64":
		res = int(v.data.(float64))
	case "bool":
		b := int(0)
		if v.data.(bool) {
			b = int(1)
		}
		res = b
	case "string":
		b, err := strconv.ParseInt(v.data.(string), 10, 0)
		if err != nil {
			return 0
		}
		res = int(b)
	}

	return res
}

// convert data to string
func (v Value) String() string {
	var res string

	switch reflect.TypeOf(v.data).String() {
	case "float64":
		res = strconv.FormatFloat(v.data.(float64), 'g', -1, 64)
	case "bool":
		b := "false"
		if v.data.(bool) {
			b = "true"
		}
		res = b
	case "string":
		res = v.data.(string)
	}

	return res
}

// convert data to float64
func (v Value) Float() float64 {
	var res float64

	switch reflect.TypeOf(v.data).String() {
	case "float64":
		res = v.data.(float64)
	case "bool":
		b := 0.
		if v.data.(bool) {
			b = 1.
		}
		res = b
	case "string":
		res, _ = strconv.ParseFloat(v.data.(string), 64)
	}

	return res
}

// convert data to boolean
func (v Value) Bool() bool {
	var res bool

	switch reflect.TypeOf(v.data).String() {
	case "float64":
		b := false
		if v.data.(float64) > 0 {
			b = true
		}
		res = b
	case "bool":
		res = v.data.(bool)
	case "string":
		b := true
		if v.data.(string) == "" {
			b = false
		}
		res = b
	}

	return res
}
