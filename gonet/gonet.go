package gonet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

//++++++++++++++++++++++++++++++++++++++++++++
// | 	 make get JSON request
//++++++++++++++++++++++++++++++++++++++++++++

// make get request
func GET(url string, defaultTimeout int, headers ...map[string]string) (string, error) {
	//defaultTimeout := 30 // time out request if no response
	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * time.Duration(defaultTimeout), // Timeout after 5 seconds
	}

	req, err := http.NewRequest("GET", url, nil)
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

// make multiple duplicated request to a server
func MakeMultipleGET(noOfReuests uint, url string, defaultTimeout int, headers map[string]string) (map[int]map[int]string, error) {

	requests := make([]uint, 0, noOfReuests)
	for noOfReuests > 0 {
		val := noOfReuests
		requests = append(requests, val)
		noOfReuests--
	}

	myChannel := make(chan map[int]string, noOfReuests)
	for _, v := range requests {

		go func(v uint) {

			getRes, err := GET(url, defaultTimeout, headers)
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

// ++++++++++++++++++++++++++++++++++++++++++++
// | 	 make post request
// ++++++++++++++++++++++++++++++++++++++++++++
// make post request - json body
func POST(url string, defaultTimeout int, data interface{}, headers ...map[string]string) (string, error) {
	//defaultTimeout := 30 // time out request if no response
	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * time.Duration(defaultTimeout), // Timeout after 5 seconds
	}

	getJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	payload := strings.NewReader(string(getJson))

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

// ++++++++++++++++++++++++++++++++++++++++++++
// | 	 make multiple post request
// ++++++++++++++++++++++++++++++++++++++++++++
// make multiple duplicated request to a server
func MakeMultiplePOST(noOfReuests uint, url string, defaultTimeout int, data interface{}, headers map[string]string) (map[int]map[int]string, error) {

	requests := make([]uint, 0, noOfReuests)
	for noOfReuests > 0 {
		val := noOfReuests
		requests = append(requests, val)
		noOfReuests--
	}

	myChannel := make(chan map[int]string, noOfReuests)
	for _, v := range requests {

		go func(v uint) {

			getRes, err := POST(url, defaultTimeout, data, headers)
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

// ++++++++++++++++++++++++++++++++++++++++++++
// | make post request - form data
// ++++++++++++++++++++++++++++++++++++++++++++
// make post request - form data
func PostFormData(url string, defaultTimeout int, formData map[string]string, headers ...map[string]string) (string, error) {
	// Encode form data
	var encodedFormData []string
	for key, value := range formData {
		encodedFormData = append(encodedFormData, key+"="+value)
	}
	formDataString := strings.Join(encodedFormData, "&")
	//defaultTimeout := 30 // time out request if no response
	// Create an HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * time.Duration(defaultTimeout), // Timeout after 5 seconds
	}

	payload := strings.NewReader(formDataString)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}

	if len(headers) > 0 {

		for k, v := range headers[0] {

			req.Header.Add(k, v)

		}

	}
	// Set the Content-Type header to application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

//++++++++++++++++++++++++++++++++++++++++++++
// | multiple post form data request
//++++++++++++++++++++++++++++++++++++++++++++

// make multiple duplicated request to a server
func MakeMultiplePostFormData(noOfReuests uint, url string, defaultTimeout int, data map[string]string, headers map[string]string) (map[int]map[int]string, error) {

	requests := make([]uint, 0, noOfReuests)
	for noOfReuests > 0 {
		val := noOfReuests
		requests = append(requests, val)
		noOfReuests--
	}

	myChannel := make(chan map[int]string, noOfReuests)

	for _, v := range requests {

		go func(v uint) {

			getRes, err := PostFormData(url, defaultTimeout, data, headers)
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
	//defer close(myChannel)
	for _, v := range requests {
		//fmt.Println("##test", mp, <-myChannel)
		resR[int(v)] = <-myChannel

	}
	//fmt.Println("##test-########THIS", resR)

	return resR, nil

}

// !++++++++++++++++++++++++++++++++++++++++++++
// | make unique multiple post - json type request
// ++++++++++++++++++++++++++++++++++++++++++++
func MakeMultiplePostUnique(numberOfConcurrentRequests uint, url string, timeouts int, dataToSend []map[string]string, headers map[string]string) (map[int]map[int]map[int]string, error) {
	myChannel := make(chan map[int]map[int]string)
	for k, dataMap := range dataToSend {

		go func(dataMap map[string]string, k int) {
			result, _ := MakeMultiplePOST(numberOfConcurrentRequests, url, timeouts, dataMap, headers)

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

// !++++++++++++++++++++++++++++++++++++++++++++
// | make unique multiple post form type request
// ++++++++++++++++++++++++++++++++++++++++++++
func MakeMultiplePostFormUnique(numberOfConcurrentRequests uint, url string, timeouts int, dataToSend []map[string]string, headers map[string]string) (map[int]map[int]map[int]string, error) {
	myChannel := make(chan map[int]map[int]string)
	for k, dataMap := range dataToSend {

		go func(dataMap map[string]string, k int) {
			result, _ := MakeMultiplePostFormData(numberOfConcurrentRequests, url, timeouts, dataMap, headers)

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

// !++++++++++++++++++++++++++++++++++++++++++++
// | make unique multiple get request 
// ++++++++++++++++++++++++++++++++++++++++++++
func MakeMultipleGetUnique(numberOfConcurrentRequests uint, url string, timeouts int, dataToSend []string, headers map[string]string) (map[int]map[int]map[int]string, error) {
	myChannel := make(chan map[int]map[int]string)
	for k, dataUrl := range dataToSend {

		go func(dataUrl string, k int) {
			result, _ := MakeMultipleGET(numberOfConcurrentRequests, dataUrl, timeouts, headers)

			myChannel <- result

		}(dataUrl, k)

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
