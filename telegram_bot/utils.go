package telegram_bot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pandaexpress/models"
	"pandaexpress/tgbotapi"
	"path/filepath"
	"regexp"
)

func fetchImageAndSave(imageURL string) (string, error) {
	// Perform GET request to fetch the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("error fetching image: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	// Create the downloads directory if it doesn't exist
	downloadsDir := "./downloads"
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		if err := os.Mkdir(downloadsDir, 0755); err != nil {
			return "", fmt.Errorf("error creating downloads directory: %v", err)
		}
	}

	// Extract the file name from the URL
	fileName := filepath.Base(imageURL)

	// Create a file in the downloads directory
	filePath := filepath.Join(downloadsDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error copying image to file: %v", err)
	}

	return filePath, nil
}

// Function to save image from Telegram to local directory
func saveTelegramImageToLocal(bot *tgbotapi.BotAPI, fileID string) (string, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get file from Telegram: %v", err)
	}

	// Construct the file path where you want to save the image locally
	uploadsDir := "./uploads"
	filePath := filepath.Join(uploadsDir, file.FileID+".jpg") // Example: Save as JPG

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Create a new file to write to
	localFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	// Download the file content
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Copy the file content to the local file
	_, err = io.Copy(localFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

func formatTable(cartItems map[int64]models.CartItem) string {
	var table string
	var total float64
	// Rows
	for _, item := range cartItems {
		var m string
		i := 0
		for _, c := range item.Quantity {
			m = fmt.Sprintf("%s\n\nVariation: %d\nQuantity:%d", m, i, c)
		}
		itemTotal := 1.0 //item.Price * float64(item.Quantity)
		//total += itemTotal
		table += fmt.Sprintf(
			"\n\n🆔 *Item ID:* %d\n📦 *Item Title:* %s\n💲 \n\nVariations\n%s*Item Price:* %.5f USDT\n🔢 *Item Quantity:* %d\n\n📝 *Item total:* %.5f * %d = %.5f USDT",
			item.ItemID, item.ItemTitle, m, item.Price, item.Quantity, item.Price, item.Quantity, itemTotal,
		)
	}
	table += fmt.Sprintf("\n\n💰 *Total Order Price:*\n%.5f USDT", total)
	return table
}

// ContainsString checks if a string exists in a list of strings.
func ContainsString(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// IsValidEmail validates an email address using a regular expression.
func IsValidEmail(email string) bool {
	// Define a regex pattern for validating email addresses.
	const emailPattern = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

// IsValidPhone validates a phone number using a regular expression.
// This example assumes a general international phone number format.
func IsValidPhone(phone string) bool {
	// Define a regex pattern for validating phone numbers.
	const phonePattern = `^\+?[1-9]\d{1,14}$`
	re := regexp.MustCompile(phonePattern)
	return re.MatchString(phone)
}

func getShippingPrice(settings models.Settings, country, city string) (float64, error) {
	if cities, ok := settings.ShippingPrices[country]; ok {
		for _, shippingPrice := range cities {
			if shippingPrice.City == city {
				return shippingPrice.Price, nil
			}
		}
	}
	return 0, fmt.Errorf("shipping price not found for city: %s, country: %s", city, country)
}
