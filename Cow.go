package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type SmartContract struct {
}

type Wallet struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

type Cow struct {
	Name     string `json:"name"`
	Maker    string `json:"maker"`
	Price    string `json:"price"`
	WalletID string `json:"walletid"`
}

type CowKey struct {
	Key string
	Idx int
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "initWallet" {
		return s.initWallet(APIstub)
	} else if function == "getWallet" {
		return s.getWallet(APIstub, args)
	} else if function == "setCow" {
		return s.setCow(APIstub, args)
	} else if function == "getAllCow" {
		return s.getAllCow(APIstub)
	} else if function == "purchaseCow" {
		return s.purchaseCow(APIstub, args)
	}
	fmt.Println("Please check your function : " + function)
	return shim.Error("Unknown function")
}

func (s *SmartContract) initWallet(APIstub shim.ChaincodeStubInterface) pb.Response {
	seller := Wallet{Name: "Hyper", ID: "hyper", Token: "1000"}
	customer := Wallet{Name: "Ledger", ID: "ledger", Token: "2000"}

	SellerasJSONBytes, _ := json.Marshal(seller)
	err := APIstub.PutState(seller.ID, SellerasJSONBytes)
	if err != nil {
		return shim.Error("Failed to create asset " + seller.Name)
	}

	CustomerasJSONBytes, _ := json.Marshal(customer)
	err = APIstub.PutState(customer.ID, CustomerasJSONBytes)
	if err != nil {
		return shim.Error("Failed to create asset " + customer.Name)
	}
	return shim.Success(nil)
}

func generateKey(stub shim.ChaincodeStubInterface) []byte {
	var isFirst bool = false
	cowkeyAsBytes, err := stub.GetState("latestKey")
	if err != nil {
		fmt.Println(err.Error())
	}
	cowkey := CowKey{}
	json.Unmarshal(cowkeyAsBytes, &cowkey)
	var tempIdx string
	tempIdx = strconv.Itoa(cowkey.Idx)
	fmt.Println(cowkey)
	fmt.Println("Key is " + strconv.Itoa(len(cow.Key)))
	if len(cowkey.Key) == 0 || cowkey.Key == "" {
		isFirst = true
		cowkey.Key = "MS"
	}
	if !isFirst {
		cowkey.Idx = cowkey.Idx + 1
	}
	fmt.Println("Last CowKey is " + cowkey.Key + " : " + tempIdx)
	returnValueBytes, _ := json.Marshal(cowkey)

	return returnValueBytes
}

func (s *SmartContract) setCow(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var cowkey = CowKey{}
	json.Unmarshal(generateKey(APIstub), &cowkey)
	keyidx := strconv.Itoa(cowkey.Idx)
	fmt.Println("Key : " + cowkey + ", Idx : " + keyidx)

	var cow = Cow{Name: args[0], Maker: args[1], Price: args[2], WalletID: args[3]}
	cowAsJSONBytes, _ := json.Marshal(cow)
	var keyString = cowkey.Key + keyidx
	fmt.Println("cowkey is " + keyString)
	err := APIstub.PutState(keyString, cowAsJSONBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record cow catch: %s", cowkey))
	}
	cowkeyAsBytes, _ := json.Marshal(cowkey)
	APIstub.PutState("latestKey", cowkeyAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getAllCow(APIstub shim.ChaincodeStubInterface) pb.Response {
	cowkeyAsBytes, _ := APIstub.GetState("latestKey")
	cowkey := CowKey{}
	json.Unmarshal(cowkeyAsBytes, &cowkey)
	idxStr := strconv.Itoa(cowkey.Idx + 1)

	var startKey = "MS0"
	var endKey = cowkey.Key + idxStr

	resultsIter, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIter.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIter.HasNext() {
		queryResponse, err := resultsIter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

func (t *SmartContract) purchaseCow(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorret number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	cowAsBytes, err := APIstub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	cowprice, _ := strconv.Atoi(string(cow.Price))

	AAsBytes, err := APIstub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if AAsBytes == nil {
		return shim.Error("Entity not found")
	}
	walletA := Wallet{}
	json.Unmarshal(AAsBytes, &walletA)
	tokenA, _ := strconv.Atoi(string(walletA.Token))

	BAsBytes, err := APIstub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if BAsBytes == nil {
		return shim.Error("Entity not found")
	}
	walletB := Wallet{}
	json.Unmarshal(BAsBytes, &walletB)
	tokenB, _ := strconv.Atoi(string(walletB.Token))

	walletA.Token = strconv.Itoa(tokenA - cowprice)
	walletB.Token = strconv.Itoa(tokenB - cowprice)
	updatedAAsBytes, _ := json.Marshal(walletA)
	updatedBAsBytes, _ := json.Marshal(walletB)
	APIstub.PutState(args[0], updatedAAsBytes)
	APIstub.PutState(args[1], updatedBAsBytes)

	fmt.Printf("A Token = %d, B Token = %d\n", walletA.Token, walletB.Token)
	return shim.Success(nil)
}

func (s *SmartContract) getWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	walletAsBytes, err := stub.GetState(args[0])
	if err != nil {
		fmt.Println(err.Error())
	}

	wallet := Wallet{}
	json.Unmarshal(walletAsBytes, &wallet)

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	if bArrayMemberAlreadyWritten == true {
		buffer.WriteString(",")
	}
	buffer.WriteString("{\"Name\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.Name)
	buffer.WriteString("\"")

	buffer.WriteString(", \"ID\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.ID)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Token\":")
	buffer.WriteString("\"")
	buffer.WriteString(wallet.Token)
	buffer.WriteString("\"")

	buffer.WriteString("}")
	bArrayMemberAlreadyWritten = true
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode : %s", err)
	}
}
