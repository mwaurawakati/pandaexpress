package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"pandaexpress/taobao"
	"pandaexpress/tgbotapi"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func AddCartKeyboard(bot *tgbotapi.BotAPI, chatID int64) int {
	cartButton := tgbotapi.NewKeyboardButton("🛒 Cart")
	checkoutButton := tgbotapi.NewKeyboardButton("♻️ Checkout ✅")
	clearButton := tgbotapi.NewKeyboardButton("⭕ Clear Cart ❌")
	transactionsButton := tgbotapi.NewKeyboardButton("♐ Transactions ♐")
	ordersButton := tgbotapi.NewKeyboardButton("🔰 Orders 🔰")
	refferals := tgbotapi.NewKeyboardButton("🌐 Refferals 🌐")
	checkButton := tgbotapi.NewInlineKeyboardButtonData("♻️ Checkout ✅", "♻️ Checkout ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("🛒 Cart", "🛒 Cart")
	// Create reply keyboard markup
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(cartButton, checkoutButton),
		tgbotapi.NewKeyboardButtonRow(clearButton),
		tgbotapi.NewKeyboardButtonRow(transactionsButton, ordersButton),
		tgbotapi.NewKeyboardButtonRow(refferals),
	)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(checkButton, cButton),
	)
	// Send message to update the reply keyboard
	replyMsg := tgbotapi.NewMessage(chatID, "You have items in the Cart. Click \n 🛒 Cart to view \n ♻️ Checkout ✅ to place an order")
	replyMsg.ReplyMarkup = struct {
		tgbotapi.InlineKeyboardMarkup
		tgbotapi.ReplyKeyboardMarkup
	}{
		inlineKeyboard,
		replyKeyboard,
	}
	m, err := bot.Send(replyMsg)

	if err != nil {
		slog.Info("Error sending reply message:", "error", err)
		return 0
	}
	return m.MessageID
}

