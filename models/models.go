package models

import (
	"context"
	"fmt"
	"os"
	"pandaexpress/taobao"
	"sync"

	"encoding/csv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	User struct {
		ID                int64  `bson:"id,,omitempty"`
		UserName          string `bson:"username,omitempty"`
		FirstName         string `bson:"firstname,omitempty"`
		LastName          string `bson:"lastname,omitempty"`
		PreferredLanguage `bson:"preferred_language,omitempty"`
		ShippingDetails   `bson:"shipping_details,omitempty"`
		Referee           string   `bson:"referee,omitempty"`
		Referrals         []string `bson:"referrals,omitempty"`
		ReferralCode      string   `bson:"referralcode,omitempty"`
	}
	Wallet struct {
		Balance int64
		UserID  int64
		Address string
	}

	Order struct {
		ID                       primitive.ObjectID `bson:"_id,omitempty"`
		Items                    []CartItem
		UserID                   int64
		TransactionIDs           []string // This makes it easier to track mutliple transactions related to an order
		ShippingFee              float64
		ShippingFeeTransactionID string
		TotalPrice               float64
		Refferee                 string
	}

	Session struct {
		UserID              int64              `bson:"user_id"`
		LastCommand         string             `bson:"last_command"`
		User                User               `bson:"user"`
		PreferredLanguage   PreferredLanguage  `bson:"preferred_language"`
		Cart                map[int64]CartItem `bson:"cart"`
		CartButtonMessageID int                `bson:"cart_button_message_id"`
	}

	Sessions struct {
		mu         sync.Mutex
		sessions   map[int64]Session
		collection *mongo.Collection
	}

	Transaction struct {
		Type          TransactionType
		Amount        int64
		To            string
		From          any
		TransactionID string `bson:"transaction_id"`
	}
	TransactionType int

	PreferredLanguage struct {
		Code        string `bson:"code,omitempty"`
		EnglishName string `bson:"english_name,omitempty"`
		Name        string `bson:"name,omitempty"`
	}

	CartItem struct {
		ItemID    int64
		ItemTitle string
		ItemURL   string
		Image     string
		Price     float64
		Quantity  map[int64]int
		MessageID int64
		Details   taobao.Response
	}

	Carts struct {
		mu    sync.Mutex
		carts map[int64][]CartItem
	}

	ShippingDetails struct {
		Name             string   `bson:"name,omitempty"`
		Continent        string   `bson:"continent,omitempty"`
		Country          string   `bson:"country,omitempty"`
		City             string   `bson:"city,omitempty"`
		Email            string   `bson:"email,omitempty"`
		Phone            string   `bson:"phone,omitempty"`
		Street           string   `bson:"street,omitempty"`
		AppartmentNumber string   `bson:"an,omitempty"`
		Addresses        []string `bson:"addresses,omitempty"`
	}
)
type SAT int

const (
	Done SAT = iota
	Continent
	Country
	City
	Email
	Phone
	Street
	AN
)

func CheckSA(u User) SAT {
	if u.Continent == "" {
		return Continent
	} else if u.Country == "" {
		return Country
	} else if u.City == "" {
		return City
	} else if u.Email == "" {
		return Email
	} else if u.Phone == "" {
		return Phone
	} else if len(u.Addresses) == 0 {
		return AN
	} else {
		return Done
	}
}

const (
	Deposit TransactionType = iota
	Withdraw
	Purchase
	ShippingFee
	RefferalEarning
)

