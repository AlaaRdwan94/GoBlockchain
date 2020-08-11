package ledger

import (
	"encoding/json" //read and send json data through api
	"fmt"
	"net/http" // using API request
	"os"

	"github.com/BurntSushi/toml"

	"../logfunc"

	//"time"

	"../accountdb"
	"../service"

	account "../account" //use accounts in the ledger structure
	"../admin"
	block "../block" //use blocks in the ledger structure

	//  write an error on the json file

	"../globalPkg"
	"../heartbeat"
	"../logpkg"
	"../systemupdate"
	"../token"
	"../tokenModule"
	transaction "../transaction" // use transaction in ledger structure
	validator "../validator"     //use validator in ledger structure

	//alaa
	globalfinctiontransaction "../globalfinctiontransaction"
)

/*-----------------structure of the ledger-----------*/
type Ledger struct {
	AccountsLstObj        []accountdb.AccountStruct
	ValidatorsLstObj      []validator.ValidatorStruct
	UnconfirmedValidators []validator.TempValidator
	ResetPassArray        []account.ResetPasswordData
	UserObjects           []account.User
	TransactionLstObj     []transaction.Transaction
	BlockchainObj         []block.BlockStruct
	AdminObj              []admin.AdminStruct
	TokenObj              []token.StructToken
	ServiceTmp            []service.ServiceStruct
	PurchasedService      []service.ServiceStruct
	LogDB                 []logpkg.LogStruct
	UserPK                []account.SavePKStruct

	//alaa
	ValidatorMap map[string]string
}

type MixedObjStruct struct {
	AdminObject  admin.Admin1
	LedgerObject Ledger
}
type TmpAccount struct {
	SessionDB []accountdb.AccountSessionStruct
	EmailDB   []accountdb.AccountEmailStruct
	NameDB    []accountdb.AccountNameStruct
	PhoneDB   []accountdb.AccountPhoneNumberStruct
}

var AdminObjec admin.Admin1

// GetLedger function to get the ledger of the miner
func GetLedger() Ledger {
	ledgerObj := Ledger{}
	ledgerObj.AccountsLstObj = accountdb.GetAllAccounts()

	ledgerObj.ValidatorsLstObj = validator.ValidatorsLstObj
	ledgerObj.UnconfirmedValidators = validator.TempValidatorlst
	// for index := range ledgerObj.ValidatorsLstObj {
	// 	if validator.CurrentValidator.ValidatorPublicKey == ledgerObj.ValidatorsLstObj[index].ValidatorPublicKey {
	// 		ledgerObj.ValidatorsLstObj[index].ValidatorPrivateKey = ""
	// 	}
	// }
	ledgerObj.BlockchainObj = block.GetBlockchain()
	ledgerObj.TransactionLstObj = transaction.GetPendingTransactions()
	ledgerObj.UserObjects = account.GetUserObjLst()
	ledgerObj.ResetPassArray = account.GetResetPasswordData()
	ledgerObj.AdminObj = admin.GetAllAdmins()
	ledgerObj.TokenObj = token.GetAllTokens()
	ledgerObj.ServiceTmp = service.GetAllservice()
	ledgerObj.PurchasedService = service.GetAllPurchusedservice()
	//alaa
	ledgerObj.ValidatorMap = globalfinctiontransaction.GetTransactionIndexTemMap()
	ledgerObj.LogDB = logfunc.GetAllLogs()
	ledgerObj.UserPK = account.GetAllsavepksave()

	return ledgerObj

}

