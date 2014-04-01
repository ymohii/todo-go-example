package example

import (
    "time"
    "appengine"
    "appengine/datastore"
)

type Todo struct {
    ID int64 `datastore:"-"`
    Description string `datastore:",noindex"`
    Deadline time.Time
    Status bool
    Created time.Time
    Updated time.Time
}

type TodoView struct {
    ID int64
    Description string
    DescriptionPreview string
    Deadline string
    Status bool
}

func (this *Todo) key(c appengine.Context) *datastore.Key {
    if this.ID == 0 {
        return datastore.NewIncompleteKey(c, "Todo", todoAncestor(c))
    }
    return datastore.NewKey(c, "Todo", "", this.ID, todoAncestor(c))
}

func (this *Todo) save(c appengine.Context) (*Todo, error) {
    this.Updated = time.Now()
    if this.ID == 0 {
        this.Created = this.Updated
    }
    k, err := datastore.Put(c, this.key(c), this)
    if err != nil {
        return nil, err
    }
    this.ID = k.IntID()
    return this, nil
}

func (this *Todo) remove(c appengine.Context) error {
    return datastore.Delete(c, this.key(c))
}

func (this *Todo) getFormattedTodo(formatFor string) (*TodoView){
    todoView := new(TodoView)
    todoView.ID = this.ID
    const layout = "Jan 2, 2006 at 3:04pm (MST)"
    todoView.Deadline = this.Deadline.Format(layout)
    todoView.Status = this.Status

    if len(this.Description) > 50{
        todoView.DescriptionPreview = this.Description[0:50] + "..."
    }else {
        todoView.DescriptionPreview = this.Description
    }

    todoView.Description = this.Description

    return todoView
}

func todoAncestor(c appengine.Context) *datastore.Key{
    return datastore.NewKey(c, "Todo", "default", 0, nil)
}

func getTodoList(c appengine.Context) ([]*TodoView, error) {
    todos := []Todo{}
    keys, err := datastore.NewQuery("Todo").Ancestor(todoAncestor(c)).Order("-Updated").GetAll(c, &todos)
    if err != nil {
        return nil, err
    }
    todoViews := []*TodoView{}
    for i, key := range keys {
        todos[i].ID = key.IntID()
        todoViews = append(todoViews, todos[i].getFormattedTodo("list"))
    }
    return todoViews, nil
}

func getTodo(c appengine.Context, ID int64) (*Todo, error) {
    var todo Todo
    todo.ID = ID
    err := datastore.Get(c, todo.key(c), &todo)

    return &todo, err
}

