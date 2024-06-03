package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// create a structure for the input and id
type Task struct {
	Text string
	Id   int
}

// create slice for the struct
var allTask []Task
var id int

// create function for the html form to get the page
func form(w http.ResponseWriter, r *http.Request) {
	files := []string{"./ui/html/forms.page.tmpl"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, allTask)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error...", 500)
	}
}

// create function for adding the task and appending it to the slice
func addtask(w http.ResponseWriter, r *http.Request) {
	id++
	todo := Task{
		Text: r.FormValue("text"),
		Id:   id,
	}
	allTask = append(allTask, todo)
//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
//create a function to remove the specified task from the slice
func deletetask(w http.ResponseWriter, r *http.Request) {
	id, _ = strconv.Atoi(r.FormValue("id"))
	for i, value := range allTask {
		if value.Id == id {
			allTask = append(allTask[:i], allTask[i+1:]...)
		}
	}
//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
