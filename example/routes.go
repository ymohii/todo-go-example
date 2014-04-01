package example

import (
    "net/http"
    "github.com/codegangsta/martini"
)

func init() {
    m := martini.Classic()

    //Todo Example
    m.Get("/todo/list", todoListHandler)
    m.Post("/todo/list", todoListHandler)
    m.Get("/todo/edit/:ID", todoEditHandler)
    m.Post("/todo/edit/:ID", todoEditHandler)
    m.Post("/todo/edit", todoEditPostHandler)
    m.Get("/todo/delete/:ID", todoDeleteHandler)
    m.Post("/todo/delete", todoDeletePostHandler)
    m.Get("/todo/create", todoCreateHandler)
    m.Post("/todo/create", todoCreatePostHandler)
    m.Get("/todo/:parentID", todoViewHandler)
    m.Post("/todo/:parentID", todoViewHandler)


    m.Get("/todo/:parentID/edit/:ID", subtaskEditHandler)
    m.Post("/subtask/edit", subtaskEditPostHandler)
    m.Get("/todo/:parentID/delete/:ID", subtaskDeleteHandler)
    m.Post("/subtask/delete", subtaskDeletePostHandler)
    m.Get("/todo/:parentID/createSubtask", subtaskCreateHandler)
    m.Post("/subtask/create", subtaskCreatePostHandler)
    // Handle this all
    http.Handle("/", m)
}

