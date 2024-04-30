package main

import (
	"bufio"
	"os/exec"
	"strings"
)

func readlog(user string) (mac string) {

	apilog := configuration.ApiLog
	readCmd := "tac " + apilog + " | grep PRO31457 | grep -m 1 " + user
	cmd, output := exec.Command("sh", "-c", readCmd), new(strings.Builder)
	cmd.Stdout = output
	cmd.Run()
	userProfile := strings.Split(output.String(), ",")

	if len(userProfile) >= 4 {
		device_attN1 := strings.Split(userProfile[3], ")")
		device_attN2 := strings.Split(device_attN1[0], "(")
		mac = device_attN2[1]
		mac = strings.Replace(mac, "-", ":", -1)
	}

	return mac
}

func userlog(user, mac string) (userLog []UserLogs) {
	apilog := configuration.ApiLog
	//readCmd := "tac " + apilog + " | grep -i -m 15 -E \"" + user + "|ADM31591.*" + mac + "\" | awk '{$1=\"\";$2=\"\";$3=\"\";$4=\"\";$5=\"\";print}'"
	readCmd := "tac " + apilog + " | grep -i -m 15 -E \"" + user + "|ADM31591.*" + mac + "\" | awk '{$1=\"\";$2=\"\";$3=\"\";print}'"
	cmd, output := exec.Command("sh", "-c", readCmd), new(strings.Builder)
	cmd.Stdout = output
	cmd.Run()

	logs := output.String()
	return convertLogToStruct(logs)
}

type UserLogs struct {
	Date string `json:"date"`
	Code string `json:"code"`
	Log  string `json:"log"`
}

func convertLogToStruct(log string) (userLog []UserLogs) {

	scanner := bufio.NewScanner(strings.NewReader(log))
	for scanner.Scan() {
		logline := scanner.Text()
		var uLog UserLogs
		logN := strings.Split(logline, "[HIDE]")
		if len(logN) > 1 {
			logN = strings.Split(logN[1], " - [API-LOG]")
			logline = logN[0]
		} else {
			logN = strings.Split(logN[0], " - [API-LOG]")
			logline = logN[0]
		}

		if strings.Contains(logline, " - ") {
			logN = strings.SplitN(logline, " - ", 3)
		} else {
			logN = strings.SplitN(logline, ",", 3)
		}

		if len(logN) < 3 {
			uLog.Log = logN[len(logN)-1]
		} else {
			uLog.Date = logN[0]
			uLog.Code = logN[1]
			uLog.Log = logN[2]
		}
		userLog = append(userLog, uLog)
	}

	return userLog
}
