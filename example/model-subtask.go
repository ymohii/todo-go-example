package example

import (
    "appengine"
    "appengine/datastore"
)

type Subtask struct {
    ID int64 `datastore:"-"`
    Description string `datastore:",noindex"`
    Status bool
    Parent *datastore.Key `datastore:"-"`
    ParentID int64 `datastore:"-"`
}

type SubtaskView struct {
    ID int64
    Description string
    DescriptionPreview string
    Status bool
    ParentID int64
}

func (this *Subtask) key(c appengine.Context) *datastore.Key {
    if this.Parent == nil {
        todo,_ := getTodo(c, this.ParentID)
        this.Parent = todo.key(c)
    }
    if this.ID == 0 {
        return datastore.NewIncompleteKey(c, "Subtask", this.Parent)
    }
    return datastore.NewKey(c, "Subtask", "", this.ID, this.Parent)
}

func (this *Subtask) save(c appengine.Context) (*Subtask, error) {
    k, err := datastore.Put(c, this.key(c), this)
    if err != nil {
        return nil, err
    }
    this.ID = k.IntID()
    return this, nil
}

func (this *Subtask) remove(c appengine.Context) error {
    return datastore.Delete(c, this.key(c))
}

func (this *Subtask) getFormattedSubtask() (*SubtaskView){
    subtaskView := new(SubtaskView)
    subtaskView.ID = this.ID
    subtaskView.Status = this.Status

    if len(this.Description) > 50{
        subtaskView.DescriptionPreview = this.Description[0:50] + "..."
    }else {
        subtaskView.DescriptionPreview = this.Description
    }

    subtaskView.Description = this.Description
    subtaskView.ParentID = this.ParentID

    return subtaskView
}

func getSubtaskList(c appengine.Context, parent *Todo) ([]*SubtaskView, error) {
    subtasks := []Subtask{}
    keys, err := datastore.NewQuery("Subtask").Ancestor(parent.key(c)).GetAll(c, &subtasks)
    if err != nil {
        return nil, err
    }
    subtaskViews := []*SubtaskView{}
    for i, key := range keys {
        subtasks[i].ID = key.IntID()
        subtasks[i].ParentID = parent.ID
        subtaskViews = append(subtaskViews, subtasks[i].getFormattedSubtask())
    }
    return subtaskViews, nil
}

func getSubtask(c appengine.Context, ID int64, parent int64) (*Subtask, error) {
    todo, _ := getTodo(c, parent)
    var subtask Subtask
    subtask.ID = ID
    subtask.Parent = todo.key(c)
    subtask.ParentID = parent
    err := datastore.Get(c, subtask.key(c), &subtask)

    return &subtask, err
}

