package top

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func createGlobalPath(path ...string) string {
	// Получение текущего рабочего каталога
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	// Создание пути до директории на основе аргументов программы
	dirPath := filepath.Join(append([]string{exPath}, path...)...)

	// Получение абсолютного пути до директории
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Println("Ошибка при получении абсолютного пути до директории:", err)
		os.Exit(1)
	}

	return filepath.Clean(absPath)
}

func isProcRunning(name string) (int, error) {
	count := 0

	cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return count, err
	}

	return bytes.Count(out, []byte(name)), nil
}

// cmdProgramm - функция для создания команды программы
func cmdProgramm(cfg *ProgrammData) (*exec.Cmd, error) {
	// Create program

	if cfg.path == "" {
		return nil, fmt.Errorf("no path %s", cfg.path)
	}

	// Create the command with the parameters
	return exec.Command(cfg.path, cfg.args), nil
}
