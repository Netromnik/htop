package top

import (
	"fmt"
	"os/exec"
)

// Program - структура, описывающая программу
type Program struct {
	cmd      *exec.Cmd
	cfg      *ProgrammData
	filepath string
}

// newProgram - конструктор для создания новой программы
func newProgram(file string) (*Program, error) {
	parser, ok := getParser(file)
	if ok == false {
		return nil, fmt.Errorf("File %s is not corect format", file)
	}

	cfg, err := parser(file)

	if err != nil {
		return nil, fmt.Errorf("File parse %s err %s", file, err)
	}

	return &Program{
		cfg:      &cfg,
		filepath: file,
		cmd:      nil,
	}, nil
}

// getKey - метод для получения ключа программы
func (self *Program) getKey() string {
	return self.filepath
}

// InterfaceProgramPool - интерфейс для пула программ
type InterfaceProgramPool interface {
	register(file string) error
	reload(file string) error
	unregister(file string) error
	close() error
}

// ProgramPool - структура, описывающая пул программ
type ProgramPool struct {
	InterfaceProgramPool
	list_app map[string]*Program
}

// newProgramPool - конструктор для создания нового пула программ
func newProgramPool() *ProgramPool {
	poll := ProgramPool{
		list_app: make(map[string]*Program),
	}
	return &poll
}

// getWithP - метод для получения программы по ключу
func (self *ProgramPool) getWithP(pm *Program) (*Program, bool) {
	pm, ok := self.list_app[pm.getKey()]
	return pm, ok
}

// getWithFile - метод для получения программы по имени файла
func (self *ProgramPool) getWithFile(file string) (*Program, error) {
	f, err := newProgram(file)
	if err != nil {
		return nil, err
	}
	app, ok := self.getWithP(f)

	if ok {
		return app, nil
	}

	return nil, nil

}

// setWithP - метод для добавления программы в пул
func (self *ProgramPool) setWithP(pm *Program) {
	self.list_app[pm.getKey()] = pm
}

// register - метод для регистрации новой программы в пуле
func (self *ProgramPool) register(file string) error {
	// Check if the file is already registered
	prm, err := self.getWithFile(file)
	if err == nil && prm == nil {
		// Is not in register
		prm, err = newProgram(file)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Start the command and add it to the list of registered apps
	prm.cmd, err = cmdProgramm(prm.cfg)

	if err != nil {
		return fmt.Errorf("%s: %s", prm.filepath, err)
	}

	if prm.cmd != nil {
		pr, ok := self.getWithP(prm)
		if ok && pr != nil {
			pr.cmd.Process.Kill()
		}

		err := prm.cmd.Start()
		if err != nil {
			return err
		}
		self.setWithP(prm)
	} else {
		return fmt.Errorf("Program not work in file %s", file)
	}

	Observer.Notify(prm)
	return nil
}

// unregister - метод для удаления программы из пула
func (self *ProgramPool) unregister(file string) error {
	// Check if the file is registered
	app, err := self.getWithFile(file)
	if err != nil {
		return fmt.Errorf("File %s is not registered", file)
	}

	cmd := app.cmd

	// Check if the command is running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		// Stop the command and remove it from the list of registered apps
		err = cmd.Process.Kill()
	}

	if err != nil {
		return err
	}
	delete(self.list_app, file)

	Observer.Notify(nil)
	return nil
}

// reload - метод для перезагрузки программы в пуле
func (self *ProgramPool) reload(file string) error {
	// Check if the file is registered
	newPr, err := newProgram(file)
	if err != nil {
		return err
	}
	old, ok := self.getWithP(newPr)

	if ok && *old.cfg != *newPr.cfg {
		self.unregister(file)
		self.register(file)
	}

	return nil
}

// close - метод для остановки всех зарегистрированных программ в пуле и очистки списка зарегистрированных программ
func (self *ProgramPool) close() error {
	for _, pm := range self.list_app {
		err := pm.cmd.Process.Kill()
		if err != nil {
			return err
		}
	}
	self.list_app = make(map[string]*Program)

	return nil
}
