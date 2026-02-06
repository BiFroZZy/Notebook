package main

import (
	"log"
	"github.com/joho/godotenv"
	h "notebook/cmd/handlers"
)
func main(){
	if err := godotenv.Load(); err != nil{
		log.Fatalf("Can't load env file: %v\n", err)
	}
	h.Handlers()
}