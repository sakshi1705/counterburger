package main

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"net"
	"strings"
	"fmt"
	"database/sql"
	"log"
	"encoding/json"
)

var mysql_connect = "root:root@tcp(127.0.0.1:3306)/cmpe281"


func UserServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	router := mux.NewRouter()
	initRoutes(router, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD","DELETE", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	n.UseHandler(handlers.CORS(allowedHeaders,allowedMethods , allowedOrigins)(router))
	return n
}

func init() {
	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func initRoutes(router *mux.Router, formatter *render.Render) {
	router.HandleFunc("/users/test/ping", checkPing(formatter)).Methods("GET")
	router.HandleFunc("/users/signup", CreateUser).Methods("POST")
	router.HandleFunc("/users/signin", UserSignIn).Methods("POST")
}

func checkPing(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		message := "Burger Users API Server Working on machine: " + getSystemIp()
		formatter.JSON(w, http.StatusOK, struct{ Test string }{message})
	}
}
func CreateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	defer db.Close()
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
	unqueId := uuid.Must(uuid.NewV4())
	person.Id = unqueId.String()
	rows, err := db.Query("select count(1) as count from user where email = ?", person.Email)
	if err != nil {
		log.Fatal(err)
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		//fmt.Println("count")
		//fmt.Println(count)
		if err != nil {
			log.Fatal(err)
			message := struct {Message string}{"Some error occured while scanning results!!"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(message)
			return
		} else {
			if count > 0 {
				//fmt.Println("user already exists")
				message := struct {Message string}{"User already exists!!"}
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(message)
				return
			} else {
				//fmt.Println("inserting new user")
				insert, err := db.Query("insert into user (userid, fname, lname, email, password) values (?, ?, ?, ? ,?)", person.Id, person.Firstname, person.Lastname, person.Email, person.Password);
				if err != nil {
				panic(err.Error())
				message := struct {Message string}{"Some error occured while querying to database!!"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(message)
				return
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(person)
				defer insert.Close()
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func UserSignIn(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		message := struct {Message string}{"Some error occured while connecting to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	defer db.Close()
	var person User
	_ = json.NewDecoder(req.Body).Decode(&person)
	rows, err := db.Query("select count(1) as count from user where email = ?", person.Email)
	if err != nil {
		log.Fatal(err)
		message := struct {Message string}{"Some error occured while querying to database!!"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		//fmt.Println("count")
		//fmt.Println(count)
		if err != nil {
			log.Fatal(err)
			message := struct {Message string}{"Some error occured while scanning results!!"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(message)
			return
		} else {
			if count == 0 {
				//fmt.Println("user already exists")
				message := struct {Message string}{"User does not exist. Please sign up!!"}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(message)
				return
			} else {
				//fmt.Println("inserting new user")
				check, err := db.Query("select count(1) as count from user where email = ?  and password = ?", person.Email, person.Password);
				if err != nil {
				panic(err.Error())
				message := struct {Message string}{"Some error occured while querying to database!!"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(message)
				return
				}
				defer check.Close()
				var checkcount int
				for check.Next() {
					err := check.Scan(&checkcount)
					if err != nil {
						log.Fatal(err)
						message := struct {Message string}{"Some error occured while scanning results!!"}
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(message)
						return
					} else {
						if checkcount == 0 {
							//fmt.Println("user already exists")
							message := struct {Message string}{"Invalid credentials!!"}
							w.WriteHeader(http.StatusUnauthorized)
							json.NewEncoder(w).Encode(message)
							return
						} else {
							details, err := db.Query("select * from user where email = ?",  person.Email);
							if err != nil {
							panic(err.Error())
							message := struct {Message string}{"Some error occured while querying to database!!"}
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(message)
							return
							}
							defer details.Close()
							var loginperson User
							for details.Next() {
								err := details.Scan(&loginperson.Id, &loginperson.Firstname, &loginperson.Lastname, &loginperson.Email, &loginperson.Password)
								if err != nil {
									log.Fatal(err)
								}
							}
							w.WriteHeader(http.StatusCreated)
							json.NewEncoder(w).Encode(loginperson)
							
						}

					}

				}
			}
		}
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