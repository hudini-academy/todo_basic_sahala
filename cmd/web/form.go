package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"todo/pkg/models/mysql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions" 
)

type Config struct{
	Addr string
	StaticDir string
	Dsn string
	Session string
}


type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config *Config
	todos    *mysql.TodoModel
	session       *sessions.Session 
}

func main() {
	//created a new config so that the default values can be used if we haven't given the port address, server route and db
	config := new(Config)
	flag.StringVar(&config.Addr,"addr",":4000", "Default HTTP network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static","Path to static access")
	flag.StringVar(&config.Dsn, "dsn","root:root@/todo?parseTime=true", "MySQL database")
	flag.StringVar(&config.Session,"secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Session")
	flag.Parse()
 
	//creating a new log infoLog
	f, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	//creating a file for error and using it as the log destination
	el, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	
	//creating a new log for error
	errorLog := log.New(el, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(config.Dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte(*&config.Session)) 
	session.Lifetime = 12 * time.Hour
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session: session,
		todos:    &mysql.TodoModel{DB: db},
		config : config,
	}
	// srv := &http.Server{
	// 	Addr:     *addr,
	// 	ErrorLog: errorLog,
	// 	Handler:  app.routes(),
	// }
	
	infoLog.Printf("Starting server on %s", app.config.Addr)
	err = http.ListenAndServe(app.config.Addr, app.routes())
	log.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
