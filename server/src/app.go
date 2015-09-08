package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

type FastCGIServer struct{}

func dbAdd(db *sql.DB, name string) {
	stmt, err := db.Prepare("INSERT INTO t_user(name) VALUES(?)")
	defer stmt.Close()

	if err != nil {
		log.Println(err)
		return
	}
	stmt.Exec(name)
}

func dbQuery(db *sql.DB) {
	rows, err := db.Query("select userId, name from t_user where userId = ?", 10001)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	fmt.Println(columns)

	for rows.Next() {
		var userId int
		var name string
		err := rows.Scan(&userId, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(userId, name)
		fmt.Println(userId)
		fmt.Println(name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (s FastCGIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	fileName := "../../document/name.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}

	db, err := sql.Open("mysql", "root:1@/myscript")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	name := req.FormValue("name")
	dbAdd(db, name)
	file.WriteString(req.FormValue("name") + "\n")
	dbQuery(db)

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
