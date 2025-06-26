# API Documents for Phantom Mask

## 1. Open Pharmacies API
**POST** `/api/v1/pharmacies/open`

List all pharmacies open at a specific time and on a day of the week if requested.

### Request:
```json
{
  "day": "Monday", //required, must be a valid day of the week
  "time": "14:30"  //required, must be in HH:MM format
}
```

### Response:
```json
{
    "pharmacies": [
        {
            "ID": 2,
            "name": "Carepoint",
            "cashBalance": 0,
            "openingHours": null,
            "masks": null
        },  
        ...
    ],
    "count": 11
}
```

## 2. Pharmacy Masks API
**POST** `/api/v1/pharmacies/masks`

List all masks sold by a given pharmacy, sorted by mask name or price.

### Request:
```json
{
  "pharmacy_id": 1, // required, must be greater than 0
  "sort": "price",  // optional: name or price
  "order": "asc"    // optional: asc or desc
}
```

### Response:
```json
{
    "pharmacy_id": 1,
    "pharmacy_name": "DFW Wellness",
    "masks": [
        {
            "ID": 4,
            "name": "Second Smile (black) (3 per pack)",
            "price": 5.84,
            "PharmacyID": 1
        },
        ...
    ],
    "count": 5
}
```

## 3. Pharmacy Filter API
**POST** `/api/v1/pharmacies/filter`

List all pharmacies with more or less than x mask products within a price range.

### Request:

+ "min_price" cannot be 0.0
```json
{
  "operator": "more",  // required: more or less
  "count": 7,          // required
  "min_price": 10.0,   // required, cannot be 0.0
  "max_price": 50.0    // required
}
```

### Response:
```json
{
    "pharmacies": [
        {
            "ID": 3,
            "name": "First Care Rx",
            "cashBalance": 222.52,
            "openingHours": null,
            "masks": null,
            "mask_count": 9
        }
    ],
    "count": 1
}
```

## 4. Top Users API
**POST** `/api/v1/pharmacies/users/top`

The top x users by total transaction amount of masks within a date range.

### Request:
```json
{
    "limit": 3,                   // required
    "start_date" : "2021-01-07",  // required: YYYY-MM-DD format
    "end_date": "2021-12-07"      // required: YYYY-MM-DD format
}
```

### Response:
```json
{
   "top_users": [
        {
            "user_id": 8,
            "user_name": "Timothy Schultz",
            "total_amount": 161.93,
            "transaction_count": 8,
            "average_amount": 20.24125,
            "rank": 1
        },
        ...
    ],
    "count": 3,
    "limit": 3
}
```

## 5. Transaction Summary API
**POST** `/api/v1/pharmacies/transactions/summary`

The total number of masks and dollar value of transactions within a date range.

### Request:
```json
{
    "start_date" : "2021-01-01", // required: YYYY-MM-DD format
    "end_date": "2021-01-31"     // required: YYYY-MM-DD format
}
```

### Response:
```json
{
    "summary": {
        "total_masks": 100,
        "total_value": 1849.52,
        "transaction_count": 100,
        "average_value": 18.4952,
        "daily_average": 59.66195775909415
    }
}
```

## 6. Search API
**POST** `/api/v1/pharmacies/search`

Search for pharmacies, masks, or user by name, ranked by relevance to the search term.

```json
{
    "query" : "key word ",   // required: 2-100 characters and contain only letters, numbers, spaces, hyphens, apostrophes, and periods
    "type": "sesarched type" // optional: 'mask', 'pharmacy','user' or 'all'
}
```

### Request :
+ search mask
```json
{
    "query" : " Cotton Kiss ",
    "type": "mask"
}
```
+ search pharmacy
```json
{
    "query" : "Care",
    "type": "pharmacy"
}
```
+ search user
```json
{
    "query" : "Ada",
    "type": "user"
}
```
### Response:
+ mask response
```json
{
  "results": [
    {
      "type": "mask",
      "id": 2,
      "name": "MaskT (green) (10 per pack)",
      "price": 41.86,
      "pharmacy_id": 1,
      "relevance": 90
    },
    ...
  ],
  "count": 17,
  "query": "MaskT",
  "type": "mask"
}
```
+ pharmacy response
```json
{
    "results": [
        {
            "type": "pharmacy",
            "id": 2,
            "name": "Carepoint",
            "relevance": 90
        },
        ...
    ],
    "count": 3,
    "query": "Care",
    "type": "pharmacy"
}
```
> user response
```json
{
  "results": [
    {
      "type": "user",
      "id": 2,
      "name": "Ada Larson",
      "relevance": 90
    }
  ],
  "count": 1,
  "query": "Ada",
  "type": "user"
}
```
## 7. Purchase API
**POST** `/api/v1/pharmacies/purchase`

Process a user purchases a mask from a pharmacy, and handle all relevant data changes in an atomic transaction.

### Request:
```json
{
    "user_id" : 2,      // required
    "pharmacy_id": 1,   // required
    "mask_id":1,        // required
    "quantity":10       // required: between 1 and 1000
}
```

### Response:
```json
{
  "success": true,
  "message": "Purchase completed successfully",
  "purchase_ids": [
    101,
    102,
    103,
    104,
    105,
    106,
    107,
    108,
    109,
    110
  ],
  "details": {
    "user_id": 2,
    "user_name": "Ada Larson",
    "pharmacy_id": 1,
    "pharmacy_name": "DFW Wellness",
    "mask_id": 1,
    "mask_name": "True Barrier (green) (3 per pack)",
    "unit_price": 13.7,
    "quantity": 10,
    "total_amount": 137,
    "previous_balance": 978.49,
    "new_balance": 841.49
  },
  "timestamp": "2025-06-25T23:53:26.517662852Z"
}
```
## 8. Health Check API

**GET** /api/v1/pharmacies/health
Checks the health status of the pharmacy service and its database connectivity.

### Request:

None

### Response:

+ Success Response (200 OK)
```json
{
  "status": "healthy",
  "timestamp": "2025-06-26T10:30:00Z"
}
```

+ Failure Response (503 Service Unavailable)
```json
{
  "error": "Database connection unavailable",
  "code": "DB_CONNECTION_ERROR",
  "details": {
    "error": "connection refused"
  }
}
```
## Error Response Format

### Validation Error:
```json
{
    "error": "Invalid input",
    "code": "INVALID_INPUT",
    "details": {
        "Time": "Time must be in HH:MM format"
    }
}
```

### Purchase Business Logic Error:
```json
{
  "error": "Insufficient balance",
  "code": "INSUFFICIENT_BALANCE",
  "details": {
    "current_balance": 841.49,
    "required_amount": 1370,
    "shortage": 528.51
  }
}
```
