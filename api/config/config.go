package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"
)

type ConfigT struct {
	Host           string `ini:"host"`
	TokenKey       string `ini:"token_key"`
	TokenName      string `ini:"token_name"`
	TokenExpTime   int64  `ini:"token_exp_time"`
	AvatarSavePath string `ini:"avatar_save_path"`
	DBSavePath string `ini:"db_save_path"`
	IsDevString string `ini:"dev"`
	IsDev bool
	TokenKeyByte   []byte
}
type Client struct {
	Name         string `json:"name"`
	Callback     string `json:"callback"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type ClientsT struct {
	Client []Client `json:"clients"`
}

var Config ConfigT
var Clients ClientsT
var ClientsMap map[string]Client

func ReadConfig() {
	iniconf, err := ini.Load("config.ini")
	if err != nil {
		log.Panicf("Cannot open the config.ini. Error: %v", err.Error())
	}
	config_section := iniconf.Section("")
	err = config_section.MapTo(&Config)
	if err != nil {
		fmt.Println(err)
		log.Panic("Cannot load config.ini")
	}
	Config.IsDev=Config.IsDevString=="true"
	Config.TokenKeyByte = []byte(Config.TokenKey)
	err = os.MkdirAll(Config.AvatarSavePath, os.ModePerm)
	if err != nil {
		log.Panicf("Cannot make dir for avatar. Error: %v", err.Error())
	}
	file, err := os.Open("oath2-conf.json")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		log.Panic("Cannot read oath2-conf.json")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Clients)
	if err != nil {
		fmt.Println(err)
		log.Panic("Cannot read oath2-conf.json")
	}
	ClientsMap = make(map[string]Client)
	for _, value := range Clients.Client {
		ClientsMap[value.ClientID] = value
	}
}
