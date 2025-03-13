package logging

import (
	"log/slog"
	"os"
)

func createLoggerFile(logFilePath string) *os.File {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
	//defer file.Close()
}

func CreateLogger(logFilePath string, debugFlag *bool) *slog.Logger {
	var programLevel = new(slog.LevelVar)
	//stdoutHandler := slog.NewJSONHandler(os.Stdout, nil)
	fileHandler := slog.NewJSONHandler(createLoggerFile(logFilePath), &slog.HandlerOptions{Level: programLevel})
	logger := slog.New(fileHandler)
	slog.SetDefault(logger)

	if *debugFlag {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}
	return logger
}
