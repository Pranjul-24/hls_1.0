package main

import (

"fmt"
"strings"
"strconv"
"encoding/json"
"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"

) 

var EVENT_COUNTER = "event_counter"
type ManageEnrollment struct {

}

var ParticipantIndexStr = "_Participantindex"
var RetailerIndexStr = "_Retailerindex"
var ConsumerIndexStr = "_Consumerindex" 
var RegulatorIndexStr= "_Regulatorindex"
var ProsumerIndexStr = "_Prosumerindex"
var ProducerIndexStr    = "_ProducerIndex"

type Enrollment struct{

  ParticipantID string `json:"ParticipantID"`
  PrimaryRole string `json:"PrimaryRole"`
  Address   string `json:"Address"`         
  TimeStamp string `json:"TimeStamp"`
  Name string `json:"Name"`
  
}

  func main() {     
  err := shim.Start(new(ManageEnrollment))
  if err != nil {
    fmt.Printf("Error starting ManageEnrollment chaincode: %s", err)
  }
}

// Chaincode Instantiation
func (t *ManageEnrollment) Init(stub shim.ChaincodeStubInterface) peer.Response {

  args := stub.GetStringArgs()
  var msg string
  var err error
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  // Initialize the chaincode
  msg = args[0]
  fmt.Println("ManageEnrollment chaincode is deployed successfully.");
  
  // Write the state to the ledger
  err = stub.PutState("abc", []byte(msg))       //making a test var "abc", I find it handy to read/write to it right away to test the network
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
  }
  
  var empty []string
  jsonAsBytes, _ := json.Marshal(empty)               //marshal an emtpy array of strings to clear the index
  err = stub.PutState(ParticipantIndexStr, jsonAsBytes)
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset in ParticipantIndex: %s", args[0]))
  }
  err = stub.PutState(EVENT_COUNTER, []byte("1"))
  if err != nil {
    return shim.Error(fmt.Sprintf("Failed to create asset in event counter: %s", args[0]))
  }
  return shim.Success(nil)
}

// Invoke Chaincode
 func (t *ManageEnrollment) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
    //fmt.Println("invoke is running " + function)
       fn, args := stub.GetFunctionAndParameters()
  // Handle different functions
       var result string
    var err error

    if fn == "enroll" {
    	result, err = enroll(stub, args)
    } else if fn == "getParticipant_byRole"{
    	result, err = getParticipant_byRole(stub, args)
    }else if fn == "getParticipant_byID"{
    	result, err = getParticipant_byID(stub, args)
    }else if fn == "delete"{
    	result, err = delete(stub, args)
    }

    if err != nil {
            return shim.Error(err.Error())
    }
   fmt.Println("invoke did not find func: " + fn)          //error
  
  return shim.Success([]byte(result))
} 

// enroll participant  
func enroll(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var err error

  if len(args) != 5{
    return "", fmt.Errorf("Incorrect number of arguments. Expecting 5")
  }
  fmt.Println("Enrollment Started")
  
  ParticipantID := strings.ToLower(args[0])
  PrimaryRole   := args[1]

  TimeStamp     := args[2]
  Address       := args[3]
  Name          := args[4]



if EnrollmentAsBytes, err := stub.GetState(ParticipantID); err != nil || EnrollmentAsBytes != nil {
		return "", fmt.Errorf("This ID already exists.")
	} 

    ParticipantDetails :=  `{`+
    `"ParticipantID": "` + ParticipantID + `" , `+
    `"PrimaryRole": "` + PrimaryRole + `" , `+
    `"Address": "` + Address + `" , `+
    `"TimeStamp": "` + TimeStamp + `" , `+
    `"Name": "` + Name + `" , `+ 
    `}`

    fmt.Println("Participant Details in Array")
    fmt.Println(ParticipantDetails)

    err = stub.PutState(ParticipantID , []byte(ParticipantDetails))

    if err != nil {
    	return "" , err
    }

    ParticipantIndexAsBytes, err := stub.GetState(ParticipantIndexStr)

    if err != nil {
    	return "" , fmt.Errorf("Failed to get Index")
    }
    	var ParticipantIndex []string

    	json.Unmarshal(ParticipantIndexAsBytes, &ParticipantIndex)

    	ParticipantIndex = append(ParticipantIndex, ParticipantID)

    	jsonAsBytes, _ := json.Marshal(ParticipantIndex)

    	err = stub.PutState(ParticipantIndexStr, jsonAsBytes)

    	if err != nil{
    		return "" , err
    	}

    	fmt.Println("End of participant creation")
    	return string(ParticipantDetails), nil
}

