package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"ex01/types"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	candyType := flag.String("k", "", "Candy type (e.g., AA)")
	candyCount := flag.Int("c", 0, "Number of candies")
	money := flag.Int("m", 0, "Amount of money")
	flag.Parse()

	reqBody := types.Request{
		CandyType:  *candyType,
		CandyCount: *candyCount,
		Money:      *money,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}
	responseBodyBuff := bytes.NewBuffer(jsonData)

	client := CertChecker()

	resp, err := client.Post("https://server.candy.tld:3333/buy_candy", "application/json", responseBodyBuff)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var responseBody types.Response
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if responseBody.Error != "" {
		fmt.Println(responseBody.Error)
	} else {
		fmt.Printf("%s Your change is %d\n", responseBody.Thanks, responseBody.Change)
	}
}

func CertChecker() *http.Client {
	cert, err := tls.LoadX509KeyPair("../certs/client.candy.tld/cert.pem", "../certs/client.candy.tld/key.pem")
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCert, err := os.ReadFile("../certs/minica.pem")
	if err != nil {
		panic(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tr := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: tr}
}
