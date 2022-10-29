package main

import (
	"database/sql"
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

	//// Creating an admin only channel for the logs.
	//permissions := []*discordgo.PermissionOverwrite{{ID: guild.ID, Allow: 0x0000000000000008, Deny: 0x0000000000000400}}
	//session.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
	//	Name:                 "asg-logs",
	//	Type:                 0,
	//	Topic:                "These are the logs for the All Systems Go bot.",
	//	Position:             0,
	//	PermissionOverwrites: permissions,
	//	ParentID:             "",
	//	NSFW:                 false,
	//})
}

var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"register": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var memberAvatarURL string
		memberName := interaction.ApplicationCommandData().Options[0].StringValue()
		memberIdentifier := interaction.ApplicationCommandData().Options[1].StringValue()
		for _, attachment := range interaction.ApplicationCommandData().Resolved.Attachments {
			memberAvatarURL = attachment.URL
		}

		// Executing a query to register commands into the database for record keeping.
		query := fmt.Sprintf(`INSERT INTO MEMBERS(discord_id, name, proxy, avatar) VALUES(%v, "%s", "%s", "%s")`,
			interaction.Member.User.ID, memberName, memberIdentifier, memberAvatarURL)
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT PLACE MEMBER IN DATABASE:\n\t%v", Red, Reset, err.Error())
			return
		}
		log.Printf("%vSUCCESS%v - PLACED MEMBER INTO DATABASE:\n\t%v", Green, Reset, result)

		log.Println(memberName, memberIdentifier, memberAvatarURL)

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
	"proxy": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var flag bool
		var webhookID string
		var webhookToken string
		var memberName string
		var memberAvatarURL string

		memberIdentifier := interaction.ApplicationCommandData().Options[0].StringValue()
		memberMessage := interaction.ApplicationCommandData().Options[1].StringValue()

		// Checking to see if an All Systems go webhook exists for this channel.
		query := fmt.Sprintf(`SELECT webhook_id FROM webhooks WHERE guild_id = %v AND channel_id = %v`,
			interaction.GuildID, interaction.GuildID)

		err = db.QueryRow(query).Scan(&webhookID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT RETRIEVE WEBHOOK FROM DATABASE:\n\t%v", Red, Reset, err)
		}

		// Looping through the webhooks in the channel to see if one already exists.
		webhooks, err := session.ChannelWebhooks(interaction.ChannelID)
		for _, webhook := range webhooks {
			if webhook.Name == "All Systems Go Proxy Webhook" {
				flag = true
				webhookID = webhook.ID
				webhookToken = webhook.Token
			}
		}

		if flag {
			// Proxying the message using the channel's webhook.
			// Getting the information of the member from the DB.
			query = fmt.Sprintf(`SELECT name, avatar FROM members WHERE discord_id = %v AND proxy = "%v";`,
				interaction.Member.User.ID, memberIdentifier)

			err = db.QueryRow(query).Scan(&memberName, &memberAvatarURL)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT RETREIVE MEMBER FROM DATABASE:\n\t%v", Red, Reset, err.Error())
				return
			}

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			})

			WebhookExecuteSimple(session, webhookID, webhookToken, memberMessage, memberName, memberAvatarURL)

			session.InteractionResponseDelete(interaction.Interaction)

		} else if !flag {
			webhook, err := session.WebhookCreate(
				interaction.ChannelID,
				"All Systems Go Proxy Webhook",
				"https://cdn.discordapp.com/attachments/990405675022700567/1035705744470843402/unknown.png",
			)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT CREATE WEBHOOOK:\n\t%v", Red, Reset, err.Error())
				return
			}

			webhookID = webhook.ID
			webhookToken = webhook.Token

			query = fmt.Sprintf(`SELECT name, avatar FROM members WHERE discord_id = %v AND proxy = "%v";`,
				interaction.Member.User.ID, memberIdentifier)

			err = db.QueryRow(query).Scan(&memberName, &memberAvatarURL)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT RETREIVE MEMBER FROM DATABASE:\n\t%v", Red, Reset, err.Error())
				return
			}

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			})

			WebhookExecuteSimple(session, webhookID, webhookToken, memberMessage, memberName, memberAvatarURL)

			session.InteractionResponseDelete(interaction.Interaction)

		}
	},
}
