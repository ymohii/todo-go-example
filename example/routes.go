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
    m.Post("/todo/edit", todoEditPostHandler)
    m.Get("/todo/delete/:ID", todoDeleteHandler)
    m.Post("/todo/delete", todoDeletePostHandler)
    m.Get("/todo/create", todoCreateHandler)
    m.Post("/todo/create", todoCreatePostHandler)

    // Handle this all
    http.Handle("/", m)
}

