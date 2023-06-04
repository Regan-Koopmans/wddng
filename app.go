package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var resources embed.FS
var t = template.Must(template.ParseFS(resources, "templates/*"))

func run_sql(statement string, w http.ResponseWriter) {
	db, _ := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
	defer db.Close()
	_, err := db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "ACK")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/bmljawo=", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{}
		t.ExecuteTemplate(w, "nick.html.tmpl", data)
	})
	http.HandleFunc("/Y3JhaWcK", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{}
		t.ExecuteTemplate(w, "craig.html.tmpl", data)
	})
	http.HandleFunc("/bmljawo=accept", func(w http.ResponseWriter, r *http.Request) {
		run_sql("UPDATE rsvp SET response = 'ACCEPT' WHERE name='nick'", w)
	})
	http.HandleFunc("/bmljawo=reject", func(w http.ResponseWriter, r *http.Request) {
		run_sql("UPDATE rsvp SET response = 'REJECT' WHERE name='nick'", w)
	})
	http.HandleFunc("/Y3JhaWcKaccept", func(w http.ResponseWriter, r *http.Request) {
		run_sql("UPDATE rsvp SET response = 'ACCEPT' WHERE name='craig'", w)
	})
	http.HandleFunc("/Y3JhaWcKreject", func(w http.ResponseWriter, r *http.Request) {
		run_sql("UPDATE rsvp SET response = 'REJECT' WHERE name='craig'", w)
	})

	http.HandleFunc("/responses-overview", func(w http.ResponseWriter, r *http.Request) {
		db, _ := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
		defer db.Close()
		sts := "SELECT * FROM rsvp;"
		var name string
		var response string
		rows, _ := db.Query(sts)
		for rows.Next() {
			_ = rows.Scan(&name, &response)
			fmt.Fprintf(w, "%s %s\n", name, response)
		}
	})

	http.HandleFunc("/aW5pdGRiCg==", func(w http.ResponseWriter, r *http.Request) {
		db, _ := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
		sts := `
			DROP TABLE IF EXISTS rsvp;
			CREATE TABLE rsvp(name TEXT PRIMARY KEY, response TEXT default 'UNDEFINED');
			INSERT INTO rsvp(name) VALUES('nick');
			INSERT INTO rsvp(name) VALUES('craig');
		`
		_, _ = db.Exec(sts)
		defer db.Close()
		fmt.Fprintf(w, "DB INIT")
	})
	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
