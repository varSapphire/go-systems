package main

import "github.com/bwmarrin/discordgo"

var dmPermission bool = false
var manageServerPermission int64 = discordgo.PermissionManageServer

var tokenMinLength int = 64
var tokenMaxLength int = 64

var switchMinLength = 0
var switchMaxLength = 100

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
		Name:         "switch",
		Description:  "This command registers a switch with PK so that you can use the /auto_proxy_message command.",
		DMPermission: &dmPermission,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "member_name",
				Description: "The name of the member switching in.",
				Required:    true,
				MinLength:   &switchMinLength,
				MaxLength:   switchMaxLength,
			},
		},
	},
	{
		Name:         "list_members",
		Description:  "This command lists your system members.",
		DMPermission: &dmPermission,
	},
	{
		Name:         "auto_proxy_message",
		Description:  "Proxies a message using the first fronter of your system.",
		DMPermission: &dmPermission,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message that you wish to send.",
				Required:    false,
				MaxLength:   2000,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_1",
				Description: "The first optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_2",
				Description: "The second optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_3",
				Description: "The third optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_4",
				Description: "The fourth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_5",
				Description: "The fifth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_6",
				Description: "The sixth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_7",
				Description: "The seventh optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_8",
				Description: "The eighth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_9",
				Description: "The fourth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_10",
				Description: "The tenth optional attachment.",
				Required:    false,
			},
		},
	},
	{
		Name:         "manual_proxy_message",
		Description:  "Proxies a message using the specified member.",
		DMPermission: &dmPermission,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "proxy",
				Description: "The proxy of your system member.",
				Required:    true,
				MaxLength:   100,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message that you wish to send.",
				Required:    false,
				MaxLength:   2000,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_1",
				Description: "The first optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_2",
				Description: "The second optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_3",
				Description: "The third optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_4",
				Description: "The fourth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_5",
				Description: "The fifth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_6",
				Description: "The sixth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_7",
				Description: "The seventh optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_8",
				Description: "The eighth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_9",
				Description: "The fourth optional attachment.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "attachment_10",
				Description: "The tenth optional attachment.",
				Required:    false,
			},
		},
	},
}
