//go:build debug

package main

import "github.com/bennicholls/tyumi/log"

func init() {
	debug = true
	log.EnableConsoleOutput()
	log.SetMinimumLogLevel(log.LVL_DEBUG)
	log.Debug("Beginning debug mode!")
}
