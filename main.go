package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

// version 3.0.2 20240429
// version 3.1 20240514

var logger *httplog.Logger
var fpLog *os.File

func main() {

	conf := os.Args[1]
	sport, cert, certkey, logpath := Config(conf)
	//_, _, _, logpath := Config(conf)

	var err error
	fpLog, err = os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	dbconn = mongoConnect()
	defer func() {
		if err := dbconn.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()

	//Logger
	logger = httplog.NewLogger("httplog", httplog.Options{
		// JSON:             true,
		LogLevel: slog.LevelInfo,
		Concise:  true,
		//RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
		Writer: fpLog,
	})

	router := chi.NewRouter()
	routes(router)

	//http.ListenAndServe(":"+configuration.ServerPort, router)
	if err := http.ListenAndServeTLS(":"+sport, cert, certkey, router); err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
