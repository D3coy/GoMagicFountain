package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/d3coy/GoMagicFountain/dice"
)

type jokeRequest struct {
	ID     string
	Joke   string
	Status int
}

// messageCreate is a function that triggers upon a user
// sending a message in the guild with any chat that the
// bot is present in. When the message matches one of the
// programmed commands, it will execute the command
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Fields(m.Content)

	if args[0] == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if args[0] == "!roll" {
		var returnString string
		if len(args) <= 1 {
			returnString = usage
		} else {
			args = args[1:]
			returnString = fmt.Sprintf("%v's rolls: \n", m.Author.Username)
			if text, err := dice.RollDice(args); err != nil {
				returnString = fmt.Sprintln(err.Error())
				returnString += usage
			} else {
				returnString += text
			}
		}
		log.Printf(fmt.Sprintf("User %v (id:%v): \n%v", m.Author.Username, m.Author.ID, returnString))
		s.ChannelMessageSend(m.ChannelID, returnString)
		return
	}

	if args[0] == "!joke" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			log.Println("Error building request:", err)
		}

		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed get joke:", err)
		}
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read contents:", err)
		}

		jokeReq := new(jokeRequest)
		json.Unmarshal(contents, &jokeReq)

		s.ChannelMessageSend(m.ChannelID, jokeReq.Joke)
	}
}

// guildMemberAdd is a function that triggers on the event of a
// user joining the guild. This will auto assign their role to
// the specified "Unknown" Role.
func guildMemberAdd(s *discordgo.Session, gm *discordgo.GuildMemberAdd) {
	err := s.GuildMemberRoleAdd(gm.GuildID, gm.User.ID, "422221833304014863")
	if err != nil {
		log.Println("Role failed to update: ", err)
	}
	message := fmt.Sprintf("User %v (id:%v) has joind the guild!", gm.User.Username, gm.User.ID)
	log.Println(message)
}

func guildMemberUpdate(s *discordgo.Session, gm *discordgo.GuildMemberUpdate) {
	message := fmt.Sprintf("User %v (id:%v) changed their roles to %v", gm.User.Username, gm.User.ID, gm.Roles)
	log.Println(message)
}

func guildMemberRemove(s *discordgo.Session, gm *discordgo.GuildMemberRemove) {
	message := fmt.Sprintf("User %v (id:%v) was removed from the guild", gm.User.Username, gm.User.ID)
	log.Println(message)
}
