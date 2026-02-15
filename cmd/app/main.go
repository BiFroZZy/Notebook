package main

import (
	"log"

	"github.com/joho/godotenv"

	"notebook/cmd/handlers"
	_ "notebook/docs"
)

// @Title Notebook documentaion (Документация по программе)
// @Description Записки для пользователей
// @version 1.0
// @Host localhost
// @BasePath 

func main(){
	if err := godotenv.Load(); err != nil{
		log.Fatalf("Can't load env file: %v\n", err)
	}
	handlers.Handlers()
}