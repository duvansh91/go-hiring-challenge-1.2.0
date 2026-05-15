# API Documentation

Base URL: `http://localhost:{HTTP_PORT}`

All responses use `Content-Type: application/json`.

---

## Catalog

### List products

**`GET /catalog`**

Returns a paginated list of products. Supports optional filtering by category and price.

#### Query parameters

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `offset` | integer | No | `0` | Number of records to skip |
| `limit` | integer | No | `10` | Number of records to return. Min `1`, max `100` |
| `category_code` | string | No | — | Filter by category code (e.g. `CLOTHING`) |
| `price_less_than` | float | No | — | Filter products with price strictly below this value |

#### Success response `200 OK`

```json
{
  "products": [
    {
      "code": "PROD001",
      "price": 59.99,
      "category": {
        "code": "clothing",
        "name": "Clothing"
      }
    }
  ],
  "total_products": 8,
  "total_pages": 1
}
```

#### Error responses

| Status | Condition | Body |
|---|---|---|
| `400 Bad Request` | `offset` or `limit` is not a valid integer | `{"error": "Invalid offset param"}` / `{"error": "Invalid limit param"}` |
| `400 Bad Request` | `offset` is negative | `{"error": "Offset must be a non-negative integer"}` |
| `400 Bad Request` | `limit` is outside the 1–100 range | `{"error": "Limit must be an integer between 1 and 100"}` |
| `400 Bad Request` | `price_less_than` is not a valid number | `{"error": "Invalid price_less_than param"}` |
| `500 Internal Server Error` | Database error | `{"error": "Error getting products"}` |

---

### Get product detail

**`GET /catalog/{code}`**

Returns the full details of a single product, including its variants.

#### Path parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `code` | string | Yes | The unique product code (e.g. `PROD001`) |

#### Success response `200 OK`

```json
{
  "code": "PROD001",
  "price": 59.99,
  "category": {
    "code": "clothing",
    "name": "Clothing"
  },
  "variants": [
    {
      "product_id": 1,
      "name": "Small",
      "sku": "PROD001-S",
      "price": 59.99
    },
    {
      "product_id": 1,
      "name": "Large",
      "sku": "PROD001-L",
      "price": 79.99
    }
  ]
}
```

#### Error responses

| Status | Condition | Body |
|---|---|---|
| `404 Not Found` | No product matches the given code | `{"error": "Product not found"}` |
| `500 Internal Server Error` | Database error | `{"error": "Error getting product"}` |

---

## Categories

### List categories

**`GET /categories`**

Returns a list of all available product categories.

#### Success response `200 OK`

```json
[
  {
    "code": "CLOTHING",
    "name": "Clothing"
  },
  {
    "code": "SHOES",
    "name": "Shoes"
  },
  {
    "code": "ACCESSORIES",
    "name": "Accessories"
  }
]
```

#### Error responses

| Status | Condition | Body |
|---|---|---|
| `500 Internal Server Error` | Database error | `{"error": "Error getting categories"}` |

---

### Create category

**`POST /categories`**

Creates a new product category.

#### Request body

```json
{
  "code": "CLOTHING",
  "name": "Clothing"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `code` | string | Yes | Unique human-readable identifier for the category |
| `name` | string | Yes | Display name of the category |

#### Success response `200 OK`

Returns the created category.

```json
{
  "code": "CLOTHING",
  "name": "Clothing"
}
```

#### Error responses

| Status | Condition | Body |
|---|---|---|
| `400 Bad Request` | Request body is not valid JSON | `{"error": "Invalid request body"}` |
| `400 Bad Request` | `code` field is missing or empty | `{"error": "Category code is required"}` |
| `400 Bad Request` | `name` field is missing or empty | `{"error": "Category name is required"}` |
| `500 Internal Server Error` | Database error | `{"error": "Error creating category"}` |
