package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/tucnak/telebot"
	"log"
	"time"
)

var userPhotos = make(map[int64][]string)

func main() {
	pref := telebot.Settings{
		Token: "TELEGRAM_TOKEN",
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		_, err := bot.Send(m.Chat, "надо: 1. отправить фотографии 2. написать название")

		if err != nil {
			log.Fatal(err)
		}
	})

	bot.Handle(telebot.OnPhoto, func(m *telebot.Message) {
		photos := m.Photo

		userPhotos[m.Chat.ID] = append(userPhotos[m.Chat.ID], photos.FileID)

		filePath := fmt.Sprintf("./images/%s.jpg", photos.FileID)

		err := bot.Download(&photos.File, filePath)
		if err != nil {
			log.Fatal(err)
		}
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
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
		_, err = bot.Send(m.Chat, telebot.Document{
			File: telebot.File{
				FilePath: outputPath,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		delete(userPhotos, m.Chat.ID)
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
