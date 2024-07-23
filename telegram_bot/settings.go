package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/tgbotapi"
	"pandaexpress/translate"
	"strconv"
	"strings"
)

func sendSettingsMessage(bot *tgbotapi.BotAPI, chatID int64) {
	button := tgbotapi.NewInlineKeyboardButtonData("Language", "set_language")
	button1 := tgbotapi.NewInlineKeyboardButtonData("Shipping Address", "set_sa")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button, button1))

	msg := tgbotapi.NewMessage(chatID, "Please select a settings button:")
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSetEmail(bot *tgbotapi.BotAPI, message tgbotapi.Message) {
	user, err := db.Adapter.GetUser(message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue setting email. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	if !IsValidEmail(message.Text) {
		msg := tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("The email: %s you sent is invalid. Please set a valid email e.g john.doe@example.com", message.Text))

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	user.Email = message.Text
	// Extract the selected language from the command
	if err := db.Adapter.UpdateUser(message.From.ID, *user); err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue updating email. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	sessions.UpdateLastCommand(context.Background(), message.From.ID, "phone")
	msg := tgbotapi.NewMessage(message.From.ID, "Email set successfully. Please enter your phone number below.\n It should be in the form of +254799172422")
	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSetPhone(bot *tgbotapi.BotAPI, message tgbotapi.Message) {
	user, err := db.Adapter.GetUser(message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue setting email. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	if !IsValidPhone(message.Text) {
		msg := tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("The pnone: %s you sent is invalid. Please set a valid phone e.g +2341234567890", message.Text))

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	user.Phone = message.Text
	// Extract the selected language from the command
	if err := db.Adapter.UpdateUser(message.From.ID, *user); err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue updating email. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	sessions.UpdateLastCommand(context.TODO(), message.From.ID, "address")
	msg := tgbotapi.NewMessage(message.From.ID, "Phone set successfully. Please enter your address below below.")
	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSetAddress(bot *tgbotapi.BotAPI, message tgbotapi.Message) {
	user, err := db.Adapter.GetUser(message.From.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue setting email. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	user.Addresses = append(user.Addresses, message.Text)
	// Extract the selected language from the command
	if err := db.Adapter.UpdateUser(message.From.ID, *user); err != nil {
		msg := tgbotapi.NewMessage(message.From.ID, "We experienced an issue updating address. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	sessions.UpdateLastCommand(context.TODO(), message.From.ID, "address")
	msg := tgbotapi.NewMessage(message.From.ID, "Address set successfully. Do you want to enter a second address?")
	//msg := tgbotapi.NewMessage(chatID, "Do you want to enter a second address?")
	var buttons = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "another_address"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no_another_address"),
		),
	)
	msg.ReplyMarkup = buttons
	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSetAnotherAddress(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please enter another address below")
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSetFullName(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	sessions.UpdateLastCommand(context.Background(), update.Message.Chat.ID, "full-name")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "To be able to search items, kindly set shipping address.\nPlease enter your full name below")
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendNoAnotherAddress(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	m := fmt.Sprintf("Your shipping address is:\nContinent:%s\nCountry:%s\nCity:%s\nEmail:%s\nPhone:%s\n", user.Continent, user.Country, user.City, user.City, user.Phone)
	for i, address := range user.Addresses {
		m = fmt.Sprintf("%sAddress %d: %s", m, i, address)
	}
	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, m)
	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}
func sendSettingsMessageEdit(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	button := tgbotapi.NewInlineKeyboardButtonData("Language", "set_language")
	button1 := tgbotapi.NewInlineKeyboardButtonData("Shipping Address", "set_sa")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button, button1))
	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select a settings button:", keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendSettingsMessageEditSASuccess(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	sessions.UpdateLastCommand(context.TODO(), callbackQuery.Message.Chat.ID, "email")
	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Country, Continent and city updated. \n Please enter your email below❕")
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendLanguagesMessage(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Create buttons with three buttons per row
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, lang := range models.Languages {
		button := tgbotapi.NewInlineKeyboardButtonData(lang.Name, fmt.Sprintf("select_language_%s", lang.Name))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/settings")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your preferred language:", keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendContinentMessage(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Create buttons with three buttons per row
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, cont := range models.Continents {
		button := tgbotapi.NewInlineKeyboardButtonData(cont, fmt.Sprintf("select_continent_%s", cont))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/settings")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your continent:", keyboard)
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func sendContinentMessageOrder(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {

	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "You have not set your shipping address or your sipping address is incomplete. \n Use /settings command to set your sipping details")
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleSelectCountry(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Extract the selected language from the command
	selectedCountry := strings.TrimPrefix(callbackQuery.Data, "select_country_")
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, models.User{ShippingDetails: models.ShippingDetails{Country: selectedCountry, Continent: user.ShippingDetails.Continent}}); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	var cities = models.Cities[selectedCountry]
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for i, city := range cities {
		button := tgbotapi.NewInlineKeyboardButtonData(city, fmt.Sprintf("select_city_%s", city))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
		if i == 29 {
			break
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	if len(cities) > 30 {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Next List", "next_list_30_"+selectedCountry)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	} else {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your city", keyboard)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleNextCityList(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Extract the selected language from the command
	callbackQueryData := strings.Split(callbackQuery.Data, "_")
	selectedCountry := callbackQueryData[3]
	startIndex, _ := strconv.Atoi(callbackQueryData[2])
	var cities = models.Cities[selectedCountry]
	if len(cities) > startIndex+30 {
		cities = cities[startIndex : startIndex+30]
	} else {
		cities = cities[startIndex:]
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, city := range cities {
		button := tgbotapi.NewInlineKeyboardButtonData(city, fmt.Sprintf("select_city_%s", city))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	cities = models.Cities[selectedCountry]
	if len(cities) > startIndex+30 {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Prev List", fmt.Sprintf("prev_list_%d_%s", startIndex-30, selectedCountry)), tgbotapi.NewInlineKeyboardButtonData("Next List", fmt.Sprintf("next_list_%d_%s", startIndex+30, selectedCountry))),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	} else {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Prev List", fmt.Sprintf("prev_list_%d_%s", startIndex-30, selectedCountry))),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your city", keyboard)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handlePrevCityList(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Extract the selected language from the command
	callbackQueryData := strings.Split(callbackQuery.Data, "_")
	selectedCountry := callbackQueryData[3]
	startIndex, _ := strconv.Atoi(callbackQueryData[2])
	var cities = models.Cities[selectedCountry]

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, city := range cities[startIndex : startIndex+30] {
		button := tgbotapi.NewInlineKeyboardButtonData(city, fmt.Sprintf("select_city_%s", city))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	if startIndex > 0 {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Prev List", fmt.Sprintf("prev_list_%d_%s", startIndex-30, selectedCountry)), tgbotapi.NewInlineKeyboardButtonData("Next List", fmt.Sprintf("next_list_%d_%s", startIndex+30, selectedCountry))),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	} else {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Next List", fmt.Sprintf("next_list_%d_%s", startIndex+30, selectedCountry))),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_"+user.Continent)))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your city", keyboard)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleSelectCountryOrder(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Extract the selected language from the command
	selectedCountry := strings.TrimPrefix(callbackQuery.Data, "select_country_order_")
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, models.User{ShippingDetails: models.ShippingDetails{Country: selectedCountry, Continent: user.ShippingDetails.Continent}}); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	var cities = models.Cities[selectedCountry]
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, city := range cities {
		button := tgbotapi.NewInlineKeyboardButtonData(city, fmt.Sprintf("select_city_order_%s", city))
		row = append(row, button)
		if len(row) == 6 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "select_continent_order_"+user.Continent)))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your city", keyboard)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}
func handleSelectCity(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	// Extract the selected language from the command
	selectedCity := strings.TrimPrefix(callbackQuery.Data, "select_city_")
	user.ShippingDetails.City = selectedCity
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, *user); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Confirm", "confirm_settings_sa")), tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", callbackQuery.Data)))

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		fmt.Sprintf("Please confirm the following as your selection: \n\n Continent: %s\n Country: %s,\n City: %s",
			user.ShippingDetails.Continent, user.ShippingDetails.Country, user.ShippingDetails.City),
		keyboard,
	)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleSelectCityOrder(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Extract the selected language from the command
	selectedCity := strings.TrimPrefix(callbackQuery.Data, "select_city_order_")
	user.ShippingDetails.City = selectedCity
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, *user); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Confirm", "confirm_order_sa")), tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", callbackQuery.Data)))

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		fmt.Sprintf("Please confirm the following as your selection: \n\n Continent: %s\n Country: %s,\n City: %s",
			user.ShippingDetails.Continent, user.ShippingDetails.Country, user.ShippingDetails.City),
		keyboard,
	)

	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func SendShippingAddressMessage(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	_, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// get the
}

func processSelectLanguageCommand(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	button := tgbotapi.NewInlineKeyboardButtonData("Language", "set_language")
	button1 := tgbotapi.NewInlineKeyboardButtonData("Shipping Address", "set_sa")
	// Extract the selected language from the command
	selectedLanguage := strings.TrimPrefix(callbackQuery.Data, "select_language_")
	language := models.GetUsingName(selectedLanguage)
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, models.User{ID: callbackQuery.Message.Chat.ID, PreferredLanguage: language}); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	// Translate the confirmation message to the selected language
	confirmationMessage := fmt.Sprintf("Your selected language is: %s", language.EnglishName)
	translatedMessage, err := translate.TranslateTextToPreferredLanguage(confirmationMessage, language.Code)
	if err != nil {
		slog.Error("Error translating message:", "error", err)
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button, button1))
	// Send the translated message
	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, translatedMessage, keyboard)
	_, err = bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message:", "error", err)
	}
}

