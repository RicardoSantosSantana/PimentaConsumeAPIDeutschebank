package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetToken() (Token, error) {

	token, err := active_token()

	// There is no token in the database
	if err != nil || token.Access_token == "" {
		token, err_new_token := new_token()

		if err_new_token != nil {
			return Token{}, err_new_token
		}
		err_save_token := token.save_token()

		if err_save_token != nil {
			return Token{}, err_save_token
		}
		return token, nil
	}

	// There is token in the database
	if token.Expired == 1 {

		token, err_refresh_token := refresh_token()

		if err_refresh_token != nil {
			return Token{}, err_refresh_token
		}

		err_save_token := token.save_token()

		if err_save_token != nil {
			return Token{}, err_save_token
		}
		return token, nil
	}
	fmt.Println("exec database ...")
	return token, nil

}

func new_token() (Token, error) {

	fmt.Println("exec new token ...")

	settings := Settings()

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", settings.Auth.Code)
	data.Set("redirect_uri", settings.Auth.Redirect_uri)

	token, err_response := response_token(strings.NewReader(data.Encode()))

	if err_response != nil {
		return token, err_response
	}

	return token, nil
}

func refresh_token() (Token, error) {

	fmt.Println("exec refresh token ...")
	token, err := active_token()

	if err != nil {
		return Token{}, err
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", token.Refresh_token)

	token, err_response := response_token(strings.NewReader(data.Encode()))

	if err_response != nil {
		return token, err_response
	}

	return token, nil
}

func response_token(payload io.Reader) (Token, error) {

	token := Token{}
	settings := Settings()

	req, err := http.NewRequest(http.MethodPost, settings.GetToken.Uri, payload)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(settings.Auth.Client_id, settings.Auth.Client_secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return token, err
	}

	if strings.Contains(string(body), "error") {
		return token, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &token); err != nil {

		return token, errors.New(string(body))
	}

	return token, nil
}

func expires_in_to_datetime(Expires_in time.Duration) string {

	DateNow := time.Now().Add(time.Second * Expires_in)

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		DateNow.Year(), DateNow.Month(), DateNow.Day(),
		DateNow.Hour(), DateNow.Minute(), DateNow.Second())
}

func (token Token) save_token() error {

	ExpiresInDateTime := expires_in_to_datetime(time.Duration(token.Expires_in))
	token.When_Expires = ExpiresInDateTime

	db, err := openConnection()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE tokens SET access_token=?, token_type=?, expires_in=?, when_expires=?, scope=?, id_token=?, refresh_token=?")
	if err != nil {
		return err
	}

	_, error := stmt.Exec(token.Access_token, token.Token_type, token.Expires_in, token.When_Expires, token.Scope, token.Id_token, token.Refresh_token)

	if error != nil {
		return error
	}

	defer db.Close()

	return nil
}

func active_token() (Token, error) {

	db, err := openConnection()
	if err != nil {
		return Token{}, err
	}

	token := Token{}
	expired := "CASE WHEN when_expires >= NOW() THEN 0	ELSE 1 	END AS expired"

	sql := "SELECT access_token, token_type, expires_in, when_expires, scope, id_token, refresh_token, " + expired + " FROM tokens ORDER BY id DESC LIMIT 1"
	err = db.QueryRow(sql).Scan(&token.Access_token, &token.Token_type, &token.Expires_in, &token.When_Expires, &token.Scope, &token.Id_token, &token.Refresh_token, &token.Expired)

	if err != nil {
		return Token{}, err
	}

	defer db.Close()

	return token, nil

}
