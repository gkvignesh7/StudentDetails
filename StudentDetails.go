package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Studentdetails struct {
	RegisterNo string `json:"RegisterNo"`
	Name       string `json:"Name"`
	Year       string `json:"Year"`
	Department string `json:"Department"`
	CGPA       string `json:"CGPA"`
	College    string `json:"College"`
}

func (t *Studentdetails) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}
func (t *Studentdetails) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	var result string

	var err error

	if function == "GetStudentdetailsByRegisterNo" {

		result, err = GetStudentdetailsByRegisterNo(stub, args)
	}
	if function == "GetStudentdetailsByNameYearAndDept" {
		result, err = GetStudentdetailsByNameYearAndDept(stub, args)
	}
	if function == "GetAllDetails" {
		result, err = GetAllDetails(stub)
	}
	if function == "GetByYearDeptAndCGPA" {
		result, err = GetByYearDeptAndCGPA(stub, args)
	}
	if function == "CreateStudentdetails" {
		result, err = CreateStudentdetails(stub, args)
	}
	if function == "GetByNameYearDeptCGPAandCollege" {
		result, err = GetByNameYearDeptCGPAandCollege(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}
func CreateStudentdetails(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering CreateStudentdetails")

	if len(args) < 6 {
		fmt.Println("Invalid number of args")
		return "", fmt.Errorf("Incorrect arguments. Expecting 6 arguments")
	}

	fmt.Println("Number of arguments:", len(args))

	RegisterNo1 := args[0]
	Name1 := args[1]
	Year1 := args[2]
	Department1 := args[3]
	CGPA1 := args[4]
	College1 := args[5]

	//assigning to struct the variables
	StudentdetailsStruct := Studentdetails{
		RegisterNo: RegisterNo1,
		Name:       Name1,
		Year:       Year1,
		Department: Department1,
		CGPA:       CGPA1,
		College:    College1,
	}
	fmt.Println("Student details struct:::::", StudentdetailsStruct)

	StudentdetailsBytes, err := json.Marshal(StudentdetailsStruct)
	if err != nil {
		fmt.Println("Couldn't marshal data from struct", err)
		return "", fmt.Errorf("Couldn't marshal data from struct")

	}
	fmt.Println("Register No is:::", RegisterNo1)
	fmt.Println("Student data in bytes:::", StudentdetailsBytes)
	StudentdetailsErr := stub.PutState(RegisterNo1, []byte(StudentdetailsBytes))
	if StudentdetailsErr != nil {
		fmt.Println("Couldn't save StudentdetailsStruct Characterestic data to ledger", StudentdetailsErr)
		return "", fmt.Errorf("Couldn't save StudentdetailsStruct Characterestic data to ledger")
	}

	/////NameYearAndDept CompKey

	NameYearAndDept := "StudentdetailsIndexx"
	NameYearAndDeptCompKey, err := stub.CreateCompositeKey(NameYearAndDept, []string{StudentdetailsStruct.Name, StudentdetailsStruct.Year, StudentdetailsStruct.Department, StudentdetailsStruct.RegisterNo})
	fmt.Println("comp key:::", NameYearAndDeptCompKey)
	if err != nil {
		return "", fmt.Errorf("composite key not found")
	}
	fmt.Println(NameYearAndDeptCompKey)

	value := []byte{0x00}
	stub.PutState(NameYearAndDeptCompKey, value)

	//////////// YearDepCgpa

	YearDepCgpaIndex1 := "YearDepCgpaIndex"
	YearDeptCgpaCompKey, err := stub.CreateCompositeKey(YearDepCgpaIndex1, []string{StudentdetailsStruct.Year, StudentdetailsStruct.Department, StudentdetailsStruct.Department, StudentdetailsStruct.CGPA, StudentdetailsStruct.RegisterNo})
	fmt.Println("comp key:::", YearDeptCgpaCompKey)
	if err != nil {
		return "", fmt.Errorf("composite key not found")
	}
	fmt.Println(YearDeptCgpaCompKey)

	value1 := []byte{0x00}
	stub.PutState(YearDeptCgpaCompKey, value1)
	///////NameYearDeptCGPAandCollege

	NameYearDeptCGPAandCollegeIndex1 := "NameYearDeptCGPAandCollegendex"
	NameYearDeptCGPAandCollegeCompKey, err := stub.CreateCompositeKey(NameYearDeptCGPAandCollegeIndex1, []string{StudentdetailsStruct.Name, StudentdetailsStruct.Year, StudentdetailsStruct.Department, StudentdetailsStruct.CGPA, StudentdetailsStruct.RegisterNo})
	fmt.Println("comp key:::", NameYearDeptCGPAandCollegeCompKey)
	if err != nil {
		return "", fmt.Errorf("composite key not found")
	}
	fmt.Println(NameYearDeptCGPAandCollegeCompKey)

	value2 := []byte{0x00}
	stub.PutState(NameYearDeptCGPAandCollegeCompKey, value2)
	/////////

	fmt.Println(args[0])

	fmt.Println("Successfully saved Studentdetails")
	return args[0], nil
}

func GetStudentdetailsByRegisterNo(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering GetStudentdetailsByRegisterNo")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return "", errors.New("Missing RegisterNo")
	}

	var RegisterNo = args[0]
	fmt.Println("data id::::", RegisterNo)
	value, err := stub.GetState(RegisterNo)
	if err != nil {
		fmt.Println("Couldn't get data for "+RegisterNo+" from ledger", err)
		return "", errors.New("Missing data")
	}

	return string(value), nil
}
func GetStudentdetailsByNameYearAndDept(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering get data by Name, Year and Dept")
	var Name string
	var Year string
	var Department string
	var err error
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting  Name, Year and Dept  to query")
	}

	Name = args[0]
	Year = args[1]
	Department = args[2]

	fmt.Println("Name::::", Name)
	fmt.Println("Year::::", Year)
	fmt.Println("Dept::::", Department)
	mdMapItr, err := stub.GetStateByPartialCompositeKey("StudentDetailsIndexx", []string{Name, Year, Department})
	if err != nil {
		return "", fmt.Errorf("Could not get composite key for this Name, Year, Dept")
	}
	fmt.Println("Iterator:::: ", mdMapItr)
	defer mdMapItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for mdMapItr.HasNext() {
		queryResponse, err := mdMapItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next record")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite key parts")
		}

		returnedName := compositeKeyParts[0]
		returnedYear := compositeKeyParts[1]
		returnedDepartment := compositeKeyParts[2]
		returnedRegisterNo := compositeKeyParts[3]

		fmt.Printf("- found a from index:%s Name:%s Year:%s Dept:%s data id:%s", objectType, returnedName, returnedYear, returnedDepartment, returnedRegisterNo)

		value, err := stub.GetState(returnedRegisterNo)
		if err != nil {
			return "", fmt.Errorf("Missing RegisterNo in ledger")
		}

		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Println("Data for deal, month and year::", buffer.String())

	return buffer.String(), nil
}
func GetAllDetails(stub shim.ChaincodeStubInterface) (string, error) {
	etAllPoolsItr, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Could not get all student data")
	}
	defer etAllPoolsItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("{\"Results\":")

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for etAllPoolsItr.HasNext() {
		queryResponseValue, err := etAllPoolsItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next Data")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponseValue.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	buffer.WriteString("}")

	return buffer.String(), nil
}
func GetByYearDeptAndCGPA(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering get data by Year, Dept and CGPA")
	var Year string
	var Department string
	var CGPA string
	var err error
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting  Year, Dept and CGPA  to query")
	}

	Year = args[0]
	Department = args[1]
	CGPA = args[2]

	fmt.Println("Year::::", Year)
	fmt.Println("Department::::", Department)
	fmt.Println("CGPA:::", CGPA)
	mdMapItr, err := stub.GetStateByPartialCompositeKey("YearDepCgpaIndex", []string{Year, Department, CGPA})
	if err != nil {
		return "", fmt.Errorf("Could not get composite key for this Year Dept, CGPA")
	}
	fmt.Println("Iterator:::: ", mdMapItr)
	defer mdMapItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for mdMapItr.HasNext() {
		queryResponse, err := mdMapItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next record")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite key parts")
		}

		returnedYear := compositeKeyParts[0]
		returnedDepartment := compositeKeyParts[1]
		returnedCGPA := compositeKeyParts[2]
		returnedRegisterNo := compositeKeyParts[3]

		fmt.Printf("- found a from index:%s Year:%s Department:%s CGPA:%s data id:%s", objectType, returnedYear, returnedDepartment, returnedCGPA, returnedRegisterNo)

		value, err := stub.GetState(returnedRegisterNo)
		if err != nil {
			return "", fmt.Errorf("Missing RegisterNo in ledger")
		}

		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Println("Data for deal, month and year::", buffer.String())

	return buffer.String(), nil
}
func GetByNameYearDeptCGPAandCollege(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("Entering get data by Name,Year,Dept,CGPA and College")
	var Name string
	var Year string
	var Department string
	var CGPA string
	var College string
	var err error
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting  Name,Year,Dept,CGPA and College  to query")
	}
	Name = args[0]
	Year = args[1]
	Department = args[2]
	CGPA = args[3]
	College = args[4]

	fmt.Println("Name::::", Name)
	fmt.Println("Year::::", Year)
	fmt.Println("Dept::::", Department)
	fmt.Println("CGPA::::", CGPA)
	fmt.Println("College::::", College)
	mdMapItr, err := stub.GetStateByPartialCompositeKey("StudenDetailsIndexx", []string{Name, Year, Department, CGPA, College})
	if err != nil {
		return "", fmt.Errorf("Could not get composite key for this Name , Year ,Dept,CGPA,College")
	}
	fmt.Println("Iterator:::: ", mdMapItr)
	defer mdMapItr.Close()

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for mdMapItr.HasNext() {
		queryResponse, err := mdMapItr.Next()
		if err != nil {
			return "", fmt.Errorf("Could not get next record")
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		objectType, compositeKeyParts, err := stub.SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return "", fmt.Errorf("Could not get Composite key parts")
		}
		returnedName := compositeKeyParts[0]
		returnedYear := compositeKeyParts[1]
		returnedDepartment := compositeKeyParts[2]
		returnedCGPA := compositeKeyParts[3]
		returnedCollege := compositeKeyParts[4]
		returnedRegisterNo := compositeKeyParts[5]

		fmt.Printf("- found a from index:%s Name:%s Year:%s Department:%s CGPA:%s College:%s data id:%s", objectType, returnedName, returnedYear, returnedDepartment, returnedCGPA, returnedCollege, returnedRegisterNo)

		value, err := stub.GetState(returnedRegisterNo)
		if err != nil {
			return "", fmt.Errorf("Missing RegisterNo in ledger")
		}

		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Println("Data for Name,Year,Dept,CGPA,College", buffer.String())

	return buffer.String(), nil
}

func main() {
	server := &shim.ChaincodeServer{
		CCID:    os.Getenv("CHAINCODE_CCID"),
		Address: os.Getenv("CHAINCODE_ADDRESS"),
		CC:      new(Studentdetails),
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	err := server.Start()

	if err != nil {
		fmt.Printf("Error starting Studentdetails chaincode: %s", err)
	}

}
