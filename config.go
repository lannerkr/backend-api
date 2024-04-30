package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	TotpServer     string
	PulseUri       string
	PulseApiKey    string
	AuthStore      string
	AuthPartner    string
	AuthEmp        string
	PPSUri         string
	PPSApiKey      string
	PPSAuthStore   string
	PPSAuthPartner string
	PPSAuthEmp     string
	DBUri          string
	ApiLog         string
	MongoDBS       string
	Cert           string
	CertKey        string
	IPaccess       []string
	LogPath        string
	ServerPort     string
	SRC_dir        string
	AllowedOrigins []string
}

var configuration Configuration

func Config(conf string) (sport, cert, certkey, logpath string) {
	file, _ := os.Open(conf)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(configuration)

	return configuration.ServerPort, configuration.Cert, configuration.CertKey, configuration.LogPath
}
