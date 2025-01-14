package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	MessageMap map[string]map[string]string
)

type MessageModel struct {
	Path     string
	FileName string
}

type IMessageConfig interface {
	Setup() *error
}

func NewMessageConfig(model MessageModel) IMessageConfig {
	return MessageModel{
		Path:     model.Path,
		FileName: model.FileName,
	}
}

func (m MessageModel) Setup() *error {

	fullPath := filepath.Join(m.Path, m.FileName)
	file, err := os.Open(fullPath)
	if err != nil {
		return &err
	}

	byteJson, err := ioutil.ReadAll(file)
	if err != nil {
		return &err
	}

	if err := json.Unmarshal(byteJson, &MessageMap); err != nil {
		return &err
	}

	return nil
}
