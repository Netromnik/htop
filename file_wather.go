package top

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Структура, представляющая проверку целостности файлов
type FileHealchek struct {
	wath_dir string       // Путь к директории, которую нужно отслеживать
	pool     *ProgramPool // Пул программ, которые нужно запускать/останавливать
	done     chan int     // Канал для завершения работы

}

// Функция для создания новой проверки целостности файлов
func newFileHealchek(wath_dir string, pool *ProgramPool) *FileHealchek {
	obj := &FileHealchek{
		wath_dir: wath_dir,
		done:     make(chan int),
		pool:     pool,
	}
	return obj
}

// Инициализация проверки целостности файлов
func (self *FileHealchek) init() {
	self.load()    // Загрузка программ
	go self.wath() // Запуск отслеживания изменений в файловой системе
	self.done <- 1 // Отправка сообщения о завершении работы
}

// Загрузка программ из директории
func (self *FileHealchek) load() error {
	path := createGlobalPath(self.wath_dir) // Получение пути к директории
	files, err := ioutil.ReadDir(path)      // Получение списка файлов в директории
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		_, ok := getParser(file.Name()) // Проверка, что файл можно запустить
		if ok {
			self.pool.register(filepath.Join(path, file.Name())) // Регистрация программы в пуле
		}
	}
	return nil
}

// Отслеживание изменений в файловой системе
func (self *FileHealchek) wath() error {
	watcher, err := fsnotify.NewWatcher() // Создание нового объекта для отслеживания изменений
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				_, ok := getParser(event.Name) // Проверка, что файл можно запустить
				if ok {
					if event.Op&fsnotify.Create == fsnotify.Create { // Обработка события создания файла
						self.pool.register(event.Name)
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove { // Обработка события удаления файла
						self.pool.unregister(event.Name)
					}
					if event.Op&fsnotify.Write == fsnotify.Write { // Обработка события записи в файл
						self.pool.reload(event.Name)
					}
					if event.Op&fsnotify.Rename == fsnotify.Rename { // Обработка события переименования файла
						self.pool.reload(event.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
				self.done <- 1
			}
		}
	}()

	err = watcher.Add(createGlobalPath(self.wath_dir)) // Добавление директории для отслеживания изменений
	if err != nil {
		return err
	}
	<-self.done // Ожидание завершения работы
	return nil
}
