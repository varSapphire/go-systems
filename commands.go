package main

import "github.com/bwmarrin/discordgo"

var dmPermission bool = false

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "register",
		Description: "This command registers you as a member of your system so that you can proxy messages.",

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Your name.",
				Required:    true,
				MaxLength:   64,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "The unique identifier used to proxy messages.",
				Required:    true,
				MaxLength:   16,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "avatar",
				Description: "The avatar used by you.",
				Required:    true,
			},
		},
	},
	{
		Name:        "proxy",
		Description: "This command will proxy a member of your system using their identifier.",

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "Your identifier.",
				Required:    true,
				MaxLength:   16,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message that you wish to send.",
				Required:    true,
				MaxLength:   2000,
			},
		},
	},
}
