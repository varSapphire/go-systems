package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
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
	"proxy_member": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var messageFlag bool
		var webhookFlag bool
		var webhookID string
		var webhookToken string
		var memberName string
		var memberIdentifier string
		var memberMessage string
		var memberAvatarURL string
		var loggingFlag bool
		var loggingChannelID string

		channels, err := session.GuildChannels(interaction.GuildID)
		for _, channel := range channels {
			if channel.Name == "asg-logs" {
				loggingChannelID = channel.ID
				loggingFlag = true
			}
		}

		for _, option := range interaction.ApplicationCommandData().Options {
			if option.Type == discordgo.ApplicationCommandOptionString {
				memberMessage = option.StringValue()
				messageFlag = true
			}
		}

		memberIdentifier = interaction.ApplicationCommandData().Options[0].StringValue()
		//if len(interaction.ApplicationCommandData().Options) > 0 {
		//	if interaction.ApplicationCommandData().Options[1].Type == discordgo.ApplicationCommandOptionString {
		//		memberMessage = interaction.ApplicationCommandData().Options[1].StringValue()
		//	}
		//}

		// Looping through the webhooks in the channel to see if one already exists.
		webhooks, _ := session.ChannelWebhooks(interaction.ChannelID)
		for _, webhook := range webhooks {
			if webhook.Name == "All Systems Go Proxy Webhook" {
				webhookFlag = true
				webhookID = webhook.ID
				webhookToken = webhook.Token
			}
		}

		if webhookFlag {
			// Proxying the message using the channel's webhook.
			// Getting the information of the member from the DB.
			query := fmt.Sprintf(`SELECT name, avatar FROM members WHERE discord_id = %v AND proxy = "%v";`,
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

			if messageFlag {
				WebhookExecuteSimple(session, webhookID, webhookToken, memberMessage, memberName, memberAvatarURL)
			}
			if interaction.ApplicationCommandData().Resolved != nil {
				for _, attachment := range interaction.ApplicationCommandData().Resolved.Attachments {
					WebhookExecuteSimple(session, webhookID, webhookToken, attachment.URL, memberName, memberAvatarURL)
				}
			}

			embeds := []*discordgo.MessageEmbed{}
			embedAuthor := discordgo.MessageEmbedAuthor{
				Name: fmt.Sprintf("%v#%v || %v",
					interaction.Member.User.Username, interaction.Member.User.Discriminator, memberName),
			}
			embedThumbnail := discordgo.MessageEmbedThumbnail{
				URL: interaction.Member.User.AvatarURL(""),
			}

			if memberMessage != "" {
				embeds = append(embeds, &discordgo.MessageEmbed{
					Description: memberMessage,
					Color:       0,
					Image:       nil,
					Thumbnail:   &embedThumbnail,
					Author:      &embedAuthor,
				})
			}

			if interaction.ApplicationCommandData().Resolved != nil {
				for _, attachment := range interaction.ApplicationCommandData().Resolved.Attachments {
					embedImage := discordgo.MessageEmbedImage{
						URL: attachment.URL,
					}

					embeds = append(embeds, &discordgo.MessageEmbed{
						Description: memberMessage,
						Color:       0,
						Image:       &embedImage,
						Thumbnail:   &embedThumbnail,
						Author:      &embedAuthor,
					})
				}
			}

			if loggingFlag {
				session.ChannelMessageSendComplex(loggingChannelID, &discordgo.MessageSend{
					Embeds: embeds,
				})

				session.ChannelMessageSend(loggingChannelID, fmt.Sprintf("Sent on: <t:%d:D> at <t:%d:T>",
					time.Now().Unix(), time.Now().Unix()))
			}

			session.InteractionResponseDelete(interaction.Interaction)
		} else if !webhookFlag {
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

			query := fmt.Sprintf(`SELECT name, avatar FROM members WHERE discord_id = %v AND proxy = "%v";`,
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

			if messageFlag {
				WebhookExecuteSimple(session, webhookID, webhookToken, memberMessage, memberName, memberAvatarURL)
			}
			if interaction.ApplicationCommandData().Resolved != nil {
				for _, attachment := range interaction.ApplicationCommandData().Resolved.Attachments {
					WebhookExecuteSimple(session, webhookID, webhookToken, attachment.URL, memberName, memberAvatarURL)
				}
			}

			embeds := []*discordgo.MessageEmbed{}
			embedAuthor := discordgo.MessageEmbedAuthor{
				Name: fmt.Sprintf("%v#%v || %v",
					interaction.Member.User.Username, interaction.Member.User.Discriminator, memberName),
			}
			embedThumbnail := discordgo.MessageEmbedThumbnail{
				URL: interaction.Member.User.AvatarURL(""),
			}

			if memberMessage != "" {
				embeds = append(embeds, &discordgo.MessageEmbed{
					Description: memberMessage,
					Color:       0,
					Image:       nil,
					Thumbnail:   &embedThumbnail,
					Author:      &embedAuthor,
				})
			}

			if interaction.ApplicationCommandData().Resolved != nil {
				for _, attachment := range interaction.ApplicationCommandData().Resolved.Attachments {
					embedImage := discordgo.MessageEmbedImage{
						URL: attachment.URL,
					}

					embeds = append(embeds, &discordgo.MessageEmbed{
						Description: memberMessage,
						Color:       0,
						Image:       &embedImage,
						Thumbnail:   &embedThumbnail,
						Author:      &embedAuthor,
					})
				}
			}

			if loggingFlag {
				session.ChannelMessageSendComplex(loggingChannelID, &discordgo.MessageSend{
					Embeds: embeds,
				})

				session.ChannelMessageSend(loggingChannelID, fmt.Sprintf("Sent on <t:%d:D> at <t:%d:T>",
					time.Now().Unix(), time.Now().Unix()))
			}

			session.InteractionResponseDelete(interaction.Interaction)

		}
	},
}
