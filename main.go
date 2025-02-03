package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt & Item Structs
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Store receipts in memory (thread-safe)
var receipts = make(map[string]Receipt)
var mutex = &sync.RWMutex{}

// Process Receipt
func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	id := uuid.New().String()
	mutex.Lock()
	receipts[id] = receipt
	mutex.Unlock()

	log.Printf("Stored receipt with ID: %s", id)

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get Points
func getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mutex.RLock()
	receipt, exists := receipts[id]
	mutex.RUnlock()

	if !exists {
		http.Error(w, "No Receipt Found", http.StatusNotFound)
		return
	}

	points := calculatePoints(receipt)
	log.Printf("Returning points for receipt ID: %s, Points: %d", id, points)

	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Calculate Points
func calculatePoints(receipt Receipt) int {
	points := countAlphanumeric(receipt.Retailer)

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if strings.HasSuffix(receipt.Total, ".00") {
			points += 50
		}
		if math.Mod(total, 0.25) == 0 {
			points += 25
		}
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	day, _ := strconv.Atoi(strings.Split(receipt.PurchaseDate, "-")[2])
	if day%2 == 1 {
		points += 6
	}

	timeParts := strings.Split(receipt.PurchaseTime, ":")
	hour, _ := strconv.Atoi(timeParts[0])
	if hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}

// Count Alphanumeric Chars
func countAlphanumeric(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	log.Println("Server running on Port 8080...")
	http.ListenAndServe(":8080", router)
}
