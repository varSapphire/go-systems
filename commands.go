package main

import "github.com/bwmarrin/discordgo"

var dmPermission bool = false
var manageServerPermission int64 = discordgo.PermissionManageServer
var tokenMinLength int = 64
var tokenMaxLength int = 64

var commands = []*discordgo.ApplicationCommand{
	{
		Name:                     "setup",
		Description:              "This command creates a channel for logging.",
		DefaultMemberPermissions: &manageServerPermission,
		DMPermission:             &dmPermission,
	},
	{
		Name:         "register",
		Description:  "This command registers your PluralKit token into the database.",
		DMPermission: &dmPermission,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "token",
				Description: "Your PluralKit token.",
				Required:    true,
				MinLength:   &tokenMinLength,
				MaxLength:   tokenMaxLength,
			},
		},
	},
	{
		Name:         "auto_proxy_message",
		Description:  "Proxies a message using the current fronter of your system.",
		DMPermission: &dmPermission,
	},
}
