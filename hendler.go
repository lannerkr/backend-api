package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("welcome"))

}

func TableHandler(w http.ResponseWriter, r *http.Request) {

	realm := chi.URLParam(r, "realm")
	v := r.URL.Query().Get("update")
	log.Printf("URL Query : %v \n", v)
	var upcheck bool = false
	if v == "true" {
		upcheck = true
	}
	updateNewDBtable(realm, upcheck)
	data := APITable{User: getUserTableNewDB(realm)}

	// var data1user []UserData
	// for i := 0; i < 2; i++ {
	// 	log.Println("data.User: ", data.User[i].Username)
	// 	data1user = append(data1user, data.User[i])
	// }
	// //data1user = append(data1user, data.User...)

	// data1 := APITable{User: data1user}

	// log.Printf("TableHander response data : \n%v", data1)
	respondwithJSON(w, 200, data)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

	realm := chi.URLParam(r, "realm")
	user := chi.URLParam(r, "user")

	data := getSingleUserPCS(realm, user)
	if len(data.Username) < 1 {
		respondwithJSON(w, 404, "USER NOT FOUND")
		return
	}

	data.UserRecord = getActiveUsers(user, "")[0]
	data.DBTable = getSingleUserTableNewDB(realm, user)

	mac := readlog(user)
	data.UserLog = userlog(user, mac)

	respondwithJSON(w, 200, data)
}

func UserStatusHandler(w http.ResponseWriter, r *http.Request) {

	realm := chi.URLParam(r, "realm")
	user := chi.URLParam(r, "user")
	status := chi.URLParam(r, "status")
	var st bool = false
	if status == "true" {
		st = true
	}

	rescode, _ := updateStatus(realm, user, st)
	respondwithJSON(w, 200, rescode)
}

func UserPWresetHandler(w http.ResponseWriter, r *http.Request) {
	realm := chi.URLParam(r, "realm")
	user := chi.URLParam(r, "user")

	rescode := resetPW(realm, user)
	respondwithJSON(w, 200, rescode)
}

func UserPWunlockHandler(w http.ResponseWriter, r *http.Request) {
	realm := chi.URLParam(r, "realm")
	user := chi.URLParam(r, "user")
	var rescode string

	if rescode = deleteUser(realm, user); rescode != "200" {
		respondwithJSON(w, 200, rescode)
		return
	} else if rescode = addUserPulse(realm, user); rescode != "200" {
		respondwithJSON(w, 200, rescode)
		return
	}

	respondwithJSON(w, 200, rescode)
}

func UserOTPresetHandler(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")

	rescode := totpReset(user)
	respondwithJSON(w, 200, rescode)
}

func UserOTPunlockHandler(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")

	rescode := totpUnlock(user)
	respondwithJSON(w, 200, rescode)
}

func UserDisconnectHandler(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "sid")

	rescode := deleteActiveUser(sid)
	respondwithJSON(w, 200, rescode)
}

func ActiveUsersHandler(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")

	rescode := getActiveUsers("all", number)
	respondwithJSON(w, 200, rescode)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	var cluster []ClusterMember
	var sysStatus SystemStatus
	var sysdate, uptime, configdate, licensed, current, maxlast, cpu, swap, disk string
	var err error

	if sysdate, uptime, configdate, cluster, err = getStatus(); err != nil {
		log.Println(err)
	}
	if licensed, current, maxlast, err = getStats(); err != nil {
		log.Println(err)
	}
	if cpu, swap, disk, err = getHealth(); err != nil {
		log.Println(err)
	}

	sysStatus = SystemStatus{sysdate, uptime, configdate, licensed, current, maxlast, cpu, swap, disk, cluster}

	respondwithJSON(w, 200, sysStatus)
}

func ApproveHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	realm := chi.URLParam(r, "realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		respondwithJSON(w, 200, warning)
	} else if realm == "EMP-GOTP" {
		if err := EMPmacApprove(user_mac); err != nil {
			respondwithJSON(w, 200, err.Error())
		} else {
			//inform := "mac-address : " + user_mac + " 승인처리가 완료 되었습니다!"
			respondwithJSON(w, 200, "200")
		}
	} else if err := macApprove(user_mac); err != nil {
		respondwithJSON(w, 200, err.Error())
	} else {
		//inform := "mac-address : " + user_mac + " 승인처리가 완료 되었습니다!"
		respondwithJSON(w, 200, "200")
	}

}
func UnapproveHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	//realm := chi.URLParam(r, "realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		respondwithJSON(w, 200, warning)
	} else if err := macUnApprove(user_mac); err != nil {
		respondwithJSON(w, 200, err.Error())
	}
	//inform := "mac-address : " + user_mac + " 미승인처리 되었습니다!"

	respondwithJSON(w, 200, "200")
}
func PermitHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	//realm := chi.URLParam(r, "realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		respondwithJSON(w, 200, warning)
	} else if err := macPermit(user_mac); err != nil {
		respondwithJSON(w, 200, err.Error())
	}
	//inform := "mac-address : " + user_mac + " USB 사용이 허용 되었습니다!"

	respondwithJSON(w, 200, "200")
}
func ProtectHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")
	//realm := chi.URLParam(r, "realm")
	user_mac := readlog(user_id)
	if user_mac == "" {
		warning := "사용자 MAC 을 찾을 수 없습니다. 사용자 재로그인 후 다시 사용해 주십시오."
		respondwithJSON(w, 200, warning)
	} else if err := macProtect(user_mac); err != nil {
		respondwithJSON(w, 200, err.Error())
	}
	//inform := "mac-address : " + user_mac + " USB 사용이 차단 되었습니다!"

	respondwithJSON(w, 200, "200")
}

type NewUser struct {
	Userid    string `json:"userid"`
	Userrealm string `json:"userrealm"`
	Userip    string `json:"userip"`
	Usermac   string `json:"usermac"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {

	//userid := chi.URLParam(r, "userid")
	var u NewUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if rescode := addUserPulse(u.Userrealm, u.Userid); rescode != "200" {
		respondwithJSON(w, 201, rescode)
		return
	}

	if u.Usermac != "" {
		if rescode := addDevicePPS(u.Userrealm, u.Usermac); rescode != "200" {
			reponse := "User [" + u.Userid + "] is successfully created, BUT " + rescode
			respondwithJSON(w, 201, reponse)
			return
		}
	}

	respondwithJSON(w, http.StatusOK, http.StatusOK)
}
