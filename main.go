package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getData(url *string) {
	if res, err := http.Get(*url); err == nil {
		defer res.Body.Close()
		if body, err := ioutil.ReadAll(res.Body); err == nil {
			users := [](map[string]interface{}){}
			if err := json.Unmarshal([]byte(body), &users); err != nil {
				log.Fatal(err)
			} else {
				for _, user := range users {
					fmt.Println(user["guild"].(map[string]interface{})["avatar_url"])
				}
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func formHandler(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(rw, "An error occured: %v", err)
		return
	}

	fmt.Fprintf(rw, "POST received")
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Println(name, address)
}

func DoneAsync() chan int {
	r := make(chan int)
	fmt.Println("Warming up ...")
	go func() {
		time.Sleep(3 * time.Second)
		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func runAsync() {
	fmt.Println("Let's start ...")
	val := DoneAsync()
	fmt.Println("Done is running ...")
	fmt.Println(<-val)
}

func startServer() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello, world!"))
	})

	http.HandleFunc("/form", formHandler)

	const PORT = 3030
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)

	fmt.Println("Server started on port " + strconv.Itoa(PORT))
}

func main() {
	url := "https://api.vimeworld.ru/user/name/ineedmoreofyou"

	go getData(&url)

	startServer()

	fmt.Scanln()

	// runAsync()
}
