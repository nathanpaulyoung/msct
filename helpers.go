package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func startServer(servername string) {
	args := buildInvocation(servername)
	cmd := exec.Command("tmux", args...)
	cmd.Dir = buildServerDir(servername)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if serverExists(servername) && !serverIsRunning(servername) {
		if err := cmd.Run(); err != nil {
			glog.Fatalf("Error in function startServer(); please inform the developer.\n%s", err)
		}
	} else {
		glog.Errorf("Server '%s' does not exist or is already running.", servername)
	}
}

func tmuxSendKeys(servername string, window int, pane int, command string) {
	cmd := exec.Command("tmux", "send-keys", "-t", buildTmuxName(servername)+":"+strconv.Itoa(window)+"."+strconv.Itoa(pane), command, "Enter")
	if serverIsRunning(servername) {
		if err := cmd.Run(); err != nil {
			glog.Fatalf("Error in function tmuxSendKeys(); please inform the developer.\n%s", err)
		}
	} else {
		glog.Errorf("Server '%s' is not running.", servername)
	}
}

func serverExists(servername string) bool {
	if _, err := os.Stat(buildServerDir(servername) + getJarFile()); err != nil {
		return false
	}
	return true
}

func serverIsRunning(servername string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", buildTmuxName(servername))
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func buildTmuxName(servername string) string {
	//Load from config and set base tmux prefix, if not set in config, default to "msct-"
	tmuxbasename, err := cfg.String("tmuxbasename")
	if err != nil {
		tmuxbasename = "msct-"
	}

	return tmuxbasename + servername
}

func buildInvocation(servername string) []string {
	//Load from config and set whether to start tmux attached or not, if not set in config, default to attached
	startTmuxAttached, err := cfg.Bool("startTmuxAttached")
	if err != nil {
		startTmuxAttached = true
	}
	var tmuxParams []string
	if startTmuxAttached == true {
		tmuxParams = append(tmuxParams, "new", "-s", buildTmuxName(servername))
	} else {
		tmuxParams = append(tmuxParams, "new", "-d", "-s", buildTmuxName(servername))
	}

	//Load from config and set java parameters, if not set in config, set reasonable defaults
	ramMin, err := cfg.String("ramMin")
	if err != nil {
		ramMin = "2048M"
	}

	ramMax, err := cfg.String("ramMax")
	if err != nil {
		ramMax = "4096M"
	}

	//Load from config and set java parameters, if not set in config, set reasonable defaults
	javaParams, err := cfg.String("javaParams")
	if err != nil {
		javaParams = "-XX:+UseConcMarkSweepGC -XX:+UseParNewGC -XX:+CMSParallelRemarkEnabled -XX:ParallelGCThreads=3 -XX:+DisableExplicitGC -XX:MaxGCPauseMillis=500 -XX:SurvivorRatio=16 -XX:TargetSurvivorRatio=90"
	}
	javaParamsArray := strings.Fields(javaParams)

	//Create full server path of the form /opt/minecraft/<servername>/server.jar
	fullpath := buildServerDir(servername) + getJarFile()

	var args []string
	args = append(args, tmuxParams...)
	args = append(args, fmt.Sprintf("java -server -Xms%sM -Xmx%sM %s -jar %s", ramMin, ramMax, strings.Join(javaParamsArray, " "), fullpath))

	if debugIsEnabled() {
		println(strings.Join(args, " "))
	}

	return args
}

func buildServerDir(servername string) string {
	//Load from config and set msct root directory, if not set in config, default to /opt/minecraft/
	rootdir, err := cfg.String("paths.root")
	if err != nil {
		rootdir = "/opt/minecraft/"
	}

	return rootdir + servername + "/"
}

func getJarFile() string {
	//Load from config and set server jar filename, if not set in config, default to server.jar
	jarFile, err := cfg.String("paths.jarFile")
	if err != nil {
		jarFile = "server.jar"
	}

	return jarFile
}

func debugIsEnabled() bool {
	debug, _ := cfg.Bool("debug")
	return debug
}
