package main

import (
	"net/http"
	"fmt"
	"time"
	b64 "encoding/base64"
	"strings"
	"encoding/json"
)

type responseSimple struct {
	Success bool
	Data map[string]string
}

func login (w http.ResponseWriter, req *http.Request) []byte {
	response := (responseSimple{Success:false})
	account := req.FormValue("account")
	fmt.Println(account)
	
	DBhandler := initDB()
	var args []interface{}
	args = append(args, account)
	success, resRow := DBhandler.selectRow("SELECT password FROM users where account = ?", args,1)
	DBhandler.closeDB()
	if success == false {
		json, err := json.Marshal(response)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte("something is wrong"+err.Error()))
			return nil
		}
		return json
	}
	if resRow[0]!=req.FormValue("password"){
		json, err := json.Marshal(response)
		if err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte("something is wrong"+err.Error()))
			return nil
		}
		return json
	}
	encodeValue := "mkhashsucc"+account+"incredablehashsalt"
	cookieValue := b64.StdEncoding.EncodeToString([]byte(encodeValue))
	expire := time.Now().Add(24 * time.Hour)
	credential := http.Cookie{Name:"mkgclogincookie",Value:cookieValue,Expires:expire}
	credential.HttpOnly = true
	http.SetCookie(w, &credential )

	response.Data = map[string]string{"account":account}
	response.Success = true
	json, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("something is wrong"+err.Error()))
		return nil
	}
	return []byte(json)
}

func loginCheck (w http.ResponseWriter, req *http.Request) (bool, string){
	loginCookie, err := req.Cookie("mkgclogincookie")
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("something is wrong"+err.Error()))
		return false, ""
	}
	decodeValue, err := b64.StdEncoding.DecodeString(loginCookie.Value)	
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("something is wrong cookie decode failed"))
		return false, ""
	}
	decodeString := string(decodeValue)
	if decodeString[0:10] != "mkhashsucc" {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("something is wrong cookie is not clean"))
		return false, ""
	}
	account := decodeString[10:strings.Index(decodeString, "incredablehashsalt")]
	fmt.Println(account)
	return true, account
}

func logout (w http.ResponseWriter, req *http.Request) []byte {
	response := (responseSimple{Success:false})
	credential := http.Cookie{Name:"mkgclogincookie",MaxAge:-1}
	http.SetCookie(w, &credential )
	response.Success = true
	json, err := json.Marshal(response)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("something is wrong"+err.Error()))
		return nil
	}
	return []byte(json)
}
