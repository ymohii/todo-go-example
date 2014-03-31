package example

import (
    "io"
    "time"
    "encoding/json"
    "appengine"
    "appengine/datastore"
)

// A Request
type Request struct {
    ID int64 `datastore:"-"`
    DeviceId string
    Created time.Time
}

func (this *Request) key(c appengine.Context) *datastore.Key {
    if this.ID == 0 {
        //this.Created = time.Now()
        return datastore.NewIncompleteKey(c, "Request", requestList(c))
    }
    return datastore.NewKey(c, "Request", "", this.ID, requestList(c))
}

func (this *Request) save(c appengine.Context) (*Request, error) {
    this.Created = time.Now()
    k, err := datastore.Put(c, this.key(c), this)
    if err != nil {
        return nil, err
    }
    this.ID = k.IntID()
    return this, nil
}

func requestList(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Request", "default", 0, nil)
}

func decodeRequest(r io.ReadCloser) (*Request, error) {
    defer r.Close()
    var request Request
    err := json.NewDecoder(r).Decode(&request)
    return &request, err
}

func getAllRequests(c appengine.Context) ([]Request, error) {
    requests := []Request{}
    keys, err := datastore.NewQuery("Request").Ancestor(requestList(c)).Order("-Created").GetAll(c, &requests)
    if err != nil {
        return nil, err
    }
    for i, key := range keys {
        requests[i].ID = key.IntID()
    }
    return requests, nil
}
