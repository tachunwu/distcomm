package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/pflag"
)

var dataDir = pflag.String("data-dir", "", "Directory to store persistent data")
var addr = pflag.String("addr", "", "Address used by pebble server")
var streamAddr = pflag.String("stream-addr", "", "Address point to messaging stream system")

func init() {
	// Seed rng.
	rand.Seed(time.Now().UTC().UnixNano())

	// Setup logging.
	log.SetOutput(os.Stderr)
	log.SetPrefix("pebble | ")
	log.SetFlags(log.LstdFlags)

	// Parse flags
	pflag.Parse()
}
