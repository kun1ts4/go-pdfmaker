package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	"gopkg.in/telebot.v4"
	"log"
	"time"
)

var userPhotos = make(map[int64][]string)

func main() {
	pref := telebot.Settings{
		Token: "7662946517:AAEnRCVDN7t6UK9VxGxk_SSDNswvP2vfExw",
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("надо: 1. отправить фотографии 2. написать название")
	})

	bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
		m := c.Message()
		photos := m.Photo

		userPhotos[m.Chat.ID] = append(userPhotos[m.Chat.ID], photos.FileID)

		filePath := fmt.Sprintf("./images/%s.jpg", photos.FileID)

		err := bot.Download(&photos.File, filePath)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		m := c.Message()
		photos := userPhotos[m.Chat.ID]

		photoPaths := make([]string, len(photos))
		for i, photo := range photos {
			photoPaths[i] = fmt.Sprintf("./images/%s.jpg", photo)
		}
		userPhotos[m.Chat.ID] = nil

		outputPath := fmt.Sprintf("./output/%s.pdf", m.Text)

		err := createPDF(photoPaths, outputPath)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(outputPath)

		delete(userPhotos, m.Chat.ID)

		return c.Send(&telebot.Document{
			File:     telebot.FromDisk(outputPath),
			FileName: fmt.Sprintf("%s.pdf", m.Text),
		})
	})

	bot.Start()
}

func createPDF(photoPaths []string, outputPath string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, photoPath := range photoPaths {
		pdf.AddPage()

		err := pdf.Image(photoPath, 10, 10, nil)
		if err != nil {
			return err
		}
	}

	return pdf.WritePdf(outputPath)
}
