package telegram_bot

import (
	"fmt"
	"log/slog"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/tgbotapi"
)

func generateReferralCode(chatID int64) string {
	return fmt.Sprintf("ref%d", chatID)
}

func SendRefferalMessage(bot *tgbotapi.BotAPI, user models.User) {
	totalEarned, err := db.Adapter.GetTotalRefferalsEarning(user.ReferralCode)
	if err != nil {
		slog.Error("error getting total earnings", "error", err)
	}
	m := fmt.Sprintf("Congratulations! You have a referral.\n\n\nTotal referrals: %d\nTotalReferralEarnings: %.3f", len(user.Referrals), totalEarned)
	message := tgbotapi.NewMessage(user.ID,m)
	if _, err =bot.Send(message); err != nil{
		slog.Error("error sending message", "error", err)
	}
}

// send earning message
func SendRefferalTransaction(bot *tgbotapi.BotAPI, user models.User, transaction models.Transaction) {
	totalEarned, err := db.Adapter.GetTotalRefferalsEarning(user.ReferralCode)
	if err != nil {
		slog.Error("error getting total earnings", "error", err)
	}
	m := fmt.Sprintf("Congratulations! You have earnings from referrals.\nTransaction ID:%s\nAmount Earned: %.3f\n\n\nTotal referrals: %d\nTotalReferralEarnings: %.3f", transaction.TransactionID, float64(transaction.Amount)/1000000, len(user.Referrals), totalEarned)
	message := tgbotapi.NewMessage(user.ID,m)
	if _, err =bot.Send(message); err != nil{
		slog.Error("error sending message", "error", err)
	}
}

func HandleRefferal(bot *tgbotapi.BotAPI, message tgbotapi.Message){
	user, err := db.Adapter.GetUser(message.Chat.ID)
	if err != nil {
		slog.Info("empty cart")
		newCaption := "We were unable to process your order. Please retry or reach out to support"
		addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Retry order ♻️", "Confirm ✅")
		contactSupoortButton := tgbotapi.NewInlineKeyboardButtonData("Contact Support 📧📝", "Contact Support 📧📝")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton),
			tgbotapi.NewInlineKeyboardRow(contactSupoortButton),
		)
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, newCaption, keyboard)
		_, err := bot.Send(editMsg)
		if err != nil {
			slog.Error("Error sending message", "error", err)
		}
		return
	}
	totalEarned, err := db.Adapter.GetTotalRefferalsEarning(user.ReferralCode)
	if err != nil {
		slog.Error("error getting total earnings", "error", err)
	}
	m := fmt.Sprintf("Earn from referrals. Get a percentage of the order total (not inclusive of shipping fee) of all the orders placed by your referrals\n\nYour refferal code is: https://t.me/%s?start=%s.\n\n\nTotal referrals: %d\nTotalReferralEarnings: %.3f",
	bot.Self.UserName,user.ReferralCode, len(user.Referrals), totalEarned)
	messageNew := tgbotapi.NewMessage(user.ID,m)
	if _, err =bot.Send(messageNew); err != nil{
		slog.Error("error sending message", "error", err)
	}
}