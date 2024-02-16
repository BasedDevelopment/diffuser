package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	configPath = "config.toml"
)

var (
	k      = koanf.New(".")
	parser = toml.Parser()
	token  = ""
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
}

func main() {
	dg, err := discordgo.New("Bot " + k.String("discord.token"))
	if err != nil {
		log.Fatal().Err(err).Msg("error creating Discord session")
	}

	dg.Identify.Intents |= discordgo.IntentsGuildMessages

	dg.AddHandler(newMsg)

	err = dg.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("error opening connection")
		return
	}

	log.Info().Msg("Bot up")

	//go getToken()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}
