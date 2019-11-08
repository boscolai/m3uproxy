package main

import (
	"log"
	"net/http"
	"os"

	"github.com/boscolai/m3uproxy/config"
	"github.com/boscolai/m3uproxy/handler"
	"github.com/spf13/viper"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalln("no config file provided")
	}
	cfgFile := os.Args[1]
	if err := config.SetupViper(cfgFile); err != nil {
		log.Fatalf("failed to initialize configuration using: %s: %s", cfgFile, err)
	}
	listenAddr := viper.GetString("server.bind_address")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	handler.InitHandlers(viper.GetViper())

	http.HandleFunc("/m3u", handler.HandleM3U)
	http.HandleFunc("/epg", handler.HandleEPG)
	http.HandleFunc("/m3u-groups", handler.HandleGroups)
	http.ListenAndServe(listenAddr, nil)
}
