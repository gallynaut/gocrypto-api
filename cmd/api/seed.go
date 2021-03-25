package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getJsonExchanges() (Exchanges, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("./data/exchanges.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Println("Successfully Opened exchanges.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	var allExchanges Exchanges
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &allExchanges)
	// fmt.Printf("Unmarshall: %s to \n%s\nLen: %d\n", byteValue, allExchanges, len(allExchanges))

	if len(allExchanges) == 0 {
		return allExchanges, fmt.Errorf("no exchanges loaded")
	}

	return allExchanges, nil
}
