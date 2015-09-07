package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

type FastCGIServer struct{}

func (s FastCGIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	fileName := "../../document/name.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}

	file.WriteString(req.FormValue("name") + "\n")

	fmt.Printf("Hello, %q\n", html.EscapeString(req.FormValue("name")))
	fmt.Printf("Hello, %q\n", html.EscapeString(req.URL.Path))
	fmt.Printf("Hello, %q\n", html.EscapeString(req.RemoteAddr))

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	resp.Write(buf)
}

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:9001")
	srv := new(FastCGIServer)
	fcgi.Serve(listener, srv)
}
