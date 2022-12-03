package main

type Token struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Refresh_token string `json:"refresh_token"`
	Expires_in    int    `json:"expires_in"`
	When_Expires  string `json:"expires_in_datetime"`
	Scope         string `json:"scope"`
	Id_token      string `json:"id_token"`
	Expired       int    `json:"expired"`
}
