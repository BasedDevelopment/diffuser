package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func lsCkpt(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info().Str("user", m.Author.Username).Msg("lsCkpt")
	reqUrl := k.String("api.url") + "/sdapi/v1/sd-models"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response")
		return
	}
	type respJsonStr struct {
		Title string `json:"title"`
	}
	var respJson []respJsonStr
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling response")
		return
	}
	outMsg := "```"
	for _, model := range respJson {
		outMsg += model.Title + "\n"
	}
	outMsg += "```"
	s.ChannelMessageSendReply(m.ChannelID, outMsg, m.Reference())
}
