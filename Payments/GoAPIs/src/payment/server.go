package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	mongo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongodb_server = os.Getenv("AWS_MONGODB")
var mongodb_database = os.Getenv("MONGODB_DBNAME")
var mongodb_collection = os.Getenv("MONGODB_COLLECTION")
var mongodb_username = os.Getenv("MONGODB_USERNAME")
var mongodb_password = os.Getenv("MONGODB_PASSWORD")

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(mx))
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/payments/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/payments", getAllPaymentDetails(formatter)).Methods("GET")
	mx.HandleFunc("/payments", addPaymentsDetails(formatter)).Methods("POST")
	mx.HandleFunc("/payment/{id}", getPaymentDetailsOfOne(formatter)).Methods("GET")
	mx.HandleFunc("/payment/{id}", deletePaymentDetailsOfOne(formatter)).Methods("DELETE")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"vCloud9.0 Payments API version 1.0 alive!"})
	}
}

// API Payments Handler
func addPaymentsDetails(formatter *render.Render) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		session, _ := mongo.Dial(mongodb_server)
		err := session.DB("admin").Login(mongodb_username, mongodb_password)
		if err != nil {
			formatter.JSON(writer, http.StatusInternalServerError, "Mongo Connection Error ")
			return
		}
		defer session.Close()
		session.SetMode(mongo.Monotonic, true)
		collection := session.DB(mongodb_database).C(mongodb_collection)
		fmt.Println(req.Body)

		var payment Payment
		_ = json.NewDecoder(req.Body).Decode(&payment)
		fmt.Printf("", payment)

		uuid := uuid.NewV4()
		payment.PaymentID = uuid.String()
		t := time.Now()
		payment.PaymentDate = t.Format("2006-01-02 15:04:05")
		payment.OrderStatus = true

		err = collection.Insert(payment)
		if err != nil {
			formatter.JSON(writer, http.StatusNotFound, "Create Payment Error")
			return
		}
		fmt.Println("Create new payment:", payment)
		formatter.JSON(writer, http.StatusOK, payment)
	}
}

//API to get all Payment
func getAllPaymentDetails(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := mongo.Dial(mongodb_server)
		err := session.DB("admin").Login(mongodb_username, mongodb_password)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Mongo Connection Error")
			return
		}
		defer session.Close()
		session.SetMode(mongo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		var result []bson.M
		err = c.Find(nil).All(&result)
		if err != nil {
			fmt.Println("error:" + err.Error())
			formatter.JSON(w, http.StatusNotFound, "Get All Payment Error")
			return
		}
		fmt.Println("getAllPaymentDetails:", result)
		formatter.JSON(w, http.StatusOK, result)
	}
}

//API to get 1 Payment
func getPaymentDetailsOfOne(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := mongo.Dial(mongodb_server)
		err := session.DB("admin").Login(mongodb_username, mongodb_password)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Mongo Connection Error")
			return
		}
		defer session.Close()
		session.SetMode(mongo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		var result bson.M
		params := mux.Vars(req)
		err = c.Find(bson.M{"orderid": params["id"]}).One(&result)
		fmt.Println("", err)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, "Get a Payment Error")
			return
		}
		fmt.Println("getAPaymentDetail:", result)
		formatter.JSON(w, http.StatusOK, result)
	}
}

//API to remove 1 Payment
func deletePaymentDetailsOfOne(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := mongo.Dial(mongodb_server)
		err := session.DB("admin").Login(mongodb_username, mongodb_password)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Mongo Connection Error")
			return
		}
		defer session.Close()
		session.SetMode(mongo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		var result Payment
		params := mux.Vars(req)
		err = c.Find(bson.M{"paymentid": params["id"]}).One(&result)
		fmt.Println("", err)
		if err != nil {
			fmt.Println("error:" + err.Error())
			formatter.JSON(w, http.StatusNotFound, "Delete a Payment Error")
			return
		} else {
			err = c.Remove(bson.M{"paymentid": result.PaymentID})
			if err != nil {
				fmt.Println("error:" + err.Error())
				formatter.JSON(w, http.StatusNotFound, "Delete Payment: deletion Error")
				return
			}
		}
		fmt.Println("DeletePaymentDetails:", result)
		formatter.JSON(w, http.StatusOK, result)
	}
}