var (
	Languages = []PreferredLanguage{
		{Name: "Afrikaans", Code: "af", EnglishName: "Afrikaans"},
		{Name: "አማርኛ", Code: "am", EnglishName: "Amharic"},
		{Name: "Български", Code: "bg", EnglishName: "Bulgarian"},
		{Name: "Català", Code: "ca", EnglishName: "Catalan"},
		/*{Name: "中文（香港）", Code: "zh-HK", EnglishName: "Chinese (Hong Kong)"},
		{Name: "中文（简体）", Code: "zh-CN", EnglishName: "Chinese (PRC)"},
		{Name: "中文（繁體）", Code: "zh-TW", EnglishName: "Chinese (Taiwan)"},*/
		{Name:"Chinese (Literary)", Code: "lzh", EnglishName: "Chinese (Literary)"},
		{Name:"Chinese Simplified", Code: "zh-Hans", EnglishName: "Chinese Simplified"},
		{Name:"Chinese Traditional", Code: "zh-Hant", EnglishName: "Chinese Traditional"},

		{Name: "Hrvatski", Code: "hr", EnglishName: "Croatian"},
		{Name: "Čeština", Code: "cs", EnglishName: "Czech"},
		{Name: "Dansk", Code: "da", EnglishName: "Danish"},
		{Name: "Nederlands", Code: "nl", EnglishName: "Dutch"},
		/*{Name: "English (UK)", Code: "en-GB", EnglishName: "English (UK)"},
		{Name: "English (US)", Code: "en-US", EnglishName: "English (US)"},*/
		{Name:"English", Code: "en", EnglishName: "English"},

		{Name: "Eesti", Code: "et", EnglishName: "Estonian"},
		{Name: "Filipino", Code: "fil", EnglishName: "Filipino"},
		{Name: "Suomi", Code: "fi", EnglishName: "Finnish"},
		{Name: "Français (Canada)", Code: "fr-ca", EnglishName: "French (Canada)"},
		{Name: "Français (France)", Code: "fr", EnglishName: "French (France)"},
		{Name: "Deutsch", Code: "de", EnglishName: "German"},
		{Name: "Ελληνικά", Code: "el", EnglishName: "Greek"},
		{Name: "עברית", Code: "he", EnglishName: "Hebrew"},
		{Name: "हिन्दी", Code: "hi", EnglishName: "Hindi"},
		{Name: "Magyar", Code: "hu", EnglishName: "Hungarian"},
		{Name: "Íslenska", Code: "is", EnglishName: "Icelandic"},
		{Name: "Bahasa Indonesia", Code: "id", EnglishName: "Indonesian"},
		{Name: "Italiano", Code: "it", EnglishName: "Italian"},
		{Name: "日本語", Code: "ja", EnglishName: "Japanese"},
		{Name: "한국어", Code: "ko", EnglishName: "Korean"},
		{Name: "Latviešu", Code: "lv", EnglishName: "Latvian"},
		{Name: "Lietuvių", Code: "lt", EnglishName: "Lithuanian"},
		{Name: "Bahasa Melayu", Code: "ms", EnglishName: "Malay"},
		{Name: "Norsk", Code: "nb", EnglishName: "Norwegian"},
		{Name: "Polski", Code: "pl", EnglishName: "Polish"},
		{Name: "Português (Brasil)", Code: "pt", EnglishName: "Portuguese (Brazil)"},
		{Name: "Português (Portugal)", Code: "pt-pt", EnglishName: "Portuguese (Portugal)"},
		{Name: "Română", Code: "ro", EnglishName: "Romanian"},
		{Name: "Русский", Code: "ru", EnglishName: "Russian"},
		{Name: "Српски", Code: "sr-Latn", EnglishName: "Serbian"},
		{Name: "Slovenčina", Code: "sk", EnglishName: "Slovak"},
		{Name: "Slovenščina", Code: "sl", EnglishName: "Slovenian"},
		//{Name: "Español (Latinoamérica)", Code: "es-419", EnglishName: "Spanish (Latin America)"},
		{Name: "Español (España)", Code: "es", EnglishName: "Spanish (Spain)"},
		{Name: "Kiswahili", Code: "sw", EnglishName: "Swahili"},
		{Name: "Svenska", Code: "sv", EnglishName: "Swedish"},
		{Name: "ไทย", Code: "th", EnglishName: "Thai"},
		{Name: "Türkçe", Code: "tr", EnglishName: "Turkish"},
		{Name: "Українська", Code: "uk", EnglishName: "Ukrainian"},
		{Name: "Tiếng Việt", Code: "vi", EnglishName: "Vietnamese"},
		{Name: "IsiZulu", Code: "zu", EnglishName: "Zulu"},
	}
)

func NewSessions(collection *mongo.Collection) *Sessions {
	return &Sessions{collection: collection}
}

func (s *Sessions) UpdateLastCommand(ctx context.Context, userID int64, command string) error {
	s.mu.Lock()
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"last_command": command}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	s.mu.Unlock()
	return err
}

func (s *Sessions) AddItemToCart(ctx context.Context, userID int64, item CartItem) error {
	s.mu.Lock()
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{fmt.Sprintf("cart.%d", item.ItemID): item}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	s.mu.Unlock()
	return err
}

func (s *Sessions) ClearCart(ctx context.Context, userID int64) error {
	s.mu.Lock()
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"cart": make(map[int64]CartItem)}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	s.mu.Unlock()
	return err
}

func (s *Sessions) IsCartEmpty(ctx context.Context, userID int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return false, err
	}
	return len(session.Cart) == 0, nil
}

func (s *Sessions) GetCart(ctx context.Context, userID int64) (map[int64]CartItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session.Cart, nil
}