func RemoveDatabase() {
	if accountdb.Open {
		accountdb.DB.Close()
		// account.DBEmail.Close()
		// account.DBName.Close()
		// account.DBPublicKey.Close()
		// account.DBPhoneNo.Close()
		// account.DBLastUpdateTime.Close()
		accountdb.Open = false
	}

	if block.Open {
		block.DB.Close()
		block.Open = false
	}
	if account.Opensave {
		account.DBsave.Close()
		account.Opensave = false
	}

	if admin.Open {
		admin.DB.Close()
		admin.Open = false
	}

	if token.Open {
		token.DB.Close()
		token.Open = false
	}
	if service.Open {
		service.DB.Close()
		service.Open = false
	}

	fmt.Println("tessssssssssssssssssst")
	//err := os.RemoveAll("AccountStruct")
	os.RemoveAll("Database/AccountStruct")
	os.RemoveAll("Database/AdminStruct")
	os.RemoveAll("Database/BlockStruct")
	//os.RemoveAll("Database/TempAccount")
	os.RemoveAll("Database/TokenStruct")
	os.RemoveAll("Database/SessionStruct")
	os.RemoveAll("Database/Service")
	os.RemoveAll("Database/SavePKStruct")
	// fmt.Println("tessssssssssssssssssst2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

// PostLegderAPI API to set the ledger of the miner-----------*/
func PostLegderAPI(w http.ResponseWriter, req *http.Request) {
	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "PostLegderAPI", "Ledger", "_", "_", "_", 0}

	recievedObj := MixedObjStruct{}
	adminReqObj := admin.Admin1{}
	ledgerReqObj := Ledger{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&recievedObj)
	if err != nil {
		globalPkg.SendError(w, "Error at reading object, please recheck it!")
		globalPkg.WriteLog(logobj, "error at reading object,please recheck it!", "failed")
		return
	}
	adminReqObj = recievedObj.AdminObject
	ledgerReqObj = recievedObj.LedgerObject

	if adminReqObj == AdminObjec {
		fmt.Println("nice")
	}

	RemoveDatabase()
	for _, accountObj := range ledgerReqObj.AccountsLstObj {
		accountObj.BlocksLst = nil
		account.AddAccount(accountObj)
	}

	for index := range ledgerReqObj.ValidatorsLstObj {
		if validator.CurrentValidator.ValidatorPublicKey == ledgerReqObj.ValidatorsLstObj[index].ValidatorPublicKey {
			ledgerReqObj.ValidatorsLstObj[index].ValidatorPrivateKey = validator.CurrentValidator.ValidatorPrivateKey
		} else {
			ledgerReqObj.ValidatorsLstObj[index].ValidatorPrivateKey = ""
		}
		validator.ValidatorCreate(ledgerReqObj.ValidatorsLstObj[index])
	}

	validator.ValidatorsLstObj = ledgerReqObj.ValidatorsLstObj
	for _, transactionObj := range ledgerReqObj.TransactionLstObj {
		transaction.AddTransaction(transactionObj)
	}
	for _, validatorObject := range ledgerReqObj.UnconfirmedValidators {
		validator.TempValidatorlst = append(validator.TempValidatorlst, validatorObject)
	}
	for _, blockObj := range ledgerReqObj.BlockchainObj {
		block.AddBlock(blockObj, true)
	}

	account.SetResetPasswordData(ledgerReqObj.ResetPassArray)
	account.SetUserObjLst(ledgerReqObj.UserObjects)

	for _, adminObj := range ledgerReqObj.AdminObj {
		admin.CreateAdmin(adminObj)
	}

	for _, tokenobj := range ledgerReqObj.TokenObj {
		tokenModule.AddToken(tokenobj)
	}
	for _, serviceObj := range ledgerReqObj.PurchasedService {
		service.AddAndUpdateServiceObj(serviceObj)
	}
	for _, UserPKObj := range ledgerReqObj.UserPK {
		account.SavePKAddress(UserPKObj)
	}
	fmt.Println("UserPKObj from ledger :", ledgerReqObj.UserPK)
	service.SetserviceTemp(ledgerReqObj.ServiceTmp)

	//alaa
	globalfinctiontransaction.SetTransactionIndexTemMap(ledgerReqObj.ValidatorMap) //globalfinctiontransaction.GetTransactionIndexTemMap()
	service.SetserviceTemp(ledgerReqObj.ServiceTmp)

	globalPkg.SendResponseMessage(w, "Ledger posted successfully")
	globalPkg.WriteLog(logobj, "Ledger posted successfully", "success")
}

