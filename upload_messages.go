package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	//"log"
	"net/http"

	"github.com/nats-io/go-nats"
)

type Status struct {
	Timestamp    int  `json:"timestamp"`
	AirplaneMode bool `json:"is_airplane_mode"`
	Telephony    Telephony
	Telephonies  []Telephony
}

type Telephony struct {
	NetworkRoaming      bool   `json:"is_network_roaming"`
	SimState            string `json:"sim_state"`
	NetworkOperatorName string `json:"network_operator_name"`
	DisplayName         string `json:"display_name"`
	SimSlot             int    `json:"sim_slot"`
}

func main() {

	//mesgs := []byte("Hello World")

	fmt.Print("Enter IP Address: ")
	var ipaddr string
	fmt.Scanln(&ipaddr)
	fmt.Print("Enter Password: ")
	var token string
	fmt.Scanln(&token)
	//fmt.Print(input)

	if CheckStatus(ipaddr) != "ready" {
		fmt.Println("Mobile Device: Offline")
	} else {
		fmt.Println("Mobile Device: ready")
		mesgs, err := GetMessages(ipaddr, "1000")
		if err != nil {
			fmt.Println("Error retrieving messages.")
			fmt.Println("Upload FAILED")
		} else {
			fmt.Println("Messages retrieved")
			UploadMessages(token, mesgs)
		}
		//fmt.Println("Messages Uploaded")
	}
}

func GetMessages(ipaddr string, limit string) ([]byte, error) {
	res, err := http.Get("http://" + ipaddr + ":8080/v1/sms/?limit=" + limit)
	if err != nil {
		//log.Fatal(err)
		return []byte("error"), err
	}

	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Fatal(err)
		return []byte(response), err
	}

	return []byte(response), err

}

func UploadMessages(token string, messages []byte) {
	fmt.Println("Uploading Messages")

	natsConnection, err := nats.Connect("nats://" + token + "@localhost:4222")
	if err != nil {
		fmt.Println("Unable to connect to NATS Server")
		fmt.Println("Upload FAILED")
	} else {
		defer natsConnection.Close()
		fmt.Println("Connected to NATS server: " + nats.DefaultURL)

		// Msg structure
		msg := &nats.Msg{
			Subject: "foo",
			Reply:   "bar",
			Data:    messages,
		}
		natsConnection.PublishMsg(msg)
		fmt.Println("Messages Uploaded Successfully")
	}
	//log.Println("Published msg.Subject = "+msg.Subject, "| msg.Data = "+string(msg.Data))
}

func CheckStatus(ip_address string) string {

	res, err := http.Get("http://" + ip_address + ":8080/v1/device/status")
	if err != nil {
		//log.Fatal(err)
		return "Offline"
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Fatal(err)
		return "Offline"
	}

	status := Status{}
	json.Unmarshal([]byte(response), &status)

	if status.Telephonies[0].SimState == "ready" {
		return status.Telephonies[0].SimState
	} else {
		return "Offline"
	}

	return "Offline"
}
