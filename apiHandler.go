package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"pandaexpress/db"
	"pandaexpress/models"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
)

// User represents the structure of user data received from the Telegram Web App
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON data
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the user data (you can save it to a database, etc.)
	fmt.Printf("Received user data: %+v\n", user)

	// create user if not exist which is unlikely since one must start with /start command

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func continentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Continents)
}

func getCountriesByContinent(w http.ResponseWriter, r *http.Request) {
	continent := r.URL.Query().Get("continent")
	if continent == "" {
		http.Error(w, "Continent query parameter is missing", http.StatusBadRequest)
		return
	}

	countries, ok := models.ContinentCountries[continent]
	if !ok {
		http.Error(w, "Continent not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func getCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		http.Error(w, "Country query parameter is missing", http.StatusBadRequest)
		return
	}

	cities, ok := models.Cities[country]
	if !ok {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func submitShippingAddress(w http.ResponseWriter, r *http.Request) {
	var address models.ShippingAddress
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error setting preffered language", "error": err.Error()})
		return
	}
	user, err := db.Adapter.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error setting preffered language", "error": err.Error()})
		return
	}
	// Decode the request body into the ShippingAddress struct
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the form fields
	if address.FirstName == "" || address.LastName == "" || address.Address == "" ||
		address.City == "" || address.Continent == "" || address.Zip == "" ||
		address.Country == "" || address.Phone == "" || address.Email == "" {
		http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		return
	}
	sa := models.ShippingDetails{
		Name:      address.FirstName + " " + address.LastName,
		Continent: address.Continent,
		Country:   address.Country,
		City:      address.City,
		Email:     address.Email,
		Phone:     address.Phone,
		Addresses: []string{address.Address},
	}
	if address.Address2 != "" {
		sa.Addresses = append(sa.Addresses, address.Address2)
	}
	user.ShippingDetails = sa
	if err := db.Adapter.UpdateUser(userID, *user); err != nil {
		http.Error(w, "Failed to update preferred language", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Form submitted successfully!"}`))
}

func SetPreferredLanguageHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PreferredLanguage
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error setting preffered language", "error": err.Error()})
		return
	}
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		slog.Error("error decoding json on set lang handler", "error", err)
		return
	}
	pl := models.GetUsingCode(req.Code)
	user, err := db.Adapter.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error setting preffered language", "error": err.Error()})
		return
	}
	user.PreferredLanguage = pl
	// Update the preferred language in the database

	if err := db.Adapter.UpdateUser(userID, *user); err != nil {
		http.Error(w, "Failed to update preferred language", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Preferred language set successfully"})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		slog.Error("error parsing userid","error", err)
		return
	}
	slog.Info("user", "id", userID)
	user, err := db.Adapter.GetUser(userID)
	if err != nil {
		slog.Error("error retrieving user", "error", err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}


// Handler function for getting wallet information
func getWalletHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userIDStr := queryParams.Get("user_id")

	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	wallet, err := db.Adapter.GetWallet(userID)
	if err != nil {
		slog.Error("Error getting wallet:","error", err)
		http.Error(w, "Error getting wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userIDStr := queryParams.Get("user_id")
	numStr := queryParams.Get("num")
	offsetStr := queryParams.Get("offset")

	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	num := 50 // Default number of transactions
	if numStr != "" {
		num, err = strconv.Atoi(numStr)
		if err != nil {
			http.Error(w, "invalid num parameter", http.StatusBadRequest)
			return
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	transactions, err := db.Adapter.GetTransactions(userID, num, offset)
	if err != nil {
		slog.Error("Error getting transactions:","error", err)
		http.Error(w, "Error getting transactions", http.StatusInternalServerError)
		return
	}
	resp := struct{
		Transactions []models.Transaction
		Offset int 
	}{
		transactions,
		offset + num,
	}
	if len(transactions) < num {
		resp.Offset = 0
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userIDStr := queryParams.Get("user_id")
	numStr := queryParams.Get("num")
	offsetStr := queryParams.Get("offset")

	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	num := 50 // Default number of transactions
	if numStr != "" {
		num, err = strconv.Atoi(numStr)
		if err != nil {
			http.Error(w, "invalid num parameter", http.StatusBadRequest)
			return
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	orders, err := db.Adapter.GetOrders(userID, num, offset)
	if err != nil {
		slog.Error("Error getting transactions:","error", err)
		http.Error(w, "Error getting transactions", http.StatusInternalServerError)
		return
	}
	resp := struct{
		Orders []models.Order
		Offset int 
	}{
		orders,
		offset + num,
	}
	if len(orders) < num {
		resp.Offset = 0
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}