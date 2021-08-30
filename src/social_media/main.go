package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"social_media/controller"
	"social_media/util"
	"social_media/wrapper"
)

func main() {
	// log to log file AND stdout
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// create database connection
	database, err := util.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login/", wrapper.Log(controller.Login(database)))
	http.HandleFunc("/logout/", wrapper.Log(controller.Logout(database)))
	http.HandleFunc("/posts/", wrapper.Log(controller.Posts(database)))
	http.HandleFunc("/profile/", wrapper.Log(controller.Profile(database)))
	http.HandleFunc("/media/", wrapper.Log(controller.Media(database)))
	http.HandleFunc("/admin/", wrapper.Log(controller.Admin(database)))
	http.HandleFunc("/", wrapper.Log(controller.Home(database)))
	http.HandleFunc("/robots.txt", wrapper.Log(controller.Robots()))
	http.HandleFunc("/static/", wrapper.Log(controller.Static()))

	log.Println("Server started ...")
	log.Fatal(http.ListenAndServe(":9092", nil))
}
