package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"pandaexpress/tgbotapi"
)

func sendWalletMessage(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat) {
	sessions.UpdateLastCommand(context.Background(),chat.ID, "/wallet")
	wallet, err := db.Adapter.GetWallet(chat.ID)
	if err != nil {
		slog.Error("Error getting wallet", "error", err)
		a, _ := payments.GetAddress(chat.ID)
		wallet = &models.Wallet{Balance: 0, Address: a, UserID: chat.ID}
	}
	withdrawButton := tgbotapi.NewInlineKeyboardButtonData("Withdraw Balance", "/withdraw_balance")
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(withdrawButton),
	)
	msg := tgbotapi.NewMessage(chat.ID,
		fmt.Sprintf("💲 Top up your balance 💲\n\n ❗️ In order to top up your balance, you need to transfer USDT to a wallet below. \nThe transfer is realized automatically.\n\n USDT Balance: %.3f \n\n Your TRC-20 wallet address is:\n `%v`\n\n\n(To copy, click on the wallet👆)", float64(wallet.Balance)/1000000, wallet.Address))
	msg.ReplyMarkup = inlineKeyboard

	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err = bot.Send(msg)
	if err != nil {
		slog.Info("Error sending message:", "error", err)
	}
}

func sendDepositMessage(bot *tgbotapi.BotAPI, chatID int64) {
	wallet, err := db.Adapter.GetWallet(chatID)
	if err != nil {
		a, _ := payments.GetAddress(chatID)
		wallet = &models.Wallet{Balance: zr, Address: a, UserID: chatID}
	}

	withdrawButton := tgbotapi.NewInlineKeyboardButtonData("Withdraw Balance", "/withdraw_balance")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(withdrawButton),
	)

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("💹 Congratutations❕💹 \n 💲 You've got money 💲 \n\n You have topped up your balance✅ 💲\n\n Your Deposit have been received.\n\n USDT Balance: %.3f \n\n Your TRC-20 wallet address is:\n `%v`\n\n\n(To copy, click on the wallet👆)", float64(wallet.Balance)/1000000, wallet.Address))
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err = bot.Send(msg)

	if err != nil {
		slog.Info("Error sending message:", "error", err)
	}
}