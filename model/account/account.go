package account

import (
	"fmt"
	"model/cryptogrpghy"
	"regexp"
	"strings"
	"time"

	error "../errorpk"
	"../validator"

	//globalPkg "github.com/tensor-programming/mining/../globalPkg"
	"../accountdb"
)

type errorMess struct {
	ErrorMessage string
}

type User struct {
	Account           accountdb.AccountStruct
	Oldpassword       string
	Reson             string
	CurrentTime       time.Time
	Confirmation_code string
	TextSearch        string
	Method            string
	PathApi           string
}

type searchResponse struct {
	UserName  string
	PublicKey string
}

/*----------------- function to get Account by public key -----------------*/
func GetAccountByAccountPubicKey(AccountPublicKey string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountPublicKey(AccountPublicKey)
}

/*----------------- function to check if account exists or not -----------------*/
func ifAccountExistsBefore(AccountPublicKey string) bool {
	if (accountdb.FindAccountByAccountPublicKey(AccountPublicKey)).AccountPublicKey == "" {
		return false //not exist
	}
	return true
}

/*----------------- function to save an account on json file -----------------*/
func AddAccount(accountObj accountdb.AccountStruct) string {

	if !(ifAccountExistsBefore(accountObj.AccountPublicKey)) && validateAccount(accountObj) {
		if accountdb.AccountCreate(accountObj) {
			return ""
		} else {
			return error.AddError("AddAccount account package", "Check your path or object to Add AccountStruct", "logical error")
		}
	}
	return error.AddError("AddAccount account package", "The account is already exists "+accountObj.AccountPublicKey, "hack error")

}

func getLastIndex() string {

	var Account accountdb.AccountStruct
	Account = accountdb.GetLastAccount()
	//if Account.AccountPublicKey == "" {
	//	return "-1"
	//}
	if Account.AccountName == "" {
		return "-1"
	}

	return Account.AccountIndex

}

/*----------------- function to update an account on json file -----------------*/
func UpdateAccount(accountObj accountdb.AccountStruct) string {
	if (ifAccountExistsBefore(accountObj.AccountPublicKey)) && validateAccount(accountObj) {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++")
		if accountdb.AccountUpdateUsingTmp(accountObj) {
			return ""
		} else {
			return error.AddError("UpdateAccount account package", "Check your path or object to Update AccountStruct", "logical error")

		}
		fmt.Println("iam update")
	}
	return error.AddError("FindjsonFile account package", "Can't find the account obj "+accountObj.AccountPublicKey, "hack error")

}

/*----------------- function to validate the account before register  -----------------*/
func validateAccount(accountObj accountdb.AccountStruct) bool {
	if len(accountObj.AccountName) < 8 || len(accountObj.AccountName) > 30 || (len(accountObj.AccountPassword) != 64) || len(accountObj.AccountAddress) < 5 || len(accountObj.AccountAddress) > 100 {
		fmt.Println("1")
		return false
	}

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(accountObj.AccountEmail) && accountObj.AccountEmail != "" {
		fmt.Println("2")
		return false
	}
	return true
}

/////////////////*******func Add by Aya to get publickey using Any string ----

func getPublicKeyUsingString(Key string) string {

	existingAccountUsingName := accountdb.FindAccountByAccountName(Key)

	existingAccountusingEmail := accountdb.FindAccountByAccountEmail(Key)
	existingAccountusingPhoneNumber := accountdb.FindAccountByAccountPhoneNumber(Key)
	if existingAccountUsingName.AccountPublicKey != "" {
		return existingAccountUsingName.AccountPublicKey
	}
	if existingAccountusingEmail.AccountPublicKey != "" {
		return existingAccountusingEmail.AccountPublicKey
	}
	if existingAccountusingPhoneNumber.AccountPublicKey != "" {
		return existingAccountusingPhoneNumber.AccountPublicKey
	}

	return ""

}

/*-------------FUNCTION TO CHECK Acoount----*/

func checkAccount(userAccountObj accountdb.AccountStruct) string {

	existingAccountUsingName := accountdb.FindAccountByAccountName(userAccountObj.AccountName)
	existingAccountusingEmail := accountdb.FindAccountByAccountEmail(userAccountObj.AccountEmail)
	existingAccountusingPhoneNumber := accountdb.FindAccountByAccountPhoneNumber(userAccountObj.AccountPhoneNumber)

	if existingAccountUsingName.AccountPublicKey != "" && existingAccountUsingName.AccountPublicKey != userAccountObj.AccountPublicKey {

		return "UserName Found"
	}
	if existingAccountusingEmail.AccountPublicKey != "" && existingAccountusingEmail.AccountPublicKey != userAccountObj.AccountPublicKey && userAccountObj.AccountEmail != "" {
		fmt.Println("Email found", existingAccountusingEmail.AccountEmail, "  ", userAccountObj.AccountEmail)
		return "Email found"
	}
	if existingAccountusingPhoneNumber.AccountPublicKey != "" && existingAccountusingPhoneNumber.AccountPublicKey != userAccountObj.AccountPublicKey && userAccountObj.AccountPhoneNumber != "" {
		return "Phone Found "
	}
	return ""
}

