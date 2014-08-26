package main

import(
  "github.com/diebels727/spyglass"
  "github.com/gorilla/mux"
  "flag"
  "net/http"
  "fmt"
  // "time"
)


var server string
var port string
var nick string
var username string
var password string
var command_and_control string

func init() {
  flag.StringVar(&server,"server","irc.freenode.org","IRC server FQDN")
  flag.StringVar(&port,"port","6667","IRC server port number")
  flag.StringVar(&nick,"nick","","Name of the bot visible on IRC channel")
  flag.StringVar(&username,"username","logbot","Username to login with to IRC")
  flag.StringVar(&password,"password","","Password for the IRC server")
  flag.StringVar(&command_and_control,"command_and_control","#spyglass-c&c","Command and control IRC channel")
}

func main() {
  flag.Parse()
  var bot *spyglass.Bot
  bot = spyglass.New(server,port,nick,username,password)
  conn := bot.Connect()
  defer conn.Close()

  bot.Run()

  // go func(bot *spyglass.Bot) {
  //   fmt.Println("ticking")
  //   ticker := time.NewTicker(time.Second * 5)
  //   for _ = range ticker.C {
  //     bot.Cmd("PING irc.freenode.net")
  //   }
  // }(bot)

  <- bot.Ready

  // go func(bot *spyglass.Bot) {


  // }(bot)

  bot.User()
  bot.Nick()
  bot.Join(command_and_control)

  router := mux.NewRouter()
  router.Methods("POST")
  router.HandleFunc("/{channel}",func(response http.ResponseWriter,request *http.Request) {
    params := mux.Vars(request)
    channel := "#" + params["channel"]
    fmt.Println("Joining channel ",channel)
    bot.Join(channel)
  })

  http.Handle("/",router)
  http.ListenAndServe(":9000",nil)

  <- bot.Stopped
}






