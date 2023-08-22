package top

import (
	"fmt"
	"testing"
	"time"
)

func TestNewProgramPool(t *testing.T) {
	poll := newProgramPool()

	if poll == nil {
		t.Error("Expected non-nil ProgramPool, but got nil")
	}

	if poll.list_app == nil {
		t.Error("Expected non-nil list_app, but got nil")
	}
}

func TestNewProgram(t *testing.T) {

	testCasesOne := []struct {
		file     string
		expected error
	}{
		{"test.txt", fmt.Errorf("File %s is not corect format", "test.txt")},
	}

	for _, tc := range testCasesOne {
		t.Run(tc.file, func(t *testing.T) {
			_, err := newProgram(tc.file)
			if err == nil {
				t.Errorf("Expected error case %s", tc.file)
			}
		})
	}

	testCasesTwo := []struct {
		file     string
		expected error
	}{
		{"config.ini", nil},
		{"data.json", nil},
	}

	for _, tc := range testCasesTwo {
		t.Run(tc.file, func(t *testing.T) {
			_, err := newProgram(tc.file)
			if err == nil {
				t.Errorf("for %s, but got %v", tc.file, err)
			}
		})
	}
	testCasesTree := []struct {
		file     string
		expected error
	}{
		{"test/config.ini", nil},
		{"test/data.json", nil},
	}

	for _, tc := range testCasesTree {
		t.Run(tc.file, func(t *testing.T) {
			_, err := newProgram(tc.file)
			if err != nil {
				t.Errorf("for %s, but got %v", tc.file, err)
			}
		})
	}

}
func TestCmdProgramm(t *testing.T) {
	// TODO: prepare test data
	testCases := []struct {
		cfg      ProgrammData
		expected error
	}{
		{ProgrammData{}, fmt.Errorf("no path %s", "")},
		{ProgrammData{path: "dir.exe"}, nil},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc.cfg), func(t *testing.T) {
			_, err := cmdProgramm(&tc.cfg)
			if fmt.Sprintln(err) != fmt.Sprintln(tc.expected) {
				t.Errorf("Expected error %v for %+v, but got %v", tc.expected, tc.cfg, err)
			}
		})
	}
}

func TestProgramPool_RegisterNeg(t *testing.T) {
	// TODO: prepare test data

	testCases := []struct {
		file     string
		expected error
	}{
		{"config.ini", nil},
		{"data.json", nil},
		{"test.txt", fmt.Errorf("File %s is not corect format", "test.txt")},
	}

	Observer = NewMultiChanObserver()
	poll := newProgramPool()

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			err := poll.register(tc.file)
			if err == nil {
				t.Errorf("Expected error %v for %s, but got %v", tc.expected, tc.file, err)
			}
		})
	}
}

func TestProgramPool_RegisterPos(t *testing.T) {

	testCases := []struct {
		file     string
		name     string
		expected error
	}{
		{"test/config.ini", "My App", nil},
		{"test/data.json", "My App", nil},
	}

	poll := newProgramPool()

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			err := poll.register(tc.file)
			if err != nil {
				t.Errorf("Expected error %v for %s, but got %v", tc.expected, tc.file, err)
			}
			_, ok := poll.list_app[tc.name]
			if ok != false {
				t.Errorf("Poll pg not found %s: %s", tc.file, tc.name)
			}
		})
	}
}

func TestProgramPool_UnregisterNeg(t *testing.T) {
	// TODO: prepare test data

	testCases := []struct {
		file     string
		expected error
	}{
		{"config.ini", fmt.Errorf("File %s is not registered", "config.ini")},
		{file: "data.json", expected: fmt.Errorf("File %s is not registered", "data.json")},
	}

	poll := newProgramPool()

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			err := poll.unregister(tc.file)
			if fmt.Sprintln(err) != fmt.Sprintln(tc.expected) {
				t.Errorf("Expected error %v for %s, but got %v", tc.expected, tc.file, err)
			}
		})
	}
}

func TestProgramPool_UnregisterPos(t *testing.T) {
	// TODO: prepare test data

	testCases := []struct {
		file     string
		expected error
	}{
		{"test/config.ini", nil},
		{file: "test/data.json", expected: nil},
	}

	poll := newProgramPool()

	poll.register("test/config.ini")
	poll.register("test/data.json")

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			err := poll.unregister(tc.file)
			if err != tc.expected {
				t.Errorf("for %s, but %v", tc.file, err)
			}
		})
	}
}

func TestProgramPool_Reload(t *testing.T) {
	// TODO: prepare test data

	testCases := []struct {
		file     string
		expected error
	}{
		// TODO: add test cases
	}

	poll := newProgramPool()

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			err := poll.reload(tc.file)
			if err != tc.expected {
				t.Errorf("Expected error %v for %s, but got %v", tc.expected, tc.file, err)
			}
		})
	}
}

func TestProgramPool_Close(t *testing.T) {
	// TODO: prepare test data

	poll := newProgramPool()

	err := poll.close()
	if err != nil {
		t.Errorf("Expected nil error for close(), but got %v", err)
	}
}

func TestProgramPoolNotify(t *testing.T) {
	Observer = NewMultiChanObserver()
	poll := newProgramPool()
	poll.register("test/config.ini")

	test := Observer.AddChannel("test")

	time.Sleep(100 * time.Millisecond)

	if len(Observer.events) != 0 {
		t.Error("event not 0 is ", len(Observer.events))
	}
	if len(*test) != 1 {
		t.Error("event not 1 is ", len(*test))
	}
	poll.unregister("test/config.ini")
	time.Sleep(120 * time.Millisecond)

	if len(Observer.events) != 0 {
		t.Error("event not 0 is ", len(Observer.events))
	}
	if len(*test) != 2 {
		t.Error("event not 2 is ", len(*test))
	}

}

func TestProgramPoolNotifyLoads(t *testing.T) {
	poll := newProgramPool()
	fh := newFileHealchek("test", poll)
	fh.load()
	test := Observer.AddChannel("test")

	time.Sleep(100 * time.Millisecond)

	if len(Observer.events) == 0 {
		t.Error("event not 0 is ", len(Observer.events))
	}
	if len(*test) == 1 {
		t.Error("event not 1 is ", len(*test))
	}
	time.Sleep(100 * time.Millisecond)

	if len(Observer.events) == 0 {
		t.Error("event not 0 is ", len(Observer.events))
	}
	if len(*test) == 2 {
		t.Error("event not 2 is ", len(*test))
	}

}
