package main

import (
	"errors"
	"log"
	"time"
	"math/rand"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Answer struct {
	choice bool
}

var (
	ErrRandNumber   = errors.New("[randAnswer]: failed to generate a response")
	ErrConstructBot = errors.New("[construct]: failed to construct a BotAPI")
)

func main() {
	bot, err := newBot()
	if errors.Is(err, ErrConstructBot) {
		log.Panic(ErrConstructBot)
	}

	ch := newUpdateConfig(bot)
	for update := range ch {
		log.Println()
		if update.Message != nil {
			checkMessage(update, bot)
		}
	}
}

// newBot - инийциализирует указатель на экземпляр объекта типа *tgbot.BotAPI
func newBot() (*tgbot.BotAPI, error) {
	bot, err := tgbot.NewBotAPI("7101656518:AAHWr6OiasFE1wo8gsyxf0uKbUuh9TKd1ss")
	if err != nil {
		return nil, ErrConstructBot
	}
	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)
	return bot, nil
}

// newUpdateConfig - инициализирует экземпляр tgbot.UpdateConfig, присваивает ему таймаут,
// который поддерживает постоянный опрос бота на наличие новых сообщений от клиента и возвращает
// канал типа tgbot.UpdatesChannel.
//
// Основной поток блокируется на период чтения из канала, в ожидании обновлений от бота
func newUpdateConfig(bot *tgbot.BotAPI) tgbot.UpdatesChannel {
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	return bot.GetUpdatesChan(u)
}

// checkMessage - проверяет сообщение от клиента.
// Не пропускаются сообщения меньше 5 символов (минимальное ограничение, чтобы было...)
func checkMessage(update tgbot.Update, bot *tgbot.BotAPI) {
	if update.Message.Text == "/start" {
		bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Магический шар знает все на свете. Задай любой вопрос и получи ответ: 'Да', или 'Нет'. Поэтому формулируй вопрос правильно."))
		return
	}
	if len([]rune(update.Message.Text)) < 5 {
		bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Напиши развернутый вопрос. Так результат будет точнее."))
		return
	}

	answer := randAnswer()
	if answer == "" {
		log.Panic(ErrRandNumber)
	}
	msg := tgbot.NewMessage(update.Message.Chat.ID, answer)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}

// randAnswer - случайныйм образом генерирует ответ 'Да', или 'Нет'
func randAnswer() string {
	randNumGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	switch randNumGen.Intn(2) {
	case 0:
		return "Нет"
	case 1:
		return "Да"
	}
	return ""
}
