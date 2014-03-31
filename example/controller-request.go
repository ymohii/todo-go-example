package example

import (
    "net/http"
    "html/template"
    "appengine"
    "github.com/codegangsta/martini"

//    "fmt"
)

func requestQueryHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {

    context := appengine.NewContext(r)
    deviceId := r.FormValue("device_id")

    request, _ := decodeRequest(r.Body)
    request.DeviceId = deviceId
    request, _ = request.save(context)

    t, _ := template.ParseFiles("templates/index.html", "templates/request/query.html")
    if err := t.Execute(w, map[string]interface{}{"deviceId": deviceId}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

}

func requestListHandler(w http.ResponseWriter, r *http.Request) {

    c := appengine.NewContext(r)

    requests, err := getAllRequests(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    t, _ := template.ParseFiles("templates/index.html", "templates/request/list.html")
    if err := t.Execute(w, map[string]interface{}{"requests": requests}); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    // debug only
    //count := len(requests)
    //fmt.Fprint(w, count)

}