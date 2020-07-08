package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	fibgen "github.com/maxLogvynyuk/firstGo/package-fibgen"
)

var tpl *template.Template
var db *sql.DB

type Lesson struct {
	ID    int
	Lname string
}

func init() {
	var err error
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	connString := "dbname=gotoeleven password=password user=teacher port=5432 host=localhost sslmode=disable"

	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("****DB not open****", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("****Ping falar****", err)
	} else {
		fmt.Println("Ping seccessful!!!")
	}
}

func main() {
	defer db.Close()

	http.HandleFunc("/", hom)
	http.HandleFunc("/add-lesson", addLesson)
	http.HandleFunc("/set-lesson", setLesson)
	http.HandleFunc("/update-lesson", updateLesson)
	http.HandleFunc("/delete-lesson", deleteLesson)
	http.HandleFunc("/fibgen", fibgenPage)
	http.HandleFunc("/fibstart", fibgenTest)
	http.Handle("/stuff/", http.StripPrefix("/stuff", http.FileServer(http.Dir("./assets/"))))
	http.ListenAndServe(":8070", nil)
}

func hom(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select * from lessons;")
	defer rows.Close()

	if err != nil {
		log.Fatal("Can`t get data", err)
	}

	ls := []Lesson{}
	for rows.Next() {
		l := Lesson{}
		rows.Scan(&l.ID, &l.Lname)
		ls = append(ls, l)
	}
	tpl.ExecuteTemplate(w, "index.gohtml", ls)
}

func addLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fn := r.FormValue("l-name")
	fmt.Println(fn)

	result, err := db.Exec("insert into lessons (Lname) values ($1);", fn)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println("rows affacted", n)

	http.Redirect(w, r, "/", http.StatusFound)
}

func setLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	id := r.FormValue("lid")

	row, err := db.Query("select * from lessons where lid = $1", id)
	defer row.Close()

	fmt.Println("Selected lesson", row)

	if err != nil {
		http.Error(w, "Lesson not found", http.StatusInternalServerError)
	}

	l := Lesson{}

	for row.Next() {
		row.Scan(&l.ID, &l.Lname)
	}

	fmt.Println("l value", l)
	tpl.ExecuteTemplate(w, "update.gohtml", l)
}

func updateLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id := r.FormValue("lid")
	fn := r.FormValue("l-name")
	fmt.Println(fn, id)

	result, err := db.Exec("update lessons set lname = $2 where lid = $1;", id, fn)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println("rows affacted", n)

	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	id := r.FormValue("lid")

	result, err := db.Exec("delete from lessons where lid = $1;", id)
	if err != nil {
		http.Error(w, "Didn`t deleted", http.StatusInternalServerError)
		return
	}

	n, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Didn`t deleted", http.StatusInternalServerError)
		return
	}
	fmt.Println("Rows affacted", n)

	http.Redirect(w, r, "/", http.StatusFound)
}

func fibgenPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	tpl.ExecuteTemplate(w, "fibgen-page.gohtml", nil)
}

func fibgenTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var fv int
	var fs []int

	go fibgen.Worker(jobs, results)
	go fibgen.Worker(jobs, results)
	go fibgen.Worker(jobs, results)
	go fibgen.Worker(jobs, results)

	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)

	for j := 0; j < 100; j++ {
		fmt.Println(<-results)
		fv = <-results
		fs = append(fs, fv)
		fmt.Println("QQQQQQQQQQQ", fs)
	}

	tpl.ExecuteTemplate(w, "fibgen-page.gohtml", fs)
}
