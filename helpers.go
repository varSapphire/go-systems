package main

import (
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

func WebhookExecuteSimple(
	session *discordgo.Session,
	webhookID string,
	webhookToken string,
	memberMessage string,
	memberName string,
	memberAvatarURL string) {
	_, err = session.WebhookExecute(webhookID, webhookToken, true, &discordgo.WebhookParams{
		Content:   memberMessage,
		Username:  memberName,
		AvatarURL: memberAvatarURL,
	})
	if err != nil {
		log.Printf("%vERROR%v - COULD NOT EXECUTE WEBHOOK:\n\t%v", Red, Reset, err.Error())
	}
}
