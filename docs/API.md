# API

This document lists the available endpoints and how to call them.

API Versioning is defined directly in the URLs. 

All unsuccessful requests return the following body:

```json
{
  "error": string
}
```

## FizzBuzz

URL: '/api/v1/fizzbuzz'
Body: none
Headers: none
Query (all required):

- int1: int (non-zero value)
- int2: int (non-zero value)
- limit: int
- str1: string
- str2: string

Successful response:

Status Code: 200 (ok)
Body:

```json
{
  "values": [
    string
  ],
  "total": int
}
```

Example:

```shell
curl "http://localhost:8080/api/v1/fizzbuzz?int1=3&int2=5&limit=10&str1=fizz&str2=buzz"
```

Other Status Codes: 400, 500

## Stats

URL: '/api/v1/stats'
Body: none
Headers: none
Query: none

Successful response:

Status Code: 200 (ok)
Body:

```json
{
  "most_frequent": string
}
```

Example:

```shell
curl "http://localhost:8080/api/v1/stats"
```

Other Status Codes: 404, 500
