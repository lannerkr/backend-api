package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func pulseReq(realm, method, url string, buff io.Reader) (resp *http.Response, err error) {
	pulseUri := configuration.PulseUri
	pulseApiKey := configuration.PulseApiKey
	var auth string

	switch realm {
	case "store":
		auth = configuration.AuthStore
	case "partner":
		auth = configuration.AuthPartner
	case "emp":
		auth = configuration.AuthEmp
	case "EMP-GOTP":
		auth = configuration.AuthEmp
	default:
		err := fmt.Errorf("realm %v is not available", err)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := pulseApiKey
	pUri := pulseUri + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + url

	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func pulseSysReq(method, url string, buff io.Reader) (resp *http.Response, err error) {
	pulseUri := configuration.PulseUri
	pulseApiKey := configuration.PulseApiKey

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := pulseApiKey

	pUri := pulseUri + "/api/v1/" + url

	req, err := http.NewRequest(method, pUri, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(apikey, "")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getUsersPCS(realm string) (pulseUsers []UserData) {

	url := "/local/users"

	resp, err := pulseReq(realm, "GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := PulseTable{}
	if err := decoder.Decode(&user); err != nil {
		log.Println(err)
	}
	pulseUsers = user.User

	return pulseUsers
}

func getSingleUserPCS(realm, user string) (pulseUser UserData) {
	//fmt.Println("getSingleUserPCS")

	url := "/local/users/user/" + user

	resp, err := pulseReq(realm, "GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&pulseUser); err != nil {
		log.Println(err)
	}

	return pulseUser
}

func getActiveUsers(user, count string) (record []Records) {
	//fmt.Println("getActiveUsers")

	var url string
	if user == "all" {
		url = "system/active-users?number=" + count
	} else {
		url = "system/active-users?name=" + user
	}

	resp, err := pulseSysReq("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	respRecords := Response{}
	if err := decoder.Decode(&respRecords); err != nil {
		log.Println(err)
	}
	record = respRecords.ActiveUsers.ActiveUserRecord.UserRecord

	if len(record) == 0 {
		record = append(record, Records{"", "", "", "", "", "", ""})
	}
	return record
}

func updateStatus(realm, user string, status bool) (int, error) {

	st := strconv.FormatBool(status)
	stup := map[string]string{
		"enabled": st,
	}
	pbytes, _ := json.Marshal(stup)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/" + user + "/enabled"
	resp, err := pulseReq(realm, "PUT", url, buff)
	if err != nil {
		log.Println(err)
		return 404, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("user: " + user + " enabled status is changed to " + st)
		return 200, nil
	} else {
		err := fmt.Errorf("user: " + user + " status change is failed")
		return 404, err
	}
}

func resetPW(realm, user string) string {

	updata := map[string]string{
		"password-cleartext":        "Fashion2022!",
		"change-password-at-signin": "true",
	}
	pbytes, _ := json.Marshal(updata)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/" + user + "/password-cleartext"
	resp, err := pulseReq(realm, "PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("user: " + user + " password reset is done.")
		return "200"
	} else {
		err := fmt.Errorf("user: " + user + " password reset is failed")
		return fmt.Sprintf("%v : %v", resp.StatusCode, err.Error())
	}
}

func addUserPulse(realm, user string) string {
	newuser := map[string]string{
		"enabled":                   "true",
		"change-password-at-signin": "true",
		"fullname":                  user,
		"password-cleartext":        "Fashion2022!",
		"username":                  user,
	}
	pbytes, _ := json.Marshal(newuser)
	buff := bytes.NewBuffer(pbytes)

	url := "/local/users/user/"

	resp, err := pulseReq(realm, "POST", url, buff)
	if err != nil {
		log.Println(err)
		return "addeUser Error: " + err.Error()
	}
	defer resp.Body.Close()

	respStrings := respString(resp, "POST")
	log.Printf("PCS response : %v\n", respStrings)

	if resp.StatusCode == 201 {
		if strings.Contains(respStrings, PostOK) {
			log.Printf("user: %v has been added successfully\n", user)
			return "200"
		} else if strings.Contains(respStrings, PostExist) {
			return fmt.Sprint("user: " + user + " already exist")
		} else {
			return fmt.Sprint("user: " + user + " status change is failed")
		}
	}
	return fmt.Sprint(resp.StatusCode, ": addeUser Error")
}

func deleteUser(realm, userid string) string {
	url := "/local/users/user/" + userid
	resp, err := pulseReq(realm, "DELETE", url, nil)
	if err != nil {
		log.Println("deleteUser Error: ", err)
		return "deleteUser Error: " + err.Error()
	} else if resp.StatusCode == 202 || resp.StatusCode == 200 || resp.StatusCode == 204 {
		defer resp.Body.Close()
		respGet, errGet := pulseReq(realm, "GET", url, nil)
		if respGet.StatusCode == 404 {
			defer respGet.Body.Close()
			return "200"
		}
		defer respGet.Body.Close()
		log.Println("deleteUser Error: ", errGet)
		return "deleteUser Error: " + errGet.Error()
	}
	defer resp.Body.Close()

	return fmt.Sprint(resp.StatusCode, ": deleteUser Error")
}

func totpReset(user string) string {
	totp := configuration.TotpServer
	url := "totp/" + totp + "/users/" + user + "?operation=reset"
	resp, err := pulseSysReq("PUT", url, nil)
	if err != nil {
		log.Println(err)
		return "totpReset Error: " + err.Error()
	}
	defer resp.Body.Close()

	respStrings := respString(resp, "OTP")
	if resp.StatusCode == 200 && strings.Contains(respStrings, "has been reset") {
		log.Println(respStrings)
		return "200"
	}

	return "totpReset Error: " + respStrings
}

func totpUnlock(user string) string {
	totp := configuration.TotpServer
	url := "totp/" + totp + "/users/" + user + "?operation=unlock"
	resp, err := pulseSysReq("PUT", url, nil)
	if err != nil {
		log.Println(err)
		return "totpUnlock Error: " + err.Error()
	}
	defer resp.Body.Close()
	respStrings := respString(resp, "OTP")
	if resp.StatusCode == 200 && strings.Contains(respStrings, "has been unlocked") {
		log.Println(respStrings)
		return "200"
	}

	return "totpUnlock Error: " + respStrings
}

func deleteActiveUser(sid string) string {
	url := "system/active-users/session/" + sid
	resp, err := pulseSysReq("DELETE", url, nil)
	if err != nil {
		log.Println(err)
		return "del ActivceUser Error: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Sprint(resp.StatusCode) + ": del ActivceUser Error"
	}

	return "200"
}

// func userCreate(realm, user, ip, mac string) (err1, err2, err3 error) {
// 	if user != "" && realm != "" {
// 		if err := addUserPulse(realm, user); err != nil {
// 			err1 = err
// 		}
// 	}

// 	if ip != "" && realm != "emp" {
// 		if err := addUserPPS(realm, user, ip); err != nil {
// 			err2 = err
// 		}
// 	} else {
// 		err2 = fmt.Errorf("static-ip is empty")
// 	}

// 	if mac != "" {
// 		if err := addDevicePPS(realm, mac); err != nil {
// 			err3 = err
// 		}

// 	} else {
// 		err3 = fmt.Errorf("mac-address is empty")
// 	}
// 	return err1, err2, err3
// }
