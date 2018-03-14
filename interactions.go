package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
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

type diceFlag struct {
	Action int
	IsTrue bool
}

type dieRoll struct {
	DieType int
	Value   int
}

type rolls struct {
	Dice []dieRoll
}

func (r rolls) RollAll() {
	for _, die := range r.Dice {
		rand.Seed(time.Now().Unix())
		die.Value = rand.Intn(die.DieType)
	}
}

func (r rolls) OnlyTop(num int) {
	r.Sort()
	if num < len(r.Dice) {
		r.Dice = r.Dice[:num]
	}
}

func (r rolls) String() string {
	var returnText string
	for _, die := range r.Dice {
		returnText += fmt.Sprintf("\t d%v: %v\n", die.DieType, die.Value)
	}
	return returnText
}

func (r rolls) Sort() {
	sort.Slice(r.Dice, func(i, j int) bool {
		return r.Dice[i].Value > r.Dice[j].Value
	})
}

func (r rolls) GetType(dieType int) []dieRoll {
	var dice []dieRoll
	for _, die := range r.Dice {
		if die.DieType == dieType {
			dice = append(dice, die)
		}
	}
	sort.Slice(r.Dice, func(i, j int) bool {
		return dice[i].Value > dice[j].Value
	})
	return dice
}

func (r rolls) Total() int {
	total := 0
	for _, die := range r.Dice {
		total += die.Value
	}
	return total
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
			returnString = fmt.Sprintf("%v's rolls: \n", m.Author.Username)
			tFlag, sFlag, errFlag := false, false, false
			args = args[1:]
			var tValue int
			currentRolls := new(rolls)

			for len(args) > 0 {
				cmd := args[0]
				if cmd == "-t" {
					if v, err := strconv.Atoi(args[1]); err == nil {
						tFlag = true
						tValue = v
						args = append(args[:1], args[2:]...)
					} else {
						returnString = "-t flag must have a number after!\n"
						errFlag = true
						break
					}
				} else if cmd == "-s" {
					sFlag = true
				} else if strings.HasPrefix(cmd, "d") {
					rollQ := 1
					if 1 < len(args) {
						if i, err := strconv.Atoi(args[1]); err == nil {
							if i < 10 {
								rollQ = i
								args = append(args[:1], args[2:]...)
							} else if i >= 10 {
								returnString = "You can only have up to 9 rolls at a time!\n"
								errFlag = true
								break
							}
						}
					}

					if i, err := strconv.Atoi(cmd[1:]); err == nil {
						for count := rollQ; count > 0; count-- {
							roll := dieRoll{i, rand.Intn(i) + 1}
							currentRolls.Dice = append(currentRolls.Dice, roll)
						}
					} else {
						log.Printf("Error converting string to int: %v\n", err)
						returnString = fmt.Sprintf("%v is not a number!~\n", args[1][1:])
						errFlag = true
						break
					}

				} else {
					returnString = fmt.Sprintf("Invalid option: %v\n", cmd)
					errFlag = true
					break
				}

				args = args[1:]
			}
			if !errFlag && (len(currentRolls.Dice) != 0) {
				currentRolls.RollAll()
				if tFlag {
					currentRolls.OnlyTop(tValue)
				} else if sFlag {
					currentRolls.Sort()
				}
				returnString += fmt.Sprint(currentRolls)
				returnString += fmt.Sprintf("Total: %v", currentRolls.Total())
			} else if !errFlag {
			} else if len(currentRolls.Dice) == 0 {
				returnString = "There are no dice to roll!\n"
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
