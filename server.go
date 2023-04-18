package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Load server certificate and private key
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load client CA certificate
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		fmt.Println(err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Set up TLS configuration with client CA certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	// Start HTTP server for ws
	go func() {
		http.HandleFunc("/ws", handleWebSocket)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Start HTTPS server for wss
	go func() {
		http.HandleFunc("/wss", handleWebSocket)
		server := &http.Server{
			Addr:      ":8443",
			Handler:   nil,
			TLSConfig: tlsConfig,
		}
		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			fmt.Println(err)
		}
	}()

	select {}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// read client send msg
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received message: %s\n", string(message))

		// send msg to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
