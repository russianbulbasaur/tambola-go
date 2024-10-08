package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TambolaLogger struct {
	gameId     string
	file       *os.File
	writer     *bufio.Writer
	LogChannel chan string
	wg         *sync.WaitGroup
}

func NewTambolaLogger(gameContext context.Context, wg *sync.WaitGroup) *TambolaLogger {
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
	return &TambolaLogger{
		gameId:     gameId.(string),
		file:       file,
		writer:     writer,
		LogChannel: make(chan string),
		wg:         wg,
	}
}

func (tl *TambolaLogger) close() {
	err := tl.writer.Flush()
	if err != nil {
		log.Fatalln(err)
	}
	err = tl.file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(tl.gameId, ":", "Killing logger")
	tl.wg.Done()
}

func (tl *TambolaLogger) StartLogging(gameContext context.Context) {
	defer tl.close()
	gameId := tl.gameId
	for {
		select {
		case text := <-tl.LogChannel:
			_, err := tl.writer.Write([]byte(formatLogEntry(text, gameId)))
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(gameId, ":", text)
		case _ = <-gameContext.Done():
			return
		}
	}
}

func formatLogEntry(text string, gameId any) string {
	return fmt.Sprintf(
		"%v \n%s \n----------- \n", time.TimeOnly, text)
}
