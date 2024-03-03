package ndaxClient

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type MsgFrame struct {
	M int    `json:"m"`
	I int    `json:"i"`
	N string `json:"n"`
	O any    `json:"o"`
}

func sendMessage(sendChannel chan MsgFrame, methodName string, payload map[string]string) error {
	fmt.Printf("Sending Message: %s, %+v\n", methodName, payload)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	frame := MsgFrame{
		M: 0,
		I: 2,
		N: methodName,
		O: jsonData,
	}

	sendChannel <- frame
	return nil
}

func Start(apiKey string, secretKey string) {

	// Connect to NDAX WS
	u := url.URL{Scheme: "wss", Host: "api.ndax.io", Path: "/WSGateway"}
	fmt.Printf("Connecting to %s\n", u.String())
	socketContext, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to NDAX:", err)
		return
	}
	fmt.Println("Connected to NDAX")

	// Read and Write channels for NDAX WS
	sendChannel := make(chan MsgFrame)
	receiveChannel := make(chan MsgFrame)

	// Read from NDAX WS
	go func() {
		defer close(receiveChannel)
		for {
			_, message, err := socketContext.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message from NDAX:", err)
				return
			}

			recievedMsg := MsgFrame{}
			err = json.Unmarshal(message, &recievedMsg)
			if err != nil {
				fmt.Println("Error unmarshalling message from NDAX:", err)
				return
			}
			receiveChannel <- recievedMsg
		}
	}()

	// Temp: Go Routine to print messages from NDAX WS
	go func() {
		receivedMessage := <-receiveChannel
		fmt.Printf("Received Message:%+v", receivedMessage)
	}()

	// Write to NDAX WS
	go func() {
		defer close(sendChannel)
		for sendMessage := range sendChannel {
			data, err := json.Marshal(sendMessage)
			if err != nil {
				fmt.Println("Error marshalling message to NDAX:", err)
				continue
			}

			err = socketContext.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("Error writing message to NDAX:", err)
				return
			}
		}
	}()

	// Authenticate
	if apiKey == "" || secretKey == "" {
		fmt.Println("No API key or secret key provided")
		return
	}
	payload := map[string]string{
		"APIKey":    apiKey,
		"SecretKey": secretKey,
	}
	sendMessage(sendChannel, "authenticateuser", payload)
}
