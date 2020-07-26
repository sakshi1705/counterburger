package main

type User struct {
	Id        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Email   string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
}
