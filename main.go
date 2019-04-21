package main

import (
	"os"
	"os/signal"
)

func main() {
	configFile := "config.json"
	if !fileExists(configFile) {
		err := SaveConfigToFile(configFile, new(Config))
		if err != nil {
			panic(err)
		}
		return
	}
	config, err := LoadConfigFromFile(configFile)
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	interrupt2 := make(chan struct{})
	go func() {
		<-interrupt
		close(interrupt2)
	}()

	server := NewPublicAddressServer(config.PublicAddressName, config.Email, config.Password, config.Port, config.CallbackURL, config.TokenStorageEndpoint, config.MessageCallbackEndpoint)
	err = server.Run(interrupt2)

	if err != nil {
		panic(err)
	}
}
