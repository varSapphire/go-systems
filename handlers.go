package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func clientJoin(session *discordgo.Session, guild *discordgo.GuildCreate) {
	// Looping through the commands array and registering them.
	// https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.ApplicationCommandCreate
	for _, command := range commands {
		registeredCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, command)
		if err != nil {
			log.Printf("CANNOT CREATE '%v' COMMAND: %v", command.Name, err)
		}

		// Executing a query to register commands into the database for record keeping.
		query := fmt.Sprintf(`INSERT INTO commands(guild_id, registered_commands) VALUES(%v, %v)`,
			guild.ID, registeredCommand.ID)
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT PLACE COMMAND IN DATABASE:\n\t%v", Red, Reset, err)
			return
		}
		log.Printf("%vSUCCESS%v - PLACED COMMAND INTO DATABASE:\n\t%v", Green, Reset, result)
	}

	// Creating an admin only channel for the logs.
	permissions := []*discordgo.PermissionOverwrite{{ID: guild.ID, Allow: 0x0000000000000008, Deny: 0x0000000000000400}}
	session.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:                 "asg-logs",
		Type:                 0,
		Topic:                "These are the logs for the All Systems Go bot.",
		Position:             0,
		PermissionOverwrites: permissions,
		ParentID:             "",
		NSFW:                 false,
	})
}

var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"hello": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello, world!",
			},
		})
	},
}