func (s *Sessions) GetCartButtonMessageID(ctx context.Context, userID int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return 0, err
	}
	return session.CartButtonMessageID, nil
}

func (s *Sessions) GetCartItem(ctx context.Context, userID, itemID int64) (CartItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return CartItem{}, err
	}
	item, exists := session.Cart[itemID]
	if !exists {
		return CartItem{}, mongo.ErrNoDocuments
	}
	return item, nil
}

func (s *Sessions) ClearItemFromCart(ctx context.Context, userID, itemID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filter := bson.M{"user_id": userID}
	update := bson.M{"$unset": bson.M{"cart." + string(itemID): ""}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *Sessions) UpdatePreferredLanguage(ctx context.Context, userID int64, pr PreferredLanguage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"preferred_language": pr}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *Sessions) GetPreferredLanguage(ctx context.Context, userID int64) (PreferredLanguage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return PreferredLanguage{}, err
	}
	return session.PreferredLanguage, nil
}

func (s *Sessions) GetLastCommand(ctx context.Context, userID int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var session Session
	filter := bson.M{"user_id": userID}
	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return "", err
	}
	return session.LastCommand, nil
}

func (s *Sessions) AddSession(ctx context.Context, session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.collection.InsertOne(ctx, session)
	return err
}
func NewCarts(s map[int64][]CartItem) *Carts {
	return &Carts{carts: s}
}

// GetUsingCode retrieves language details using the language code
func GetUsingCode(code string) PreferredLanguage {
	for _, lang := range Languages {
		if lang.Code == code {
			return lang
		}
	}
	return GetUsingCode("en-US") // Return empty struct if code not found
}

// GetUsingEnglishName retrieves language details using the English name
func GetUsingEnglishName(englishName string) PreferredLanguage {
	for _, lang := range Languages {
		if lang.EnglishName == englishName {
			return lang
		}
	}
	return GetUsingEnglishName("English (US)") // Return empty struct if English name not found
}

// GetUsingName retrieves language details using the name
func GetUsingName(name string) PreferredLanguage {
	for _, lang := range Languages {
		if lang.Name == name {
			return lang
		}
	}
	return GetUsingName("English (US)") // Return empty struct if name not found
}

// AddToCart adds an item to the user's cart in a thread-safe manner
func (c *Carts) AddToCart(userID int64, item CartItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the user's cart exists
	//if _, exists := c.carts[userID]; exists {
	// If the cart exists, append the item to the cart
	c.carts[userID] = append(c.carts[userID], item)
	//} else {
	// If the cart does not exist, create a new cart with the item
	c.carts[userID] = []CartItem{item}
	//}
}

// RemoveFromCart removes an item from the user's cart in a thread-safe manner
func (c *Carts) RemoveFromCart(userID int64, itemID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the user's cart exists
	if cart, exists := c.carts[userID]; exists {
		for i, item := range cart {
			if item.ItemID == itemID {
				// Remove the item from the cart
				c.carts[userID] = append(cart[:i], cart[i+1:]...)
				return true // Item successfully removed
			}
		}
	}
	return false // Item not found in the cart
}

const (
	Africa       = "Africa"
	Antarctica   = "Antarctica"
	Asia         = "Asia"
	Australia    = "Australia"
	Europe       = "Europe"
	NorthAmerica = "North America"
	SouthAmerica = "South America"
)

