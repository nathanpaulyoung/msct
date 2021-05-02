package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/urfave/cli/v2"
)

func main() error {
	// make sure we're on linux
	if runtime.GOOS != "Linux" {
		glog.Fatal("msct currently only supports linux")
	}

	// initialize config struct to default values
	cfg := newConfig()
	// check default global config location for config file
	if _, err := os.Stat("/etc/msct.conf"); err == nil {
		// found, so load it into our config struct
		err := cfg.load("/etc/msct.conf")
		if err != nil {
			glog.Fatal("error loading /etc/msct.conf; check the file for syntax errors, or delete the file and a default will be generated")
		}
	} else {
		// not found, save default values to file
		err := cfg.save("/etc/msct.conf")
		if err != nil {
			glog.Fatal("could not create config file at default location /etc/msct.conf; are you running as root?")
		} else {
			glog.Warning("generated default config file at /etc/msct.conf")
		}
	}

	sc := newServerController(cfg)

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start a minecraft server",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startCommand() cli.Command {
	command := cli.Command{
		Name:  "s",
		Usage: "start a server",
		Action: func(c *cli.Context) {
			servername := c.Args()[0]
			startServer(servername)
		},
	}
	return command
}

func haltCommand() cli.Command {
	command := cli.Command{
		Name:  "h",
		Usage: "halt a server",
		Action: func(c *cli.Context) {
			servername := c.Args()[0]
			tmuxSendKeys(servername, 0, 0, "stop")
		},
	}
	return command
}

func resumeCommand() cli.Command {
	command := cli.Command{
		Name:  "r",
		Usage: "resume a server's tmux session",
		Action: func(c *cli.Context) {
			servername := c.Args()[0]
			args := []string{"a", "-t", buildTmuxName(servername)}
			cmd := exec.Command("tmux", args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if serverIsRunning(servername) {
				if err := cmd.Run(); err != nil {
					glog.Fatalf("Error in function resumeCommand(); please inform the developer.\n%s", err)
				}
			} else {
				glog.Errorf("Server '%s' is not running.", servername)
			}
		},
	}
	return command
}

func keepAliveCommand() cli.Command {
	command := cli.Command{
		Name:  "keepalive",
		Usage: "restart a server's tmux session if server detected as stopped",
		Action: func(c *cli.Context) {
			servername := c.Args()[0]
			keepAliveFreq, _ := cfg.Int("keepAliveFreq")
			for {
				if !serverIsRunning(servername) {
					startServer(servername)
					time.Sleep(time.Second * time.Duration(keepAliveFreq))
				}
			}
		},
	}
	return command
}

func commandCommand() cli.Command {
	command := cli.Command{
		Name:  "cmd",
		Usage: "send a command to the tmux session",
		Action: func(c *cli.Context) {
			servername := c.Args()[0]
			var args []string
			for i := 1; i < len(c.Args()); i++ {
				args = append(args, c.Args()[i])
			}
			command := strings.Join(args, " ")
			tmuxSendKeys(servername, 0, 0, command)
		},
	}
	return command
}
