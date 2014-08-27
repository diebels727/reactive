package main

import(
  "github.com/diebels727/spyglass"
  // "github.com/gorilla/mux"
  "github.com/diebels727/faker"
  "flag"
  // "net/http"
  "fmt"
  "strings"
  // "net"
  // "net/textproto"
  "time"
)


var server string
var port string
var nick string
var username string
var password string
var command_and_control string

var fake *faker.Faker

func init() {
  flag.StringVar(&server,"server","irc.freenode.net","IRC server FQDN")
  flag.StringVar(&port,"port","6667","IRC server port number")
  flag.StringVar(&nick,"nick","","Name of the bot visible on IRC channel")
  flag.StringVar(&username,"username","","Username to login with to IRC")
  flag.StringVar(&password,"password","","Password for the IRC server")
  flag.StringVar(&command_and_control,"command_and_control","#spyglass-c&c","Command and control IRC channel")
}

type Clients [](*spyglass.Bot)

type Client *spyglass.Bot

func NewClient(server string,port string) (*spyglass.Bot) {
  return spyglass.New(server,port,fake.Username(),fake.Username(),"")
}

func main() {
  fake = faker.New()

  flag.Parse()

  clients := make(Clients,5)

  clients = ([](*spyglass.Bot))(clients)
  for i:=0;i<len(clients);i++ {
    clients[i] = NewClient(server,port)
    clients[i].Connect()
    time.Sleep(time.Duration(time.Second * 3))
    defer clients[i].Conn.Close()
  }

  var channels map[string]bool
  channels = make(map[string]bool)
  var channel string

  master := clients[0]

  master.Run()

  //react to event 322, which is each listed channel
  master.RegisterEventHandler("322",func(event *spyglass.Event) {
    arguments := event.RawArguments
    args := strings.Split(arguments," ")
    channel = args[1]
    channels[channel] = false
  })

  // master.RegisterEventHandler("323",func(event *spyglass.Event) {
  //   for name,_ := range channels {
  //     channels[name] = true
  //     bot.Join(name)
  //     time.Sleep(time.Duration(time.Millisecond * 250))
  //   }
  // })

  //405 events are triggered when a client has joined too many channels
  // bot.RegisterEventHandler("405",func(event *spyglass.Event) {
  // })

  <- master.Ready

  master.User()
  master.Nick()
  master.Join(command_and_control)
  master.List()

  for id,client := range clients[1:] {
    fmt.Println("Starting up client #",id)
    client.Run()
    <- client.Ready
    client.User()
    client.Nick()
    client.Join(command_and_control)
  }


  for _,client := range clients {
    <- client.Stopped
  }
}






