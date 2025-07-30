# API

This sections lists the available endpoints and how to call them.
All unuseful requests return the following body:

```json
{
  "error": string
}
```

## FizzBuzz

URL: '/fizzbuzz'
Body: none
Headers: none
Query (all required):

- int1 (non-zero value)
- int2 (non-zero value)
- limit
- str1
- str2

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
curl "http://localhost:8080/fizzbuzz?int1=3&int2=5&limit=10&str1=fizz&str2=buzz"
```

Other Status Codes: 400, 500

## Stats

URL: '/stats'
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
curl "http://localhost:8080/stats"
```

Other Status Codes: 404, 500
