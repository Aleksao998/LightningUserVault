package main

import (
	"github.com/Aleksao998/LightingUserVault/core/command/root"
	"github.com/joho/godotenv"
)

func main() {
	// if .env does not exist, cobra commands will be used
	godotenv.Load()

	root.NewRootCommand().Execute()
}
