package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "fmt"
  "path"
  "io/ioutil"
  "database/sql"
)

func ServerHandler(response http.ResponseWriter,request *http.Request) {
  response.Header().Set("Content-Type", "application/json")
  params := mux.Vars(request)
  db := "db"

  //open directory
  db_path := path.Join(db,params["server"])
  files,err := ioutil.ReadDir(db_path)
  if err != nil {
    http.Error(response, http.StatusText(404), 404)
    return
  }

  dbs := make([]*sql.DB,len(files))
  //find dbfiles
  for id,file := range files {
    full_path := path.join(db_path,file)
    dbs[id],err := sql.Open("sqlite3",full_path)
    if err != nil {
      http.Error(response,http.StatusTest(500),500)
    }
  }



  fmt.Fprint(response,"")
}

func main() {
  port := ":8080"

  router := mux.NewRouter()
  router.HandleFunc("/{server}",ServerHandler)
  http.Handle("/",router)
  http.ListenAndServe(port,nil)
  fmt.Println("done!")
}