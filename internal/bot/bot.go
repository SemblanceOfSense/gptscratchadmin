package bot

import (
	"fmt"
	"gptscratchadmin/internal/flagcomment"
	"gptscratchadmin/internal/getcomments"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
    appID = "1222735923796512779"
    guildID = "785618905665896478"
)

func Run(BotToken string, openaikey string) {
    discord, err := discordgo.New(("Bot " + BotToken))
    if err != nil { log.Fatal(err) }

    _, err = discord.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand {
        {
            Name: "scan-comments",
            Description: "scan comments from project",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionString,
                    Name: "project-owner-username",
                    Description: "the username of the project owner",
                    Required: true,
                },
                {
                    Type: discordgo.ApplicationCommandOptionInteger,
                    Name: "project-id",
                    Description: "the id of the project",
                    Required: true,
                },
                {
                    Type: discordgo.ApplicationCommandOptionInteger,
                    Name: "hours-to-scan",
                    Description: "the amount of hours since present time to scan",
                    Required: true,
                },
            },
        },
    })
    if err != nil { fmt.Println("1"); log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i *discordgo.InteractionCreate,
    ) {
        data := i.ApplicationCommandData()
        switch data.Name {
        case "scan-comments":
            if i.Interaction.Member.User.ID == s.State.User.ID { return }

            err = s.InteractionRespond(
                i.Interaction,
                &discordgo.InteractionResponse {
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData {
                        Content: "Wait please!",
                    },
                },
            )
            if err != nil { log.Fatal(err) }

            var projectownerusername string
            var projectid, hours int
            _ = projectownerusername
            _ = projectid
            _ = hours
            for _, v := range i.Interaction.ApplicationCommandData().Options {
                switch v.Name {
                case "project-owner-username":
                    projectownerusername = v.StringValue()
                case "project-id":
                    projectid = int(v.IntValue())
                case "hours-to-scan":
                    hours = int(v.IntValue())
                }
            }

            comments, err := getcomments.GetComments(projectownerusername, projectid, hours)
            if err != nil { log.Fatal(err) }

            flaggedComments := []string{}
            for _, val := range comments {
                flag, err := flagcomment.FlagComment(val.Content, openaikey)
                if err != nil { log.Fatal(err) }
                if (flag) { flaggedComments = append(flaggedComments, val.Content + " by " + val.Author.Username) }
                time.Sleep(time.Millisecond * 100)
            }

            discord.ChannelMessageSend(i.ChannelID, strings.Join(flaggedComments, "\n"))
        }
    })

    err = discord.Open()
    if err != nil { fmt.Println("3"); log.Fatal(err) }

    stop := make (chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    log.Println("Press Ctrl+C to Exit")
    <-stop

    err = discord.Close()
    if err != nil { fmt.Println("4"); log.Fatal(err) }
}
