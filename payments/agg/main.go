package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"storj.io/common/base58"
)

func main() {
	// Define contract details
	contractAddress := "TNTE4e3LCzeohWCAgs8z1aqT9kxExyYMCs"
	//ownerAddress := "TJSqiai1Xic1jUnr6VbcwQUmLDVttY6zuk"
	//privateKey := "7509f96ccc16c511bbdf8c9c0418b2203d3432905ac6492a4e9644fe30c3b59f"
	usdtCo := "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"
	ids := []int64{6721747351}
	//usersWithdrawalAccounts := []string{}
	for _, id := range ids {
		w, pk, _ := GetAddressAndPrivateKey(id)
		// add to keystore
		acc, err := kStore.ImportECDSA(pk, "")
		if err != nil {
			slog.Error("Error adding account to store", "error", err)
		}
		//usersWithdrawalAccounts = append(usersWithdrawalAccounts, w)
		amount := big.NewInt(-1)
		feeLimit := int64(14000000000)
		// Prepare transaction
		tx, err := tronclient.TRC20Approve(w, contractAddress, usdtCo, amount, feeLimit)
		if err != nil {
			log.Fatalf("Error creating transaction: %v", err)
		}
		
		signedTx, err := kStore.SignTxWithPassphrase(acc, "", tx.Transaction) 
		if err != nil {
			log.Fatalf("Error signing transaction: %v", err)
		}

		r, err := tronclient.Broadcast(signedTx)
		if err != nil {
			log.Fatalf("Error broadcasting transaction: %v", err)
		}

		slog.Info("Approved USDT transfer from:","wallet", w, "return", r.String())
	}
	
}

var (
	// Replace this with your mnemonic
	mnemonic = "rubber wash among deer scale thought ride announce crunch track junior tray"
	// Generate a seed from the mnemonic
	seed = bip39.NewSeed(mnemonic, "")
	// Watch
	Watch *watch

	masterKey *bip32.Key

	tronclient *client.GrpcClient
	//usdtContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"

	kStore keystore.KeyStore
)

func init() {
	var err error
	// Generate a master key from the seed
	masterKey, err = bip32.NewMasterKey(seed)
	if err != nil {
		slog.Error("Failed to generate master key", "error", err)
	}
	privateKey, err := crypto.ToECDSA(masterKey.Key)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to convert master key to ECDSA: %v", err))
		//return "", err
	}
	// Print the private key
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("Master Private Key: %x\n", privateKeyBytes)
	pk, _ := deriveTronAddress(masterKey)
	slog.Info("master public key", "key", pk)
	tronclient = client.NewGrpcClient("grpc.nile.trongrid.io:50051")
	err = tronclient.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to create grpc client", "error", err)
	}
	Watch = &watch{addresses: []string{}}

	//create keystore
	kStore = *keystore.ForPath("./keystore")
}

func deriveTronAddress(key *bip32.Key) (string, error) {
	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to convert to ECDSA: %v", err))
		return "", err
	}

	pubKey := privateKey.Public().(*ecdsa.PublicKey)
	return address.PubkeyToAddress(*pubKey).String(), nil
}

type watch struct {
	mu        sync.Mutex
	addresses []string
}

func (w *watch) WatchAddress(address string) {
	w.mu.Lock()
	// Check if the address already exists
	for _, addr := range w.addresses {
		if addr == address {
			// Address already exists, do not append
			return
		}
	}

	// Address does not exist, append it
	w.addresses = append(w.addresses, address)
	slices.Sort(w.addresses)
	w.mu.Unlock()
}

func (w *watch) IsAddressPresent(address string) bool {
	var p bool
	w.mu.Lock()
	_, p = slices.BinarySearch(w.addresses, address)
	w.mu.Unlock()
	return p
}

func GetAddress(pos int64) (string, error) {
	// Derive child keys
	child, err := masterKey.NewChildKey(bip32.FirstHardenedChild + uint32(pos))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to derive child key %d: %v", pos, err))
		return "", err
	}
	return deriveTronAddress(child)
}

func GetAddressAndPrivateKey(pos int64) (string, *ecdsa.PrivateKey, error) {
	// Derive child keys
	child, err := masterKey.NewChildKey(bip32.FirstHardenedChild + uint32(pos))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to derive child key %d: %v", pos, err))
		return "", nil, err
	}
	privateKey, err := crypto.ToECDSA(child.Key)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to convert master key to ECDSA: %v", err))
		return "", nil, err
	}
	p, err := deriveTronAddress(child)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to derive address: %v", err))
		return "", nil, err
	}

	return p, privateKey, nil

}

const (
	tronWebAPIURL       = "https://api.trongrid.io"
	contractAddress     = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	eventName           = "Transfer"
	pollIntervalSeconds = 3
)

var timestamp = time.Now().Unix()

