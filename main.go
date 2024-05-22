package main

import (
	"encoding/json"
	"fmt"

	"github.com/emcodest/emcode-go-net/gonet"
)

func main() {

	//!++++++++++++++++++++++++++++++++++++++++++++
	// | make a get request with nested body
	//++++++++++++++++++++++++++++++++++++++++++++
	//! sample - create a struct to match the sample
	// {
	// 	"method": "create_account",
	// 	"params": {
	// 		"name": "test"
	// 	},
	// 	"jsonrpc": "2.0",
	// 	"id": 1
	// }
	url := "https://webhook.site/f884a569-4ff9-4ca2-b25a-bc7933545315"
	type MyParams struct {
		Name string `json:"name"`
	}
	body := struct {
		Method  string   `json:"method"`
		Params  MyParams `json:"params"`
		JsonRpc string   `json:"jsonrpc"`
		Id      uint     `json:"id"`
	}{
		Method:  "create_account",
		Params:  MyParams{Name: "test"},
		JsonRpc: "2.0",
		Id:      1,
	}

	jsonStr, _ := json.Marshal(body)
	jsonStrBody := string(jsonStr)
	fmt.Println("baddest ", jsonStrBody)

	header := map[string]string{
		"content-type": "application/json",
	}

	res, _ := gonet.GetWithBody(url, 10, jsonStrBody, header)
	fmt.Println("##RESULT", res)

}
