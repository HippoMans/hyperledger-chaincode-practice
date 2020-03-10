package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type BasicChain struct {
}

func (t *BasicChain) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Error Incorrect arguments")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
		return shim.Success(nil)

	}
}

func (t *BasicChain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "transfer" {
		result, err = transfer(stub, args)
	} else {
		result, err = get(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Error Incorrect arguments")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key, %s", args[0])
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error : %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found : %s", args[0])
	}
	return string(value), nil
}

func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var A, B string
	var Aval, Bval int
	var X int
	var err error

	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	Avalbytes, err := stub.GetState(A)
	Aval, _ = strconv.Atoi(string(Avalbytes))
	Bvalbytes, err := stub.GetState(b)
	Bval, _ = strconv.Atoi(string(Bvalbytes))
	X, err = strconv.Atoi(args[2])

	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return args[2], nil
}

func main() {
	if err := shim.Start(new(BasicChain)); err != nil {
		fmt.Printf("Error starting BasicChain : %s", err)
	}
}