func PostLedger(ledgerReqObj Ledger) {

	RemoveDatabase()
	for _, accountObj := range ledgerReqObj.AccountsLstObj {
		accountObj.BlocksLst = nil
		account.AddAccount(accountObj)
	}

	for index := range ledgerReqObj.ValidatorsLstObj {
		if validator.CurrentValidator.ValidatorPublicKey == ledgerReqObj.ValidatorsLstObj[index].ValidatorPublicKey {
			ledgerReqObj.ValidatorsLstObj[index].ValidatorPrivateKey = validator.CurrentValidator.ValidatorPrivateKey
		} else {
			ledgerReqObj.ValidatorsLstObj[index].ValidatorPrivateKey = ""
		}
		validator.ValidatorCreate(ledgerReqObj.ValidatorsLstObj[index])
	}
	for _, UserPKObj := range ledgerReqObj.UserPK {
		account.SavePKAddress(UserPKObj)
	}
	fmt.Println("UserPKObj from ledger :", ledgerReqObj.UserPK)
	validator.ValidatorsLstObj = ledgerReqObj.ValidatorsLstObj

	for _, transactionObj := range ledgerReqObj.TransactionLstObj {
		transaction.AddTransaction(transactionObj)
	}
	for _, blockObj := range ledgerReqObj.BlockchainObj {
		block.AddBlock(blockObj, true)
	}

	account.SetResetPasswordData(ledgerReqObj.ResetPassArray)
	account.SetUserObjLst(ledgerReqObj.UserObjects)
	for _, adminObj := range ledgerReqObj.AdminObj {
		admin.CreateAdmin(adminObj)
	}
	for _, tokenobj := range ledgerReqObj.TokenObj {
		tokenModule.AddToken(tokenobj)
	}
	for _, serviceObj := range ledgerReqObj.PurchasedService {
		service.AddAndUpdateServiceObj(serviceObj)
	}
	service.SetserviceTemp(ledgerReqObj.ServiceTmp)

	for _, validatorObject := range ledgerReqObj.UnconfirmedValidators {
		validator.TempValidatorlst = append(validator.TempValidatorlst, validatorObject)
	}

	globalfinctiontransaction.SetTransactionIndexTemMap(ledgerReqObj.ValidatorMap)

}

//GetLegderAPI API to get the ledger from the miner-----------*/
func GetLegderAPI(w http.ResponseWriter, req *http.Request) {
	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "GetLegderAPI", "ledger", "_", "_", "_", 0}

	Adminobj := admin.Admin{}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&Adminobj)

	if err != nil {
		globalPkg.SendError(w, "please enter your correct request")
		globalPkg.WriteLog(logobj, "please enter your correct request", "failed")
		return
	}

	if admin.ValidationAdmin(Adminobj) {
		sendJSON, _ := json.Marshal(GetLedger())
		globalPkg.SendResponse(w, sendJSON)
		globalPkg.WriteLog(logobj, "get ledger success", "success")

		var UpdateDataObj systemupdate.UpdateData
		toml.DecodeFile("./config.toml", &UpdateDataObj)
		heartbeat.SendUpdateHeartBeat(UpdateDataObj.Updatestruct.Updateversion, UpdateDataObj.Updatestruct.Updateurl)
	} else {
		globalPkg.SendError(w, "you are not admin")
		globalPkg.WriteLog(logobj, "you are not admin to get ledger", "failed")
	}

}

func GetAllTmpAccountDB() TmpAccount {
	ledgerObj := TmpAccount{}
	ledgerObj.EmailDB = accountdb.GetAllEmails()
	ledgerObj.NameDB = accountdb.GetAllNames()
	ledgerObj.PhoneDB = accountdb.GetAllPhones()
	ledgerObj.SessionDB = accountdb.GetAllSessions()
	return ledgerObj
}

/////////////////////
func GetTmpAccountDB(w http.ResponseWriter, req *http.Request) {
	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"", now, userIP, "macAdress", "GetLegderAPI", "ledger", "", "", "", 0}

	Adminobj := admin.Admin{}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&Adminobj)
	if err != nil {
		globalPkg.SendError(w, "please enter your correct request")
		globalPkg.WriteLog(logobj, "please enter your correct request", "failed")
		return
	}

	if admin.ValidationAdmin(Adminobj) {
		sendJSON, _ := json.Marshal(GetAllTmpAccountDB())
		globalPkg.SendResponse(w, sendJSON)
		globalPkg.WriteLog(logobj, "get GetAllTmpAccountDB success", "success")
	} else {
		globalPkg.SendError(w, "you are not admin")
		globalPkg.WriteLog(logobj, "you are not admin to get tmpAccountDB", "failed")
	}

}

func PostLedgerTmpAccount(ledgerReqObj TmpAccount) {

	RemoveDatabase()

	//validator.CurrentValidator.ValidatorRegisterTime, _ = time.Parse("2006-01-02 03:04:05 PM -0000", time.Now().UTC().Format("2006-01-02 03:04:05 PM -0000"))

	for _, sessionObj := range ledgerReqObj.SessionDB {
		accountdb.AddSessionIdStruct(sessionObj)
	}
	// for _, sessionObj := range ledgerReqObj.NameDB {
	// 	accountdb.AddSessionIdStruct(sessionObj)
	// }

}
