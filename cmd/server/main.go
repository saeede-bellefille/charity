package main

import (
	"SABKAD/pkg/handler"
	"SABKAD/pkg/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	db := createConnection()
	//db.Migrator().DropTable(&model.CommonBaseType{}, &model.CommonBaseData{}, &model.CharityAccount{}, &model.Personal{}, &model.NeedyAccount{}, &model.Plan{}, &model.AssignNeedyToPlan{}, &model.CashAssistanceDetail{}, &model.Payment{})
	db.AutoMigrate(&model.CommonBaseType{}, &model.CommonBaseData{}, &model.CharityAccount{}, &model.Personal{}, &model.NeedyAccount{}, &model.Plan{}, &model.AssignNeedyToPlan{}, &model.CashAssistanceDetail{}, &model.Payment{})
	//defer db.Close()

	r := handler.Router(db)

	server := http.Server{
		Addr:         ":9091",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		fmt.Println("Starting server on port 9091")

		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}

func createConnection() *gorm.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	path := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open(postgres.Open(path), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}
	database, err := db.DB()

	err = database.Ping()

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected!")

	return db

}
