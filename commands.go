package main

import "github.com/bwmarrin/discordgo"

var dmPermission bool = false

var commands = []*discordgo.ApplicationCommand{
	// Hello, world! command.
	{
		Name:         "hello",
		Description:  "Hello, world!",
		DMPermission: &dmPermission,
	},
}
