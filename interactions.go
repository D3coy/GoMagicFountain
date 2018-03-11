package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "ping") {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.HasPrefix(m.Content, "pong") {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
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
