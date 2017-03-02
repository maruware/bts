package main

import (
	"bufio"
	"fmt"

	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	"github.com/nlopes/slack"
	// "github.com/nlopes/slack"
	"io"
	"os"
	"os/exec"
	"strings"
)

// bts "Changed C Param, It may improve accuracy" -- python some_job.py --input 001.json

var opts struct {
	Args struct {
		Memo    string
		Execute []string
	} `positional-args:"yes" required:"yes"`
}

func printScanner(r io.Reader, colorFn func(string, ...interface{})) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		colorFn(scanner.Text())
	}
}

func execCommand(cmdWithArgs []string, onStart func(), onDone func(), onFailed func()) error {
	var cmd *exec.Cmd
	switch len(cmdWithArgs) {
	case 1:
		cmd = exec.Command(cmdWithArgs[0])
	default:
		cmd = exec.Command(cmdWithArgs[0], cmdWithArgs[1:]...)
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}
	color.Cyan("Executing: " + strings.Join(cmdWithArgs, " "))
	onStart()

	white := color.New(color.FgWhite).PrintfFunc()
	printScanner(stdout, white)

	red := color.New(color.FgRed).PrintfFunc()
	printScanner(stderr, red)

	if err := cmd.Wait(); err != nil {
		onFailed()
		return err
	}
	onDone()
	return nil
}

func findChannel(api *slack.Client, name string) (slack.Channel, error) {
	var channel slack.Channel
	channels, err := api.GetChannels(false)
	if err != nil {
		return channel, err
	}
	for _, c := range channels {
		if c.Name == name {
			return c, nil
		}
	}
	return channel, fmt.Errorf("Not found channel %s", name)
}

func postSlack(api *slack.Client, channelName string, message string, params slack.PostMessageParameters) {
	channel, err := findChannel(api, channelName)
	if err != nil {
		color.Red("%s\n", err)
		return
	}

	api.PostMessage(channel.ID, message, params)
	if err != nil {
		color.Red("%s\n", err)
		return
	}
	color.Cyan("Successfully sent to channel")
}

func postStartToSlack(api *slack.Client, channelName string, memo string, cmdWithArgs []string) {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text:  strings.Join(cmdWithArgs, " "),
		Color: "#00A388",
	}
	params.Attachments = []slack.Attachment{attachment}
	postSlack(api, channelName, memo, params)
}

func postDoneToSlack(api *slack.Client, channelName string, cmdWithArgs []string) {
	params := slack.PostMessageParameters{}
	postSlack(api, channelName, "*Done*", params)
}

func postFailedToSlack(api *slack.Client, channelName string, cmdWithArgs []string) {
	params := slack.PostMessageParameters{}
	postSlack(api, channelName, "*Failed*", params)
}

func main() {
	_, err := flags.Parse(&opts)

	if err != nil {
		color.Red("%s\n", err)
		return
	}

	token := os.Getenv("SLACK_TOKEN")
	channelName := os.Getenv("BTS_SLACK_CHANNEL_NAME")
	api := slack.New(token)

	memo := opts.Args.Memo
	c := opts.Args.Execute

	onStart := func() {
		postStartToSlack(api, channelName, "*memo*: "+memo, c)
	}
	onDone := func() {
		postDoneToSlack(api, channelName, c)
	}
	onFailed := func() {
		postFailedToSlack(api, channelName, c)
	}
	execCommand(c, onStart, onDone, onFailed)
}
