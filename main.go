package main

import (
	"database/sql"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

// Config Creating a structure to hole the Discord token along with anything else needed for configuration.
type Config struct {
	DiscordToken  string
	MySQLUsername string
	MySQLPassword string
}

// Creating a global variable to hold that configuration structure.
var config Config

// Creating a global variable to hold the database connection.
var db *sql.DB
var err error

// Main function that will hold the entire logic loop.
func main() {
	// Retrieve the tokens from the tokens.json file.
	var tokens []byte
	tokens, err = os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT READ 'config.json' FILE:\n\t%v", Red, Reset, err)
	}

	// Unmarshal the tokens from the config file.
	err = json.Unmarshal(tokens, &config)
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT UNMARSHAL 'config.json' FILE:\n\t%v", Red, Reset, err)
	}

	log.Printf("%vSTARTING DATABASE CONNECTION.%v", Blue, Reset)

	// Opening a connection to the database.
	// Set up the parameters for the database connection.
	sqlConfiguration := mysql.Config{
		User:   config.MySQLUsername,
		Passwd: config.MySQLPassword,
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "systems",
	}

	// Open a connection to the database.
	db, err = sql.Open("mysql", sqlConfiguration.FormatDSN())
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT CONNECT TO DATABASE:\n\t%v", Red, Reset, err)
	}

	log.Printf("%vBOT IS STARTING UP.%v", Blue, Reset)

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatalf("%vERROR%v - PROBLEM CREATING DISCORD SESSION:\n\t%v", Red, Reset, err)
	}

	// Identify that we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatalf("%vERROR%v - PROBLEM OPENING WEBSOCKET:\n\t%v", Red, Reset, err)
	}

	session.AddHandler(clientJoin)

	// Looping through the array of interaction handlers and adding them to the session.
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Printf("%vBOT IS NOW RUNNING.%v", Blue, Reset)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Closing the session.
	err = session.Close()
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT CLOSE SESSION:\n\t%v", Red, Reset, err)
	}

	// Closing the connection to the database.
	err = db.Close()
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT CLOSE DATABASE:\n\t%v", Red, Reset, err)
	}
}
