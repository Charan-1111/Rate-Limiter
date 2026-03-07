package main

import (
	"fmt"

	"goapp/server"
)

func main() {
	filePath := "manifest/config.json"
	app, err := server.NewApplication(filePath)
	if err != nil {
		fmt.Println("Error creating the application, retrying...")
		return
	}

	err = app.StartServer()
	if err != nil {
		
	}
}