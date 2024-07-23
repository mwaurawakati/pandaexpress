package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var continents = map[string][]string{
	"Africa":   {"Nigeria", "Kenya", "South Africa"},
	"Asia":     {"China", "Japan", "India"},
	"Europe":   {"Germany", "France", "UK"},
	"America":  {"USA", "Canada", "Brazil"},
	"Oceania":  {"Australia", "New Zealand"},
}

type UserAddress struct {
	Continent        string
	Country          string
	City             string
	Address          string
	Address2         string
	Address3         string
	Email            string
	PhoneNumber      string
	AskingForAddress int
}

var userAddresses = make(map[int64]*UserAddress)

func main() {
	bot, err := tgbotapi.NewBotAPI("7114413650:AAEsC0PDIJ9lvPXg5-9k8DHXXVvutbODrOM")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
		} else if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress, exists := userAddresses[chatID]

	if !exists || userAddress.Continent == "" {
		sendContinentPrompt(bot, chatID)
		return
	}

	if userAddress.Country == "" {
		handleCountryInput(bot, message)
		return
	}

	if userAddress.City == "" {
		handleCityInput(bot, message)
		return
	}

	if userAddress.Address == "" {
		handleAddressInput(bot, message)
		return
	}

	if userAddress.Address2 == "" && userAddress.AskingForAddress == 1 {
		handleAddress2Input(bot, message)
		return
	}

	if userAddress.Address3 == "" && userAddress.AskingForAddress == 2 {
		handleAddress3Input(bot, message)
		return
	}

	if userAddress.Email == "" {
		handleEmailInput(bot, message)
		return
	}

	if userAddress.PhoneNumber == "" {
		handlePhoneNumberInput(bot, message)
		return
	}

	confirmAddress(bot, chatID)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	userAddress := getUserAddress(chatID)

	if userAddress.Continent == "" {
		userAddress.Continent = callbackQuery.Data
		userAddresses[chatID] = userAddress
		sendCountryPrompt(bot, chatID, userAddress.Continent)
	} else if userAddress.Country == "" {
		userAddress.Country = callbackQuery.Data
		userAddresses[chatID] = userAddress
		sendCityPrompt(bot, chatID)
	} else if userAddress.AskingForAddress == 0 {
		if callbackQuery.Data == "Yes" {
			userAddress.AskingForAddress = 1
			sendAddress2Prompt(bot, chatID)
		} else {
			sendEmailPrompt(bot, chatID)
		}
		userAddresses[chatID] = userAddress
	} else if userAddress.AskingForAddress == 1 {
		if callbackQuery.Data == "Yes" {
			userAddress.AskingForAddress = 2
			sendAddress3Prompt(bot, chatID)
		} else {
			sendEmailPrompt(bot, chatID)
		}
		userAddresses[chatID] = userAddress
	}
}

func sendContinentPrompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please select your continent:")
	var buttons []tgbotapi.InlineKeyboardButton
	for continent := range continents {
		button := tgbotapi.NewInlineKeyboardButtonData(continent, continent)
		buttons = append(buttons, button)
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	bot.Send(msg)
}

func sendCountryPrompt(bot *tgbotapi.BotAPI, chatID int64, continent string) {
	msg := tgbotapi.NewMessage(chatID, "Please select your country:")
	var buttons []tgbotapi.InlineKeyboardButton
	for _, country := range continents[continent] {
		button := tgbotapi.NewInlineKeyboardButtonData(country, country)
		buttons = append(buttons, button)
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	bot.Send(msg)
}

func sendCityPrompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please enter your city:")
	bot.Send(msg)
}

func sendAddressPrompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please enter your address:")
	bot.Send(msg)
}

func sendAddress2Prompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Do you want to enter a second address?")
	var buttons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "Yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "No"),
		),
	)
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func sendAddress3Prompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Do you want to enter a third address?")
	var buttons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "Yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "No"),
		),
	)
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func sendEmailPrompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please enter your email address:")
	bot.Send(msg)
}

func sendPhoneNumberPrompt(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please enter your phone number:")
	bot.Send(msg)
}

func confirmAddress(bot *tgbotapi.BotAPI, chatID int64) {
	userAddress := userAddresses[chatID]
	msg := tgbotapi.NewMessage(chatID, "Please confirm your address:\n" +
		"Continent: " + userAddress.Continent + "\n" +
		"Country: " + userAddress.Country + "\n" +
		"City: " + userAddress.City + "\n" +
		"Address: " + userAddress.Address + "\n" +
		"Address 2: " + userAddress.Address2 + "\n" +
		"Address 3: " + userAddress.Address3 + "\n" +
		"Email: " + userAddress.Email + "\n" +
		"Phone Number: " + userAddress.PhoneNumber)
	bot.Send(msg)
}

func handleCountryInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.Country = message.Text
	userAddresses[chatID] = userAddress
	sendCityPrompt(bot, chatID)
}

func handleCityInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.City = message.Text
	userAddresses[chatID] = userAddress
	sendAddressPrompt(bot, chatID)
}

func handleAddressInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.Address = message.Text
	userAddresses[chatID] = userAddress
	sendAddress2Prompt(bot, chatID)
}

func handleAddress2Input(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.Address2 = message.Text
	userAddresses[chatID] = userAddress
	sendAddress3Prompt(bot, chatID)
}

func handleAddress3Input(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.Address3 = message.Text
	userAddresses[chatID] = userAddress
	sendEmailPrompt(bot, chatID)
}

func handleEmailInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.Email = message.Text
	userAddresses[chatID] = userAddress
	sendPhoneNumberPrompt(bot, chatID)
}

func handlePhoneNumberInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userAddress := getUserAddress(chatID)
	userAddress.PhoneNumber = message.Text
	userAddresses[chatID] = userAddress
	confirmAddress(bot, chatID)
}

func getUserAddress(chatID int64) *UserAddress {
	userAddress, exists := userAddresses[chatID]
	if !exists {
		userAddress = &UserAddress{}
		userAddresses[chatID] = userAddress
	}
	return userAddress
}
