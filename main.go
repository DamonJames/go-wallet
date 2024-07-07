package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

func main() {
	fmt.Printf("Hello world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		getToken()
		fmt.Printf("token %s\n", token.Access)
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

var (
	token TokenResponse
)

type TokenResponse struct {
	Access          string
	Access_expires  int
	Refresh         string
	Refresh_expires int
}

func getToken() (string, error) {
	jsonBody := []byte(`{ "secret_id": "$SECRET_ID", "secret_key": "$SECRET_KEY" }`)
	bodyReader := bytes.NewReader(jsonBody)
	res, err := http.Post("https://bankaccountdata.gocardless.com/api/v2/token/new/", "application/json", bodyReader)
	if err != nil {
		fmt.Printf("got err getting token: %s", err.Error())
		return "", err
	}
	body, _ := io.ReadAll(res.Body)

	bodyString := string(body)
	var tokenResponse TokenResponse
	json.Unmarshal([]byte(bodyString), &tokenResponse)
	fmt.Printf("token: %s refresh: %s\n", tokenResponse.Access, tokenResponse.Refresh)
	token = tokenResponse
	return tokenResponse.Access, nil
}
