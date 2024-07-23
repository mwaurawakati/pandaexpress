package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"pandaexpress/tgbotapi"
	"pandaexpress/translate"
	"strings"
)

func sendStartMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// get user preffered language
	code := message.From.LanguageCode
	user := models.User{
		ID:                message.Chat.ID,
		FirstName:         message.Chat.FirstName,
		LastName:          message.Chat.LastName,
		UserName:          message.Chat.UserName,
		PreferredLanguage: models.GetUsingCode(code),
		ReferralCode:      generateReferralCode(message.From.ID),
	}
	sessions.AddSession(context.TODO(), models.Session{UserID: user.ID, User: user, LastCommand: "full-name", PreferredLanguage: models.GetUsingCode(code), Cart: make(map[int64]models.CartItem)})
	var wallet *models.Wallet
	if err := db.Adapter.UserCreate(&user); err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			// User aleady exist. They had either uninstalled Telegram or deleted chats
			if wallet, err = db.Adapter.GetWallet(message.Chat.ID); err != nil {
				address, _ := payments.GetAddress(message.Chat.ID)
				wallet = &models.Wallet{
					Address: address,
					Balance: zr,
					UserID:  message.Chat.ID,
				}
			}
		} else {
			address, _ := payments.GetAddress(message.Chat.ID)
			wallet = &models.Wallet{
				Address: address,
				Balance: zr,
				UserID:  message.Chat.ID,
			}
		}
	} else {
		// Create wallet
		address, _ := payments.GetAddress(message.Chat.ID)
		wallet = &models.Wallet{
			Address: address,
			Balance: zr,
			UserID:  message.Chat.ID,
		}
		if err := db.Adapter.CreateWallet(*wallet); err != nil {
			if strings.Contains(err.Error(), "duplicate key error") {
				if wallet, err = db.Adapter.GetWallet(message.Chat.ID); err != nil {
					address, _ := payments.GetAddress(message.Chat.ID)
					wallet = &models.Wallet{
						Address: address,
						Balance: zr,
						UserID:  message.Chat.ID,
					}
				}
			}
		}
	}
	var keyboard tgbotapi.InlineKeyboardMarkup
	// Define buttons
	button := tgbotapi.NewInlineKeyboardButtonURL("Shop", "http://t.me/panda_express_test_bot/shop/about")
	if models.CheckSA(user) != models.Done {
		button = tgbotapi.NewInlineKeyboardButtonData("Set shipping Address", "/set-shipping-address")
		keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	}

	//payments.Watch.WatchAddress(wallet.Address)
	transactionsButton := tgbotapi.NewKeyboardButton("♐ Transactions ♐")
	ordersButton := tgbotapi.NewKeyboardButton("🔰 Orders 🔰")
	refferals := tgbotapi.NewKeyboardButton("🌐 Refferals 🌐")

	// Create reply keyboard markup
	_ = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(transactionsButton, ordersButton),
		tgbotapi.NewKeyboardButtonRow(refferals),
	)
	messageText := fmt.Sprintf("Hello %s %s (@%s) \n", message.Chat.FirstName, message.Chat.LastName, message.Chat.UserName)
	messageText = fmt.Sprintf("%sWelcome to Panda Express Bot!\nPanda Express helps you to have a seamless easy shopping. ", messageText)
	messageText = fmt.Sprintf("%sTo be able to have great experience, make sure to fund your wallet. Use /wallet command to access wallet.\n", messageText)
	messageText = fmt.Sprintf("%sYour TRC20 wallet address is: `%v` \n\nUSDT Balance: %.3f", messageText, wallet.Address, float64(wallet.Balance)/1000000)
	messageText = fmt.Sprintf("%s\n\n\nTo be able to use Panda Express, you MUST set your shipping Adress. ", messageText)
	messageText = fmt.Sprintf("%s\n\n\nEarn by refering your friend. \nUser /referral command for more infomation about referrals.\nYour referal link is:\nhttps://t.me/%s?start=%s",
		messageText, bot.Self.UserName, user.ReferralCode)
	messageText = fmt.Sprintf("%s\n\n. Kindly set your your shipping details to be able to use this bot. Waht is your full name?", messageText)
	var translatedMessage string
	if user.PreferredLanguage.Code != "en" {
		var err error
		translatedMessage, err = translate.TranslateTextToPreferredLanguage(messageText, user.PreferredLanguage.Code)
		if err != nil {
			slog.Error("Error translating message:", "error", err)
			translatedMessage = messageText
		}
	} else {
		translatedMessage = messageText
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, translatedMessage)
	_ = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		slog.Info("Error sending message:", "error", err)
	}
}

// Function to send all photos in a single message with captions and an inline button
func sendPhotosWithCaptionsAndButton(bot *tgbotapi.BotAPI, message *tgbotapi.Message, imageURLs []string) {
	var mediaGroup []interface{}
	for i, url := range imageURLs {
		fullURL := "http:" + url
		caption := fmt.Sprintf("Image %d", i+1)
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(fullURL))
		photo.Caption = caption
		//photo.
		mediaGroup = append(mediaGroup, photo)

	}

	// Create the media group message
	mediaMsg := tgbotapi.NewMediaGroup(message.Chat.ID, mediaGroup)
	// Define the inline button
	button := tgbotapi.NewInlineKeyboardButtonURL("View All Images", "https://t.me/your_bot_url")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))

	// Send the media group message
	msg := tgbotapi.NewMessage(message.Chat.ID, "Here are the images:")
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		slog.Error("Error sending message", "error", err)
	}
	if _, err := bot.Send(mediaMsg); err != nil {
		slog.Error("Error sending media group", "error", err)
	}
}

