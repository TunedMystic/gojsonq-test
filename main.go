package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

//go:embed "users.json"
var UsersJson string

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func HandleUser(jq *gojsonq.JSONQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the url parameter.
		//  "/users/bob/" -> "bob"
		//
		userName, found := strings.CutPrefix(r.URL.Path, "/user/")
		userName = strings.Trim(userName, "/")

		// If no userName, then return 404.
		if !found {
			fmt.Println("Not Found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Query from json file.
		users := []User{}
		jq.Copy().From("users").Where("name", "=", userName).Out(&users)

		if len(users) == 0 {
			fmt.Println("Not Found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Marshal json response.
		out, err := json.Marshal(users[0])
		if err != nil {
			http.Error(w, "error with struct marshalling", http.StatusInternalServerError)
		}

		w.Write(out)
	}
}

func HandleIndex(jq *gojsonq.JSONQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// If not root page, then return 404.
		if r.URL.Path != "/" {
			fmt.Println("Not Found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Query from json file.
		res, err := jq.Copy().From("users").PluckR("name")
		if err != nil {
			fmt.Println(err)
		}

		// Convert to string slice.
		names, err := res.StringSlice()
		if err != nil {
			fmt.Printf("⚠️ %v\n", err)
		}

		// Construct URLs for response.
		for i, name := range names {
			names[i] = fmt.Sprintf("http://localhost:3000/user/%s", name)
		}

		// Marshal json response.
		out, err := json.Marshal(names)
		if err != nil {
			fmt.Println("Server Error")
			http.Error(w, "error with struct marshalling", http.StatusInternalServerError)
		}

		w.Write([]byte(out))
	}
}

func main() {
	jq := gojsonq.New().FromString(UsersJson)

	http.Handle("/", HandleIndex(jq))
	http.Handle("/user/", HandleUser(jq))
	fmt.Println("Starting server on :3000")
	http.ListenAndServe(":3000", nil)
}
