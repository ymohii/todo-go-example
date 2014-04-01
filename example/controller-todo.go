package example

import (
    "net/http"
    "html/template"
    "appengine"
    "github.com/codegangsta/martini"
    "time"
    "strconv"
//    "fmt"
)

func todoListHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)

    todos, err := getTodoList(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    t, _ := template.ParseFiles("templates/todo/index.html", "templates/todo/list.html")
    if err := t.Execute(w, map[string]interface{}{"todos": todos}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}

func todoEditHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,err := strconv.Atoi(params["ID"])
    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    todo,err := getTodo(c, (int64)(id))

    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    t, err := template.ParseFiles("templates/todo/index.html", "templates/todo/edit.html")


    if err := t.Execute(w, map[string]interface{}{"todo": todo.getFormattedTodo("edit")} ); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func todoEditPostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    context := appengine.NewContext(r)
    description := r.FormValue("description")
    const layout = "Jan 2, 2006 at 3:04pm (MST)"
    deadline,err := time.Parse(layout, r.FormValue("deadline"))
    status := r.FormValue("status")

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    id,err := strconv.Atoi(r.FormValue("ID"))

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    todo,err := getTodo(context, (int64)(id))
    todo.Description = description
    todo.Deadline = deadline
    if status == "completed" {
        todo.Status = true
    }else {
        todo.Status = false
    }


    todo.save(context)

    http.Redirect(w, r, "/todo/list", http.StatusTemporaryRedirect)

}

func todoCreateHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    t, _ := template.ParseFiles("templates/todo/index.html", "templates/todo/create.html")
    now := time.Now()
    const layout = "Jan 2, 2006 at 3:04pm (MST)"

    if err := t.Execute(w, map[string]interface{}{"timeNow": now.Format(layout)}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func todoDeleteHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,_ := strconv.Atoi(params["ID"])
    todo,err := getTodo(c, (int64)(id))

    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    t, err := template.ParseFiles("templates/todo/index.html", "templates/todo/delete.html")

    if err := t.Execute(w, map[string]interface{}{"todo": todo.getFormattedTodo("delete")}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}

func todoDeletePostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,_ := strconv.Atoi(r.FormValue("ID"))
    var todo Todo
    todo.ID = (int64)(id)
    err := todo.remove(c)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    http.Redirect(w, r, "/todo/list", http.StatusTemporaryRedirect)

}

func todoCreatePostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    context := appengine.NewContext(r)
    description := r.FormValue("description")
    const layout = "Jan 2, 2006 at 3:04pm (MST)"
    deadline,_ := time.Parse(layout, r.FormValue("deadline"))
    status := r.FormValue("status")
    var todo Todo
    todo.Description = description
    todo.Deadline = deadline
    if status == "completed" {
        todo.Status = true
    }else {
        todo.Status = false
    }

    todo.save(context)

    http.Redirect(w, r, "/todo/list", http.StatusTemporaryRedirect)
}