func checkAccountbeforeRegister(userAccountObj accountdb.AccountStruct) string {

	existingAccountUsingName := accountdb.FindAccountByAccountName(userAccountObj.AccountName)
	existingAccountusingEmail := accountdb.FindAccountByAccountEmail(userAccountObj.AccountEmail)
	existingAccountusingPhoneNumber := accountdb.FindAccountByAccountPhoneNumber(userAccountObj.AccountPhoneNumber)

	if existingAccountUsingName.AccountName != "" && existingAccountUsingName.AccountName == userAccountObj.AccountName {

		return "UserName Found"
	}
	if existingAccountusingEmail.AccountEmail != "" && existingAccountusingEmail.AccountEmail == userAccountObj.AccountEmail && userAccountObj.AccountEmail != "" {
		fmt.Println("Email found", existingAccountusingEmail.AccountEmail, "  ", userAccountObj.AccountEmail)
		return "Email found"
	}
	if existingAccountusingPhoneNumber.AccountPhoneNumber != "" && existingAccountusingPhoneNumber.AccountPhoneNumber == userAccountObj.AccountPhoneNumber && userAccountObj.AccountPhoneNumber != "" {
		return "Phone Found "
	}
	return ""
}

/*-------------integrat account module with miner------*/
/*-------------check Befor Add------*/

func checkingIfAccountExixtsBeforeRegister(accountObj accountdb.AccountStruct) string {

	/*if IfAccountExistsBefore(accountObj.AccountPublicKey) {
		return "public key exists before"
	}*/
	Error := checkAccountbeforeRegister(accountObj)
	if Error != "" {
		return Error

	}

	if !validateAccount(accountObj) {
		return "please, check your data"
	}
	return ""

}
func checkingIfAccountExixtsBeforeAdd(accountObj accountdb.AccountStruct) string {

	//IfAccountExistsBefore(accountObj.AccountPublicKey)
	if ifAccountExistsBefore(accountObj.AccountPublicKey) {
		return "the publick key exist before "
	}

	return ""

}

/*-------------check Befor update------*/
func checkingIfAccountExixtsBeforeUpdating(accountObj accountdb.AccountStruct) string {
	if !(ifAccountExistsBefore(accountObj.AccountPublicKey)) {
		return "Please Check Your data to help me to find your account"
	}
	s := checkAccount(accountObj)
	if s != "" {
		return s
	}
	if !validateAccount(accountObj) {
		return "please, check your data"
	}
	return ""

}

/*-------------getAccountPassword-------*/
func getAccountPassword(AccountPublicKey string) string {
	return accountdb.FindAccountByAccountKey(AccountPublicKey).AccountPassword
}

/*-------------get AccountStruct using publicKey-------*/
// func getAccountByPublicKey(AccountPublicKey string) AccountStruct {
// 	return findAccountByAccountPublicKey(AccountPublicKey)
// }

// func GetAccountByPublicKey(AccountPublicKey string) AccountStruct {
// 	return findAccountByAccountPublicKey(AccountPublicKey)
// }

/*-------------get AccountStruct using email-------*/
func getAccountByEmail(AccountEmail string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountEmail(AccountEmail)
}
func getAccountByPhone(AccountPhoneNumber string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountPhoneNumber(AccountPhoneNumber)
}

/*-------------get AccountStruct using user name-------*/
func getAccountByName(AccountName string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountName(AccountName)
}

/*-------------get AccountStruct using user name-------*/
func GetAccountByName(AccountName string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountName(AccountName)
}

/*----------------AddBlockToAnaccount-----*/
func AddBlockToAccount(AccountPublicKey string, blockIndex string, tokenID string) {
	accountObj := accountdb.FindAccountByAccountPublicKey(AccountPublicKey)
	hashedIndex := cryptogrpghy.KeyEncrypt(validator.CurrentValidator.ValidatorPrivateKey, blockIndex)
	accountObj.BlocksLst = append(accountObj.BlocksLst, hashedIndex)

	containid := ContainstokenID(accountObj.AccountTokenID, tokenID)
	if !containid {
		accountObj.AccountTokenID = append(accountObj.AccountTokenID, tokenID)
	}
	UpdateAccount2(accountObj)
}

//ContainstokenID Contains tells whether a contains x.
func ContainstokenID(AccountTokenID []string, tokenid string) bool {
	for _, n := range AccountTokenID {
		if tokenid == n {
			return true
		}
	}
	return false
}

func GetAccountByIndex(index string) accountdb.AccountStruct {
	return accountdb.FindAccountByAccountKey(index)
}

func UpdateAccount2(accountObj accountdb.AccountStruct) string {
	if (ifAccountExistsBefore(accountObj.AccountPublicKey)) && validateAccount(accountObj) {
		if accountdb.AccountUpdate2(accountObj) {
			return ""
		} else {
			return error.AddError("UpdateAccount account package", "Check your path or object to Update AccountStruct", "logical error")

		}
		fmt.Println("iam update22")
	}
	return error.AddError("FindjsonFile account package", "Can't find the account obj "+accountObj.AccountPublicKey, "hack error")

}

//SetPublicKey update the public key into the database
func SetPublicKey(accountObjc accountdb.AccountStruct) {
	if accountdb.AddBKey(accountObjc) {
		fmt.Println("public key added successfully")
	} else {
		fmt.Println("failed to add public key")
	}
}

//convertStringTolowerCaseAndtrimspace approve username , email is lowercase and trim spaces
func convertStringTolowerCaseAndtrimspace(stringphrase string) string {
	stringphrase = strings.ToLower(stringphrase)
	stringphrase = strings.TrimSpace(stringphrase)
	return stringphrase
}

//----------to convert User Object to Account Object----
//convertUserTOAccount Deleted
