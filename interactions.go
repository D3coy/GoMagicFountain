package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
			returnString = "Usage: !roll [d#] [quantity]\n\nExamples:\n\t- !roll d4\n\t- !roll d10 3"
		} else {
			rand.Seed(time.Now().Unix())
			returnString += fmt.Sprintf("%v's rolls: \n", m.Author.Username)
			cmdList := args[1:]

			for index, cmd := range cmdList {
				if strings.HasPrefix(cmd, "d") {
					rollQ := 1
					if nextIndex := index + 1; nextIndex < len(cmdList) {
						if i, err := strconv.Atoi(cmdList[nextIndex]); err == nil {
							if i < 10 {
								rollQ = i
							} else if i > 10 {
								returnString = "You can only have between 1 and 9 rolls at a time!"
								break
							}
						}
					}

					i, err := strconv.Atoi(cmd[1:])
					if err != nil {
						log.Printf("Error converting string to int: %v\n", err)
						returnString = fmt.Sprintf("%v is not a number!~", args[1][1:])
						break
					} else {
						total := 0
						for count := rollQ; count > 0; count-- {
							result := rand.Intn(i) + 1
							total += result
							returnString += fmt.Sprintf("\t%v: %v\n", cmd, result)
						}
						if rollQ != 1 {
							returnString += fmt.Sprintf("Total (%v): %v\n\n", cmd, total)
						}
					}

				} else {
					if _, err := strconv.Atoi(cmd); err == nil {
						continue
					} else {
						returnString = fmt.Sprintf("Invalid option: %v", cmd)
						break
					}
				}
			}
		}

		log.Printf(fmt.Sprintf("User %v (id:%v): \n%v", m.Author.Username, m.Author.ID, returnString))
		s.ChannelMessageSend(m.ChannelID, returnString)
	}

	if args[0] == "!joke" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			log.Printf("Error building request: %v\n", err)
		}

		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed get joke: %v\n", err)
		}
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read contents: %v\n", err)
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
