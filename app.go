package iptables_logs_valve

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/moznion/go-iptables-logs-parser"
)

func Run(ctx context.Context, inputReader io.Reader, outputWriter io.Writer, bufferingDuration time.Duration, jsonArrayMode bool) {
	logEmitter := func(l *iptables.Log) {
		serializedLog, err := json.Marshal(l)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			return
		}

		if jsonArrayMode {
			_, _ = fmt.Fprintf(outputWriter, "[%s]\n", serializedLog)
			return
		}
		_, _ = fmt.Fprintf(outputWriter, "%s\n", serializedLog)
	}

	if bufferingDuration > 0 {
		var mu sync.Mutex
		buff := make([]string, 0, 100000)

		logEmitter = func(l *iptables.Log) {
			serializedLog, err := json.Marshal(l)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
				return
			}

			mu.Lock()
			buff = append(buff, string(serializedLog))
			mu.Unlock()
		}

		go func() {
			ticker := time.NewTicker(bufferingDuration)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					mu.Lock()

					if len(buff) > 0 {
						if jsonArrayMode {
							_, _ = fmt.Fprintf(outputWriter, "[%s]\n", strings.Join(buff, ","))
						} else {
							_, _ = fmt.Fprintf(outputWriter, "%s\n", strings.Join(buff, "\n"))
						}
					}

					buff = buff[:0]
					mu.Unlock()
				}
			}
		}()
	}

	lineCh := make(chan string)
	scannerClosedCh := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(inputReader)
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		}

		scannerClosedCh <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-scannerClosedCh:
			return
		case line := <-lineCh:
			parsedLog, err := iptables.Parse(line)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
				continue
			}
			logEmitter(parsedLog)
		}
	}
}
