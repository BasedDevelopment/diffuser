package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func draw(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	s.MessageReactionAdd(m.ChannelID, m.Reference().MessageID, "⌛")

	log.Info().
		Str("user", m.Author.Username).
		Str("prompt", msg).
		Msg("draw")

	reqUrl := k.String("api.url") + "/sdapi/v1/txt2img"

	var prompt, checkpoint string
	msgParts := strings.Split(msg, "?")

	if len(msgParts) < 2 {
		prompt = msgParts[0]
		checkpoint = "wd-v1-3-full-opt.ckpt [3e1a125f]"
	} else {
		prompt = msgParts[0]
		checkpoint = msgParts[1]
	}

	reqBody := map[string]interface{}{
		"prompt": prompt,
		"override_settings": map[string]string{
			"sd_model_checkpoint": "\"" + checkpoint + "\"",
		},
	}

	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		log.Error().
			Err(err).
			Msg("json marshal")
		return
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		log.Error().
			Err(err).
			Msg("http req")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("cookie jar")
		return
	}

	jar.SetCookies(req.URL, []*http.Cookie{
		{
			Name:  "access-token",
			Value: token,
		},
	})

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("http req")
		return
	}

	defer resp.Body.Close()

	// make resp body into string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBody := buf.String()

	var respBodyParsed map[string]interface{}
	err = json.NewDecoder(strings.NewReader(respBody)).Decode(&respBodyParsed)
	if err != nil {
		log.Error().
			Err(err).
			Str("resp", respBody).
			Str("req", string(reqJson)).
			Msg("json decode")
		return
	}

	// print index 0 of images
	img := (respBodyParsed["images"].([]interface{})[0])

	// Convirt image from bae64 to bytes data:image/png;base64
	imgData := strings.Split(img.(string), ",")[1]
	unbased, _ := base64.StdEncoding.DecodeString(imgData)

	outMsg := "Prompt: " + prompt + ", Checkpoint: " + checkpoint + respBodyParsed["info"].(string)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: strings.NewReader(string(unbased)),
			},
		},
		Content: outMsg,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	})
}
