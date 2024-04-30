package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// func catch(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// respondwithError return error message
// func respondWithError(w http.ResponseWriter, code int, msg string) {
// 	respondwithJSON(w, code, map[string]string{"message": msg})
// }

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Logger return log message
// func Logger() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println(time.Now(), r.Method, r.URL)
// 		router.ServeHTTP(w, r) // dispatch the request
// 	})
// }

type PostResp struct {
	Results Result `json:"result"`
}
type Result struct {
	Warnings []Warning `json:"warnings"`
}
type Warning struct {
	Messages string `json:"message"`
}

const (
	PostOK    string = "The configuration has been implicitly changed"
	PostExist string = "already exist"
)

func respString(resp *http.Response, method string) (respString string) {

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	jsonResp := strings.NewReader(bodyString)

	if strings.Contains(bodyString, "admin_op") {
		type DeviceResp struct {
			Status  string `json:"status"`
			Mac     string `json:"macaddr"`
			AdminOp string `json:"admin_op"`
		}
		var postResp DeviceResp
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			log.Println(err)
		}
		respString = "mac-address: " + postResp.Mac + " , status: " + postResp.Status + " , operation: " + postResp.AdminOp
		return respString
	}

	if method == "POST" {
		var postResp PostResp
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			log.Println(err)
		}
		respString = postResp.Results.Warnings[0].Messages
	} else if method == "PPS" {
		var postResp Warning
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			log.Println(err)
		}
		respString = postResp.Messages
	} else if method == "OTP" {
		type TypeError struct {
			Messages string `json:"message"`
		}
		type Info struct {
			Messages string `json:"message"`
		}
		type Result struct {
			ErrorsR []TypeError `json:"errors"`
			Infos   []Info      `json:"info"`
		}
		type OTPresp struct {
			Results Result `json:"result"`
		}

		var postResp OTPresp
		if err := json.NewDecoder(jsonResp).Decode(&postResp); err != nil {
			log.Println(err)
		}
		if postResp.Results.Infos != nil {
			respString = postResp.Results.Infos[0].Messages
		} else if postResp.Results.ErrorsR != nil {
			respString = postResp.Results.ErrorsR[0].Messages
		}
	}

	return respString
}

// // Chi by default doesnt show message for 405 error text
// // This middleware checks the gloobal router context and
// // displays an error message of "Method not Allowed"
// func methodCheck(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// get the context from the given router
// 		rctx := chi.RouteContext(r.Context())

// 		// Temporary context
// 		ctx := chi.NewRouteContext()
// 		// Matching the router context with the current request.
// 		// If there is no method matching for the current
// 		// request route, throw method not allowed (405) err.
// 		if !rctx.Routes.Match(ctx, r.Method, r.URL.Path) {
// 			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 			return
// 		}
// 		next.ServeHTTP(w, r.WithContext(r.Context()))
// 	})
// }

// func hashresult(input string) string {
// 	h := sha256.New()
// 	h.Write([]byte(input))
// 	hashValue := h.Sum(nil)
// 	//fmt.Printf("SHA-256 hash of '%s' is: %s\n", input, hex.EncodeToString(hashValue))
// 	return hex.EncodeToString(hashValue)
// }

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