func sendStartMessageRefferal(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// get user preffered language
	code := message.From.LanguageCode
	user := models.User{
		ID:                message.Chat.ID,
		FirstName:         message.Chat.FirstName,
		LastName:          message.Chat.LastName,
		UserName:          message.Chat.UserName,
		PreferredLanguage: models.GetUsingCode(code),
		ReferralCode:      generateReferralCode(message.From.ID),
		Referee:           strings.Split(message.Text, " ")[1],
	}
	sessions.AddSession(context.TODO(),models.Session{UserID: user.ID, User: user, LastCommand: "/start", PreferredLanguage: models.GetUsingCode(code), Cart: make(map[int64]models.CartItem)})
	var wallet *models.Wallet
	if err := db.Adapter.UserCreate(&user); err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			// User aleady exist. They had either uninstalled Telegram or deleted chats
			if wallet, err = db.Adapter.GetWallet(message.Chat.ID); err != nil {
				address, _ := payments.GetAddress(message.Chat.ID)
				wallet = &models.Wallet{
					Address: address,
					Balance: zr,
					UserID:  message.Chat.ID,
				}
			}
		} else {
			address, _ := payments.GetAddress(message.Chat.ID)
			wallet = &models.Wallet{
				Address: address,
				Balance: zr,
				UserID:  message.Chat.ID,
			}
		}
	} else {
		// Create wallet
		address, _ := payments.GetAddress(message.Chat.ID)
		wallet = &models.Wallet{
			Address: address,
			Balance: zr,
			UserID:  message.Chat.ID,
		}
		if err := db.Adapter.CreateWallet(*wallet); err != nil {
			if strings.Contains(err.Error(), "duplicate key error") {
				if wallet, err = db.Adapter.GetWallet(message.Chat.ID); err != nil {
					address, _ := payments.GetAddress(message.Chat.ID)
					wallet = &models.Wallet{
						Address: address,
						Balance: zr,
						UserID:  message.Chat.ID,
					}
				}
			}
		}
	}
	var keyboard tgbotapi.InlineKeyboardMarkup
	// Define buttons
	button := tgbotapi.NewInlineKeyboardButtonURL("Shop", "http://t.me/panda_express_test_bot/shop/about")
	if models.CheckSA(user) != models.Done {
		button = tgbotapi.NewInlineKeyboardButtonData("Set shipping Address", "/set-shipping-address")
		keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	}

	//payments.Watch.WatchAddress(wallet.Address)
	transactionsButton := tgbotapi.NewKeyboardButton("♐ Transactions ♐")
	ordersButton := tgbotapi.NewKeyboardButton("🔰 Orders 🔰")
	refferals := tgbotapi.NewKeyboardButton("🌐 Refferals 🌐")

	// Create reply keyboard markup
	_ = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(transactionsButton, ordersButton),
		tgbotapi.NewKeyboardButtonRow(refferals),
	)
	ref, err := db.Adapter.GetUserByRefID(user.Referee)
	if err != nil {
		slog.Error("error querying db", "error", err)
		sendErrorMessageToOwner(bot, err.Error())
	}
	err = db.Adapter.UpdateReferral(ref.ReferralCode, user.ReferralCode)
	if err != nil {
		slog.Error("error querying db", "error", err)
		sendErrorMessageToOwner(bot, err.Error())
	}
	go SendRefferalMessage(bot, *ref)
	messageText := fmt.Sprintf("Hello %s %s (@%s) \n", message.Chat.FirstName, message.Chat.LastName, message.Chat.UserName)
	messageText = fmt.Sprintf("%sWelcome to Panda Express Bot!\nPanda Express helps you to have a seamless easy shopping. ", messageText)
	messageText = fmt.Sprintf("%sTo be able to have great experience, make sure to fund your wallet. Use /wallet command to access wallet.\n", messageText)
	messageText = fmt.Sprintf("%sYour TRC20 wallet address is: `%v` \n\nUSDT Balance: %.3f", messageText, wallet.Address, float64(wallet.Balance)/1000000)
	messageText = fmt.Sprintf("%s\n\n\nTo be able to use Panda Express, you MUST set your shipping Adress. ", messageText)
	messageText = fmt.Sprintf("%s\n\n\nEarn by refering your friend. \nUser /referral command for more infomation about referrals.\n You have been referred by: %s %s(%s)\nYour referal link is:\nhttps://t.me/%s?start=%s",
		messageText,
		ref.FirstName,
		ref.LastName,
		ref.UserName,
		bot.Self.UserName, user.ReferralCode)
	var translatedMessage string
	if user.PreferredLanguage.Code != "en" {
		var err error
		translatedMessage, err = translate.TranslateTextToPreferredLanguage(messageText, user.PreferredLanguage.Code)
		if err != nil {
			slog.Error("Error translating message:", "error", err)
			translatedMessage = messageText
		}
	} else {
		translatedMessage = messageText
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, translatedMessage)
	_ = keyboard
	_, err = bot.Send(msg)
	if err != nil {
		slog.Info("Error sending message:", "error", err)
	}
}
