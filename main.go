package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables required from command line parameters
var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error when creating Discord session: %b\n", err)
		return
	}
	defer discord.Close()

	// Register Handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(guildMemberAdd)
	discord.AddHandler(guildMemberUpdate)
	discord.AddHandler(guildMemberRemove)

	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v\n", err)
		return
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