func handleSelectContinentOrder(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Extract the selected language from the command
	selectedContinent := strings.TrimPrefix(callbackQuery.Data, "select_continent_order_")
	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, models.User{ID: callbackQuery.Message.Chat.ID, ShippingDetails: models.ShippingDetails{Continent: selectedContinent}}); err != nil {
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	var countries = models.ContinentCountries[selectedContinent]
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, country := range countries {
		button := tgbotapi.NewInlineKeyboardButtonData(country, fmt.Sprintf("select_country_order_%s", country))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "set_osa")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your country", keyboard)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleSelectContinent(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Extract the selected language from the command
	selectedContinent := strings.TrimPrefix(callbackQuery.Data, "select_continent_")

	if err := db.Adapter.UpdateUser(callbackQuery.Message.Chat.ID, models.User{ShippingDetails: models.ShippingDetails{Continent: selectedContinent}}); err != nil {
		slog.Error("error saving continent", "error", err)
		msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "We experienced an issue. Please contact customer support.")

		_, err := bot.Send(msg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	var countries = models.ContinentCountries[selectedContinent]
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, country := range countries {
		button := tgbotapi.NewInlineKeyboardButtonData(country, fmt.Sprintf("select_country_%s", country))
		row = append(row, button)
		if len(row) == 3 {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "set_sa")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewEditMessageTextAndMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Please select your country", keyboard)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}
