package validatorModule

import (
	"encoding/json" //read and send json data through api
	"../globalfinctiontransaction"
	"net/http" // using API request

	"../ledger"

	"../broadcastTcp"

	"../admin"
	"../errorpk"   //  write an error on the json file
	"../globalPkg" //to use send request function
	"../logpkg"
	"../validator"
)

type MixedObjec struct {
	Admn  admin.Admin
	Vldtr validator.ValidatorStruct
}

//BroadcastValidatorAPI endpoint to broadcasting adding, updating or deleting validator in the miner  ----------------- */
func BroadcastValidatorAPI(w http.ResponseWriter, req *http.Request) {
	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "BroadcastValidatorAPI", "Validator", "_", "_", "_", 0}
	var parentObjec MixedObjec
	errStr := ""

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&parentObjec)

	if err != nil {
		globalPkg.SendError(w, "Error at reading object check that admin object is first and in format Admn and validator second in formate : Vldtr ")
		globalPkg.WriteLog(logobj, "Error at reading object check that admin object is first and in format Admn and validator second in formate : Vldtr ", "failed")
		return
	}

	// converting mixed obj to 2 object
	admin2 := parentObjec.Admn                    // admin object
	validator.NewValidatorObj = parentObjec.Vldtr // validator obj

	//check if authorized :
	if admin2.UsernameAdmin != "inoadmin" && admin2.PasswordAdmin != "In0v@ti@n@dmin" {
		globalPkg.SendError(w, "Not Authorized, false password and/or false admin name")
		globalPkg.WriteLog(logobj, "Not Authorized, false password and/or false admin name", "failed")
		return
	}

	validator.NewValidatorObj.ValidatorRegisterTime = globalPkg.UTCtime()
	validator.NewValidatorObj.ValidatorLastHeartBeat = globalPkg.UTCtime()

	if req.Method == "PUT" {
		broadcastTcp.BoardcastingTCP(validator.NewValidatorObj, req.Method, "validator")
	} else {
		now := globalPkg.UTCtime()
		confCode := validator.EncodeToString(4)
		tmpvalidator := validator.TempValidator{validator.NewValidatorObj, confCode, now}

		for _, timpValidator := range validator.ValidatorsLstObj {
			if timpValidator.ValidatorPublicKey == tmpvalidator.ValidatorObjec.ValidatorPublicKey {
				globalPkg.SendResponse(w, []byte("The Validator Public Key exists before"))
				globalPkg.WriteLog(logobj, "This Validator IP exists before", "failed")
				return
			}
		}
		for _, timpValidator := range validator.TempValidatorlst {
			if timpValidator.ValidatorObjec.ValidatorPublicKey == tmpvalidator.ValidatorObjec.ValidatorPublicKey {
				globalPkg.SendResponse(w, []byte("This Validator is waiting for confirmation"))
				globalPkg.WriteLog(logobj, "This Validator is waiting for confirmation", "failed")
				return
			}
		}
		adminslst := admin.GetAllAdmins()
		ValidatorAdmin := adminslst[0]
		validator.SendConfMail(ValidatorAdmin, tmpvalidator)

		broadcastTcp.BoardcastingTCP(tmpvalidator, req.Method, "validator")
		globalPkg.SendResponse(w, []byte("validator broadcasted successfully"))
		globalPkg.WriteLog(logobj, "validator broadcasted successfully", "success")

		//create tempvalidator
		// now := globalPkg.UTCtime()
		// confCode := validator.EncodeToString(4)
		// tmpvalidator := validator.TempValidator{validator.NewValidatorObj, confCode, now}
		// broadcastTcp.BoardcastingTCP(tmpvalidator, req.Method, "validator")
	}

	if errStr == "" {
		globalPkg.SendResponse(w, []byte("validator broadcasted successfully"))
		globalPkg.WriteLog(logobj, "validator broadcasted successfully", "success")
		return
	} else {
		globalPkg.SendError(w, errStr)
		globalPkg.WriteLog(logobj, errStr, "failed")
		return
	}
}

//ValidatorAPI endpoint to add, update or delete validator in the miner  ----------------- */
func ValidatorAPI(w http.ResponseWriter, req *http.Request) {

	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "ValidatorAPI", "Validator", "_", "_", "_", 0}
	validatorObj := validator.ValidatorStruct{}
	errorStr := ""

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&validatorObj)

	if err != nil {
		globalPkg.SendError(w, "please enter your correct request")
		globalPkg.WriteLog(logobj, "please enter your correct request", "failed")
		return
	}

	switch req.Method {
	case "POST":
		errorStr = validator.AddValidator(validatorObj)
	case "PUT":
		errorStr = validator.UpdateValidator(validatorObj)
	case "DELETE":
		errorStr = validator.DeleteValidator(validatorObj)
	default:
		errorStr = errorpk.AddError("Validator API validator package "+req.Method, "wrong method ", "logical error")

	}

	if errorStr == "" {
		sendJSON, _ := json.Marshal(validator.ValidatorsLstObj)
		globalPkg.SendResponse(w, sendJSON)
		globalPkg.WriteLog(logobj, "boardcast validator success to add or register validator", "success")

	} else {
		globalPkg.SendError(w, errorStr)
		globalPkg.WriteLog(logobj, errorStr, "failed")
	}
}

