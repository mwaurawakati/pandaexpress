package telegram_bot

import (
	"context"
	"log/slog"
	"os"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"strconv"
	"strings"

	"pandaexpress/tgbotapi"
)

var (
	sessions *models.Sessions
	zr       int64 = 0
	//carts              *models.Carts
)

func RunBot(c chan *[]payments.L) {
	sessions = models.NewSessions(db.Adapter.GetCollection("sessions"))
	bot, err := tgbotapi.NewBotAPI("7114413650:AAEsC0PDIJ9lvPXg5-9k8DHXXVvutbODrOM")
	if err != nil {
		sendErrorMessageToOwner(bot, err.Error())
		slog.Error("error creating bot", "error", err)
		slog.Warn("Exiting")
		os.Exit(1)
	}
	go func(c chan *[]payments.L) {
		for {
			data := <-c
			adds, err := db.Adapter.GetAllAddresses()
			if err != nil {
				sendErrorMessageToOwner(bot, err.Error())
			}
			for _, t := range *data {
				to, _ := payments.EthereumToTronAddress(t.Result.To)
				from, _ := payments.EthereumToTronAddress(t.Result.From)
				if ContainsString(adds, to) {
					wallet, err := db.Adapter.GetWalletByAddress(to)
					sendErrorMessageToOwner(bot, err.Error())
					transactionAmount, _ := strconv.ParseInt(t.Result.Value, 10, 64)
					db.Adapter.UpdateWallet(wallet.UserID, &transactionAmount)

					transaction := models.Transaction{
						TransactionID: t.TransactionID,
						Amount:        int64(transactionAmount),
						To:            to,
						Type:          models.Deposit,
						From:          from,
					}
					// Normally, if the transaction is not saved (maybe a duplicate)
					if err := db.Adapter.CreateTransaction(&transaction); err == nil {
						sendDepositMessage(bot, wallet.UserID)
					} else {
						sendErrorMessageToOwner(bot, err.Error())
					}
				}
			}
		}
	}(c)
	//bot.Debug = true

	slog.Info("Authorized on account", "username", bot.Self.UserName)

	// Set up command handlers
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// handle in a routine to free the loop to handle other messages
		go func(update tgbotapi.Update) {
			//slog.Info("update", "update", update)
			if update.CallbackQuery != nil {
				slog.Info(update.CallbackData())
				if strings.HasPrefix(update.CallbackQuery.Data, "select_language_") {
					processSelectLanguageCommand(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "/add_to_cart_") {
					handleAddToCart(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "/remove_from_cart_") {
					handleRemoveFromCart(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "/clear_from_cart_") {
					handleClearFromCart(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "♻️ Checkout ✅" {
					handleCheckout(bot, update.CallbackQuery.Message.Chat.ID)
				} else if update.CallbackQuery.Data == "🛒 Cart" {
					handleCart(bot, update.CallbackQuery.Message.Chat.ID)
				} else if update.CallbackQuery.Data == "Confirm ✅" {
					handleOrder(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "Calcel Order ❌" {
					handleClearCart(bot, update.CallbackQuery.Message.Chat.ID)
				} else if update.CallbackQuery.Data == "Contact Support 📧📝" {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Contact Support feature is coming soon"))
				} else if update.CallbackQuery.Data == "set_language" {
					sendLanguagesMessage(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "set_sa" {
					sendContinentMessage(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "set_osa" {
					sendContinentMessageOrder(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "settings" || update.CallbackQuery.Data == "/settings" || update.CallbackQuery.Data == "⚙️ Settings" {
					sendSettingsMessageEdit(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "confirm_settings_sa" {
					sendSettingsMessageEditSASuccess(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_continent_order_") {
					handleSelectContinentOrder(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_country_order_") {
					handleSelectCountryOrder(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_city_order_") {
					handleSelectCityOrder(bot, update.CallbackQuery)
				} else if update.CallbackQuery.Data == "confirm_order_sa" {
					handleOrder(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_continent_") {
					handleSelectContinent(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_country_") {
					handleSelectCountry(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackQuery.Data, "select_city_") {
					handleSelectCity(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackData(), "next_list_") {
					handleNextCityList(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackData(), "prev_list_") {
					handlePrevCityList(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackData(), "another_address") {
					sendSetAnotherAddress(bot, update.CallbackQuery)
				} else if strings.HasPrefix(update.CallbackData(), "no_another_address") {
					sendNoAnotherAddress(bot, update.CallbackQuery)
				}

			}

			if update.Message == nil {
				return
			}
			switch update.Message.Text {
			case "/start":
				sendStartMessage(bot, update.Message)
			case "/help", "💡 Help", "💡 Help 🫂":
				sendHelpMessage(bot, update.Message.Chat.ID)
			case "/search_product":
				sendImageRequestMessage(bot, update.Message.Chat.ID)
			case "/wallet", "🏧 Wallet":
				sendWalletMessage(bot, update.Message.Chat)
			case "/settings", "⚙️ Settings":
				sendSettingsMessage(bot, update.Message.Chat.ID)
			case "/rm":
				RemoveCartButton(bot, update.Message.Chat.ID)
			case "/a":
				imageURLs := []string{
					"//img.alicdn.com/imgextra/i3/2096578813/O1CN01nifAtA2EyPOssJsgc_!!2-item_pic.png",
					"//img.alicdn.com/imgextra/i3/2096578813/O1CN01ZiNWrU2EyPK80HEPW_!!2096578813.jpg",
					"//img.alicdn.com/imgextra/i2/2096578813/O1CN01pTLeLF2EyPKuw0hYw_!!2096578813.jpg",
					"//img.alicdn.com/imgextra/i3/2096578813/O1CN0180wGEO2EyPJjzn3US_!!2096578813.jpg",
					"//img.alicdn.com/imgextra/i4/2096578813/O1CN013rNWfo2EyPKYQ32j1_!!2096578813.jpg",
				}
				sendPhotosWithCaptionsAndButton(bot, update.Message, imageURLs)
			case "🛒 Cart":
				handleCart(bot, update.Message.Chat.ID)
			case "⭕ Clear Cart ❌":
				handleClearCart(bot, update.Message.Chat.ID)
			case "🌐 Referrals 🌐", "referral", "referrals", "/referrals", "/referral", "🌐 Refferals 🌐":
				HandleRefferal(bot, *update.Message)
			case "♻️ Checkout ✅":
				handleCheckout(bot, update.Message.Chat.ID)
			default:
				// get last command
				lastCommand, _ := sessions.GetLastCommand(context.Background(), update.Message.From.ID)
				if lastCommand == "email" {
					sendSetEmail(bot, *update.Message)
				} else if lastCommand == "phone" {
					sendSetPhone(bot, *update.Message)
				} else if lastCommand == "address" {
					sendSetAddress(bot, *update.Message)
				} else if lastCommand == "full-name" {

				} else {
					if update.Message.Photo != nil {

						if user, err := db.Adapter.GetUser(update.Message.Chat.ID); err != nil {
							slog.Error("error getting user", "error", err)
							if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "We were unable to process your order. Please retry or reach out to support")); err != nil {
								slog.Error("error sending message", "error", err)
							}
						} else if user.Country == "" {
							sendSetFullName(bot, &update)
						} else if user.Country == "Turkmenistan" {
							go searchTaobao(bot, &update)
						} else {
							go search1688(bot, &update)
						}

					} else {
						switch strings.Split(update.Message.Text, " ")[0] {
						case "/start":
							sendStartMessageRefferal(bot, update.Message)
						default:
							sendUnknownCommandMessage(bot, update.Message.Chat.ID)
						}

					}
				}
			}
		}(update)
	}
}
