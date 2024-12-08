package logger

import (
	"fmt"
	"os"
)

// createLogFile создает файл с логами и возвращает его FileWriteSyncer
func createLogFile(filePathLog string) *os.File {
	// Логирование в файл
	// Если параметры не указаны, задаем значения по умолчанию
	if filePathLog == "" {
		// Путь к файлу логов по умолчанию (текущая директория)
		filePathLog = "app.log"
	}
	// Создаем директорию, если её нет
	logDir := filepath.Dir(filePathLog)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("failed to create log directory: %s\n", err.Error())
	}

	//Открываем или создаем файл для записи логов
	logFile, err := os.OpenFile(filePathLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("failed to open log file: %s\n", err.Error())
	}
	
	return logFile
}
