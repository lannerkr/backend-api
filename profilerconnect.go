package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func ppsReq(realm, method, url string, buff io.Reader) (resp *http.Response, err error) {
	ppsUri := configuration.PPSUri
	ppsApiKey := configuration.PPSApiKey
	var auth string

	switch realm {
	case "store":
		auth = configuration.PPSAuthStore
	case "partner":
		auth = configuration.PPSAuthPartner
	case "emp":
		auth = configuration.PPSAuthEmp
	case "EMP-GOTP":
		auth = configuration.PPSAuthEmp
	default:
		err := fmt.Errorf("realm %v is not available", realm)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := ppsApiKey
	pUri := ppsUri + "/api/v1/configuration/authentication/auth-servers/auth-server/" + auth + url

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

func ppsSysReq(method, url string, buff io.Reader) (resp *http.Response, err error) {
	ppsUri := configuration.PPSUri
	ppsApiKey := configuration.PPSApiKey

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apikey := ppsApiKey
	pUri := ppsUri + "/api/v1/" + url

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

func macApprove(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	setApprove := map[string]string{
		"status": "approved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Printf("macApprove response: %v\n", bodyString)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		log.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macUnApprove(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	setApprove := map[string]string{
		"status": "unapproved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		log.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macPermit(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	setApprove := map[string]string{
		"category": "permit",
		"override": "true",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		log.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func macProtect(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	setApprove := map[string]string{
		"category": "Windows",
		"override": "false",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		log.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func EMPmacApprove(mac string) error {

	url := "profiler/endpoints/simplified/" + mac
	setApprove := map[string]string{
		"notes":  "emp",
		"status": "approved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("PUT", url, buff)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Successfully updated") {
		return nil
	} else {
		bodyString = strings.ReplaceAll(bodyString, "{", "")
		bodyString = strings.ReplaceAll(bodyString, "}", "")
		log.Println(bodyString)
		return fmt.Errorf(bodyString)
	}
}

func addDevicePPS(realm, mac string) string {

	url := "profiler/endpoints?=addNewDevice"
	if realm == "EMP-GOTP" {
		realm = "emp"
	}
	setApprove := map[string]string{
		"macaddr":  mac,
		"category": "Windows",
		"notes":    realm,
		"status":   "approved",
	}
	pbytes, _ := json.Marshal(setApprove)
	buff := bytes.NewBuffer(pbytes)

	resp, err := ppsSysReq("POST", url, buff)
	if err != nil {
		return fmt.Sprint(err)

	}
	defer resp.Body.Close()

	respStrings := respString(resp, "PPS")
	if resp.StatusCode == 201 {
		if strings.Contains(respStrings, "operation: add") {
			// return fmt.Sprint("%v : mac-address for user : %v has been added successfully\n", resp.StatusCode, mac)
			return "200"
		} else {
			return fmt.Sprintf("mac-address: " + mac + " adding is failed")
		}
	} else if resp.StatusCode == 409 {
		return fmt.Sprintf("mac-address: " + mac + " already exist")
	}

	return respStrings //strconv.Itoa(resp.StatusCode)
}
