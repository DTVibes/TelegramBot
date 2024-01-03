package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI
var chatID int64
var userGreeted bool

type JokeAPIResponse struct {
	Type     string `json:"type"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}

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

func fetchRandomJoke() string {
	apiUrl := "https://v2.jokeapi.dev/joke/Any" // URL API для получения случайной шутки

	response, err := http.Get(apiUrl)
	if err != nil {
		return "Произошла ошибка при получении шутки."
	}
	defer response.Body.Close()

	var jokeResponse JokeAPIResponse
	if err := json.NewDecoder(response.Body).Decode(&jokeResponse); err != nil {
		return "Произошла ошибка при обработке ответа от API."
	}

	if jokeResponse.Type == "twopart" {
		return jokeResponse.Setup + "\n" + jokeResponse.Delivery
	} else {
		return jokeResponse.Setup
	}
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

func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	button1 := tgbotapi.NewKeyboardButton("Расскажи шутку")

	// Создаем ряды кнопок
	row1 := tgbotapi.NewKeyboardButtonRow(button1)

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewReplyKeyboard(row1)

	// Устанавливаем текст для сообщения
	keyboard.ResizeKeyboard = true
	return keyboard
}

func main() {
	connectWithTelegram()

	updateConfig := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		if update.Message != nil {
			// Проверка и установка chatID, если он еще не установлен
			if chatID == 0 {
				chatID = update.Message.Chat.ID
			}

			if update.Message.Text == "/start" && !userGreeted {
				// Greet the user and ask for their name

				// Устанавливаем клавиатуру после команды /start
				keyboard := getKeyboard()
				msg := tgbotapi.NewMessage(chatID, "Привет! Я твой телеграм-бот советчик. Я здесь, чтобы помочь тебе с твоим вопросом. Ты можешь обращаться ко мне как \"Советчик\" или \"Приятель\". Как тебя зовут? ")
				msg.ReplyMarkup = keyboard
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Error sending message with keyboard: %v", err)
				}

				userGreeted = true
			} else if isMessageForBot(&update) {
				sendAnswer(&update)
			} else if update.Message.Text == "Расскажи шутку" {
				// Загрузка случайной шутки и отправка пользователю
				joke := fetchRandomJoke()
				sendMessage(joke)
			} else if userGreeted {
				// Assuming the next message after the greeting is the user's name
				// You may need to implement a more sophisticated name retrieval mechanism
				userName := update.Message.Text
				sendMessage("Приятно познакомиться, " + userName + "!")
				sendMessage("Какой вопрос ты мне хочешь задать?")
				userGreeted = false
			}
		}
	}
}
