package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todo/pkg/models"
	"unicode/utf8"
)

// create a structure for the input and id
type Task struct {
	Text string
	Id   int
}

// create slice for the struct
// var allTask []Task
// var id int

// create function for the html form to get the page
func (app *application) form(w http.ResponseWriter, r *http.Request) {
	s, err := app.todos.Latest()
	if err != nil {
		http.Error(w, "Internal Server Error...", 500)
		return
	}

	//panic("An Error Happened")

	files := []string{"./ui/html/forms.page.tmpl"}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, struct {
        Tasks []*models.Todo
        Flash string
    }{
        Tasks: s,
        Flash: app.session.PopString(r, "flash"),
    })
	
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error...", 500)
	}
}

// create function for adding the task and appending it to the slice
func (app *application) addtask(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	//     w.Header().Set("Allow", "POST")
	//     app.clientError(w, http.StatusMethodNotAllowed)
	//     return
	// }
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.

	// title := "O snail"
	//content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi
	// expires := "7"

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.

	// Redirect the user to the relevant page for the snippet.
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	app.session.Put(r, "flash", "Todo successfully created!")

	err := r.ParseForm()
	if err != nil {
		// app.clientError(w, http.StatusBadRequest)
		http.Error(w, "Method not allowed123", 400)
		return
	}

	title := r.PostForm.Get("text")
	log.Println(title)
	//expires := r.PostForm.Get("expires")
	//errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		//errors["title"] = "This field cannot be blank"
		app.session.Put(r, "flash", "This field cannot be blank!")
	} else if utf8.RuneCountInString(title) > 50 {
		//errors["title"] = "This field is too long (maximum is 50 characters)"
		app.session.Put(r, "flash", "This field cannot be blank!")
	}else{
		_, errInsert := app.todos.Insert(title)
	if errInsert != nil {
		//app.serverError(w, err)

		http.Error(w, "Method not allowed...", 400)
		return
	}else{
		app.session.Put(r, "flash", "Todo successfully created!")
	}
	}

	// if strings.TrimSpace(expires) == "" {
	// 	errors["expires"] = "This fied cannot be blank"
	// 	} else if expires != "365" && expires != "7" && expires != "1" {
	// 	errors["expires"] = "This field is invalid"
	// 	}
	// if len(errors) > 0 {
	// 	fmt.Fprint(w, errors)
	// 	return
	// }
	http.Redirect(w, r, "/", http.StatusSeeOther) //fmt.Sprintf("/addtask/%d", id)
}

// id++
// todo := Task{
// 	Text: r.FormValue("text"),
// 	Id:   id,
// }
// allTask = append(allTask, todo)
// //redirecting to the same page
// _, err := app.todos.Insert(todo.Text)
// if err != nil {
// 	app.errorLog.Println(err.Error())
// 	return
// }

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

// create a function to remove the specified task from the slice
func (app *application) deletetask(w http.ResponseWriter, r *http.Request) {
	ids, _ := strconv.Atoi(r.FormValue("ID"))
	_, err := app.todos.Delete(ids)
	// if err != nil {
	// 	http.Error(w, "Internal Server Error...", 500)
	// 	return
	// }
	if err != nil {
		//app.serverError(w, err)

		http.Error(w, "Method not allowed...", 400)
		return
	}else{
		app.session.Put(r, "flash", "Todo successfully deleted!")
	}
	
	//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) getsingle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	s, err := app.todos.Get(id)
	if err == models.ErrNoRecord {
		http.Error(w, "Internal Server Error", 500)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	fmt.Fprintf(w, "%v", s)
}

//redirecting to the same page
//http.Redirect(w, r, "/", http.StatusSeeOther)

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("ID"))
	name := r.FormValue("update")
	
	// if err != nil {
	// 	http.Error(w, "Internal Server Error...", 500)
	// 	return
	// }
	if strings.TrimSpace(name) == "" {
		//errors["title"] = "This field cannot be blank"
		app.session.Put(r, "flash", "This field cannot be blank!")
	}else{
		err := app.todos.Update(id, name)
	if err != nil {
		//app.serverError(w, err)

		http.Error(w, "Method not allowed...", 400)
		return
	}else{
		app.session.Put(r, "flash", "Todo successfully updated!")
	}
	}
	//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
