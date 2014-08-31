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
  "os"
  "path"
  "path/filepath"

  "code.google.com/p/go-sqlite/go1/sqlite3"
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

func NewClient(server string,port string,db *sqlite3.Conn) (bot *spyglass.Bot) {
  bot = spyglass.New(server,port,fake.Username(),fake.Username(),"")
  bot.DB = db
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

  log_path,err := logPath(server)
  if err != nil {
    panic("Cannot create log path!")
  }
  err = makeLogPath(server)
  if err != nil {
    panic("Cannot create log path!")
  }

  var channels map[string]Channel
  channels = make(map[string]Channel)
  var channel Channel

  num_clients,err := strconv.Atoi(n)
  if err != nil || num_clients < 1 {
    panic("Error client.")
  }
  clients := make(Clients,num_clients)

  clients = ([](*spyglass.Bot))(clients)

  var db *sqlite3.Conn

  for i:=0;i<len(clients);i++ {

    //need to handle better  ... but just want it to work right now
    db_name := fmt.Sprintf("_db%d.sqlite3",i)
    db_path := path.Join(log_path,db_name)
    db,_ = sqlite3.Open(db_path)
    db.Exec(`CREATE TABLE events(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      bot VARCHAR(32),
      timestamp INTEGER,
      source VARCHAR(255),
      command VARCHAR(32),
      target VARCHAR(64),
      message TEXT
      );`)
    defer db.Close()

    clients[i] = NewClient(server,port,db)
    client := clients[i]
    client.Connect()
    fmt.Println("[DEBUG] Starting up client #",i)
    client.Run()
    fmt.Println("[DEBUG] Running client: ",client.GetNick())
    <- client.Ready
    fmt.Println("[DEBUG] Client is ready: ",client.GetNick())
    client.User()
    client.Nick()
    client.Join(command_and_control)
    time.Sleep(time.Duration(time.Second * 3))
    defer client.Conn.Close()

  }

  master := clients[0]

  master.Run()

  //react to event 322, which is each listed channel
  oscillator := 0
  minimum,err := strconv.Atoi(m)
  if err != nil {
    panic("couldn't convert minimum to integer!")
  }
  if err != nil {
    panic("couldn't convert sleep duration to integer!")
  }

  master.RegisterEventHandler("322",func(event *spyglass.Event) {
    arguments := event.RawArguments

    args := strings.Split(arguments," ")

    if len(args) <= 1 {
      fmt.Println("[DEBUG] Expected arguments to be length 1, but got ",len(args))
      return
    }

    name := args[1]

    var users int

    s := fmt.Sprintf("[DEBUG 322] RawArguments: %s",arguments)
    fmt.Println(s)

    if len(args) <= 2 {
      fmt.Println("[DEBUG] Expected arguments to be length 2, but got ",len(args))
      return
    }



    // if len(args) > 1 {
      users,err := strconv.Atoi(args[2])
      if err != nil {
        fmt.Println("[DEBUG] Cannot handle event. Args: ",args," users: ",users)
        return
      }
    // } else {
    //   fmt.Println("[DEBUG] Arguments are not long enough.")
    //   return
    // }

    channel = Channel{name,users,false}

    s = fmt.Sprintf("[DEBUG 322] Will join name: %s, users: %d",name,users)
    fmt.Println(s)

    if channel.users > minimum {
      debug_str := fmt.Sprintf("[DEBUG] Joining channel: %s,%d",args[1],args[2])
      fmt.Println(debug_str)
      client := clients[oscillator % len(clients)]
      debug_str = fmt.Sprintf("[DEBUG] client: %s channel.name: %s",client.GetNick(),channel.name)
      fmt.Println(debug_str)
      client.Join(channel.name)
      debug_str = fmt.Sprintf("[DEBUG] %s has joined %d channels.",client.GetNick(),len(client.JoinedChannels) )
      fmt.Println(debug_str)

      // time.Sleep(time.Duration(time.Millisecond * 100))

      oscillator++
    }

    channels[channel.name] = channel
  })

  //Event 263: Server load too heavy.
  // master.RegisterEventHandler("263",func(event *spyglass.Event)) {
  //
  // }

  master.RegisterEventHandler("323",func(event *spyglass.Event) {
    debug_str := fmt.Sprintf("[DEBUG] Done listing channels!")
    fmt.Println(debug_str)
    debug_str = fmt.Sprintf("[DEBUG] Found %d channels!",len(channels))
    fmt.Println(debug_str)
  })

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