var Continents = []string{Africa, Antarctica, Asia, Australia, Europe, NorthAmerica, SouthAmerica}
var ContinentCountries = map[string][]string{
	"Africa": {
		"Algeria", "Angola", "Benin", "Botswana", "Burkina Faso", "Burundi",
		"Cabo Verde", "Cameroon", "Central African Republic", "Chad", "Comoros",
		"Congo", "Democratic Republic of the Congo", "Djibouti", "Egypt",
		"Equatorial Guinea", "Eritrea", "Eswatini", "Ethiopia", "Gabon", "Gambia",
		"Ghana", "Guinea", "Guinea-Bissau", "Ivory Coast", "Kenya", "Lesotho",
		"Liberia", "Libya", "Madagascar", "Malawi", "Mali", "Mauritania",
		"Mauritius", "Morocco", "Mozambique", "Namibia", "Niger", "Nigeria",
		"Rwanda", "Sao Tome and Principe", "Senegal", "Seychelles", "Sierra Leone",
		"Somalia", "South Africa", "South Sudan", "Sudan", "Swaziland", // (use Eswatini)
		"Tanzania", "Togo", "Tunisia", "Uganda", "Zambia", "Zimbabwe",
	},
	"Antarctica": {}, // No countries in Antarctica
	"Asia": {
		"Afghanistan", "Armenia", "Azerbaijan", "Bahrain", "Bangladesh", "Bhutan",
		"Brunei", "Cambodia", "China", "Cyprus", "East Timor", "Egypt", "Georgia",
		"India", "Indonesia", "Iran", "Iraq", "Israel", "Japan", "Jordan", "Kazakhstan",
		"Kuwait", "Kyrgyzstan", "Laos", "Lebanon", "Malaysia", "Maldives", "Mongolia",
		"Myanmar", "Nepal", "North Korea", "Oman", "Pakistan", "Palestine", "Philippines",
		"Qatar", "Russia", "Saudi Arabia", "Singapore", "South Korea", "Sri Lanka",
		"Syria", "Tajikistan", "Thailand", "Turkey", "Turkmenistan", "United Arab Emirates",
		"Uzbekistan", "Vietnam", "Yemen",
	},
	"Australia": {
		"Australia", // Mainland continent
		"Fiji",
		"Kiribati",
		"Marshall Islands",
		"Micronesia",
		"Nauru",
		"New Zealand",
		"Palau",
		"Papua New Guinea",
		"Samoa",
		"Solomon Islands",
		"Tonga",
		"Tuvalu",
		"Vanuatu",
	},
	"Europe": {
		"Albania", "Andorra", "Austria", "Azerbaijan", "Belarus", "Belgium", "Bosnia and Herzegovina",
		"Bulgaria", "Croatia", "Cyprus", "Czech Republic", "Denmark", "Estonia", "Finland", "France",
		"Georgia", "Germany", "Greece", "Hungary", "Iceland", "Ireland", "Italy", "Kazakhstan", "Kosovo",
		"Latvia", "Liechtenstein", "Lithuania", "Luxembourg", "Malta", "Moldova", "Monaco", "Montenegro",
		"Netherlands", "North Macedonia", "Norway", "Poland", "Portugal", "Romania", "Russia", "San Marino",
		"Serbia", "Slovakia", "Slovenia", "Spain", "Sweden", "Switzerland", "Turkey", "Ukraine", "United Kingdom",
		"Vatican City",
	},
	"North America": {
		"Canada",
		"Mexico",
		"United States",
		// Include these if you consider them part of North America:
		"Guatemala",
		"Belize",
		"El Salvador",
		"Honduras",
		"Nicaragua",
		"Costa Rica",
		"Panama",
	},
	"South America": {
		"Argentina",
		"Bolivia",
		"Brazil",
		"Chile",
		"Colombia",
		"Ecuador",
		"Guyana",
		"Paraguay",
		"Peru",
		"Suriname",
		"Uruguay",
		"Venezuela",
	},
}

var (
	JSONconfig []byte
	Cities     map[string][]string
)

type Location struct {
	ID          int
	Name        string
	StateID     int
	StateCode   string
	StateName   string
	CountryID   int
	CountryCode string
	CountryName string
	Latitude    float64
	Longitude   float64
	WikiDataID  string
}

func init() {
	/*if err := json.Unmarshal(JSONconfig, &Cities); err != nil {
		slog.Error("error reading cities", "error", err)
	}*/
	Cities = make(map[string][]string)
	// Open the CSV file
	file, err := os.Open("./models/cities.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header line (if present)
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read all the remaining records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate through the records and populate the structs
	//var locations []Location
	for _, record := range records {
		/*id, _ := strconv.Atoi(record[0])
		stateID, _ := strconv.Atoi(record[2])
		countryID, _ := strconv.Atoi(record[5])
		latitude, _ := strconv.ParseFloat(record[8], 64)
		longitude, _ := strconv.ParseFloat(record[9], 64)

		location := Location{
			ID:          id,
			Name:        record[1],
			StateID:     stateID,
			StateCode:   record[3],
			StateName:   record[4],
			CountryID:   countryID,
			CountryCode: record[6],
			CountryName: record[7],
			Latitude:    latitude,
			Longitude:   longitude,
			WikiDataID:  record[10],
		}*/
		Cities[record[7]] = append(Cities[record[7]], record[1])
		//locations = append(locations, location)
	}

	// Print the locations to verify
	/*for _, location := range locations {
		fmt.Printf("%+v\n", location)
	}*/
}

type Settings struct {
	ReferralCommission float64                    `bson:"referral_commission"`
	ProductCommission  float64                    `bson:"product_commission"`
	ShippingPrices     map[string][]ShippingPrice `bson:"shipping_prices"`
}

type ShippingPrice struct {
	City  string  `bson:"city"`
	Price float64 `bson:"price"`
}
