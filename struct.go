package main

type User struct {
	ID           string `json:"id"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Phone        int64  `json:"phone"`
	QQ           int64  `json:"qq"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type Response struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   struct {
		User User `json:"user"`
	} `json:"data"`
}
