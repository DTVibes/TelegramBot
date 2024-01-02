package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI
var chatID int64

var botNames = [3]string{"советчик", "приятель", "эу"}

var answers = []string{
	"Развивайтесь, но помните, что великий дуб вырастает из маленького желудя.",
	"Иногда лучшее решение — это остановиться и взглянуть на происходящее со стороны.",
	"Если вы не уверены, стоит ли что-то делать, задайте себе вопрос: 'Это приносит мне радость?'",
	"Не бойтесь перемен, они могут открыть новые горизонты.",
	"Прежде чем принимать решение, подумайте, стоит ли вам нарушать свой внутренний покой.",
	"Если что-то приносит вам силу и вдохновение, вы уже на правильном пути.",
	"Не стоит тратить энергию на то, что не приносит вам счастья или роста.",
	"Лучший способ предсказать будущее — создать его самому, но оцените, стоит ли это вашего времени и усилий.",
	"Помните, что важнее — не количество сделанных вещей, а качество вашей жизни.",
	"Иногда отказ — это ответ, который приводит к новым возможностям.",
	"Не давайте влиять на себя чужим мнениям, если это не соответствует вашим ценностям и целям.",
	"Решения, принятые в состоянии гнева или страха, редко приносят долгосрочное удовлетворение.",
	"Всегда стоит стремиться к лучшей версии себя, но не забывайте оценивать свой текущий путь.",
	"Сравнение себя с другими — это быстрый способ потерять свою уникальность и радость.",
	"Прежде чем делать что-то ради других, убедитесь, что это соответствует вашим собственным ценностям.",
	"Не бойтесь отклоняться от общественных ожиданий, если это приводит к вашему собственному развитию.",
	"Стремление к совершенству — это чудесно, но иногда просто достаточно быть достаточно хорошим.",
	"Если что-то не приносит вам удовлетворения, спросите себя, стоит ли это вашего времени и энергии.",
	"Принимайте решения с любовью и добротой, ведь они формируют ваш будущий путь.",
	"Слушайте свое внутреннее чувство — интуиция часто знает ответы на наши вопросы.",
	"Стремитесь к гармонии внутри себя, и вы найдете ответы на многие свои вопросы.",
	"Ваше время ограничено, поэтому стоит делать только то, что приносит вам настоящее счастье.",
	"Не замыкайтесь на прошлом или будущем — настоящий момент часто содержит ответы, которые мы ищем.",
	"Забудьте о страхах и сомнениях, они только мешают вашему личному и профессиональному росту.",
	"При принятии решений доверяйте своей интуиции и опыту — они часто указывают на верный путь.",
	"Не стремитесь угодить всем — важно быть верным самому себе.",
	"Счастье часто находится в простых радостях, поэтому не упускайте их из виду в поисках сложных ответов.",
	"Прежде чем что-то делать, спросите себя, приносит ли это вам радость и удовлетворение.",
	"Стремитесь к балансу в жизни — это ключ к долгосрочному благополучию.",
	"Помните, что ваш путь уникален, и именно он формирует вашу уникальную историю.",
	"Да",
	"Нет",
}

func connectWithTelegram() {
	TOKEN := os.Getenv("TELEGRAM_BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Переменная окружения TELEGRAM_BOT_TOKEN не установлена")
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatalf("Cannot connect to Telegram: %v", err)
	}
}

func sendMessage(msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	_, err := bot.Send(msgConfig)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func isMessageForBot(update *tgbotapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	msgLowerCase := strings.ToLower(update.Message.Text)
	for _, name := range botNames {
		if strings.Contains(msgLowerCase, name) {
			return true
		}
	}
	return false
}

func getBotAnswer() string {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(answers))))
	return answers[index.Int64()]
}

func sendAnswer(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(chatID, getBotAnswer())
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending answer: %v", err)
	}
}

func main() {
	connectWithTelegram()

	updateConfig := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			chatID = update.Message.Chat.ID
			sendMessage("Привет! Я твой телеграм-бот советчик. Я здесь, чтобы помочь тебе с твоим вопросом")
		}

		if isMessageForBot(&update) {
			sendAnswer(&update)
		}
	}
}