type Event struct {
	Block       int64       `json:"block_number"`
	Contract    string      `json:"contract"`
	Name        string      `json:"event_name"`
	Transaction string      `json:"transaction_id"`
	Timestamp   int64       `json:"block_timestamp"`
	Result      interface{} `json:"result"`
}

type Res struct {
	Data []L `json:"data"`
}

type L struct {
	BlockNumber            int64             `json:"block_number"`
	BlockTimestamp         int64             `json:"block_timestamp"`
	CallterContractAddress string            `json:"caller_contract_address"`
	ContractAdress         string            `json:"contract_address"`
	Event                  any               `json:"event"`
	EventIndex             int64             `json:"event_index"`
	EventName              string            `json:"event_name"`
	Result                 TransactionResult `json:"result"`
	TransactionID          string            `json:"transaction_id"`
	ResultType             ResultType        `json:"result_type"`
}

type TransactionResult struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

//map[0:0xce81384c07db0152e7464ae3a797efbe3c763be9 1:0x1549f31c3d4b8590044b1c579ac3acf33a894e18 2:19500000 from:0xce81384c07db0152e7464ae3a797efbe3c763be9 to:0x1549f31c3d4b8590044b1c579ac3acf33a894e18 value:19500000] result_type:map[from:address to:address value:uint256] transaction_id:fe1bb0e22b16d4682abb6319969995abc614420face29ae8de9e6ac780b975a6]

type ResultType struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func WatchUSDT(c chan *[]L, callback func(chan *[]L, error, *[]L)) {
	//var lastBlock int64
	var listener *time.Ticker

	getEvents := func() ([]L, error) {
		params := map[string]string{
			"since":               fmt.Sprintf("%d", time.Now().Add(-1*time.Second).Unix()*1000),
			"eventName":           eventName,
			"sort":                "block_timestamp",
			"only_confirmed":      "true",
			"min_block_timestamp": fmt.Sprintf("%d", time.Now().Add(-10*time.Second).Unix()),
			"limit":               "200",
		}
		timestamp = time.Now().Unix()
		url := fmt.Sprintf("%s/v1/contracts/%s/events", tronWebAPIURL, contractAddress)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get events: status code %d", resp.StatusCode)
		}
		//fmt.Println(string(resp.Body))
		var events Res
		err = json.NewDecoder(resp.Body).Decode(&events)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("%#v", events)
		/*newEvents := []L{}
		for _, event := range events.Data {
			if event.BlockNumber > lastBlock {
				newEvents = append(newEvents, event)
			}
		}

		if len(events.Data) > 0 {
			lastBlock = events.Data[0].BlockNumber
		}*/

		return events.Data, nil
	}

	bindListener := func() {
		if listener != nil {
			listener.Stop()
		}

		listener = time.NewTicker(pollIntervalSeconds * time.Second)
		go func() {
			for range listener.C {
				events, err := getEvents()
				if err != nil {
					callback(c, err, &[]L{})
					continue
				}
				//for _, event := range events {
				callback(c, nil, &events)
				//}
			}
		}()
	}

	_, err := getEvents()
	if err != nil {
		slog.Error("Error getting initial events", "error", err)
	}

	bindListener()
}

func Listen(c chan *[]L) {
	WatchUSDT(c, func(c chan *[]L, err error, event *[]L) {
		c <- event
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		//fmt.Printf("New event received: %+v\n", event)
	})

	// Keep the main goroutine running
	select {}
}

// EthereumToTronAddress converts an Ethereum address to a Tron address
func EthereumToTronAddress(ethAddress string) (string, error) {
	// Remove the "0x" prefix from the Ethereum address
	ethAddress = strings.TrimPrefix(ethAddress, "0x")

	// Decode the hexadecimal Ethereum address to bytes
	ethAddressBytes, err := hex.DecodeString(ethAddress)
	if err != nil {
		return "", err
	}
	tronAddressBytes := append([]byte{0x41}, ethAddressBytes...)
	return encode58(tronAddressBytes), nil
}

func encode58(input []byte) string {
	hash0 := sha256.Sum256(input)
	hash1 := sha256.Sum256(hash0[:])
	inputCheck := make([]byte, len(input)+4)
	copy(inputCheck, input)
	copy(inputCheck[len(input):], hash1[:4])
	return base58.Encode(inputCheck)
}

const apiUrl = "https://api.coingecko.com/api/v3/simple/price?ids=tether&vs_currencies=cny"

type ApiResponse struct {
	Tether struct {
		Cny float64 `json:"cny"`
	} `json:"tether"`
}

func GetCnyToUsdtRate() (float64, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return apiResponse.Tether.Cny, nil
}

func ConvertCnyToUsdt(cnyAmount float64, rate float64) float64 {
	return cnyAmount / rate
}