func handleOrder(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	settings, err := db.Adapter.GetSettings()
	if err != nil {
		slog.Error(err.Error())
	}
	// Get the original message
	message := callbackQuery.Message
	// Check if shipping details exist
	user, err := db.Adapter.GetUser(callbackQuery.Message.Chat.ID)
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
	if models.CheckSA(*user) != models.Done {
		sendContinentMessageOrder(bot, callbackQuery)
		return
	}
	order := models.Order{UserID: callbackQuery.Message.Chat.ID, ShippingFee: 100.34}

	items, _ := sessions.GetCart(context.TODO(), message.Chat.ID)
	if len(items) == 0 {
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
	rate, _ := payments.GetCnyToUsdtRate()
	for _, item := range items {
		order.Items = append(order.Items, item)
		for id, sku := range item.Quantity {
			var p string
			for _, base := range item.Details.Result.Item.Sku.Base {
				if base.SkuId == id {
					p = base.PromotionPrice
					break
				}
			}
			pFloat, _ := strconv.ParseFloat(p, 64)

			//p, _ := strconv.ParseFloat(base.PromotionPrice, 64)
			usdtPrice := payments.ConvertCnyToUsdt(pFloat, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			order.TotalPrice += (thePrice * float64(sku))
		}
	}
	var transactions []models.Transaction
	transaction := models.Transaction{
		Type:          models.Purchase,
		To:            "Panda Express",
		From:          user.ID,
		TransactionID: uuid.NewString(),
		Amount:        int64(order.TotalPrice * 1000000),
	}
	order.TransactionIDs = append(order.TransactionIDs, transaction.TransactionID)
	transactions = append(transactions, transaction)
	/*if err := db.Adapter.CreateTransaction(&transaction); err != nil {
		slog.Error("error creating transaction", "error", err)
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
	}*/

	// Handle refferal
	if user.Referee != "" {
		//referee, _ := db.Adapter.GetUserByRefID(user.Referee)
		reftransaction := models.Transaction{
			Type:          models.RefferalEarning,
			To:            user.Referee,
			From:          "Panda Express",
			TransactionID: uuid.NewString(),
			Amount:        int64(order.TotalPrice * 1000000 * settings.ReferralCommission),
		}
		order.TransactionIDs = append(order.TransactionIDs, reftransaction.TransactionID)

		transactions = append(transactions, reftransaction)
		/*if err := db.Adapter.CreateTransaction(&reftransaction); err != nil {
			slog.Error("error creating ref transactions", "error", err)
			go sendErrorMessageToOwner(bot, "error creating refferal transaction:"+err.Error())
		}else{
			go SendRefferalTransaction(bot, *referee, reftransaction)
		}*/
	}

	//order.TransactionID = transaction.TransactionID
	// Create a shipping transaction
	// get shipping fee
	shippingfee, err := getShippingPrice(settings, user.Country, user.City)
	if err != nil {
		slog.Info("error getting shipping fee")
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
	order.ShippingFee = shippingfee
	shipingTransaction := models.Transaction{
		Type:          models.ShippingFee,
		To:            "Pand Express",
		From:          user.ID,
		TransactionID: uuid.NewString(),
		Amount:        int64(shippingfee * 1000000),
	}
	order.TransactionIDs = append(order.TransactionIDs, shipingTransaction.TransactionID)

	transactions = append(transactions, shipingTransaction)

	// Amount inclusive of sipping fee
	a := -(transaction.Amount + shipingTransaction.Amount)
	// Update wallet and make sure its not empty
	if err := db.Adapter.UpdateWallet(order.UserID, &a); err != nil {
		slog.Error("error updating wallet", "error", err)
		newCaption := "We were unable to process your order. Please retry or reach out to support"
		if strings.Contains(err.Error(), "insufficient funds") {
			newCaption = err.Error() + "Kindly top up your wallet and try again. Use /wallet command for more information"
		}
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
	// Save transactions
	if err := db.Adapter.CreateTransactions(transactions); err != nil {
		slog.Error("error updating wallet", "error", err)
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

	// send referree and eanring message
	if user.Referee != "" {
		referee, _ := db.Adapter.GetUserByRefID(user.Referee)
		go SendRefferalTransaction(bot, *referee, transactions[1])

	}
	if err := db.Adapter.CreateOrder(&order); err != nil {
		slog.Error("error creating order", "error", err)
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
	balance, _ := db.Adapter.GetWallet(order.UserID)
	newCaption := fmt.Sprintf(
		"✅ *Order placed successfully!*\n\n🆔 *Order ID:* %s\n🔗 *Transaction IDs:* %v\n Order Price: %.5f\n\n *Shipping Continent:*%s \n *ShippingCountry*: %s\n *ShippingToCity:* %s \nShipping fee:%.5f \n\n\n Order Total + Shipping fee:\n %.5f💰 *Wallet Balance:* %.5f USDT",
		order.ID.Hex(),
		order.TransactionIDs,
		order.TotalPrice,
		user.ShippingDetails.Continent,
		user.ShippingDetails.Country,
		user.ShippingDetails.City,
		order.ShippingFee,
		order.ShippingFee+order.TotalPrice,
		float64(balance.Balance)/1000000,
	)
	orderButton := tgbotapi.NewInlineKeyboardButtonData("View Order details", fmt.Sprintf("/add_to_cart_%s", order.ID.String()))
	//transactionButton := tgbotapi.NewInlineKeyboardButtonData("View Transaction", fmt.Sprintf("/remove_from_cart_%s", order.TransactionID))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(orderButton),
		//tgbotapi.NewInlineKeyboardRow(),
	)
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, newCaption, keyboard)
	_, err = bot.Send(editMsg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}

func handleRemoveFromCart(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Get the original message
	message := callbackQuery.Message
	cbd := strings.Split(strings.TrimPrefix(callbackQuery.Data, "/remove_from_cart_"), "_")
	itemID, _ := strconv.ParseInt(cbd[0], 10, 64)
	skuID, _ := strconv.ParseInt(cbd[1], 10, 64)
	item, _ := sessions.GetCartItem(context.TODO(), message.Chat.ID, itemID)
	item.Quantity[skuID]--
	sessions.AddItemToCart(context.TODO(), message.Chat.ID, item)
	// Append the new information to the caption
	var messageText string
	settings, err := db.Adapter.GetSettings()
	if err != nil {
		slog.Error(err.Error())
	}
	for variation, base := range item.Details.Result.Item.Sku.Base {
		if base.SkuId == skuID {
			messageText = fmt.Sprintf("Variation: %d", variation)
			proppaths := strings.Split(base.PropPath, ";")
			var props []taobao.PropValues
			for _, path := range proppaths {
				ps := strings.Split(path, ":")
				for _, prop := range item.Details.Result.Item.Sku.Props {
					if strconv.Itoa(prop.PID) == ps[0] {
						for _, val := range prop.Values {
							if strconv.Itoa(val.VID) == ps[1] {
								props = append(props, val)
							}
						}
					}
				}
			}
			for i, p := range props {
				messageText = fmt.Sprintf("%s\n\n Variation Property %d: %s", messageText, i, p.Name)
			}
			rate, _ := payments.GetCnyToUsdtRate()
			p, _ := strconv.ParseFloat(base.PromotionPrice, 64)
			usdtPrice := payments.ConvertCnyToUsdt(p, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			messageText = fmt.Sprintf("%s\n\nVariation Price: %s\nVariation Promotion Price: %.3f USDT\nVariation Quantity: %d", messageText, base.Price, thePrice, base.Quantity)
			break
		}
	}
	newCaption := fmt.Sprintf("%s\nTitle: %s \nPrice: %.4f USDT \n\nQuantity in cart: %d", messageText, item.ItemTitle, item.Price, item.Quantity[skuID])
	addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add ➕1️⃣", fmt.Sprintf("/add_to_cart_%d_%d", item.ItemID, skuID))
	RemoveFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove ➖1️⃣", fmt.Sprintf("/remove_from_cart_%d_%d", item.ItemID, skuID))
	clearFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove item from Cart ❌", fmt.Sprintf("/clear_from_cart_%d_%d", item.ItemID, skuID))
	checkoutButton := tgbotapi.NewInlineKeyboardButtonData("♻️ Checkout ✅", "♻️ Checkout ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("🛒 Cart", "🛒 Cart")
	// Edit the message caption
	editMsg := tgbotapi.NewEditMessageText(
		message.Chat.ID,
		message.MessageID,
		newCaption,
	)
	var keyboard tgbotapi.InlineKeyboardMarkup
	if item.Quantity[skuID] == 0 {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton, RemoveFromCartButton),
			tgbotapi.NewInlineKeyboardRow(checkoutButton, cButton),
			tgbotapi.NewInlineKeyboardRow(clearFromCartButton),
		)
	}
	editMsg.ReplyMarkup = &keyboard

	// Send the edit message
	_, err = bot.Send(editMsg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}

	// Answer the callback query to acknowledge the button click
	answerCallback := tgbotapi.NewCallback(callbackQuery.ID, "Item removed from cart!")
	_, err = bot.Send(answerCallback)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
	if isempty, _ := sessions.IsCartEmpty(context.TODO(), editMsg.ChatID); isempty {
		RemoveCartButton(bot, editMsg.ChatID)
	}
}

func handleClearFromCart(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Get the original message
	message := callbackQuery.Message
	cbd := strings.Split(strings.TrimPrefix(callbackQuery.Data, "/clear_from_cart_"), "_")
	itemID, _ := strconv.ParseInt(cbd[0], 10, 64)
	skuID, _ := strconv.ParseInt(cbd[1], 10, 64)
	item, _ := sessions.GetCartItem(context.TODO(), message.Chat.ID, itemID)
	item.Quantity[skuID] = 0
	sessions.AddItemToCart(context.TODO(), message.Chat.ID, item)
	// Append the new information to the caption
	//newCaption := fmt.Sprintf("Item retrieved: \n Title: %s \n Price: %.4f USDT \n\n Quantity in cart: %d", item.ItemTitle, item.Price, item.Quantity)
	var messageText string
	settings, err := db.Adapter.GetSettings()
	if err != nil {
		slog.Error(err.Error())
	}
	for variation, base := range item.Details.Result.Item.Sku.Base {
		if base.SkuId == skuID {
			messageText = fmt.Sprintf("Variation: %d", variation)
			proppaths := strings.Split(base.PropPath, ";")
			var props []taobao.PropValues
			for _, path := range proppaths {
				ps := strings.Split(path, ":")
				for _, prop := range item.Details.Result.Item.Sku.Props {
					if strconv.Itoa(prop.PID) == ps[0] {
						for _, val := range prop.Values {
							if strconv.Itoa(val.VID) == ps[1] {
								props = append(props, val)
							}
						}
					}
				}
			}
			for i, p := range props {
				messageText = fmt.Sprintf("%s\n\n Variation Property %d: %s", messageText, i, p.Name)
			}
			rate, _ := payments.GetCnyToUsdtRate()
			p, _ := strconv.ParseFloat(base.PromotionPrice, 64)
			usdtPrice := payments.ConvertCnyToUsdt(p, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			messageText = fmt.Sprintf("%s\n\nVariation Price: %s\nVariation Promotion Price: %.3f USDT\nVariation Quantity: %d", messageText, base.Price, thePrice, base.Quantity)
			break
		}
	}
	newCaption := fmt.Sprintf("%s\nTitle: %s \nPrice: %.4f USDT \n\nQuantity in cart: %d", messageText, item.ItemTitle, item.Price, item.Quantity[skuID])
	addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add ➕1️⃣", fmt.Sprintf("/add_to_cart_%d_%d", item.ItemID, skuID))
	RemoveFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove ➖1️⃣", fmt.Sprintf("/remove_from_cart_%d_%d", item.ItemID, skuID))
	clearFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove item from cart ❌", fmt.Sprintf("/clear_from_cart_%d_%d", item.ItemID, skuID))

	// Edit the message caption
	editMsg := tgbotapi.NewEditMessageText(
		message.Chat.ID,
		message.MessageID,
		newCaption,
	)
	//editMsg.Caption = newCaption
	var keyboard tgbotapi.InlineKeyboardMarkup
	if item.Quantity[skuID] == 0 {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton, RemoveFromCartButton), tgbotapi.NewInlineKeyboardRow(clearFromCartButton),
		)
	}
	editMsg.ReplyMarkup = keyboard
	_, err = bot.Send(editMsg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
	answerCallback := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Item cleared from cart!")
	_, err = bot.Send(answerCallback)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
	if isEmpty, _ := sessions.IsCartEmpty(context.TODO(), editMsg.ChatID); isEmpty {
		RemoveCartButton(bot, editMsg.ChatID)
	}
}

func handleAddToCart(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Get the original message
	message := callbackQuery.Message
	cbd := strings.Split(strings.TrimPrefix(callbackQuery.Data, "/add_to_cart_"), "_")
	itemID, _ := strconv.ParseInt(cbd[0], 10, 64)
	skuID, _ := strconv.ParseInt(cbd[1], 10, 64)
	item, _ := sessions.GetCartItem(context.TODO(), message.Chat.ID, itemID)
	item.Quantity[skuID]++
	sessions.AddItemToCart(context.TODO(), message.Chat.ID, item)
	// Append the new information to the caption
	var messageText string
	settings, err := db.Adapter.GetSettings()
	if err != nil {
		slog.Error(err.Error())
	}
	for variation, base := range item.Details.Result.Item.Sku.Base {
		if base.SkuId == skuID {
			messageText = fmt.Sprintf("Variation: %d", variation)
			proppaths := strings.Split(base.PropPath, ";")
			var props []taobao.PropValues
			for _, path := range proppaths {
				ps := strings.Split(path, ":")
				for _, prop := range item.Details.Result.Item.Sku.Props {
					if strconv.Itoa(prop.PID) == ps[0] {
						for _, val := range prop.Values {
							if strconv.Itoa(val.VID) == ps[1] {
								props = append(props, val)
							}
						}
					}
				}
			}
			for i, p := range props {
				messageText = fmt.Sprintf("%s\n\n Variation Property %d: %s", messageText, i, p.Name)
			}
			rate, _ := payments.GetCnyToUsdtRate()
			p, _ := strconv.ParseFloat(base.PromotionPrice, 64)
			usdtPrice := payments.ConvertCnyToUsdt(p, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			messageText = fmt.Sprintf("%s\n\nVariation Price: %s\nVariation Promotion Price: %.3f USDT\nVariation Quantity: %d", messageText, base.Price, thePrice, base.Quantity)
			break
		}
	}
	newCaption := fmt.Sprintf("%s\nTitle: %s \nPrice: %.4f USDT \n\nQuantity in cart: %d", messageText, item.ItemTitle, item.Price, item.Quantity[skuID])
	addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add ➕1️⃣", fmt.Sprintf("/add_to_cart_%d_%d", item.ItemID, skuID))
	RemoveFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove ➖1️⃣", fmt.Sprintf("/remove_from_cart_%d_%d", item.ItemID, skuID))
	clearFromCartButton := tgbotapi.NewInlineKeyboardButtonData("Remove item from Cart ❌", fmt.Sprintf("/clear_from_cart_%d_%d", item.ItemID, skuID))
	checkoutButton := tgbotapi.NewInlineKeyboardButtonData("♻️ Checkout ✅", "♻️ Checkout ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("🛒 Cart", "🛒 Cart")

	// Edit the message caption
	editMsg := tgbotapi.NewEditMessageText(
		message.Chat.ID,
		message.MessageID,
		newCaption,
	)
	var keyboard tgbotapi.InlineKeyboardMarkup
	if item.Quantity[skuID] == 0 {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addToCartButton, RemoveFromCartButton),
			tgbotapi.NewInlineKeyboardRow(checkoutButton, cButton),
			tgbotapi.NewInlineKeyboardRow(clearFromCartButton),
		)
	}
	editMsg.ReplyMarkup = &keyboard

	// Send the edit message
	_, err = bot.Send(editMsg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}

	// Answer the callback query to acknowledge the button click
	answerCallback := tgbotapi.NewCallback(callbackQuery.ID, "Item added to cart!")
	_, err = bot.Send(answerCallback)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
	if item.Quantity[skuID] == 1 {
		AddCartKeyboard(bot, editMsg.ChatID)
	}
}

func RemoveCartButton(bot *tgbotapi.BotAPI, chatID int64) {
	cartButton := tgbotapi.NewKeyboardButton("⚙️ Settings")
	checkoutButton := tgbotapi.NewKeyboardButton("💡 Help 🫂")
	clearButton := tgbotapi.NewKeyboardButton("🏧 Wallet")
	transactionsButton := tgbotapi.NewKeyboardButton("♐ Transactions ♐")
	ordersButton := tgbotapi.NewKeyboardButton("🔰 Orders 🔰")
	refferals := tgbotapi.NewKeyboardButton("🌐 Refferals 🌐")
	checkButton := tgbotapi.NewInlineKeyboardButtonData("♻️ Checkout ✅", "♻️ Checkout ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("🌐 Refferals 🌐", "🏧 Wallet")
	// Create reply keyboard markup
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(cartButton, checkoutButton),
		tgbotapi.NewKeyboardButtonRow(clearButton),
		tgbotapi.NewKeyboardButtonRow(transactionsButton, ordersButton),
		tgbotapi.NewKeyboardButtonRow(refferals),
	)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(checkButton, cButton),
	)
	// Send message to update the reply keyboard
	replyMsg := tgbotapi.NewMessage(chatID, "You have items in the Cart. Click \n 🛒 Cart to view \n ♻️ Checkout ✅ to place an order")
	replyMsg.ReplyMarkup = struct {
		tgbotapi.InlineKeyboardMarkup
		tgbotapi.ReplyKeyboardMarkup
	}{
		inlineKeyboard,
		replyKeyboard,
	}
	_, err := bot.Send(replyMsg)

	if err != nil {
		slog.Info("Error sending reply message:", "error", err)
	}
}

func handleClearCart(bot *tgbotapi.BotAPI, cid int64) {
	sessions.ClearCart(context.TODO(), cid)
	m := tgbotapi.NewMessage(cid, "The cart was cleared")
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	bot.Send(m)
}

func handleCheckout(bot *tgbotapi.BotAPI, cid int64) {
	m := tgbotapi.NewMessage(cid, "The please confirm to place order:")
	checkButton := tgbotapi.NewInlineKeyboardButtonData("Confirm ✅", "Confirm ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("Calcel Order ❌", "Calcel Order ❌")
	m.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(checkButton, cButton),
	)
	bot.Send(m)
}

func handleCart(bot *tgbotapi.BotAPI, cid int64) {
	cart, _ := sessions.GetCart(context.TODO(), cid)

	var message string
	if len(cart) == 0 {
		message = "Your cart is empty"
	} else {
		message += formatTable(cart)
	}

	msg := tgbotapi.NewMessage(cid, message)
	checkButton := tgbotapi.NewInlineKeyboardButtonData("♻️ Checkout ✅", "♻️ Checkout ✅")
	cButton := tgbotapi.NewInlineKeyboardButtonData("🛒 Cart", "🛒 Cart")
	if len(cart) != 0 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(checkButton, cButton),
		)
	}
	// Send the message
	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}
