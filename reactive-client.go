package main

import(
  "fmt"
  "github.com/diebels727/spyglass"
  "time"
)


func main() {

  server := "localhost"
  m := spyglass.New(server,"6667","nick89122","jiggly101001","")
  conn := m.Connect()
  defer conn.Close()
  m.Run()

  user_cmd := fmt.Sprintf("USER %s 8 * :%s\r\n", "nick89122", "nick89122")
  nick_cmd := fmt.Sprintf("NICK %s\r\n", "nick89122")
  fmt.Println("[TestInit] Sending USER command")
  m.Send(user_cmd)
  fmt.Println("[TestInit] Sending NICK command")
  m.Send(nick_cmd)
  fmt.Println("[TestInit] Sending JOIN command")
  m.Join("#cinch-bots")
  fmt.Println("[TestInit] Sending JOIN command")
  m.Send("JOIN #foofoo\r\n")

  fmt.Println("Sleeping for 5 seconds")
  time.Sleep(time.Second * 5)

  for {
    time.Sleep(time.Second * 1)
  }
}






