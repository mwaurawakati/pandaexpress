package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"pandaexpress/db"
	"pandaexpress/models"
	"pandaexpress/payments"
	"pandaexpress/telegram_bot"
	_ "pandaexpress/translate"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("paniced", "error", err)
		}
	}()
	var jsonRaw json.RawMessage = []byte(`{"uri":"mongodb+srv://developer:YC9Ypt7ZZBYgynk0@pandaexpress.eu5f3ng.mongodb.net/?retryWrites=true&w=majority&appName=PandaExpress"}`)
	if err := db.Adapter.Open(jsonRaw); err != nil {
		slog.Error("error creating database", "error", err)
		slog.Warn("exiting")
		os.Exit(1)
	}

	// Initialize settings
	collection := db.Adapter.GetCollection("settings")
	ctx := context.TODO()
	if err := initializeSettings(ctx, collection); err != nil {
		slog.Error("Failed to initialize settings","error", err)
		slog.Warn("exiting")
		os.Exit(1)
	}

	eventsChannel := make(chan *[]payments.L)
	// listen to transaction
	go payments.Listen(eventsChannel)
	go telegram_bot.RunBot(eventsChannel)
	srv := &http.Server{
		Addr:        ":8080",
		Handler:     router(),
		IdleTimeout: time.Minute,
	}
	slog.Info("starting the ui server at :8080")
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("error listening and serving the ui", "error", err)
	}
}

func router() http.Handler {

	mux := http.NewServeMux()

	// index
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/callback", CallbackHandler)

	// static files
	httpFS := http.FileServer(http.Dir("ui/dist"))
	mux.Handle("/static/", httpFS)
	mux.HandleFunc("/uploads/", serveImageHandler)
	// api
	mux.HandleFunc("/api/v1/greeting", greetingAPI)
	mux.HandleFunc("/api/v1/url", handler)
	mux.HandleFunc("/api/v1/set-url", handleSetUrl)
	return mux
}

func handleSetUrl(w http.ResponseWriter, r *http.Request) {
	ur := r.URL.Query().Get("url")
	if ur == ""{
		fmt.Fprintf(w, "Kindle set the param using the url query param")
		return
	}

	if string(ur[len(ur)-1]) == "/"{
		ur = ur[:len(ur)-1]
	}
	telegram_bot.ServiceURL = ur
	fmt.Fprintf(w, "Image url set to: %s", ur)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Image search url is: %s",telegram_bot.ServiceURL )
}

func getCloudRunURL(service, revision string) string {
	if service == "" {
		service = "default"
	}
	if revision == "" {
		revision = "default"
	}

	// Construct the Cloud Run URL
	region := "africa-south1" // Replace with your Cloud Run region
	//projectID := "PROJECT_ID" // Replace with your Google Cloud Project ID

	return fmt.Sprintf("https://%s-%s-%s.a.run.app", service, revision, region)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, "ui/dist/favicon.ico")
		return
	}

	http.ServeFile(w, r, "ui/dist/index.html")
}

func greetingAPI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, there!"))
}

// CallbackHandler handles the OAuth callback and extracts the authorization code.
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the 'code' query parameter from the URL.
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// For demonstration purposes, just print the code to the console.
	fmt.Printf("Received OAuth code: %s\n", code)

	// Respond to the client (optional).
	fmt.Fprintf(w, "Received code: %s", code)

	// Here you can proceed with exchanging the code for an access token, etc.
}

// Handler for serving images from /uploads directory
func serveImageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image filename from the request URL
	imagePath := r.URL.Path[len("/uploads/"):] // Remove "/uploads/" prefix to get the filename
	//fmt.Println(imagePath)
	// Construct the full path to the image file
	uploadsDir := "./uploads"
	imageFilePath := filepath.Join(uploadsDir, imagePath)

	// Open the image file
	file, err := os.Open(imageFilePath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get the file's content type
	contentType := getContentType(imageFilePath)
	if contentType == "" {
		http.Error(w, "Unknown file type", http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=604800")                      // Cache for one week (adjust as needed)
	w.Header().Set("Expires", time.Now().AddDate(0, 0, 7).Format(http.TimeFormat)) // Cache for one week (adjust as needed)

	// Serve the file content
	http.ServeContent(w, r, imageFilePath, time.Now(), file)
}

// Function to get the content type of a file based on its extension
func getContentType(filename string) string {
	switch filepath.Ext(filename) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return ""
	}
}

func initializeSettings(ctx context.Context, collection *mongo.Collection) error {
	var settings models.Settings
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		defaultSettings := createDefaultSettings()
		_, err := collection.InsertOne(ctx, defaultSettings)
		return err
	}
	return err
}

func createDefaultSettings() models.Settings {
	shippingPrices := make(map[string][]models.ShippingPrice)
	for country, cities := range models.Cities {
		for _, city := range cities {
			shippingPrices[country] = append(shippingPrices[country], models.ShippingPrice{
				City:  city,
				Price: 100,
			})
		}
	}

	return models.Settings{
		ReferralCommission: 0.20,
		ProductCommission:  0.20,
		ShippingPrices:     shippingPrices,
	}
}


