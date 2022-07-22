package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	valve "github.com/moznion/iptables-logs-valve"
	"github.com/moznion/iptables-logs-valve/internal"
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `%s: A CLI tool to convert the iptables logs that come from STDIN into JSON/JSONL.
Example Usage:
  tail -F /var/log/iptables.log | %s

Options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var bufferMillisecDuration uint64
	var arrayMode bool
	var showVersion bool
	flag.Uint64Var(&bufferMillisecDuration, "buffer-millisec", 0, "The duration for buffering. This value's time unit is milliseconds. If this value is 0, it won't do buffering. (default: 0; i.e. don't buffer it)")
	flag.BoolVar(&arrayMode, "array", false, "If this option is specified, this tool output the result by wrapping it in a top-level JSON array. Else, it output the results as JSONL. (default: not array; i.e. JSONL)")
	flag.BoolVar(&showVersion, "version", false, "Show version information.")
	flag.Parse()

	if showVersion {
		fmt.Println(internal.GetVersionJSONString())
		os.Exit(0)
	}

	ctx := context.Background()
	valve.Run(ctx, os.Stdin, os.Stdout, time.Duration(bufferMillisecDuration)*time.Millisecond, arrayMode)
}
