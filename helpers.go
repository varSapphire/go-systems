package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func checkForWebhook(session *discordgo.Session, interaction *discordgo.Interaction) (webhookFlag bool, webhookID string, webhookToken string) {
	webhookFlag = false
	webhookID = ""
	webhookToken = ""

	webhooks, err := session.ChannelWebhooks(interaction.ChannelID)
	if err != nil {
		log.Printf("%vERROR%v - COULD NOT RETRIEVE LIST OF WEBHOOKS FROM CHANNEL:\n\t%v", Red, Reset, err.Error())

		content := fmt.Sprintf("ERROR - COULD NOT RETRIEVE LIST OF WEBHOOKS FROM CHANNEL:\n\t%v", err.Error())
		session.InteractionResponseEdit(interaction, &discordgo.WebhookEdit{Content: &content})

		return webhookFlag, webhookID, webhookToken
	}

	for _, webhook := range webhooks {
		if webhook.Name == "All Systems Go Proxy Webhook" {
			webhookFlag = true
			webhookID = webhook.ID
			webhookToken = webhook.Token

			return webhookFlag, webhookID, webhookToken
		}
	}

	return false, "", ""
}

func checkForMessage(interaction *discordgo.Interaction) (memberMessage string, messageFlag bool) {
	for _, option := range interaction.ApplicationCommandData().Options {
		if option.Name == "message" {
			memberMessage = option.StringValue()
			messageFlag = true

			return memberMessage, messageFlag
		}
	}

	return "", false
}

func checkforLogging(session *discordgo.Session, interaction *discordgo.Interaction) (loggingChannelID string, loggingFlag bool) {
	channels, _ := session.GuildChannels(interaction.GuildID)
	for _, channel := range channels {
		if channel.Name == "asg-logs" {
			loggingChannelID = channel.ID
			loggingFlag = true

			return loggingChannelID, loggingFlag
		}
	}

	return "", false
}
