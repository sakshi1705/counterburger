package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/gorilla/handlers"
	"gopkg.in/mgo.v2"
	"net"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// MongoDB Config
// var database_server = "ds227185.mlab.com:27185"
// var database = "counterburger"
// var collection = "menu"
// var mongo_user = "cmpe281"
// var mongo_pass = "cmpe281" 

var database_server = os.Getenv("Server")
var database = os.Getenv("Database")
var collection = os.Getenv("Collection")
var mongo_user = os.Getenv("User")
var mongo_pass = os.Getenv("Pass") 

// var database_server = "13.56.168.122:27017"
// var database = "cb"
// var collection = "menu"
// var mongo_user = "cmpe281"
// var mongo_pass = "cmpe281"

// MenuServer configures and returns a MenuServer instance.
func MenuServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	router := mux.NewRouter()
	initRoutes(router, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD","DELETE","OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders,allowedMethods , allowedOrigins)(router))
	return n
}

// Menu Service API Routes
func initRoutes(router *mux.Router, formatter *render.Render) {
	router.HandleFunc("/menu/ping", pingHandler(formatter)).Methods("GET")
	router.HandleFunc("/menu", GetMenu(formatter)).Methods("GET")
}

// Error Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Menu Serivce Health Check API 
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
        message := "Burger Menu API Server Working on machine: " + getSystemIp()
		formatter.JSON(w, http.StatusOK, struct{ Test string }{message})
	}
}

func getSystemIp() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "" 
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr).String()
    address := strings.Split(localAddr, ":")
    fmt.Println("address: ", address[0])
    return address[0]
}

// API to find an item in the menu
func GetMenu(formatter *render.Render) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("here in Get")
		session, _ := mgo.Dial(database_server)
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		fmt.Println("user pass", mongo_user, mongo_pass)
		err:= session.DB("test").Login(mongo_user, mongo_pass)
		if err!=nil{
			log.Fatalf(" %s", err)
			formatter.JSON(response, http.StatusInternalServerError, "Internal Server Error")
			return
		}
        //session.SetMode(mgo.Monotonic, true) need to check
        mongo_collection := session.DB(database).C(collection)
        var result []Item
		err = mongo_collection.Find(bson.M{}).All(&result)
		fmt.Println("Result: ", result)
        if err != nil {
			log.Fatalf(" %s", err)
            formatter.JSON(response, http.StatusNotFound, "Menu not found !!!")
            return
        }
        fmt.Println("Result: ", result)
		formatter.JSON(response, http.StatusOK, result)
	}
}
