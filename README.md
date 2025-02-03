# Receipt Processor API

A simple receipt processing service built in Go using `gorilla/mux`. It processes receipts and calculates reward points based on predefined rules.

## Features
- **Submit receipts** (`POST /receipts/process`)
- **Retrieve receipt points** (`GET /receipts/{id}/points`)
- **Optimized with `sync.RWMutex`** for concurrency

## API Endpoints

### Submit a Receipt
- **POST** `/receipts/process`
- **Request:**
  ```json
  {
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      }
    ],
    "total": "35.35"
  }
  ```
- **Response:**
  ```json
  {
    "id": "e8d90ccc-1102-4e6c-af2f-56206909aeae"
  }
  ```

### Get Receipt Points
- **GET** `/receipts/{id}/points`
- **Response:**
  ```json
  {
    "points": 12
  }
  ```
