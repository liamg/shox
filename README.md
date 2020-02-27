# Shox: Terminal Status Bar

[![Travis Build Status](https://travis-ci.org/liamg/shox.svg?branch=master)](https://travis-ci.org/liamg/shox)
[![GoReportCard](https://goreportcard.com/badge/github.com/liamg/shox)](https://goreportcard.com/report/github.com/liamg/shox)
[![Github Release](https://img.shields.io/github/release/liamg/shox.svg)](https://github.com/liamg/shox/releases)

A customisable terminal status bar with universal shell/terminal compatibility. Currently works on Mac/Linux.

![](./screenshot.png)

## Installation

**NOTE** This is still very experimental. I'm using it locally without any problems right now, but there's still a lot of testing and tweaking to do. Feel free to try it out, but get ready for some potential bugginess!

```bash
curl -s "https://raw.githubusercontent.com/liamg/shox/master/scripts/install.sh" | sudo bash
```

If you don't like to pipe to sudo - as well you shouldn't - you can remove the `sudo` above, but you'll have to add the shox dir to your `PATH` env var manually, as instructed by the installer.

## Configuration

The shox config file should be created at `$XDG_CONFIG_HOME/shox/config.yaml`, which is usually `~/.config/shox/config.yaml`. You can alternatively create it at `~/.shox.yaml`

The config file looks like the following:

```yaml
shell: /bin/bash
bar:
    format: "{time}||CPU: {cpu} MEM: {memory}"
    colours: 
      bg: red
      fg: white
    padding: 0
```

Shox will use your `SHELL` environment variable to determine the shell to run if a shell is not specified in the config file, but if your `SHELL` is set to shox, it'll default to `/bin/bash` to prevent a horrible recursive mess.

### Bar Configuration

Bar configuration is done using a simple string format. Helpers are encased in braces e.g. `{time}`, alignment is done using pipes (see below), and any other text will be written to the bar.

#### Alignment

You can use pipes to align content within the status bar. All content before the first pipe will be aligned to the left, all content between the first and second will be centre aligned, and all content after the second pipe will be right aligned.

For example, to display a bar that centre aligns the time, you could use `|{time}|` 

#### Colours

The following colours are available: `black`, `white`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`,  `darkgrey`, `lightgrey`, `lightred`, `lightgreen`, `lightyellow`, `lightblue`, `lightmagenta`, `lightcyan`.

#### Helpers

Helpers create dynamic output in your status bar. You can use one by adding it to your bar format config. The following is a list of available helpers.

| Helper  | Description                                       | Example Config   | Example Output |
|---------|---------------------------------------------------|------------------|----------------|
| time    | Show current time                                 | {time}           | 11:58:17       |
| cpu     | Show current CPU usage                            | {cpu}            | 20%            |
| memory  | Show current memory usage %                       | {memory}         | 20%            |
| battery | Show current battery charge %                     | {battery}        | 20%            |
| bash    | Run a custom bash command                         | {bash:echo hi}   | hi             |
| weather | Show current weather (provided by http://wttr.in) | {weather:1}      | 🌧 +6°C         |

Ideally this list would be much longer - please feel free to PR more helpers! You can see simple examples [here](https://github.com/liamg/shox/tree/master/pkg/helpers).

##### Weather

The configuration section of the weather helper holds the display format.
For all available display formats please visit
[chubin/wttr.in#one-line-output](https://github.com/chubin/wttr.in#one-line-output)
The default value is `1` which only shows the weather

> **_NOTE:_** You don't need to URL-encode the weather format, i.e. use `%l: %c %t` instead of `%l:+%c+%t`

## Uninstallation

### If installed with `sudo`
Remove the binary from `/usr/local/bin`
```bash
rm /usr/local/bin/shox
```

### If installed without `sudo`
Remove the binary from the shox installation dir `$HOME/bin`
```bash
rm $HOME/bin/shox
```

> **_NOTE:_** Don't forget to remove any configuration files you've created should you decide you don't need them

## Why?

I frequently needed a way to have a quick overview of several things without cramming them into my PS1, and to update those things dynamicly.

## How does it work?

Shox sits between the terminal and your shell and proxies all data sent between them. It identifies ANSI commands which contain coordinates and dimensions and adjusts them accordingly, so that the status bar can be drawn efficiently without interfering with the shell and it's child programs.
