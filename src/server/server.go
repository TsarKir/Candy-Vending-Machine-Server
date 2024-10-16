package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"ex01/types"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/buy_candy", BuyHandler)
	cert, err := tls.LoadX509KeyPair("../certs/server.candy.tld/cert.pem", "../certs/server.candy.tld/key.pem")
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    loadCA(),
		ServerName:   "server.candy.tld",
	}

	server := &http.Server{
		Addr:      ":3333",
		Handler:   http.HandlerFunc(BuyHandler),
		TLSConfig: config,
	}

	fmt.Println("Starting server on https://server.candy.tld:3333")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Printf("Server launch error: %s\n", err)
	}
}

func loadCA() *x509.CertPool {
	caCert, err := os.ReadFile("../certs/minica.pem")
	if err != nil {
		panic(err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	return caPool
}

func BuyHandler(w http.ResponseWriter, r *http.Request) {
	Prices := map[string]int{
		"AA": 15,
		"CE": 10,
		"NT": 17,
		"DE": 21,
		"YR": 23,
	}
	w.Header().Set("Content-Type", "application/json")
	var req types.Request
	var resp types.Response
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = "Decoding error"
		json.NewEncoder(w).Encode(resp)
		return
	}
	candyPrice, exists := Prices[req.CandyType]
	if !exists {
		resp.Error = fmt.Sprintf("There are no candy like %s", req.CandyType)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
	} else {
		if req.CandyCount*candyPrice <= req.Money {
			resp.Change = req.Money - req.CandyCount*candyPrice
			resp.Thanks = "Thank you!"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp.Error = fmt.Sprintf("You need %d more money!", req.CandyCount*candyPrice-req.Money)
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(resp)

		}
	}
}
