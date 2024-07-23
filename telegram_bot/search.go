package telegram_bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"pandaexpress/taobao"
	"pandaexpress/tgbotapi"
	"pandaexpress/translate"
	"strconv"
	"strings"
)

var ServiceURL string
func searchTaobao(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	photo := update.Message.Photo[len(update.Message.Photo)-1]
	newPhoto := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(photo.FileID))
	newPhoto.Caption = "Your photo was received and a search is being conducted"

	// Send the new photo message with caption
	_, err := bot.Send(newPhoto)
	if err != nil {
		slog.Error("Failed to send photo", "error", err)
	}

	fp, err := saveTelegramImageToLocal(bot, photo.FileID)
	if err != nil {
		slog.Error(err.Error())
	}
	settings, err := db.Adapter.GetSettings()
	if err != nil{
		slog.Error(err.Error())
	}
	// Construct the file URL
	fileURL := fmt.Sprintf("%s/%s",ServiceURL, fp)
	go func(message tgbotapi.Message) {
		defer os.Remove(fp)
		resp, err := taobao.SearchImageOnTaobao(fileURL, 1)
		if err != nil {
			sendErrorMessageToOwner(bot, err.Error())
			return
		}
		if len(resp.Result.ResultList) == 0 {
			newPhoto.Caption = "There are no results for this image"
			slog.Info("res", "resp", resp)
			_, err := bot.Send(newPhoto)
			if err != nil {
				sendErrorMessageToOwner(bot, err.Error())
				return
			}
		}
		rate, err := payments.GetCnyToUsdtRate()
		if err != nil {
			slog.Error(err.Error())
		}
		for _, res := range resp.Result.ResultList {
			itemID, err := strconv.ParseInt(res.Item.ItemID.(string), 10, 64)
			if err != nil {
				slog.Error(err.Error())
			}
			preferredLanguage, err := sessions.GetPreferredLanguage(context.TODO(), message.Chat.ID)
			if err != nil {
				slog.Error(err.Error())
			}
			title, err := translate.TranslateTextToPreferredLanguage(res.Item.Title, preferredLanguage.Code)
			if err != nil {
				slog.Error("error transalating tittle", "error", err, "preffered_lang", preferredLanguage)
			}
			p, err := strconv.ParseFloat(res.Item.Sku.Def.PromotionPrice, 64)
			if err != nil {
				slog.Error(err.Error())
			}
			usdtPrice := payments.ConvertCnyToUsdt(p, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			if err != nil {
				title = res.Item.Title
			}
			item := models.CartItem{
				ItemID:    itemID,
				ItemTitle: title,
				Price:     thePrice,
				Quantity:  make(map[int64]int),
				Image:     fmt.Sprintf("http:%s", res.Item.Image),
				ItemURL:   fmt.Sprintf("http:%s", res.Item.ItemURL),
			}
			//carts.AddToCart(message.Chat.ID, item)
			b, err := fetchImageAndSave("http:" + res.Item.Image)
			defer os.Remove(b)
			if err != nil {
				slog.Error(err.Error())
			}
			fpp := tgbotapi.FilePath(b)
			newPhoto := tgbotapi.NewPhoto(message.Chat.ID, fpp)
			newPhoto.Caption = fmt.Sprintf("Item retrieved: \n Title: %s \n Price: %.4f USDT \n\n Quantity in cart: %d \n Image: %s \n ItemURL: %s", item.ItemTitle, item.Price, 0, item.Image, item.ItemURL)
			/*if err != nil {
				slog.Error("error marshalling data", "error", err)
			}*/

			addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add to Cart", fmt.Sprintf("/add_to_cart_%d", item.ItemID))
			_ = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(addToCartButton),
			)
			//newPhoto.ReplyMarkup = &keyboard
			// Send the new photo message with caption
			m, err := bot.Send(newPhoto)
			if err != nil {
				slog.Error("Failed to send photo", "error", err)
			}
			item.MessageID = int64(m.MessageID)
			//sessions.AddItemToCart(message.Chat.ID, item)

			// after sending the message, we can query the item
			itemDetails, err := taobao.SearchItemDetails(res.Item.ItemID.(string))
			if err != nil {
				slog.Error("error fetching item details")
			} else {
				// send images
				sendItemDescription(bot, &m, itemDetails.Result.Item.Images)
				// send description images as reply to the first image
				sendItemDescription(bot, &m, itemDetails.Result.Item.Description.Images)
				// loop through the items
				for variation, base := range itemDetails.Result.Item.Sku.Base {
					proppaths := strings.Split(base.PropPath, ";")
					//fmt.Println(proppaths)
					// Get props
					var props []taobao.PropValues
					for _, path := range proppaths {
						ps := strings.Split(path, ":")
						//fmt.Println(ps)
						for _, prop := range itemDetails.Result.Item.Sku.Props {
							if strconv.Itoa(prop.PID) == ps[0] {
								for _, val := range prop.Values {
									if strconv.Itoa(val.VID) == ps[1] {
										props = append(props, val)
									}
								}
							}
						}
					}
					//fmt.Printf("%+v\n", props)
					// Send a reply message
					vm := fmt.Sprintf("Variation: %d\n\nTitle: %s", variation, item.ItemTitle)
					var mediaGroup []interface{}
					var filesToRemove []string
					defer func() {
						for _, file := range filesToRemove {
							os.Remove(file)
						}
					}()
					for i, p := range props {
						if p.Image != "" {
							fullURL := "http:" + p.Image
							photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(fullURL))
							mediaGroup = append(mediaGroup, photo)
						}
						vname, _ := translate.TranslateTextToPreferredLanguage(p.Name, preferredLanguage.Code)
						p.Name = vname
						vm = fmt.Sprintf("%s\n\nVariation Property %d: %s", vm, i, p.Name)
					}
					// Create the media group message
					if len(mediaGroup) != 0 {
						mediaMsg := tgbotapi.NewMediaGroup(message.Chat.ID, mediaGroup)
						mediaMsg.ReplyToMessageID = message.MessageID
						_, err := bot.Send(mediaMsg)
						if err != nil {
							slog.Error("Error sending media group", "error", err)
						}
					}
					// Add price
					p, err := strconv.ParseFloat(base.PromotionPrice, 64)
					if err != nil {
						slog.Error(err.Error())
					}
					usdtPrice := payments.ConvertCnyToUsdt(p, rate)
					thePrice := usdtPrice * (settings.ProductCommission + 1)
					vm = fmt.Sprintf("%s\n\nVariation Price: %s\nVariation Promotion Price: %.3f USDT\nVariation Quantity: %d", vm, base.Price, thePrice, base.Quantity)
					/*b, err := fetchImageAndSave("http:" + res.Item.Image)
					defer os.Remove(b)
					if err != nil {
						slog.Error(err.Error())
					}*/
					//fpp := tgbotapi.FilePath(b)
					newPhoto := tgbotapi.NewMessage(message.Chat.ID, vm)
					//newPhoto.Caption = vm
					//vMessage := tgbotapi.NewPhoto(message.Chat.ID, vm)
					//newPhoto.ReplyToMessageID = mesMes.MessageID
					addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add to Cart", fmt.Sprintf("/add_to_cart_%d_%d", item.ItemID, base.SkuId))
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(addToCartButton),
					)
					newPhoto.ReplyMarkup = keyboard
					if _, err := bot.Send(newPhoto); err != nil {
						slog.Error("Error sending normal variation message", "error", err)
					}
				}
			}
			item.Details = *itemDetails
			sessions.AddItemToCart(context.TODO(), message.Chat.ID, item)
		}

	}(*update.Message)
}

