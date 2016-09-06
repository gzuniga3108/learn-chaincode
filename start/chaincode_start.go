
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

var accountPrefix = "acct:"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Account struct {
	ID          string  `json:"id"`
	Prefix      string  `json:"prefix"`
	CashBalance float64 `json:"cashBalance"`
	AssetsIds   []string `json:"assetIds"`
}



// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err := stub.PutState("hello_world",[]byte(args[0]))
	if(err != nil){
		return nil,err
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "write" {
		return t.write(stub,args)
	}else if function == "createAccount" {
		return t.createAccount(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub,args)
	}
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query")
}


//Custom functions 
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub,args []string)([]byte, error){
	var key,value string
	var err error
	fmt.Println("running write()")
	if len(args) != 2 {
		return nil,errors.New("Incorrect number of arguments. Expecting 2")
	}

	key 	= args[0]
	value 	= args[1]
	err = stub.PutState(key, []byte(value))
	if(err != nil){
		return nil,err
	}
	return nil,nil
}

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub,args []string)([]byte,error){
	var key,jsonResp string	
	var err error
	if len(args) != 1 {
		return nil,errors.New("Incorrect number of arguments, expecting 1")
	}
	key = args[0]
	jsonResp =  args[0]
	valAsbytes,err := stub.GetState(key)
	if err != nil{
		jsonResp = "{\"Error\":\"Failed to get state for "+ key +"\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes,nil
}

func (t *SimpleChaincode) createAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    // Obtain the username to associate with the account
    if len(args) != 1 {
        fmt.Println("Error obtaining username")
        return nil, errors.New("createAccount accepts a single username argument")
    }
    username := args[0]
    
    // Build an account object for the user
    var assetIds []string
    suffix := "000A"
    prefix := username + suffix
    var account = Account{ID: username, Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
    accountBytes, err := json.Marshal(&account)
    if err != nil {
        fmt.Println("error creating account" + account.ID)
        return nil, errors.New("Error creating account " + account.ID)
    }
    
    fmt.Println("Attempting to get state of any existing account for " + account.ID)
    existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {        
        var company Account
        err = json.Unmarshal(existingBytes, &company)
        if err != nil {
            fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())
            
            if strings.Contains(err.Error(), "unexpected end") {
                fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
                err = stub.PutState(accountPrefix+account.ID, accountBytes)
                
                if err == nil {
                    fmt.Println("created account" + accountPrefix + account.ID)
                    return nil, nil
                } else {
                    fmt.Println("failed to create initialize account for " + account.ID)
                    return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
                }
            } else {
                return nil, errors.New("Error unmarshalling existing account " + account.ID)
            }
        } else {
            fmt.Println("Account already exists for " + account.ID + " " + company.ID)
		    return nil, errors.New("Can't reinitialize existing user " + account.ID)
        }
    } else {
        
        fmt.Println("No existing account found for " + account.ID + ", initializing account.")
        err = stub.PutState(accountPrefix+account.ID, accountBytes)
        
        if err == nil {
            fmt.Println("created account" + accountPrefix + account.ID)
            return nil, nil
        } else {
            fmt.Println("failed to create initialize account for " + account.ID)
            return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
        }
        
    }      
}