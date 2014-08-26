package main

import(
  "github.com/diebels727/spyglass"
  "flag"
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
  bot := spyglass.New(server,port,nick,username,password)
  conn := bot.Connect()
  defer conn.Close()

  bot.Run()

  <- bot.Ready

  bot.User()
  bot.Nick()
  bot.Join(command_and_control)

  <- bot.Stopped
}






