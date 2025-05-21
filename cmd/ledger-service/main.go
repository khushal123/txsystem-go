package main

import "github.com/joho/godotenv"

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func run() {

}

func main() {
	run()
}
