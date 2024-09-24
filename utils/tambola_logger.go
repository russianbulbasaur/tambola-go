package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TambolaLogger struct {
	gameContext context.Context
	writer      *bufio.Writer
}

func NewTambolaLogger(gameContext context.Context) *TambolaLogger {
	directory := "logs"
	err := os.Mkdir(directory, 0777)
	fileName := fmt.Sprintf("%v.log", time.DateOnly)
	filePath := filepath.Join(directory, fileName)
	file, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0700)
	if err != nil {
		log.Fatalln("Unable to create logs file", err)
	}
	writer := bufio.NewWriter(file)
	return &TambolaLogger{gameContext: gameContext,
		writer: writer}
}

func (tl *TambolaLogger) Log(text string) {
	gameId := tl.gameContext.Value("game_id")
	_, err := tl.writer.Write([]byte(formatLogEntry(text, gameId)))
	err = tl.writer.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(gameId, ":", text)
}

func formatLogEntry(text string, gameId any) string {
	return fmt.Sprintf(
		"%v \nGame %d \n%s \n----------- \n", time.TimeOnly, gameId, text)
}