func sendItemDescription(bot *tgbotapi.BotAPI, message *tgbotapi.Message, imageURLs []string) {
	batchSize := 10
	for i := 0; i < len(imageURLs); i += batchSize {
		end := i + batchSize
		if end > len(imageURLs) {
			end = len(imageURLs)
		}

		var mediaGroup []interface{}
		for _, url := range imageURLs[i:end] {
			fullURL := "http:" + url
			photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(fullURL))
			mediaGroup = append(mediaGroup, photo)
		}

		// Create the media group message
		mediaMsg := tgbotapi.NewMediaGroup(message.Chat.ID, mediaGroup)
		mediaMsg.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(mediaMsg); err != nil {
			slog.Error("Error sending media group", "error", err)
		}
	}
}

func sendImageRequestMessage(bot *tgbotapi.BotAPI, chatID int64) {
	imageRequestText := "Please upload an image to search for products"

	msg := tgbotapi.NewMessage(chatID, imageRequestText)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message:", "error", err)
	}
}

// Function to send an error message to the owner of the bot
func sendErrorMessageToOwner(bot *tgbotapi.BotAPI, errorMessage string) {
	msg := tgbotapi.NewMessage(6721747351, errorMessage)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message", "error", err)
	}
}


func search1688(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	photo := update.Message.Photo[len(update.Message.Photo)-1]
	newPhoto := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(photo.FileID))
	newPhoto.Caption = "Your photo was received and a search is being conducted"

	// Send the new photo message with caption
	_, err := bot.Send(newPhoto)
	if err != nil {
		slog.Error("Failed to send photo", "error", err)
	}

	fp, err := saveTelegramImageToLocal(bot, photo.FileID)
	if err != nil {
		slog.Error(err.Error())
	}
	settings, err := db.Adapter.GetSettings()
	if err != nil{
		slog.Error(err.Error())
	}
	// Construct the file URL
	fileURL := fmt.Sprintf("%s/%s",ServiceURL, fp)
	go func(message tgbotapi.Message) {
		defer os.Remove(fp)
		resp, err := taobao.SearchImageOn1688(fileURL, 1)
		if err != nil {
			sendErrorMessageToOwner(bot, err.Error())
			return
		}
		if len(resp.Result.ResultList) == 0 {
			newPhoto.Caption = "There are no results for this image"
			slog.Info("res", "resp", resp)
			_, err := bot.Send(newPhoto)
			if err != nil {
				sendErrorMessageToOwner(bot, err.Error())
				return
			}
		}
		rate, err := payments.GetCnyToUsdtRate()
		if err != nil {
			slog.Error(err.Error())
		}
		for _, res := range resp.Result.ResultList {
			var itemID int64
			switch t := res.Item.ItemID.(type) {
			case int:
				itemID = int64(t)
			case int8:
				itemID = int64(t)
			case int16:
				itemID = int64(t)
			case int32:
				itemID = int64(t)
			case int64:
				itemID = t
			case uint:
				itemID = int64(t)
			case uint8:
				itemID = int64(t)
			case uint16:
				itemID = int64(t)
			case uint32:
				itemID = int64(t)
			case uint64:
				itemID = int64(t)
			case float32:
				itemID = int64(t)
			case float64:
				itemID = int64(t)
			case string:
				// Try to parse the string as an int64
				if parsedInt, err := strconv.ParseInt(t, 10, 64); err == nil {
					itemID = parsedInt
				} else if parsedFloat, err := strconv.ParseFloat(t, 64); err == nil {
					itemID = int64(parsedFloat)
				} else {
					fmt.Printf("Error parsing string as int64 or float64: %v\n", err)
					return
				}
			default:
				fmt.Printf("Unsupported type: %T\n", t)
				return
			}
			
			preferredLanguage, err := sessions.GetPreferredLanguage(context.TODO(), message.Chat.ID)
			if err != nil {
				slog.Error(err.Error())
			}
			title, err := translate.TranslateTextToPreferredLanguage(res.Item.Title, preferredLanguage.Code)
			if err != nil {
				slog.Error("error transalating tittle", "error", err, "preffered_lang", preferredLanguage)
			}
			p, err := strconv.ParseFloat(res.Item.Sku.Def.PromotionPrice, 64)
			if err != nil {
				slog.Error(err.Error())
			}
			usdtPrice := payments.ConvertCnyToUsdt(p, rate)
			thePrice := usdtPrice * (settings.ProductCommission + 1)
			if err != nil {
				title = res.Item.Title
			}
			item := models.CartItem{
				ItemID:    itemID,
				ItemTitle: title,
				Price:     thePrice,
				Quantity:  make(map[int64]int),
				Image:     fmt.Sprintf("http:%s", res.Item.Image),
				ItemURL:   fmt.Sprintf("http:%s", res.Item.ItemURL),
			}
			//carts.AddToCart(message.Chat.ID, item)
			b, err := fetchImageAndSave("http:" + res.Item.Image)
			defer os.Remove(b)
			if err != nil {
				slog.Error(err.Error())
			}
			fpp := tgbotapi.FilePath(b)
			newPhoto := tgbotapi.NewPhoto(message.Chat.ID, fpp)
			newPhoto.Caption = fmt.Sprintf("Item retrieved: \n Title: %s \n Price: %.4f USDT \n\n Quantity in cart: %d \n Image: %s \n ItemURL: %s", item.ItemTitle, item.Price, 0, item.Image, item.ItemURL)
			/*if err != nil {
				slog.Error("error marshalling data", "error", err)
			}*/

			addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add to Cart", fmt.Sprintf("/add_to_cart_%d", item.ItemID))
			_ = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(addToCartButton),
			)
			//newPhoto.ReplyMarkup = &keyboard
			// Send the new photo message with caption
			m, err := bot.Send(newPhoto)
			if err != nil {
				slog.Error("Failed to send photo", "error", err)
			}
			item.MessageID = int64(m.MessageID)
			//sessions.AddItemToCart(message.Chat.ID, item)

			// after sending the message, we can query the item
			itemDetails, err := taobao.SearchItemDetails1688(fmt.Sprintf("%v",res.Item.ItemID))
			if err != nil {
				slog.Error("error fetching item details")
			} else {
				// send images
				sendItemDescription(bot, &m, itemDetails.Result.Item.Images)
				// send description images as reply to the first image
				sendItemDescription(bot, &m, itemDetails.Result.Item.Description.Images)
				// loop through the items
				for variation, base := range itemDetails.Result.Item.Sku.Base {
					proppaths := strings.Split(base.PropPath, ";")
					//fmt.Println(proppaths)
					// Get props
					var props []taobao.PropValues
					for _, path := range proppaths {
						ps := strings.Split(path, ":")
						//fmt.Println(ps)
						for _, prop := range itemDetails.Result.Item.Sku.Props {
							if strconv.Itoa(prop.PID) == ps[0] {
								for _, val := range prop.Values {
									if strconv.Itoa(val.VID) == ps[1] {
										props = append(props, val)
									}
								}
							}
						}
					}
					//fmt.Printf("%+v\n", props)
					// Send a reply message
					vm := fmt.Sprintf("Variation: %d\n\nTitle: %s", variation, item.ItemTitle)
					var mediaGroup []interface{}
					var filesToRemove []string
					defer func() {
						for _, file := range filesToRemove {
							os.Remove(file)
						}
					}()
					for i, p := range props {
						if p.Image != "" {
							fullURL := "http:" + p.Image
							photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(fullURL))
							mediaGroup = append(mediaGroup, photo)
						}
						vname, _ := translate.TranslateTextToPreferredLanguage(p.Name, preferredLanguage.Code)
						p.Name = vname
						vm = fmt.Sprintf("%s\n\nVariation Property %d: %s", vm, i, p.Name)
					}
					// Create the media group message
					if len(mediaGroup) != 0 {
						mediaMsg := tgbotapi.NewMediaGroup(message.Chat.ID, mediaGroup)
						mediaMsg.ReplyToMessageID = message.MessageID
						_, err := bot.Send(mediaMsg)
						if err != nil {
							slog.Error("Error sending media group", "error", err)
						}
					}
					// Add price
					p, err := strconv.ParseFloat(base.PromotionPrice, 64)
					if err != nil {
						slog.Error(err.Error())
					}
					usdtPrice := payments.ConvertCnyToUsdt(p, rate)
					thePrice := usdtPrice * (settings.ProductCommission + 1)
					vm = fmt.Sprintf("%s\n\nVariation Price: %s\nVariation Promotion Price: %.3f USDT\nVariation Quantity: %d", vm, base.Price, thePrice, base.Quantity)
					/*b, err := fetchImageAndSave("http:" + res.Item.Image)
					defer os.Remove(b)
					if err != nil {
						slog.Error(err.Error())
					}*/
					//fpp := tgbotapi.FilePath(b)
					newPhoto := tgbotapi.NewMessage(message.Chat.ID, vm)
					//newPhoto.Caption = vm
					//vMessage := tgbotapi.NewPhoto(message.Chat.ID, vm)
					//newPhoto.ReplyToMessageID = mesMes.MessageID
					addToCartButton := tgbotapi.NewInlineKeyboardButtonData("Add to Cart", fmt.Sprintf("/add_to_cart_%d_%d", item.ItemID, base.SkuId))
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(addToCartButton),
					)
					newPhoto.ReplyMarkup = keyboard
					if _, err := bot.Send(newPhoto); err != nil {
						slog.Error("Error sending normal variation message", "error", err)
					}
				}
			}
			item.Details = *itemDetails
			sessions.AddItemToCart(context.TODO(), message.Chat.ID, item)
		}

	}(*update.Message)
}