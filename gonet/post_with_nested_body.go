package gonet

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func PostWithNestedBody(url string, defaultTimeout int, jsonReqBodyStr string, headers ...map[string]string) (string, error) {
	//defaultTimeout := 30 // time out request if no response
	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * time.Duration(defaultTimeout), // Timeout after 5 seconds
	}

	// getJson, err := json.Marshal(data)
	// if err != nil {
	// 	return "", err
	// }
	// payload := strings.NewReader(string(getJson))

	payload := strings.NewReader(jsonReqBodyStr)

	req, err := http.NewRequest("POST", url, payload)
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

func MultiplePostWithNestedBody(noOfReuests uint, url string, defaultTimeout int, jsonReqBodyStr string, headers map[string]string) (map[int]map[int]string, error) {

	requests := make([]uint, 0, noOfReuests)
	for noOfReuests > 0 {
		val := noOfReuests
		requests = append(requests, val)
		noOfReuests--
	}

	myChannel := make(chan map[int]string, noOfReuests)
	for _, v := range requests {

		go func(v uint) {

			getRes, err := PostWithNestedBody(url, defaultTimeout, jsonReqBodyStr, headers)
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
	defer close(myChannel)
	resR := make(map[int]map[int]string, noOfReuests)

	for _, v := range requests {
		//fmt.Println("##test", mp, <-myChannel)
		resR[int(v)] = <-myChannel

	}
	//fmt.Println("##test-########THIS", resR)

	return resR, nil

}

func MultiplePostWithNestedBodyUnique(numberOfConcurrentRequests uint, url string, timeouts int, dataToSend []string, headers map[string]string) (map[int]map[int]map[int]string, error) {
	myChannel := make(chan map[int]map[int]string)
	for k, dataMap := range dataToSend {

		go func(dataMap string, k int) {
			result, _ := MultiplePostWithNestedBody(numberOfConcurrentRequests, url, timeouts, dataMap, headers)

			myChannel <- result

		}(dataMap, k)

	}
	//defer close(myChannel)
	resR := make(map[int]map[int]map[int]string, 0)
	counter := 0
	for mp := range dataToSend {
		counter++
		fmt.Println("##test-x", mp, "counter")
		resR[counter] = <-myChannel
		//resR = append(resR, mp)

	}

	return resR, nil
}
