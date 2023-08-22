package programmactive

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewJson(t *testing.T) {
	// Создание временного файла с JSON-данными
	file, err := ioutil.TempFile("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	jsonData := `{
		"name": "My App",
		"path": "/usr/local/bin/myapp",
		"args": "--config config.json"
	}`

	err = ioutil.WriteFile(file.Name(), []byte(jsonData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Тестирование функции newJson()
	data, err := newJson(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if data.name != "My App" {
		t.Errorf("Expected name to be 'My App', but got '%s'", data.name)
	}

	if data.path != "/usr/local/bin/myapp" {
		t.Errorf("Expected path to be '/usr/local/bin/myapp', but got '%s'", data.path)
	}

	if data.args != "--config config.json" {
		t.Errorf("Expected args to be '--config config.json', but got '%s'", data.args)
	}
}

func TestNewIni(t *testing.T) {
	// Создание временного файла с INI-данными
	file, err := ioutil.TempFile("", "config.ini")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	iniData := `
name = My App
path = /usr/local/bin/myapp
args = --config config.json`

	err = ioutil.WriteFile(file.Name(), []byte(iniData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Тестирование функции newIni()
	data, err := newIni(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if data.name != "My App" {
		t.Errorf("Expected name to be 'My App', but got '%s'", data.name)
	}

	if data.path != "/usr/local/bin/myapp" {
		t.Errorf("Expected path to be '/usr/local/bin/myapp', but got '%s'", data.path)
	}

	if data.args != "--config config.json" {
		t.Errorf("Expected args to be '--config config.json', but got '%s'", data.args)
	}
}

func TestGetParserOk(t *testing.T) {

	testCases := []struct {
		path     string
		expected bool
	}{
		{"config.ini", true},
		{"data.json", true},
		{"test.txt", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			_, ok := getParser(tc.path)
			if ok != tc.expected {
				t.Errorf("Expected ok=%v for %s, but got %v", tc.expected, tc.path, ok)
			}
		})
	}
}
