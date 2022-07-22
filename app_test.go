package iptables_logs_valve

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/moznion/go-iptables-logs-parser"

	"github.com/stretchr/testify/assert"
)

func TestRun_NoBuffers(t *testing.T) {
	input, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(input.Name())
	}()

	output, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(output.Name())
	}()

	ctx, canceler := context.WithCancel(context.Background())
	defer canceler()

	timestamp1 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp1.UnixMicro(),
	))
	assert.NoError(t, err)

	timestamp2 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp2.UnixMicro(),
	))
	assert.NoError(t, err)

	_, err = input.Seek(0, 0)
	assert.NoError(t, err)

	go Run(ctx, input, output, 0, false)

	time.Sleep(100 * time.Millisecond)

	_, err = output.Seek(0, 0)
	assert.NoError(t, err)

	convertedLogs := make([]*iptables.Log, 0)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var l iptables.Log
		err := json.Unmarshal(scanner.Bytes(), &l)
		assert.NoError(t, err)

		convertedLogs = append(convertedLogs, &l)
	}
	err = scanner.Err()
	assert.NoError(t, err)

	assert.Len(t, convertedLogs, 2)
	assert.Equal(t, fmt.Sprintf("%d", timestamp1.UnixMicro()), convertedLogs[0].Timestamp)
	assert.Equal(t, fmt.Sprintf("%d", timestamp2.UnixMicro()), convertedLogs[1].Timestamp)
}

func TestRun_NoBuffersWithJSONArray(t *testing.T) {
	input, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(input.Name())
	}()

	output, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(output.Name())
	}()

	ctx, canceler := context.WithCancel(context.Background())
	defer canceler()

	timestamp1 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp1.UnixMicro(),
	))
	assert.NoError(t, err)

	timestamp2 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp2.UnixMicro(),
	))
	assert.NoError(t, err)

	_, err = input.Seek(0, 0)
	assert.NoError(t, err)

	go Run(ctx, input, output, 0, true)

	time.Sleep(100 * time.Millisecond)

	_, err = output.Seek(0, 0)
	assert.NoError(t, err)

	convertedLogs := make([]*iptables.Log, 0)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var logs []iptables.Log
		err := json.Unmarshal(scanner.Bytes(), &logs)
		assert.NoError(t, err)

		convertedLogs = append(convertedLogs, &logs[0])
	}
	err = scanner.Err()
	assert.NoError(t, err)

	assert.Len(t, convertedLogs, 2)
	assert.Equal(t, fmt.Sprintf("%d", timestamp1.UnixMicro()), convertedLogs[0].Timestamp)
	assert.Equal(t, fmt.Sprintf("%d", timestamp2.UnixMicro()), convertedLogs[1].Timestamp)
}

func TestRun_WithBuffers(t *testing.T) {
	input, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(input.Name())
	}()

	output, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(output.Name())
	}()

	ctx, canceler := context.WithCancel(context.Background())
	defer canceler()

	timestamp1 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp1.UnixMicro(),
	))
	assert.NoError(t, err)

	timestamp2 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp2.UnixMicro(),
	))
	assert.NoError(t, err)

	_, err = input.Seek(0, 0)
	assert.NoError(t, err)

	go Run(ctx, input, output, 500, false)

	time.Sleep(1000 * time.Millisecond)

	_, err = output.Seek(0, 0)
	assert.NoError(t, err)

	convertedLogs := make([]*iptables.Log, 0)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var l iptables.Log
		err := json.Unmarshal(scanner.Bytes(), &l)
		assert.NoError(t, err)

		convertedLogs = append(convertedLogs, &l)
	}
	err = scanner.Err()
	assert.NoError(t, err)

	assert.Len(t, convertedLogs, 2)
	assert.Equal(t, fmt.Sprintf("%d", timestamp1.UnixMicro()), convertedLogs[0].Timestamp)
	assert.Equal(t, fmt.Sprintf("%d", timestamp2.UnixMicro()), convertedLogs[1].Timestamp)
}

func TestRun_WithBuffersAndJSONArray(t *testing.T) {
	input, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(input.Name())
	}()

	output, err := os.CreateTemp("", "")
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(output.Name())
	}()

	ctx, canceler := context.WithCancel(context.Background())
	defer canceler()

	timestamp1 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp1.UnixMicro(),
	))
	assert.NoError(t, err)

	timestamp2 := time.Now()
	_, err = input.WriteString(fmt.Sprintf(
		"%d ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3\n",
		timestamp2.UnixMicro(),
	))
	assert.NoError(t, err)

	_, err = input.Seek(0, 0)
	assert.NoError(t, err)

	go Run(ctx, input, output, 500, true)

	time.Sleep(1000 * time.Millisecond)

	_, err = output.Seek(0, 0)
	assert.NoError(t, err)

	convertedLogs := make([]*iptables.Log, 0)
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		var logs []iptables.Log
		err := json.Unmarshal(scanner.Bytes(), &logs)
		assert.NoError(t, err)

		for _, l := range logs {
			convertedLogs = append(convertedLogs, &l)
		}
	}
	err = scanner.Err()
	assert.NoError(t, err)

	assert.Len(t, convertedLogs, 2)
	assert.Equal(t, fmt.Sprintf("%d", timestamp1.UnixMicro()), convertedLogs[0].Timestamp)
	assert.Equal(t, fmt.Sprintf("%d", timestamp2.UnixMicro()), convertedLogs[1].Timestamp)
}
