package model

type LoginResponse struct {
	Id            string `json:"id"`
	First_name    string `json:"first_name"`
	Role          string `json:"role"`
	Gender        string `json:"gender"`
	Last_name     string `json:"last_name`
	Email         string `json:"email`
	Date_of_birth string `json:"date_of_birth"`
}

type Tokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}