package common

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	Config       map[string]string
	ConfigString map[string]string
	ConfigInt    map[string]int
	ConfigBool   map[string]bool
)

type ConfigModel struct {
	FileName string
}

type IConfigSetup interface {
	Open() *error
}

func NewConfig(model ConfigModel) IConfigSetup {
	return ConfigModel{FileName: model.FileName}
}

func (c ConfigModel) Open() *error {

	splits := strings.Split(filepath.Base(c.FileName), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(c.FileName))

	if err := viper.ReadInConfig(); err != nil {
		return &err
	}

	if !viper.IsSet("server.mode") {
		newError := errors.New("please define server.mode in your toml file")
		return &newError
	}

	environment := viper.GetString("server.mode")
	mapString := make(map[string]string)
	mapInt := make(map[string]int)
	mapBool := make(map[string]bool)

	for k, v := range viper.AllSettings() {

		if strings.EqualFold(k, environment) || strings.EqualFold(k, "server") {
			reflectV := reflect.ValueOf(v)
			if reflectV.Kind() == reflect.Map {

				for _, v := range reflectV.MapKeys() {
					mapIndex := reflectV.MapIndex(v)
					stringValue, ok := v.Interface().(string)
					if !ok {
						newError := errors.New("cannot found type")
						return &newError
					}

					switch value := mapIndex.Interface().(type) {
					case string:
						mapString[stringValue] = value
					case int64:
						mapInt[stringValue] = int(value)
					case bool:
						mapBool[stringValue] = value
					default:
						newError := errors.New("cannot found type")
						return &newError
					}

					ConfigString = mapString
					ConfigInt = mapInt
					ConfigBool = mapBool

				}

			}
		}
	}

	return nil
}
