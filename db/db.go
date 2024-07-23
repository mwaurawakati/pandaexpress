package db

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"pandaexpress/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// adapter holds MongoDB connection data.
type adapter struct {
	conn              *mongo.Client
	db                *mongo.Database
	dbName            string
	maxResults        int
	maxMessageResults int
	version           int
	ctx               context.Context
	useTransactions   bool
}

const (
	defaultHost              = "localhost:27017"
	defaultDatabase          = "panda"
	defaultMaxResults        = 1024
	defaultMaxMessageResults = 100
	defaultAuthMechanism     = "SCRAM-SHA-256"
	defaultAuthSource        = "admin"
)

type configType struct {
	Uri                string                   `json:"uri,omitempty"`
	Addresses          any                      `json:"addresses,omitempty"`
	ConnectTimeout     int                      `json:"timeout,omitempty"`
	Database           string                   `json:"database,omitempty"`
	ReplicaSet         string                   `json:"replica_set,omitempty"`
	AuthMechanism      string                   `json:"auth_mechanism,omitempty"`
	AuthSource         string                   `json:"auth_source,omitempty"`
	Username           string                   `json:"username,omitempty"`
	Password           string                   `json:"password,omitempty"`
	UseTLS             bool                     `json:"tls,omitempty"`
	TlsCertFile        string                   `json:"tls_cert_file,omitempty"`
	TlsPrivateKey      string                   `json:"tls_private_key,omitempty"`
	InsecureSkipVerify bool                     `json:"tls_skip_verify,omitempty"`
	APIVersion         options.ServerAPIVersion `json:"api_version,omitempty"`
}

// Open initializes mongodb session
func (a *adapter) Open(jsonconfig json.RawMessage) error {
	if a.conn != nil {
		return errors.New("adapter mongodb is already connected")
	}

	if len(jsonconfig) < 2 {
		return errors.New("adapter mongodb missing config")
	}

	var err error
	var config configType
	if err = json.Unmarshal(jsonconfig, &config); err != nil {
		return errors.New("adapter mongodb failed to parse config: " + err.Error())
	}

	var opts options.ClientOptions

	if config.Addresses == nil {
		opts.SetHosts([]string{defaultHost})
	} else if host, ok := config.Addresses.(string); ok {
		opts.SetHosts([]string{host})
	} else if ihosts, ok := config.Addresses.([]any); ok && len(ihosts) > 0 {
		hosts := make([]string, len(ihosts))
		for i, ih := range ihosts {
			h, ok := ih.(string)
			if !ok || h == "" {
				return errors.New("adapter mongodb invalid config.Addresses value")
			}
			hosts[i] = h
		}
		opts.SetHosts(hosts)
	} else {
		return errors.New("adapter mongodb failed to parse config.Addresses")
	}

	if config.Database == "" {
		a.dbName = defaultDatabase
	} else {
		a.dbName = config.Database
	}

	if config.ReplicaSet != "" {
		opts.SetReplicaSet(config.ReplicaSet)
		a.useTransactions = true
	} else {
		// Retriable writes are not supported in a standalone instance.
		opts.SetRetryWrites(false)
	}

	if config.Username != "" {
		if config.AuthMechanism == "" {
			config.AuthMechanism = defaultAuthMechanism
		}
		if config.AuthSource == "" {
			config.AuthSource = defaultAuthSource
		}
		var passwordSet bool
		if config.Password != "" {
			passwordSet = true
		}
		opts.SetAuth(
			options.Credential{
				AuthMechanism: config.AuthMechanism,
				AuthSource:    config.AuthSource,
				Username:      config.Username,
				Password:      config.Password,
				PasswordSet:   passwordSet,
			})
	}

	if config.UseTLS {
		tlsConfig := tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}

		if config.TlsCertFile != "" {
			cert, err := tls.LoadX509KeyPair(config.TlsCertFile, config.TlsPrivateKey)
			if err != nil {
				return err
			}

			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}

		opts.SetTLSConfig(&tlsConfig)
	}

	if a.maxResults <= 0 {
		a.maxResults = defaultMaxResults
	}

	if a.maxMessageResults <= 0 {
		a.maxMessageResults = defaultMaxMessageResults
	}

	// Connection string URI overrides any other options configured earlier.
	if config.Uri != "" {
		opts.ApplyURI(config.Uri)
	}

	if config.APIVersion != "" {
		opts.SetServerAPIOptions(options.ServerAPI(config.APIVersion))
	}

	// Make sure the options are sane.
	if err = opts.Validate(); err != nil {
		return err
	}

	a.ctx = context.Background()
	a.conn, err = mongo.Connect(a.ctx, &opts)
	a.db = a.conn.Database(a.dbName)
	if err != nil {
		return err
	}
	collection := a.db.Collection("users")

	// Ensure the ID field is unique
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	collection = a.db.Collection("transactions")

	// Ensure the ID field is unique
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	collection = a.db.Collection("wallets")

	// Ensure the ID field is unique
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "userid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	collection = a.db.Collection("sessions")

	// Ensure the ID field is unique
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	return nil
}

// Close the adapter
func (a *adapter) Close() error {
	var err error
	if a.conn != nil {
		err = a.conn.Disconnect(a.ctx)
		a.conn = nil
		a.version = -1
	}
	return err
}

// IsOpen checks if the adapter is ready for use
func (a *adapter) IsOpen() bool {
	return a.conn != nil
}

