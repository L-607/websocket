package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Load client cert
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println(err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	dialer := websocket.Dialer{
		TLSClientConfig: tlsConfig,
	}

	// Connect to wss server
	conn, _, err := dialer.Dial("wss://emqx.ee:8443/wss", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		now := time.Now().Format("2006-01-02 15:04:05")
		message := now + " Hello, server!"

		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println(err)
			return
		}
		// Print message received from server
		fmt.Printf("Send message: %s\n", string(message))

		// Read message from server
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// // Print message received from server
		fmt.Printf("Received message: %s\n", string(p))

		time.Sleep(5 * time.Second)
	}
}
