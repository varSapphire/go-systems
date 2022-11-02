package main

import (
	"fmt"
	"github.com/starshine-sys/pkgo/v2"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
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
	// This command creates a logging channel for All Systems Go.
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
	// This command registers a users PluralKit token into a database so that they can use the commands.
	"register": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Setting up a delayed response.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
			},
		})

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

			content := fmt.Sprintf("%ERROR - COULD NOT PLACE SYSTEM TOKEN IN DATABASE:\n\t%v", err.Error())
			session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})

			return
		}
		log.Printf("%vSUCCESS%v - REGISTERED SYSTEM INTO DATABASE:\n\t%v", Green, Reset, result)

		content := fmt.Sprintf("You are successfully registered!")
		session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})
	},
	// This command proxies a message using a PluralKit proxy.
	"auto_proxy_message": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var pkToken string

		var webhook *discordgo.Webhook

		var embeds []*discordgo.MessageEmbed

		// Setting up a delayed response.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
			},
		})

		// Getting the user's PK token from the database.
		query := fmt.Sprintf(`SELECT pk_token FROM tokens WHERE discord_id = %v`, interaction.Member.User.ID)
		err := db.QueryRow(query).Scan(&pkToken)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT RETRIEVE SYSTEM TOKEN FROM DATABASE:\n\t%v", Red, Reset, err.Error())

			content := fmt.Sprintf("%ERROR - COULD NOT RETRIEVE SYSTEM TOKEN FROM DATABASE:\n\t%v", err.Error())
			session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})

			return
		}

		// Authenticating a PK session to grab member information.
		pk := pkgo.New(pkToken)
		front, err := pk.Fronters("@me")
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT RETRIEVE MEMBER INFORMATION FROM PLURALKIT:\n\t%v", Red, Reset, err.Error())

			content := fmt.Sprintf("%ERROR - COULD NOT RETRIEVE MEMBER INFORMATION FROM PLURALKIT:\n\t%v", err.Error())
			session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})

			return
		}

		system, err := pk.System("@me")
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT RETRIEVE SYSTEM INFORMATION FROM PLURALKIT:\n\t%v", Red, Reset, err.Error())

			content := fmt.Sprintf("%ERROR - COULD NOT RETRIEVE SYSTEM INFORMATION FROM PLURALKIT:\n\t%v", err.Error())
			session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})

			return
		}

		memberName := front.Members[0].Name
		memberAvatarURL := front.Members[0].AvatarURL
		systemTag := system.Tag

		webhookFlag, webhookID, webhookToken := checkForWebhook(session, interaction.Interaction)
		memberMessage, messageFlag := checkForMessage(interaction.Interaction)
		loggingChannelID, loggingFlag := checkforLogging(session, interaction.Interaction)

		embedAuthor := discordgo.MessageEmbedAuthor{
			URL:          "",
			Name:         memberName,
			IconURL:      memberAvatarURL,
			ProxyIconURL: "",
		}

		embedFields := []*discordgo.MessageEmbedField{
			{
				Name:   "User:",
				Value:  fmt.Sprintf("<@%v>", interaction.Member.User.ID),
				Inline: true,
			},
			{
				Name:   "Channel:",
				Value:  fmt.Sprintf("<#%v>", interaction.ChannelID),
				Inline: true,
			},
		}

		if !webhookFlag {
			webhook, err = session.WebhookCreate(
				interaction.ChannelID,
				"All Systems Go Proxy Webhook",
				"https://cdn.discordapp.com/attachments/990405675022700567/1035705744470843402/unknown.png",
			)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT CREATE WEBHOOOK:\n\t%v", Red, Reset, err.Error())

				content := fmt.Sprintf("%ERROR - COULD NOT CREATE WEBHOOOK:\n\t%v", err.Error())
				session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})
				return
			}

			webhookFlag = true
			webhookID = webhook.ID
			webhookToken = webhook.Token
		}

		if messageFlag {
			session.WebhookExecute(webhookID, webhookToken, true, &discordgo.WebhookParams{
				Content:         memberMessage,
				Username:        memberName + systemTag,
				AvatarURL:       memberAvatarURL,
				TTS:             false,
				Files:           nil,
				Components:      nil,
				Embeds:          nil,
				AllowedMentions: nil,
				Flags:           0,
			})

			t, _ := discordgo.SnowflakeTimestamp(interaction.ID)
			ts, _ := time.Parse("2006-01-02T15:04:05.000+0000", t.String())

			embeds = append(embeds, &discordgo.MessageEmbed{
				URL:         "",
				Type:        "",
				Title:       "",
				Description: memberMessage,
				Timestamp:   ts.String(),
				Color:       0,
				Footer:      nil,
				Image:       nil,
				Thumbnail:   nil,
				Video:       nil,
				Provider:    nil,
				Author:      &embedAuthor,
				Fields:      embedFields,
			})

		}

		if loggingFlag {
			session.ChannelMessageSendComplex(loggingChannelID, &discordgo.MessageSend{
				Content:         "",
				Embeds:          embeds,
				TTS:             false,
				Components:      nil,
				Files:           nil,
				AllowedMentions: nil,
				Reference:       nil,
				File:            nil,
				Embed:           nil,
			})
		}
	},
}
