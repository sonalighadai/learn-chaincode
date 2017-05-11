/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type ledger struct {
	Vendor string
	Time string
	Geolocation string
	Vehicleno string
	Vehicletype string
	Items []struct{
		Name string
		Desc string
		Qty int
	}
	Defects []struct{
		Name string
		Desc string
		Qty int
	}

}

type Warehouse struct {
	Vendor string
	Time string
	Geolocation string
	Vehicleno string
	Vehicletype string
	Name string
	Desc string
	ScannedItem int
	Defect int
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "enterHDWHLedgerDetails" {
		 return t.enterHDWHLedgerDetails(stub,"write", args)
	}else if function == "write" { //read a variable
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) enterHDWHLedgerDetails(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	var key, s, value string
	str := make([]string, 2)
	key = args[0]
        value = args[1]
        bytes := []byte(value)
	
	// Unmarshal string into structs.
	var languages []ledger
	json.Unmarshal(bytes, &languages)
	
	// Loop over structs and display them.
	for l := range languages {
		for item := range languages[l].Items {
			for defect := range languages[l].Defects {
				warehouse := Warehouse{
					Vendor: languages[l].Vendor,
					Time: languages[l].Time,
					Geolocation: languages[l].Geolocation,
					Vehicleno: languages[l].Vehicleno,
					Vehicletype: languages[l].Vehicletype,
					Name: languages[l].Items[item].Name,
					Desc: languages[l].Items[item].Desc,
					ScannedItem: languages[l].Items[item].Qty-languages[l].Defects[defect].Qty,
					Defect: languages[l].Defects[defect].Qty,
				}
				// Create JSON from the instance data.
				// ... Ignore errors.
				b, _ := json.Marshal(warehouse)
				// Convert bytes to string.
				s = string(b)
				str[0] = key
                                str[1] = s
				
			}
		}
	}
	
	// Handle different functions
	if function == "writ" {
		 return t.write(stub, str)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + str[1])
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
