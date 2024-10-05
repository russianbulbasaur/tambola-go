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
	gameId string
	file   *os.File
	writer *bufio.Writer
}

func NewTambolaLogger(gameContext context.Context) *TambolaLogger {
	directory := "logs"
	gameId := gameContext.Value("game_id")
	err := os.Mkdir(directory, 0777)
	fileName := fmt.Sprintf("%s.log", gameId)
	filePath := filepath.Join(directory, fileName)
	file, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0700)
	if err != nil {
		log.Fatalln("Unable to create logs file", err)
	}
	writer := bufio.NewWriter(file)
	return &TambolaLogger{gameId: gameId.(string),
		file:   file,
		writer: writer}
}

func (tl *TambolaLogger) Log(text string) {
	gameId := tl.gameId
	_, err := tl.writer.Write([]byte(formatLogEntry(text, gameId)))
	err = tl.writer.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(gameId, ":", text)
}

func formatLogEntry(text string, gameId any) string {
	return fmt.Sprintf(
		"%v \n%s \n----------- \n", time.TimeOnly, text)
}