func (a *adapter) UserCreate(usr *models.User) error {
	if _, err := a.db.Collection("users").InsertOne(a.ctx, &usr); err != nil {
		return err
	}

	return nil
}
func (a *adapter) GetCollection(col string) *mongo.Collection {
	 return a.db.Collection(col)
}
var Adapter *adapter
func (a *adapter)GetSettings() (models.Settings, error) {
	collection := a.db.Collection("settings")
	var settings models.Settings
	err := collection.FindOne(context.TODO(), bson.M{}).Decode(&settings)
	return settings, err
}

func init() {
	Adapter = &adapter{dbName: defaultDatabase}
}

func (a *adapter) GetUser(userID int64) (*models.User, error) {
	collection := a.db.Collection("users")
	filter := bson.D{{Key: "id", Value: userID}}
	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}
	return &user, nil
}

func (a *adapter) GetUserByRefID(refID string) (*models.User, error) {
	collection := a.db.Collection("users")
	filter := bson.D{{Key: "referralcode", Value: refID}}
	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with ID %s not found", refID)
		}
		return nil, err
	}
	return &user, nil
}

// updateUser updates the user's information in the MongoDB collection by their unique ID.
func (a *adapter) UpdateUser(userID int64, updatedFields models.User) error {
	collection := a.db.Collection("users")
	filter := bson.D{{Key: "id", Value: userID}}
	update := bson.D{{Key: "$set", Value: updatedFields}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}

func (a *adapter) UpdateReferral(referralCode string, newUser string) error {
	collection := a.db.Collection("users")

	// Find the user by referral code
	filter := bson.M{"referralcode": referralCode}
	var referrer models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&referrer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("referral code not found")
		}
		return err
	}

	// Update the referrer's referrals list, only add if not exists
	update := bson.M{
		"$addToSet": bson.M{"referrals": newUser},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// createWallet inserts a new wallet into the MongoDB collection.
func (a *adapter) CreateWallet(wallet models.Wallet) error {
	collection := a.db.Collection("wallets")
	_, err := collection.InsertOne(context.TODO(), wallet)
	if err != nil {
		return err
	}
	return nil
}

// getWallet retrieves a wallet from the MongoDB collection by the user's unique ID.
func (a *adapter) GetWallet(userID int64) (*models.Wallet, error) {
	collection := a.db.Collection("wallets")
	filter := bson.D{{Key: "userid", Value: userID}}
	var wallet models.Wallet
	err := collection.FindOne(context.TODO(), filter).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet with UserID %d not found", userID)
		}
		return nil, err
	}
	return &wallet, nil
}

// getWallet retrieves a wallet from the MongoDB collection by the user's unique ID.
func (a *adapter) GetWalletByAddress(address string) (*models.Wallet, error) {
	collection := a.db.Collection("wallets")
	filter := bson.D{{Key: "address", Value: address}}
	var wallet models.Wallet
	err := collection.FindOne(context.TODO(), filter).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("wallet with UserID %s not found", address)
		}
		return nil, err
	}
	return &wallet, nil
}

// updateWallet updates the wallet's balance in the MongoDB collection by the user's unique ID.
func (a *adapter) UpdateWallet(userID int64, amount *int64) error {
	collection := a.db.Collection("wallets")
	wallet, err := a.GetWallet(userID)
	if err != nil {
		return err
	}

	// Check if the new balance would be less than zero
	newBalance := wallet.Balance + *amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds: current balance %d, transaction amount %d", wallet.Balance, *amount)
	}

	// Update the balance
	wallet.Balance = newBalance

	// Update the wallet in the database
	filter := bson.D{{Key: "userid", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "balance", Value: wallet.Balance}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("wallet with UserID %d not found", userID)
	}

	return nil
}

// createWallet inserts a new wallet into the MongoDB collection.
func (a *adapter) CreateTransaction(transaction *models.Transaction) error {
	collection := a.db.Collection("transactions")
	_, err := collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		return err
	}
	return nil
}

// CreateTransactions inserts multiple transactions into the MongoDB collection.
func (a *adapter) CreateTransactions(transactions []models.Transaction) error {
	collection := a.db.Collection("transactions")

	var docs []interface{}
	for _, t := range transactions {
		docs = append(docs, t)
	}

	_, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}
// createorder inserts a new order into the MongoDB collection.
func (a *adapter) CreateOrder(order *models.Order) error {
	collection := a.db.Collection("orders")
	o, err := collection.InsertOne(context.TODO(), order)
	if err != nil {
		return err
	}
	order.ID = o.InsertedID.(primitive.ObjectID)
	return nil
}

// GetAllAddresses retrieves all wallet addresses from the MongoDB collection.
func (a *adapter) GetAllAddresses() ([]string, error) {
	collection := a.db.Collection("wallets")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var addresses []string
	for cursor.Next(context.TODO()) {
		var wallet models.Wallet
		err := cursor.Decode(&wallet)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, wallet.Address)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (a *adapter)GetTotalRefferalsEarning(refID string) (float64, error) {
	matchStage := bson.D{{"$match", bson.D{{"to", refID}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", nil}, {"total", bson.D{{"$sum", "$amount"}}}}}}
	collection := a.db.Collection("transactions")
	ctx := context.Background()
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("failed to decode result: %w", err)
		}
		return float64(result.Total)/1000000, nil
	}

	return 0, nil
}