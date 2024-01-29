package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"encoding/json"
	"path/filepath"

	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ImageJsonData struct {
	File_id        string
	File_unique_id string
	File_size      int
	File_path      string
}

type ImageJsonRoot struct {
	Ok     bool
	Result ImageJsonData
}

var token string

var dlPath = os.Getenv("HOME") + "/Slideshow"

func main() {
	if len(token) == 0 {
		fmt.Println("Telegram bot token not given, see makefile")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func getJson(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}

func downloadImage(url string) bool {
	jr := ImageJsonRoot{}
	getJson(url, &jr)

	_, fname := filepath.Split(jr.Result.File_path)

	dlUrl := "https://api.telegram.org/file/bot" + token + "/" + jr.Result.File_path
	outputPath := dlPath + "/" + fname

	file, err := os.Create(outputPath)
	if err != nil {
		return false
	}
	defer file.Close()

	response, err := http.Get(dlUrl)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)

	return err == nil
}

func deleteImageFiles() error {
	dir, err := os.Open(dlPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	fnames, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	if len(fnames) == 0 {
		return err
	}

	for _, fname := range fnames {
		err = os.RemoveAll(filepath.Join(dir.Name(), fname))
		if err != nil {
			return err
		}
	}

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var responseMessage string

	if update.Message.Photo != nil {
		responseMessage = "Photo received, but got thrown in the trash :D Please send images as files"
	} else if update.Message.Document != nil {
		url := "https://api.telegram.org/bot" + token + "/getFile?file_id=" + update.Message.Document.FileID

		if downloadImage(url) {
			responseMessage = "Successfully loaded image"
		} else {
			responseMessage = "Image load failed"
		}
	} else {
		switch update.Message.Text {
		case "/hello":
			responseMessage = "Hello there"
		case "/delete":
			deleteImageFiles()
			responseMessage = "Old files deleted"
		default:
			responseMessage = "Command not understood\n\nAvailable commands:\n"
			responseMessage += "/hello, sends 'Hello there'\n"
			responseMessage += "/delete, deletes all sent images\n\n"
			responseMessage += "Images sent as files will be added to the server"
		}
	}

	b.SendMessage(ctx, messageBody(update, responseMessage))
}

func messageBody(update *models.Update, msg string) *bot.SendMessageParams {
	return &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	}
}
