package main

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("./")
	viper.ReadInConfig()
}

func main() {
	// Open routes
	route()

	// Start http server
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Tiger server listening on 127.0.0.1:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %s", err)
	}
}

