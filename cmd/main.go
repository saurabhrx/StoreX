package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"storeX/database"
	"storeX/routes"
)

func main() {
	host := os.Getenv("DB_HOST")
	post := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	if err := database.ConnectToDB(host, post, user, password, databaseName); err != nil {
		logrus.Panicf("failed to connect to database : %+v", err)
	}
	fmt.Println("database connected")

	srv := routes.SetUpStoreXRoutes()

	if srvErr := http.ListenAndServe(":8080", srv); srvErr != nil {
		logrus.Panicf("failed to connect to server %+v", srvErr)
		return
	}
	fmt.Println("server is running on port :")
	if DBCloseErr := database.CloseDBConnection(); DBCloseErr != nil {
		logrus.Panicf("failed to close database %+v", DBCloseErr)
		return
	}

}
