package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/notrishabh/tuivia/tui"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tui.RunTui()
}
