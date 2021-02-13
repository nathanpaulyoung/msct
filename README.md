# **MSCT** - Minecraft Server Control Tool [![Build Status](https://travis-ci.org/nathanpaulyoung/msct.svg?branch=master)](https://travis-ci.org/nathanpaulyoung/msct)  
Pronounced like "musket". Written in [golang](http://golang.org). Made by a server admin for server admins.

## Installation
1. Download the [latest](https://github.com/nathanpaulyoung/msct/releases/latest) binary from the releases page.
2. Install to `/usr/bin` on the server.
3. Run any command once to generate the default config file at `/etc/msct.conf` and edit it to your satisfaction.

## Usage

### Starting a server
    msct start <servername>
    msct s <servername>

### Resuming a server
    msct resume <servername>
    msct r <servername>

### Halting a server
    msct halt <servername>
    msct h <servername>
    msct stop <servername>

### Keeping a server alive
    msct keepalive <servername>
    msct ka <servername>
Note: You probably want to run `msct keepalive` in the background. I recommend `nohup msct keepalive` for its simplicity, but there are many options for backgrounding a program available to choose from.

Alternatively, you can set up a cronjob to periodically run `msct start`, as this command will not run if the target server is already running.

### Sending arbitrary commands to a server
    msct command <servername> <command>
    msct cmd <servername> <command>
    msct c <servername> <command>

### Detaching from the server console
Press the key combination `CTRL + b`, then type `d` by itself.

## Planned Features
* Daemonized mode, with web portal and user authentication for starting and stopping servers.
* Configurable automailer via [Mailgun](http://mailgun.com), with attached crash reports and server logs.
* Self updater.
* Windows and OSX compatibility.
* Launch server as non-root user.
* World backups.
* Automatic leveraging of tmux's panes for system resource visibility while viewing the console.

## Project Information
I started writing MSCT when I started using a solid state drive on my own personal friends and family server. I made a symlink for my world directory, and [MSM](http://msmhq.com) overwrote it with its own symlink. I was not thrilled that so few things in MSM are configurable to my satisfaction, and that there was so much extra fluff. So I set out to roll my own solution.

This project is currently run by just one guy. If you want to get involved, send me a pull request, open an issue, or request a feature. I'm very open to making this project the best it can be. You can also always [PM me on reddit](https://www.reddit.com/message/compose?to=nathanpaulyoung), under the same username as here.

## Special Thanks:  
[@Merovius](http://github.com/Merovius) for helping with the issues with attaching.