func getParticipant_byRole(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  var PrimaryRole, jsonResp, errResp string
  var err error
  var valIndex Enrollment
  fmt.Println("start getParticipant_byRole")
  if len(args) != 1 {
    return "", fmt.Errorf("Incorrect number of arguments. Expecting ID of the participant to query")
  }
  // set ID
   PrimaryRole= args[0]
  EnrollmentAsBytes, err := stub.GetState(ParticipantID)                  //get the getParticipant_byRole from chaincode state
  if err != nil {
    jsonResp = "{\"Error\":\"Failed to get state for " + PrimaryRole + "\"}"
    return "",fmt.Errorf(jsonResp)
  }
 
    var ParticipantIndex []string
  json.Unmarshal(EnrollmentAsBytes, &ParticipantIndex) 

  jsonResp = "{"
  for i,val := range ParticipantIndex{

    fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getParticipant_byRole")
    valueAsBytes, err := stub.GetState(val)
    if err != nil {
      errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
      return "", fmt.Errorf(errResp)
    }
  
    var err1 error
    err1 = json.Unmarshal(valueAsBytes, &valIndex)
    if err1 != nil {
      fmt.Println(err1)
  }
      
    if valIndex.PrimaryRole == PrimaryRole{
      fmt.Println("Participant found")
      jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
      if i < len(ParticipantIndex)-1 {
        jsonResp = jsonResp + ","
      }}}
  //fmt.Print("valAsbytes : ")
  //fmt.Println(valAsbytes)
  jsonResp = jsonResp + "}"
  //fmt.Println("jsonResp : " + jsonResp)
  //fmt.Print("jsonResp in bytes : ")
  //fmt.Println([]byte(jsonResp))
  fmt.Println("end getby role")
  return string(jsonResp), nil 
}

 func getParticipant_byID(stub shim.ChaincodeStubInterface, args []string)(string, error){
 	var ParticipantID string
 	var err error 
 	fmt.Println("Started getParticipant_byID")
 	if len(args) != 1{
 		return "", fmt.Errorf("Incorrect arguments. Expecting a key (ParticipantID)")
}
 		ParticipantID = strings.ToLower(args[0])
 		valAsbytes, err := stub.GetState(ParticipantID)

 		if err != nil{
 			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
  		}
  		if valAsbytes == nil{
  			return "", fmt.Errorf("Asset not found: %s", args[0])
  		}
 	
 	return string(valAsbytes) , nil
}

func delete(stub shim.ChaincodeStubInterface, args []string)(string, error){
	fmt.Println("Delete function executed")

	if len(args) != 1 {
		return "" , fmt.Errorf("Incorrect argument. Expecting only one argument")
	}

	ParticipantID := args[0]
	err := stub.DelState(ParticipantID)

	ParticipantIndexAsBytes, err := stub.GetState(ParticipantIndexStr)

	if err != nil {
		return "" , fmt.Errorf("Failed to get Participant Index")
	}

	var ParticipantIndex []string
	json.Unmarshal(ParticipantIndexAsBytes, &ParticipantIndex)

	for i, val := range ParticipantIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + ParticipantID)
		 if val == ParticipantID{
		 	fmt.Println("found patient with matching patientID")
		 	ParticipantIndex = append(ParticipantIndex[:i], ParticipantIndex[i+1:]...)

		 	for x := range ParticipantIndex{
		 		fmt.Println(string(x) + " - " + ParticipantIndex[x])
		 	}
		 	break
		 }
	}

	jsonAsBytes , _ :=json.Marshal(ParticipantIndex)
	err = stub.PutState(ParticipantIndexStr, jsonAsBytes)
	return string(ParticipantIndexStr), nil
} 