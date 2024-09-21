package main

import notrhttp "github.com/Notr-Dev/notr-http"

func main() {
	server := notrhttp.NewServer("8080", "1.0")
	server.SetName("My Server")
	err := server.Run()
	if err != nil {
		panic(err)
	}

}
