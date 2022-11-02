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
		_, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, command)
		if err != nil {
			log.Printf("CANNOT CREATE '%v' COMMAND: %v", command.Name, err)
		}
	}

}

var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"setup": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Creating an admin only channel for the logs.
		permissions := []*discordgo.PermissionOverwrite{{ID: interaction.GuildID, Allow: 0x0000000000000008, Deny: 0x0000000000000400}}
		session.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
			Name:                 "asg-logs",
			Type:                 0,
			Topic:                "These are the logs for the All Systems Go bot.",
			Position:             0,
			PermissionOverwrites: permissions,
			ParentID:             "",
			NSFW:                 false,
		})
	},
	"register": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		token := interaction.ApplicationCommandData().Options[0].StringValue()
		query := fmt.Sprintf(`INSERT INTO tokens(discord_id, pk_token) VALUES(%v, "%s")`,
			interaction.Member.User.ID, token)
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT PLACE SYSTEM TOKEN IN DATABASE:\n\t%v", Red, Reset, err.Error())

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "%ERROR - COULD NOT PLACE SYSTEM TOKEN IN DATABASE:\n\t" + err.Error(),
					Flags:   1 << 6,
				},
			})

			return
		}
		log.Printf("%vSUCCESS%v - REGISTERED SYSTEM INTO DATABASE:\n\t%v", Green, Reset, result)

		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS:             false,
				Content:         "You are successfully registered!",
				Components:      nil,
				Embeds:          nil,
				AllowedMentions: nil,
				Files:           nil,
				Flags:           1 << 6,
				Choices:         nil,
				CustomID:        "",
				Title:           "",
			},
		})
	},
}
