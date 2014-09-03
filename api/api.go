package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "path"
  "io/ioutil"
  "code.google.com/p/go-sqlite/go1/sqlite3"
  // "database/sql"
)

type Event struct {
  Timestamp int64
  Message string
  Bot string
  Source string
  Command string
}

func ServerHandler(response http.ResponseWriter,request *http.Request) {
  response.Header().Set("Content-Type", "application/json")
  params := mux.Vars(request)
  db_root := "db"

  //open directory
  db_path := path.Join(db_root,params["server"])
  files,err := ioutil.ReadDir(db_path)
  if err != nil {
    http.Error(response, http.StatusText(404), 404)
    return
  }

  //needs to be a map (as in map-reduce map)
  //spin this up as a go routine, send to a reducer
  dbs := make([]*sqlite3.Conn,len(files))
  //find dbfiles
  for id,file := range files {
    full_path := path.Join(db_path,file.Name())
    dbs[id],err = sqlite3.Open(full_path)
    defer dbs[id].Close()
    if err != nil {
      fmt.Println(err)
      http.Error(response,http.StatusText(500),500)
    }
  }

  sql := "SELECT * FROM events;"
  row := make(sqlite3.RowMap)
  events := make([]Event,0)
  for _,db := range dbs {
    for s,err := db.Query(sql); err == nil; err = s.Next() {
      s.Scan(row)
      event := Event{
        Timestamp: row["timestamp"].(int64),
        Message: row["message"].(string),
        Bot: row["bot"].(string),
        Source: row["source"].(string),
        Command: row["command"].(string),
      }
      fmt.Println(row)
      events = append(events,event)
    }
  }

  fmt.Fprint(response,events)
}

func main() {
  port := ":8080"
  router := mux.NewRouter()
  router.HandleFunc("/{server}",ServerHandler)
  http.Handle("/",router)
  http.ListenAndServe(port,nil)
  fmt.Println("done!")
}