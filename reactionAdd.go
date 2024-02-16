package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func reactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	if r.Emoji.Name != "🗑️" {
		return
	}

	proceed := false

	channelID := r.ChannelID

	channel, err := s.State.Channel(r.ChannelID)
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		s.ChannelMessageSendReply(r.ChannelID, "Error getting channel: "+err.Error(), msg.Reference())
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

	if err != nil {
		log.Error().Err(err).Msg("Failed to get message")
		return
	}

	originalMsg, err := s.ChannelMessage(r.ChannelID, msg.MessageReference.MessageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get original message")
		return
	}
	if r.UserID == originalMsg.Author.ID {
		log.Info().
			Str("user", originalMsg.Author.Username).
			Str("msgID", originalMsg.ID).
			Str("channel", originalMsg.ChannelID).
			Msg("User deleted message")

		err = s.ChannelMessageDelete(msg.ChannelID, msg.ID)

		s.ChannelMessageSendReply(msg.ChannelID, "Deleted per your request", originalMsg.Reference())

		if err != nil {
			log.Error().Err(err).Msg("Failed to delete message")
		}
	}
}
