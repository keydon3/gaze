/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package app

import (
	"flag"
	"strings"

	"github.com/wtetsu/gaze/pkg/config"
	"github.com/wtetsu/gaze/pkg/gazer"
	"github.com/wtetsu/gaze/pkg/logger"
	"github.com/wtetsu/gaze/pkg/uniq"
)

// Start starts a gaze process
func Start(watchFiles []string, userCommand string, file string, appOptions AppOptions) error {
	theGazer, err := gazer.New(watchFiles, appOptions.MaxWatchDirs())
	if err != nil {
		return err
	}
	defer theGazer.Close()

	commandConfigs, err := createCommandConfig(userCommand, file)
	if err != nil {
		return err
	}
	err = theGazer.Run(commandConfigs, appOptions.Timeout(), appOptions.Restart())
	return err
}

func createCommandConfig(userCommand string, file string) (*config.Config, error) {
	if userCommand != "" {
		logger.Debug("userCommand: %s", userCommand)
		commandConfigs := config.New(userCommand)
		return commandConfigs, nil
	}

	if file != "" {
		return config.LoadConfig(file)
	}

	return config.InitConfig()
}

// ParseArgs parses command arguments.
func ParseArgs(osArgs []string, usage func()) *Args {
	flagSet := flag.NewFlagSet(osArgs[0], flag.ExitOnError)

	flagSet.Usage = func() {
		if usage != nil {
			usage()
		}
	}

	help := flagSet.Bool("h", false, "")
	restart := flagSet.Bool("r", false, "")
	userCommand := flagSet.String("c", "", "")
	timeout := flagSet.Int64("t", 1<<50, "")
	yaml := flagSet.Bool("y", false, "")
	quiet := flagSet.Bool("q", false, "")
	verbose := flagSet.Bool("v", false, "")
	file := flagSet.String("f", "", "")
	color := flagSet.Int("color", 1, "")
	debug := flagSet.Bool("debug", false, "")
	version := flagSet.Bool("version", false, "")
	maxWatchDirs := flagSet.Int("w", 100, "")

	files := []string{}
	optionStartIndex := len(osArgs)
	for i, a := range osArgs[1:] {
		if strings.HasPrefix(a, "-") {
			optionStartIndex = i + 1
			break
		}
		files = append(files, a)
	}
	err := flagSet.Parse(osArgs[optionStartIndex:])
	if err != nil {
		return nil
	}

	u := uniq.New()
	u.AddAll(files)
	u.AddAll(flagSet.Args())

	args := Args{
		help:         *help,
		restart:      *restart,
		userCommand:  *userCommand,
		timeout:      *timeout,
		yaml:         *yaml,
		quiet:        *quiet,
		verbose:      *verbose,
		debug:        *debug,
		file:         *file,
		color:        *color,
		version:      *version,
		targets:      u.List(),
		maxWatchDirs: *maxWatchDirs,
	}

	return &args
}
