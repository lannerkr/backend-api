package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// var allowOrigin string = "http://192.168.0.157:5173"
//var allowOrigins []string = []string{"http://192.168.0.157:5173", "https://192.168.0.157:3434"}

func routes(router *chi.Mux) {

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(logger))
	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	log.Println(configuration.AllowedOrigins)
	router.Use(cors.Handler(cors.Options{
		//AllowedOrigins: []string{allowOrigin}, // Use this to allow specific origin hosts
		AllowedOrigins: configuration.AllowedOrigins,
		AllowedMethods: []string{"POST", "PUT", "DELETE", "GET"},
		AllowedHeaders: []string{"Origin", "content-type", "Authorization"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))
	//router.Use(methodCheck)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(cAuthenticator(tokenAuth))

		// isAdmin check need
		router.Post("/login/users/register", putUserHandler)
		router.Delete("/login/users/{userid}", delUserHandler)

		// different athority for admin need
		router.Get("/login/users/{userid}", getUsersHandler)
		router.Put("/login/users/{userid}", putUserHandler)

		router.Get("/table/{realm}", TableHandler)

		router.Get("/users/{realm}/{user}/pwreset", UserPWresetHandler)
		router.Get("/users/{realm}/{user}/pwunlock", UserPWunlockHandler)
		router.Get("/users/{realm}/{user}/otpreset", UserOTPresetHandler)
		router.Get("/users/{realm}/{user}/otpunlock", UserOTPunlockHandler)
		router.Get("/users/{realm}/{user}/disconnect/{sid}", UserDisconnectHandler)
		router.Get("/users/{realm}/{user}/{status}", UserStatusHandler)
		router.Get("/users/{realm}/{user}", UserHandler)

		router.Get("/user/approve/{realm}/{user_id}", ApproveHandler)
		router.Get("/user/unapprove/{realm}/{user_id}", UnapproveHandler)
		router.Get("/user/permit/{realm}/{user_id}", PermitHandler)
		router.Get("/user/protect/{realm}/{user_id}", ProtectHandler)

		router.Post("/user/create", createUserHandler)

		router.Get("/activeusers/{number}", ActiveUsersHandler)

		router.Get("/system/status", DashboardHandler)
	})

	// Public routes
	router.Group(func(router chi.Router) {
		router.Get("/", HomeHandler)
		router.Post("/login/users/authenticate", authHandler)

		// router.Put("/login/users/{userid}", putUserHandler)
		// router.Delete("/login/users/{userid}", delUserHandler)

	})
}

func cAuthenticator(*jwtauth.JWTAuth) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				log.Println("cAuth Error: ", err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// exptime, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", fmt.Sprint(rClaim["exp"]))
			// if err != nil {
			// 	log.Println(err)
			// }
			// if time.Until(exptime) <= time.Minute*5 {
			// 	log.Println("exptime: ", exptime)
			// 	token.Expiration().Add(time.Minute * 15)
			// 	updateAuthToken(w, r, token)
			// }

			if token == nil || jwt.Validate(token) != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
