package gonet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// make get request with a body
func GetWithBody(url string, defaultTimeoutInSeconds int, jsonReqBodyStr string, headers ...map[string]string) (string, error) {
	//defaultTimeout := 30 // time out request if no response
	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * time.Duration(defaultTimeoutInSeconds), // Timeout after 5 seconds
	}
	// getJson, err := json.Marshal(reqBody)
	// if err != nil {
	// 	return "", err
	// }
	// payload := strings.NewReader(string(getJson))
	getJson, err := json.Marshal(jsonReqBodyStr)
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(getJson))

	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return "", err
	}

	if len(headers) > 0 {

		for k, v := range headers[0] {

			req.Header.Add(k, v)

		}

	}

	res, err := client.Do(req)
	if err != nil {

		return "", err

	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {

		return "", err

	}

	return string(body), err

}

// make multiple duplicated get request to a server
func MultipleGetWithBody(noOfReuests uint, url string, defaultTimeout int, jsonReqBodyStr string, headers map[string]string) (map[int]map[int]string, error) {

	requests := make([]uint, 0, noOfReuests)
	for noOfReuests > 0 {
		val := noOfReuests
		requests = append(requests, val)
		noOfReuests--
	}

	myChannel := make(chan map[int]string, noOfReuests)
	for _, v := range requests {

		go func(v uint) {

			getRes, err := GetWithBody(url, defaultTimeout, jsonReqBodyStr, headers)
			data := make(map[int]string, 0)
			if err != nil {
				data[int(v)] = err.Error()
				myChannel <- data
				//close(myChannel)
			}
			data[int(v)] = getRes
			myChannel <- data

		}(v)

	}
	resR := make(map[int]map[int]string, noOfReuests)
	defer close(myChannel)
	for _, v := range requests {
		//fmt.Println("##test", mp, <-myChannel)
		resR[int(v)] = <-myChannel

	}

	return resR, nil

}

// make multiple get request to different urls and unique data to send
func MultipleUrlsGetWithBodyUnique(numberOfConcurrentRequests uint, timeouts int, urlDataToSend map[string]string, headers map[string]string) (map[int]map[int]map[int]string, error) {
	myChannel := make(chan map[int]map[int]string)
	for url, datastr := range urlDataToSend {

		go func(url string, datastr string) {
			result, _ := MultipleGetWithBody(numberOfConcurrentRequests, url, timeouts, datastr, headers)

			myChannel <- result

		}(url, datastr)

	}
	//defer close(myChannel)
	resR := make(map[int]map[int]map[int]string, 0)
	counter := 0
	for mp := range urlDataToSend {
		counter++
		fmt.Println("##test-x", mp, "counter")
		resR[counter] = <-myChannel
		//resR = append(resR, mp)

	}

	return resR, nil
}
