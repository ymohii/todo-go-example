package example

import (
    "net/http"
    "html/template"
    "appengine"
    "github.com/codegangsta/martini"
    "strconv"
//    "fmt"
)

func subtaskEditHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,err := strconv.Atoi(params["ID"])
    parentID, _ := strconv.Atoi(params["parentID"])
    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    subtask,err := getSubtask(c, (int64)(id), (int64)(parentID))


    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    t, err := template.ParseFiles("templates/subtask/index.html", "templates/subtask/edit.html")


    if err := t.Execute(w, map[string]interface{}{"subtask": subtask.getFormattedSubtask() } ); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func subtaskEditPostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    context := appengine.NewContext(r)
    description := r.FormValue("description")
    status := r.FormValue("status")


    id,err := strconv.Atoi(r.FormValue("ID"))
    parentID,err := strconv.Atoi(r.FormValue("parent"))

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    subtask,err := getSubtask(context, (int64)(id), (int64)(parentID))
    subtask.Description = description
    if status == "completed" {
        subtask.Status = true
    }else {
        subtask.Status = false
    }


    subtask.save(context)

    http.Redirect(w, r, "/todo/" + r.FormValue("parent"), http.StatusTemporaryRedirect)

}

func subtaskCreateHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    t, _ := template.ParseFiles("templates/subtask/index.html", "templates/subtask/create.html")

    if err := t.Execute(w, map[string]interface{}{"parentID": params["parentID"]}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func subtaskDeleteHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,_ := strconv.Atoi(params["ID"])
    parentID, _ := strconv.Atoi(params["parentID"])
    subtask,err := getSubtask(c, (int64)(id), (int64)(parentID))

    if(err != nil) {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    t, err := template.ParseFiles("templates/subtask/index.html", "templates/subtask/delete.html")

    if err := t.Execute(w, map[string]interface{}{"subtask": subtask.getFormattedSubtask()}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}

func subtaskDeletePostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    c := appengine.NewContext(r)
    id,_ := strconv.Atoi(r.FormValue("ID"))
    parentID,_ := strconv.Atoi(r.FormValue("parent"))


    subtask, _ := getSubtask(c, (int64)(id), (int64)(parentID))
    err := subtask.remove(c)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    http.Redirect(w, r, "/todo/" + r.FormValue("parent"), http.StatusTemporaryRedirect)

}

func subtaskCreatePostHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
    context := appengine.NewContext(r)
    description := r.FormValue("description")
    status := r.FormValue("status")
    parentID,err := strconv.Atoi(r.FormValue("parent"))
    if err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    var subtask Subtask
    subtask.Description = description
    subtask.ParentID = (int64)(parentID)
    todo,err := getTodo(context, subtask.ParentID)
    if err != nil {
         http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    subtask.Parent = todo.key(context)

    if status == "completed" {
        subtask.Status = true
    }else {
        subtask.Status = false
    }

    subtask.save(context)

    http.Redirect(w, r, "/todo/" + r.FormValue("parent"), http.StatusTemporaryRedirect)
}

