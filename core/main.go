package main

import (
	"github.com/Aleksao998/LightingUserVault/core/command/root"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	// If the .env file is not present, fallback to using cobra commands
	godotenv.Load()

	// Initialize and execute the root command
	root.NewRootCommand().Execute()
}
