package main

import(
  "github.com/diebels727/spyglass"
  "github.com/diebels727/faker"
  "flag"
  "fmt"
  "strings"
  "strconv"
  "time"
  "gopkg.in/mgo.v2"
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
var mongo string

var fake *faker.Faker

type Channel struct {
  name string
  users int
  joined bool
}

type Datastore struct {
  Session *mgo.Session
  Collection *mgo.Collection
}

func slug(str string) string {
  str = strings.ToLower(str)
  return strings.Replace(str,".","-",-1)
}

func NewDatastore(host string,server string,session *mgo.Session) (*Datastore) {
  local := session.Copy()
  collection := local.DB(slug(server)).C("events")
  datastore := Datastore{local,collection}
  return &datastore
}

func (d *Datastore) Write(event *spyglass.Event) {
  collection := d.Collection
  err := collection.Insert(event)
  if err != nil {
    panic(err)
  }
}

type Clients [](*spyglass.Bot)
type Client *spyglass.Bot
func NewClient(server string,port string) (bot *spyglass.Bot) {
  bot = spyglass.New(server,port,fake.Username(),fake.Username(),"")
  return bot
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
  flag.StringVar(&mongo,"mongo","localhost","Mongo address")
}

var session *mgo.Session  //if package Datastore, move this in to that package

func main() {
  session, err := mgo.Dial(mongo)
  if err != nil {
    panic(err)
  }
  defer session.Close()

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
    datastore := NewDatastore("localhost",server,session)
    client.Datastore = datastore
    defer datastore.Session.Close()

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

  <- master.Ready

  master.User()
  master.Nick()
  master.Join(command_and_control)
  master.List()

  for _,client := range clients {
    <- client.Stopped
  }
}






