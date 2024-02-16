package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func newMsg(s *discordgo.Session, m *discordgo.MessageCreate) {

	proceed := false

	channelID := m.ChannelID

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Error getting channel: "+err.Error(), m.Reference())
		return
	}
	if channel.IsThread() {
		channelID = channel.ParentID
	}

	for _, channel := range k.Strings("discord.msgChan") {
		if channelID == channel {
			proceed = true
		}
	}

	if !proceed {
		return
	}

	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, user := range k.Strings("discord.bannedUsers") {
		if m.Author.ID == user {
			return
		}
	}

	if m.Content == "lsckpt!" {
		lsCkpt(s, m)
	}

	if len(m.Content) < 8 {
		return
	}

	if m.Content[:8] == "diffuse!" {

		msg := m.Content[8:]
		if msg[0] == ' ' {
			msg = msg[1:]
		}

		user := m.Author.Username + "#" + m.Author.Discriminator
		for _, word := range k.Strings("discord.bannedWords") {
			if strings.Contains(msg, word) {
				log.Warn().
					Str("user", user).
					Str("msg", msg).
					Str("word", word).
					Msg("Banned word detected")
				if _, err := s.ChannelMessageSendReply(m.ChannelID, "Banned word detected, dropping request", m.Reference()); err != nil {
					log.Error().Err(err).Msg("Failed to send message")
				}
				return
			}
		}
		draw(s, m, msg)
	}
}
