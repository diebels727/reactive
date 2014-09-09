package main

import(
  "github.com/diebels727/spyglass"
  "github.com/diebels727/faker"
  "flag"
  "fmt"
  "strings"
  "strconv"
  "time"
  "os"
  "path"
  "path/filepath"

  "gopkg.in/mgo.v2"
  // "gopkg.in/mgo.v2/bson"
)


var server string
var port string
var nick string
var username string
var password string
var command_and_control string
var n string
var m string
var s string

var fake *faker.Faker

type Channel struct {
  name string
  users int
  joined bool
}

type Datastore struct {
  Session *mgo.Session
}

func NewDatastore(host string) (*Datastore) {
  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()
  datastore := Datastore{session}
  return &datastore
}

func (d *Datastore) Write(event *spyglass.Event) {
  fmt.Println("[DATASTORE]:",event)
}

func init() {
  flag.StringVar(&server,"server","localhost","IRC server FQDN")
  flag.StringVar(&port,"port","6667","IRC server port number")
  flag.StringVar(&nick,"nick","","Name of the bot visible on IRC channel")
  flag.StringVar(&username,"username","","-Username to login with to IRC")
  flag.StringVar(&password,"password","","Password for the IRC server")
  flag.StringVar(&command_and_control,"command_and_control","#spyglass-c&c","Command and control IRC channel")
  flag.StringVar(&n,"n","1","Number of clients; minimum is one")
  flag.StringVar(&m,"m","50","Minimum number of users per channel")
  flag.StringVar(&s,"s","100","Amount of time to sleep between channel joins")
}

type Clients [](*spyglass.Bot)

type Client *spyglass.Bot

func NewClient(server string,port string) (bot *spyglass.Bot) {
  bot = spyglass.New(server,port,fake.Username(),fake.Username(),"")
  return bot
}

func toSlug(str string) string {
  return strings.Replace(str,".","-",-1)
}

func logPath(server string) (string,error) {
  slug := strings.Replace(server,".","-",-1)
  log_path := path.Join(slug)
  abs_log_path,err := filepath.Abs(log_path)
  return abs_log_path,err
}

func makeLogPath(server string) error {
  log_path,err := logPath(server)
  file_mode := os.FileMode(0777)
  err = os.MkdirAll(log_path,file_mode)
  return err
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
    client := clients[i]
    client.Datastore = NewDatastore("localhost")

    time.Sleep(time.Duration(time.Second * 10))
    client.Connect()
    client.Run()
    <- client.Ready
    client.User()
    client.Nick()
    client.Join(command_and_control)
    //client.Register()
    time.Sleep(time.Duration(time.Second * 10))
    defer client.Conn.Close()
  }

  master := clients[0]

  master.Run()

  //react to event 322, which is each listed channel
  minimum,err := strconv.Atoi(m)
  if err != nil {
    panic("couldn't convert minimum to integer!")
  }
  if err != nil {
    panic("couldn't convert sleep duration to integer!")
  }

  oscillator := 0
  master.RegisterEventHandler("322",func(event *spyglass.Event) {
    arguments := event.RawArguments

    args := strings.Split(arguments," ")

    if len(args) <= 1 {
      fmt.Println("[DEBUG] Expected arguments to be length 1, but got ",len(args))
      return
    }

    if len(args) <= 2 {
      fmt.Println("[DEBUG] Expected arguments to be length 2, but got ",len(args))
      return
    }

    name := args[1]

    users,err := strconv.Atoi(args[2])
    if err != nil {
      fmt.Println("[DEBUG] Cannot handle event. Args: ",args," users: ",users)
      return
    }

    channel = Channel{name,users,false}

    if channel.users > minimum {
      client := clients[oscillator % len(clients)]
      client.Join(channel.name)
      oscillator++
    }

    channels[channel.name] = channel
  })

  //Event 263: Server load too heavy.
  // master.RegisterEventHandler("263",func(event *spyglass.Event)) {
  //
  // }

  // master.RegisterEventHandler("323",func(event *spyglass.Event) {
  //
  // })

  //405 events are triggered when a client has joined too many channels
  // bot.RegisterEventHandler("405",func(event *spyglass.Event) {
  // })

  <- master.Ready

  master.User()
  master.Nick()
  master.Join(command_and_control)
  master.List()

  for _,client := range clients {
    <- client.Stopped
  }
}






