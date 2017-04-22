
package main

import (
	"errors"
	"fmt"
    "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
var logger = shim.NewLogger("CLDChaincode")
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}


//ASSET
type KYCInfo struct {
	KYC_Id         string	  `json:"kyc_id"`
	Kyc_Type       string `json:"kyc_type"`
	Cust_Id        string `json:"cust_id"`
}

//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

type User_and_eCert struct {
	Identity string `json:"identity"`
	eCert string `json:"ecert"`
}

type Kyc_Holder struct {
	KYCs 	[]string `json:"kycs"`
}

func main() {

	err := shim.Start(new(SimpleChaincode))

															if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}


//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//Args
	//				0
	//			peer_address

	var  kycs Kyc_Holder


	bytes, err := json.Marshal(kycs)

    if err != nil { return nil, errors.New("Error creating KYC Holder record") }

	err = stub.PutState("Deployed successfully", bytes)


	return nil, nil
}


//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    logger.Debug("Inside Invoke")

	if function == "create_kyc" {
        logger.Debug("Inside Invoke: calling create kyc")

        return t.create_kyc(stub, args)
	} else if function == "ping" {
        return t.ping(stub)
    } else { 																				
		return nil, errors.New("Function of the name "+ function +" doesn't exist.")

	}
}

func (t *SimpleChaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}
//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Vehicle - Creates the initial JSON for the kyc and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_kyc(stub shim.ChaincodeStubInterface, k []string) ([]byte, error) {
	var v KYCInfo
    logger.Debug("Inside create KYC")

	kyc_id         := "\"KYC_Id\":\""+k[1]+"\", "							// Variables to define the JSON
	kyc_type       := "\"Kyc_Type\"\""+k[2]+"\","
    cust_id        := "\"Kyc_Type\"\""+k[0]+"\","

	fmt.Printf(+k[0]);
	fmt.Printf(+k[1]);
	fmt.Printf(+k[2]);
	kyc_json := "{"+kyc_id+kyc_type+cust_id+"}" 	// Concatenates the variables to create the total JSON object

	logger.Debug("kyc_json: ", kyc_json)

	err := json.Unmarshal([]byte(kyc_json), &v)							// Convert the JSON defined above into a vehicle object for go

	if err != nil { return nil, errors.New("Invalid JSON object....") }

	//record, err := stub.GetState(v.KYC_Id) 								// If not an error then a record exists so cant create a new car with this V5cID as it must be unique


	bytes, err := stub.GetState("KYCs")
																		if err != nil { return nil, errors.New("Unable to get KYCs") }
	var kycs Kyc_Holder

	err = json.Unmarshal(bytes, &kycs)

																		if err != nil {	return nil, errors.New("Corrupt KYC record") }


	kycs.KYCs = append(kycs.KYCs, kyc_id)


	bytes, err = json.Marshal(kycs)

															if err != nil { fmt.Print("Error creating Kyc_Holder record") }

	err = stub.PutState("kycs", bytes)

															if err != nil { return nil, errors.New("Unable to put the state") }

	return nil, nil

}



func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	
    logger.Debug("function: ", function)
    logger.Debug("customer: ", args[0])
 

	if function == "get_kyc_details" {
		if len(args) != 1 { fmt.Printf("Incorrect number of arguments passed"); return nil, errors.New("QUERY: Incorrect number of arguments passed") }
		//if err != nil { fmt.Printf("QUERY: Error retrieving cust_id: %s", err); return nil, errors.New("QUERY: Error retrieving cust_id "+err.Error()) }
	return t.get_kyc_details(stub,args[0] )
	//) else if function == "get_ecert" {
		//return t.get_ecert(stub, args[0])
	} else if function == "ping" {
		return t.ping(stub)
	}

	return nil, errors.New("Received unknown function invocation " + function)

}

func (t *SimpleChaincode) get_kyc_details(stub shim.ChaincodeStubInterface, cust_id string) ([]byte, error) {
	bytes, err := stub.GetState("kycs")

																			if err != nil { return nil, errors.New("Unable to get kycs") }
	var kycs Kyc_Holder
	err = json.Unmarshal(bytes, &kycs)
																			if err != nil {	return nil, errors.New("Corrupt V5C_Holder") }
	result := "["
	var temp []byte
	var v KYCInfo

	for _, kyc := range kycs.KYCs {

		v, err = t.retrieve_v5c(stub, kyc)

		if err != nil {return nil, errors.New("Failed to retrieve V5C")}
temp, err = t.get_kyc(stub, v,cust_id)

		if err == nil {
			result += string(temp) + ","
		}
	}
	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}
	return []byte(result), nil

}

func (t *SimpleChaincode) retrieve_v5c(stub shim.ChaincodeStubInterface, kyc_id string) (KYCInfo, error) {

	var v KYCInfo

	bytes, err := stub.GetState(kyc_id);
    fmt.Printf("KYC_ID------" +kyc_id);
	if err != nil {	fmt.Printf("RETRIEVE_V5C: Failed to invoke vehicle_code: %s", err); return v, errors.New("RETRIEVE_V5C: Error retrieving vehicle with v5cID = " + kyc_id) }

	err = json.Unmarshal(bytes, &v);

    if err != nil {	fmt.Printf("RETRIEVE_V5C: Corrupt vehicle record "+string(bytes)+": %s", err); return v, errors.New("RETRIEVE_V5C: Corrupt vehicle record"+string(bytes))	}

	return v, nil
}
func (t *SimpleChaincode) get_kyc(stub shim.ChaincodeStubInterface, v KYCInfo, cust_id string) ([]byte, error) {

	bytes, err := json.Marshal(v)

																if err != nil { return nil, errors.New("GET_VEHICLE_DETAILS: Invalid vehicle object") }

	if 		v.Cust_Id				== cust_id		{

					return bytes, nil
	} else {
																return nil, errors.New("Permission Denied. get_vehicle_details")
	}

}