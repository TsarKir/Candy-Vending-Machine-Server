# Technical Specification: Candy Vending Machine API

## Overview
This document outlines the requirements for implementing a backend server for a candy vending machine. The server will accept orders via HTTP and JSON and respond based on the validity of the input data and the amount of money provided.

## API Specification

### Swagger Definition
```yaml
swagger: '2.0'
info:
  version: 1.0.0
  title: Candy Server
paths:
  /buy_candy:
    post:
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: order
          description: summary of the candy order
          schema:
            type: object
            required:
              - money
              - candyType
              - candyCount
            properties:
              money:
                description: amount of money put into vending machine
                type: integer
              candyType:
                description: kind of candy (valid types: "CE", "AA", "NT", "DE", "YR")
                type: string
              candyCount:
                description: number of candy (must be a positive integer)
                type: integer
      operationId: buyCandy
      responses:
        201:
          description: purchase successful
          schema:
              type: object
              properties:
                thanks:
                  type: string
                change:
                  type: integer
        400:
          description: input data error (invalid candyType or negative candyCount)
          schema:
              type: object
              properties:
                error:
                  type: string
        402:
          description: not enough money provided
          schema:
              type: object
              properties:
                error:
                  type: string
```

# Business Logic

## Input Validation
- Validate `candyCount` to ensure it is a positive integer.
- Validate `candyType` to ensure it matches one of the following valid types: "CE", "AA", "NT", "DE", or "YR".

## Price Calculation
- Calculate the total price based on the `candyType` and `candyCount`.
- If `money` provided is less than the total price, respond with HTTP 402 and an error message indicating how much more is needed.
- If money is sufficient, respond with HTTP 201, including a thank you message and the change to be returned.

## Example Requests

### Valid Purchase Request
```bash
curl -XPOST -H "Content-Type: application/json" -d '{"money": 20, "candyType": "AA", "candyCount": 1}' http://127.0.0.1:3333/buy_candy 
# Response: {"change":5,"thanks":"Thank you!"}
```

## Insufficient Funds Request
```bash
curl -XPOST -H "Content-Type: application/json" -d '{"money": 46, "candyType": "YR", "candyCount": 2}' http://127.0.0.1:3333/buy_candy 
# Response: {"error":"You need {amount} more money!"}
```

# Security Requirements

## Mutual TLS Authentication

### Certificate Generation
- Use Minica to generate two certificate/key pairs—one for the server and one for the client. Ensure that a CA file (`minica.pem`) is created.

### Server Configuration
- Update the server to support TLS using secure URLs and accept command-line parameters for CA file, key file, and cert file.

### Client Configuration
- Implement a test client that supports the following flags:
  - `-k`: Two-letter abbreviation for the candy type.
  - `-c`: Count of candy to buy.
  - `-m`: Amount of money provided.

### Example Client Command
```bash
./candy-client -k AA -c 2 -m 50 
# Expected Output: Thank you! Your change is 20.
```
