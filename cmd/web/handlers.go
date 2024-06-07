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
		http.Error(w, "Internal Server Error..", 500)
		return
	}

	//panic("An Error Happened")

	files := []string{
		"./ui/html/forms.page.tmpl",
	}
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
	app.session.Put(r, "flash", "Todo successfully created!")

	err := r.ParseForm()
	if err != nil {
		// app.clientError(w, http.StatusBadRequest)
		http.Error(w, "Method not allowed12", 400)
		return
	}

	title := r.PostForm.Get("text")

	if strings.TrimSpace(title) == "" {
		//errors["title"] = "This field cannot be blank"
		app.session.Put(r, "flash", "This field cannot be blank!")
	} else if utf8.RuneCountInString(title) > 50 {
		//errors["title"] = "This field is too long (maximum is 50 characters)"
		app.session.Put(r, "flash", "This field cannot be blank!")
	} else if strings.Contains(title, "special:") {
		_, err := app.specials.Insert(title)
		_, errInsert := app.todos.Insert(title)
		if errInsert != nil {
			http.Error(w, "Method not allowed..", 400)
			return
		}else if err !=nil{
			http.Error(w, "Method not allowed", 400)
		} else {
			app.session.Put(r, "flash", "Todo successfully created!")
		}
	}else{
		_, errInsert := app.todos.Insert(title)
		if errInsert != nil {
			http.Error(w, "Method not allowed..", 400)
			return
		}else {
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

	if err != nil {
		http.Error(w, "Method not allowed...", 400)
		return
	} else {
		app.session.Put(r, "flash", "Todo successfully deleted!")
	}
	//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// func (app *application) getsingle(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.URL.Query().Get("id"))
// 	if err != nil || id < 1 {
// 		http.Error(w, "Internal Server Error", 500)
// 		return
// 	}
// 	s, err := app.todos.Get(id)
// 	if err == models.ErrNoRecord {
// 		http.Error(w, "Internal Server Error", 500)
// 		return
// 	} else if err != nil {
// 		http.Error(w, "Internal Server Error", 500)
// 		return
// 	}
// 	fmt.Fprintf(w, "%v", s)
// }

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("ID"))
	name := r.FormValue("update")

	// if err != nil {
	// 	http.Error(w, "Internal Server Error...", 500)
	// 	return
	// }
	if strings.TrimSpace(name) == "" {
		//errors["title"] = "This field cannot be blank"
		app.session.Put(r, "flash", "This field cannot be blank123")
	} else {
		err := app.todos.Update(id, name)
		if err != nil {
			//app.serverError(w, err)

			http.Error(w, "Method not allowed...", 400)
			return
		} else {
			app.session.Put(r, "flash", "Todo successfully updated!")
		}
	}
	//redirecting to the same page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	ts.Execute(w, nil)

}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Method not allowed123", 400)
		return
	}

	useremail := r.PostForm.Get("email")
	userpassword := r.PostForm.Get("password")

	isUser, err := app.users.Authenticate(useremail, userpassword)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
	}
	if isUser {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "Flash", "Login successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Flash", "Login failed")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.session.Put(r, "Authenticated", false)
	}

	//fmt.Fprintln(w, "Create a new user...")
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	log.Println(ts)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	ts.Execute(w, nil)

}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// app.clientError(w, http.StatusBadRequest)
		http.Error(w, "Method not allowed123", 400)
		return
	}
	email := r.PostForm.Get("email")
	//log.Println(email)
	password := r.PostForm.Get("password")
	//log.Println(password)

	isUser, err := app.users.Authenticate(email, password)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server error", 500)
	}
	if isUser {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "Flash", "Login successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Flash", "Login failed")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.session.Put(r, "Authenticated", false)
	}
	//fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

// creating a form named specialForm to display it in the front end
func (app *application) specialForm(w http.ResponseWriter, r *http.Request) {
	s, err := app.specials.Latest()
	if err != nil {
		http.Error(w, "Internal Server Error..", 500)
		return
	}

	//panic("An Error Happened")

	files := []string{
		"./ui/html/special.page.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, struct {
		Tasks []*models.Special
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
// created function to take it into db named special
func (app *application) specialadd(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Method not allowed12", 400)
		return
	}

	title := r.PostForm.Get("text")

	if strings.TrimSpace(title) == "" {
		//errors["title"] = "This field cannot be blank"
		app.session.Put(r, "flash", "This field cannot be blank!")
	} else if utf8.RuneCountInString(title) > 50 {
		//errors["title"] = "This field is too long (maximum is 50 characters)"
		app.session.Put(r, "flash", "This field cannot be blank!!!!")
		//checking whether the name contains the substring "special:"
		//if yes added it to the specials table and todo table
	} else if strings.Contains(title, "special:") {
		_, errInsert := app.specials.Insert(title)
		_, err := app.todos.Insert(title)
		if err != nil {
			http.Error(w, "Method not allowed", 400)
		}
		if errInsert != nil {
			http.Error(w, "Method not allowed..", 400)
			return
		}
		app.session.Put(r, "flash", "Special successfully created!")

	}
	http.Redirect(w, r, "/user/special", http.StatusSeeOther)
}
//created a delete function to delete from the special table
func (app *application) specialdelete(w http.ResponseWriter, r *http.Request) {

	//delete by title from specials table using the delete function
	name := r.FormValue("title")
	_, err := app.specials.Delete(name)
	if err != nil {
		http.Error(w, "Method not allowed...", 400)
		return
	} else {
		app.session.Put(r, "flash", "Successfully deleted!")
	}
	//redirecting to the same page
	http.Redirect(w, r, "/user/special", http.StatusSeeOther)

}
