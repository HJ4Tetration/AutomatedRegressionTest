package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var c *http.Client

type myClient struct {
	d  websocket.Dialer
	ip string
}

func getGatewayIP() string {
	return ""
}

func getContent() string {
	return ""
}

func getWebsocketIP() string {
	return ""
}

func getRegistrationInfo(swn int) SwitchRegistration {
	var regist SwitchRegistration
	return regist
}

func singleSwitchRegistration(swn int, wg *sync.WaitGroup) bool { //true:registration succeds,false:registration fails
	defer wg.Done()
	url := getGatewayIP()
	content := getContent()
	regist := getRegistrationInfo(swn)

	jsonRegist, err := json.Marshal(regist)
	if err != nil {
		glog.Errorf("Can't marshal switchRegistration (switch %d)\n", swn)
		return false
	}

	resp, err := c.Post(url, content, bytes.NewReader(jsonRegist))
	if err != nil {
		glog.Errorf("Error getting response for registration request (switch %d)\n", swn)
		return false
	}
	defer resp.Body.Close()

	jsonRespb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Error reading from registration response body (switch %d)\n", swn)
		return false
	}
	var registResponse SwitchRegistrationResponse
	err = json.Unmarshal(jsonRespb, &registResponse)
	if err != nil {
		glog.Errorf("Can't unmarshal switchRegistrationResponse (switch %d)\n", swn)
		return false
	}

	/*if registResponse. ==  {
		glog.Info("Registration failed before building websocket")
		return
	}*/
	glog.Infof("Switch registration (http request) done, ready to build websocket and send config information\n")

	mc := myClient{ip: getWebsocketIP()}
	proxyConnection, _, err := mc.d.Dial(mc.ip, nil)
	if err != nil {
		glog.Errorf("Error making websocket connection to proxy (switch %d)\n", swn)
		return false
	}
	glog.Infof("Switch websocket with proxy established (switch %d)\n", swn)

	rc := make(chan ServerMessage)
	sc := make(chan SwitchMessage)

	//Pump from proxy to Switch
	go func(chan<- ServerMessage) {
		for {
			messageType, message, err := proxyConnection.ReadMessage()
			if messageType == websocket.CloseMessage {
				glog.Infof("Switch received websocket close information from proxy, websocket closed (switch %d)\n", swn)
				defer proxyConnection.Close()
				return
			}
			if err != nil {
				glog.Errorf("Error reading websocket message from proxy (switch %d)\n", swn)
				return
			}

			var jsonMessage ServerMessage
			err = json.Unmarshal(message, &jsonMessage)
			if err != nil {
				glog.Errorf("Cannot unmarshal ServerMessage (switch %d)\n", swn)
				return
			}
			rc <- jsonMessage
		}
	}(rc)

	//Pump from Switch to proxy
	go func(<-chan SwitchMessage) {
		for {
			message := <-sc
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				glog.Errorf("Cannot marshal SwitchMessage (switch %d)\n", swn)
				return
			}

			messageType := 2
			proxyConnection.WriteMessage(messageType, jsonMessage)
			glog.Infof("Write message from switch to proxy (switch %d)\n", swn)
		}
	}(sc)

	go func(<-chan ServerMessage, chan<- SwitchMessage) {
		for {
			//messageR := <-rc

			//do something based on command in messageR

			//var messageS SwitchMessage
			//sc <- messageS
		}
	}(rc, sc)
	return true
}
