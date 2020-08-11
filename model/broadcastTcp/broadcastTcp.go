package broadcastTcp

import (
	"encoding/json"
	"fmt"
	"net"

	"../cryptogrpghy"
	"../validator"
)

var TempData []TCPData

//TCPData struct contain data about object,package name ,method
type TCPData struct {
	Obj         []byte
	ValidatorIP string
	Method      string
	PackageName string // CurrentTime        []string
	Signature   string
}
type NetStruct struct {
	Encryptedkey  string
	Encrypteddata string
}

type TxBroadcastResponse struct {
	TxID  string
	Valid bool
}

//TempNotRecieving IS store tcpdata an ip
// type TempNotRecieving struct {
// 	TCPData
// 	ValidatorSoketIP string
// }

// var temp []TempNotRecieving

//BoardcastingTCP Object
func BoardcastingTCP(obj interface{}, Method, PackageName string) TxBroadcastResponse {
	var res TxBroadcastResponse

	for _, validatorObj := range validator.ValidatorsLstObj {
		if !validatorObj.ValidatorRemove {
			if validatorObj.ValidatorIP == validator.CurrentValidator.ValidatorIP {
				if PackageName == "transaction" && Method == "addTransaction" {
					_, res = SendObject(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validator.CurrentValidator.ValidatorSoketIP)
					fmt.Println("\n @#########@ validatorIP", validatorObj.ValidatorIP, " and the res", res)
				} else {
					SendObject(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validator.CurrentValidator.ValidatorSoketIP)
				}
			} else {
				if PackageName == "transaction" && Method == "addTransaction" {
					_, res = SendObject(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validatorObj.ValidatorSoketIP)
					fmt.Println("\n @#########@ validatorIP", validatorObj.ValidatorIP, " and the res", res)
				} else {
					SendObject(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validatorObj.ValidatorSoketIP)
				}
			}
		}
	}
	return res
}

//SendObject to spacific miner
func SendObject(obj interface{}, Validatorpublickey, Method, PackageName, ValidatorSoketIP string) (TCPData, TxBroadcastResponse) {
	jsonObj, _ := json.Marshal(obj)
	//fmt.Println("\n **************** validator.CurrentValidator.ValidatorPrivateKey", validator.CurrentValidator.ValidatorPrivateKey)
	signature := cryptogrpghy.SignPKCS1v15(string(jsonObj), *cryptogrpghy.ParsePEMtoRSAprivateKey(validator.CurrentValidator.ValidatorPrivateKey))

	objTCP := TCPData{jsonObj, validator.CurrentValidator.ValidatorIP, Method, PackageName, signature}
	RemoteAddr, err := net.ResolveUDPAddr("udp", ValidatorSoketIP)
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	defer conn.Close()
	netObj := NetStruct{}

	/*---------------------------*/
	hashedkey := cryptogrpghy.CreateSHA1(Validatorpublickey)
	// netObj.Encryptedkey = cryptogrpghy.RSAENC(Validatorpublickey, []byte(hashedkey))
	netObj.Encryptedkey, _ = cryptogrpghy.PublicEncrypt(Validatorpublickey, hashedkey)
	byteData, _ := json.Marshal(objTCP)
	strofdata := string(byteData)
	netObj.Encrypteddata = cryptogrpghy.KeyEncrypt(hashedkey, strofdata)

	var responseObj TxBroadcastResponse

	//conn, err := net.Dial("tcp", ValidatorSoketIP)
	if err != nil {
		fmt.Println("error at at net.dial")
		fmt.Println(err)

	} else {
		byteData, _ := json.Marshal(netObj)
		_, err := conn.Write(byteData)

		if err != nil {
			fmt.Println("broadcastTcp Write data error:", err)
		}

		if objTCP.PackageName == "transaction" && Method == "addTransaction" {
			responseObj = ReadTxResponseData(conn)

		}
	}
	fmt.Println("\n responseObj := ReadTxResponseData(conn)", responseObj)
	return objTCP, responseObj
}

func ReadTxResponseData(conn *net.UDPConn) TxBroadcastResponse {
	var txResponse TxBroadcastResponse

	buffer := make([]byte, 1024)
	n, _, err1 := conn.ReadFromUDP(buffer)

	fmt.Println("UDP Server read bytes count: ", n)

	err2 := json.Unmarshal(buffer[:n], &txResponse)
	fmt.Println("\n json.Unmarshal(buffer[:n], &txResponse) error: ", err2)

	fmt.Println("the Response from broadcast handle:", txResponse)

	if err1 != nil {
		fmt.Println("broadcastTcp read data error1:", err1)
	}
	return txResponse
}

func SendTokenImg(obj string, Validatorpublickey, Method, PackageName, ValidatorSoketIP string) TCPData {
	jsonObj, _ := json.Marshal(obj)
	// signature := cryptogrpghy.SignPKCS1v15(string(jsonObj), *cryptogrpghy.ParsePEMtoRSAprivateKey(validator.CurrentValidator.ValidatorPrivateKey))

	objTCP := TCPData{jsonObj, validator.CurrentValidator.ValidatorIP, Method, PackageName, ""}
	RemoteAddr, err := net.ResolveUDPAddr("udp", ValidatorSoketIP)
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	defer conn.Close()
	netObj := NetStruct{}

	// hashedkey := cryptogrpghy.CreateSHA1(Validatorpublickey)
	netObj.Encryptedkey = "key"
	byteDat, _ := json.Marshal(objTCP)
	strofdata := string(byteDat)
	netObj.Encrypteddata = strofdata //cryptogrpghy.KeyEncrypt(hashedkey, strofdata)

	// conn, err := net.Dial("tcp", ValidatorSoketIP)
	if err != nil {
		fmt.Println("error at at net.dial")
		fmt.Println(err)

	} else {
		byteData, _ := json.Marshal(netObj)
		conn.Write(byteData)
		// reqBodyBytes := new(bytes.Buffer)
		// json.NewEncoder(reqBodyBytes).Encode(objTCP)
		// conn.Write(reqBodyBytes.Bytes())
	}
	return objTCP
}

func BoardcastingTokenImgUDP(obj string, Method, PackageName string) {
	// if PackageName == "transaction" {
	// }
	for _, validatorObj := range validator.ValidatorsLstObj {
		if !validatorObj.ValidatorRemove {
			if validatorObj.ValidatorIP == validator.CurrentValidator.ValidatorIP {

				SendTokenImg(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validator.CurrentValidator.ValidatorSoketIP)
			} else {

				SendTokenImg(obj, validatorObj.ValidatorPublicKey, Method, PackageName, validatorObj.ValidatorSoketIP)
			}
		}
	}
}
