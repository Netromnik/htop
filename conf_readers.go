package top

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type ProgrammData struct {
	name string
	path string
	args string
}

func newJson(path string) (ProgrammData, error) {
	var data ProgrammData

	// Чтение файла конфигурации
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}

	// Декодирование JSON-данных в map[string]interface{}
	var dataClean map[string]interface{}
	err = json.Unmarshal(bytes, &dataClean)
	if err != nil {
		return data, err
	}

	// Заполнение полей структуры
	data.name = dataClean["name"].(string)
	data.path = dataClean["path"].(string)
	data.args = dataClean["args"].(string)

	return data, nil
}

func newIni(path string) (ProgrammData, error) {
	var data ProgrammData

	// Чтение файла конфигурации
	cfg, err := ini.Load(path)
	if err != nil {
		return data, err
	}

	// Чтение значений параметров из файла конфигурации
	data.name = cfg.Section("").Key("name").String()
	data.path = cfg.Section("").Key("path").String()
	data.args = cfg.Section("").Key("args").String()

	return data, nil
}

var FILE_PARSER = map[string]func(string) (ProgrammData, error){
	".ini":  newIni,
	".json": newJson,
}

// Возврат парсера
func getParser(path string) (func(string) (ProgrammData, error), bool) {
	parser, ok := FILE_PARSER[filepath.Ext(path)]
	return parser, ok
}