//GetAllValidatorAPI endpoint to get all validators from the miner  ----------------- */
func GetAllValidatorAPI(w http.ResponseWriter, req *http.Request) {
	//log
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "GetAllValidatorAPI", "ValidatorModule", "_", "_", "_", 0}

	Adminobj := admin.Admin{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&Adminobj)

	if err != nil {
		globalPkg.SendError(w, "  please enter your correct request")
		globalPkg.WriteLog(logobj, "please enter your correct request", "failed")
		return
	}

	if admin.ValidationAdmin(Adminobj) {
		json.NewEncoder(w).Encode(validator.ValidatorsLstObj)
		globalPkg.WriteLog(logobj, "get all validators success", "success")

	} else {
		globalPkg.SendError(w, "you are not admin")
		globalPkg.WriteLog(logobj, "you are not admin to get all validators", "failed")

	}

}

// DeactiveNode admin can change or Update status of validator IP from validatorActive to disactive
func DeactiveNode(w http.ResponseWriter, req *http.Request) {

	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "GetAllAdminsAPI", "adminModule", "_", "_", "_", 0}

	AdminObj := admin.AdminStruct{}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&AdminObj); err != nil {

		globalPkg.SendError(w, "please enter your correct request")
		globalPkg.WriteLog(logobj, "please enter your correct request", "failed")
		return
	}
	if AdminObj.AdminUsername == "" || AdminObj.AdminPassword == "" {
		globalPkg.SendError(w, "Please enter username and password")
		globalPkg.WriteLog(logobj, "please enter userName and password", "failed")
		return
	}

	Adminexist := admin.FindAdminByid(AdminObj.AdminUsername)

	if AdminObj.AdminUsername == Adminexist.AdminUsername && AdminObj.AdminPassword == Adminexist.AdminPassword {
		listValidator := Adminexist.Validatorlst
		exist := false
		for _, validatorip := range listValidator {
			if validatorip == AdminObj.ValiatorIPtoDeactive {
				validatorObj := validator.FindValidatorByValidatorIP(validatorip)

				validatorObj.ValidatorActive = !validatorObj.ValidatorActive
				validator.UpdateValidator(validatorObj)
				exist = true
			}
		}
		if exist == false {
			globalPkg.SendError(w, "please check validator ip ")
			globalPkg.WriteLog(logobj, "please check validator ip", "failed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		globalPkg.WriteLog(logobj, "Update validator status", "success")
	} else {
		globalPkg.SendError(w, "please check password and username")
		globalPkg.WriteLog(logobj, "please check password and username ", "failed")
	}
}

func ConfirmedValidatorAPI(w http.ResponseWriter, req *http.Request) {
	now, userIP := globalPkg.SetLogObj(req)
	logobj := logpkg.LogStruct{"_", now, userIP, "macAdress", "ConfirmedValidatorAPI", "validatorModule", "", "", "_", 0}

	var validValidator validator.TempValidator

	keys, ok := req.URL.Query()["confirmationcode"] // values.Get("confirmationcode") //return parameter from url

	if !ok || len(keys) == 0 {
		globalPkg.SendNotFound(w, "Please check your parameters")
		globalPkg.WriteLog(logobj, "Please check your parameters", "failed")
		return
	}

	//i is the index refer to location of the current confirmed validator in he tempvalidator array
	i := 0
	var flag bool
	for _, Valid := range validator.TempValidatorlst {
		if Valid.ConfirmationCode == keys[0] {
			validValidator = Valid
			flag = true
			break
		}
		i++
	}

	if flag != true {
		globalPkg.SendError(w, "please,check Your verification Code")
		globalPkg.WriteLog(logobj, "please , check your verification code", "failed")
		return
	}
	if now.Sub(validValidator.CurrentTime).Seconds() > globalPkg.GlobalObj.DeleteAccountTimeInseacond {
		globalPkg.SendError(w, "Time out")
		globalPkg.WriteLog(logobj, "Timeout", "failed")
		return

	}
	validator.AddValidator(validValidator.ValidatorObjec)
	//alaa
	M := make(map[string]string)
	M[validValidator.ValidatorObjec.ValidatorIP] = "00000"
	globalfinctiontransaction.SetTransactionIndexTemMap(M)
	broadcastTcp.BoardcastingTCP(validValidator.ValidatorObjec, req.Method, "confirmedvalidator")            // broadcast the validator
	validator.TempValidatorlst = append(validator.TempValidatorlst[:i], validator.TempValidatorlst[i+1:]...) // delete validator from temp list

	ledObj := ledger.GetLedger()
	// ledObj := ledger.GetLedgerForBroadcasting()
	broadcastTcp.SendObject(ledObj, validValidator.ValidatorObjec.ValidatorPublicKey, "POST", "ledger", validValidator.ValidatorObjec.ValidatorSoketIP)
	globalPkg.WriteLog(logobj, "sending success as response", "success")
	globalPkg.SendResponse(w, []byte("Validator addedd successfully"))
}
