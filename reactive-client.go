package main

import(
  "github.com/diebels727/spyglass"
  // "github.com/gorilla/mux"
  "github.com/diebels727/faker"
  "flag"
  // "net/http"
  "fmt"
  "strings"
  "strconv"
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
var n string
var m string

var fake *faker.Faker

type Channel struct {
  name string
  users int
  joined bool
}

func init() {
  flag.StringVar(&server,"server","irc.freenode.net","IRC server FQDN")
  flag.StringVar(&port,"port","6667","IRC server port number")
  flag.StringVar(&nick,"nick","","Name of the bot visible on IRC channel")
  flag.StringVar(&username,"username","","Username to login with to IRC")
  flag.StringVar(&password,"password","","Password for the IRC server")
  flag.StringVar(&command_and_control,"command_and_control","#spyglass-c&c","Command and control IRC channel")
  flag.StringVar(&n,"n","1","Number of clients; minimum is one")
  flag.StringVar(&m,"m","50","Minimum number of users per channel")
}

type Clients [](*spyglass.Bot)

type Client *spyglass.Bot

func NewClient(server string,port string) (*spyglass.Bot) {
  return spyglass.New(server,port,fake.Username(),fake.Username(),"")
}

func main() {
  fake = faker.New()
  flag.Parse()

  var channels map[string]Channel
  channels = make(map[string]Channel)
  var channel Channel

  num_clients,err := strconv.Atoi(n)
  if err != nil || num_clients < 1 {
    panic("Error client.")
  }
  clients := make(Clients,num_clients)

  clients = ([](*spyglass.Bot))(clients)
  for i:=0;i<len(clients);i++ {
    clients[i] = NewClient(server,port)
    clients[i].Connect()
    time.Sleep(time.Duration(time.Second * 3))
    defer clients[i].Conn.Close()
  }

  master := clients[0]

  master.Run()

  //react to event 322, which is each listed channel
  master.RegisterEventHandler("322",func(event *spyglass.Event) {
    arguments := event.RawArguments
    args := strings.Split(arguments," ")
    users,err := strconv.Atoi(args[2])
    if err != nil {
      panic("Num Users conversion error")
    }
    channel = Channel{args[1],users,false}
    channels[channel.name] = channel
  })

  //Event 263: Server load too heavy.
  // master.RegisterEventHandler("263",func(event *spyglass.Event)) {
  //
  // }

  master.RegisterEventHandler("323",func(event *spyglass.Event) {
    minimum,err := strconv.Atoi(m)
    if err != nil {
      panic("couldn't convert minimum to integer!")
    }

    for _,channel := range channels {
      if channel.users >
      if channel.users > 100 {
        num_100++
      }
      if channel.users > 500 {
        num_500++
      }
      if channel.users > 50 {
        num_50++
      }
    }
    fmt.Println("[NUM] num_channels: ",num_channels)
    fmt.Println("[NUM] num_100: ",num_100)
    fmt.Println("[NUM] num_500: ",num_500)
    fmt.Println("[NUM] num_50: ",num_50)

    // for name,_ := range channels {
    //   channels[name] = true
    //   bot.Join(name)
    //   time.Sleep(time.Duration(time.Millisecond * 250))
    // }
  })

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






