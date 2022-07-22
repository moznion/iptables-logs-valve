# iptables-logs-valve [![.github/workflows/check.yml](https://github.com/moznion/iptables-logs-valve/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/iptables-logs-valve/actions/workflows/check.yml)

A CLI tool to convert the iptables logs that come from STDIN into JSON/JSONL.

## Usage

```
iptables-logs-valve: A CLI tool to convert the iptables logs that come from STDIN into JSON/JSONL.
Example Usage:
  tail -F /var/log/iptables.log | iptables-logs-valve

Options:
  -array
        If this option is specified, this tool output the result by wrapping it in a top-level JSON array. Else, it output the results as JSONL. (default: not array; i.e. JSONL)
  -buffer-millisec uint
        The duration for buffering. This value's time unit is milliseconds. If this value is 0, it won't do buffering. (default: 0; i.e. don't buffer it)
  -version
        Show version information.
```

## Example

### Simple transformation

```
$ echo '2022-07-12T09:01:27.345918+00:00 ubuntu-jammy kernel: [14879.600492] OUT-LOG: IN= OUT=enp0s3 SRC=10.0.2.15 DST=8.8.8.8 LEN=84 TOS=0x00 PREC=0x00 TTL=64 ID=6495 DF PROTO=ICMP TYPE=8 CODE=0 ID=1 SEQ=3' | iptables-logs-valve
{"timestamp":"2022-07-12T09:01:27.345918+00:00","hostname":"ubuntu-jammy","kernelTimestamp":14879.600492,"prefix":"OUT-LOG:","inputInterface":"","outputInterface":"enp0s3","macAddress":"","source":"10.0.2.15","destination":"8.8.8.8","length":84,"tos":0,"precedence":0,"ttl":64,"id":6495,"congestionExperienced":false,"doNotFragment":true,"moreFragmentsFollowing":false,"frag":0,"ipOptions":"","protocol":"ICMP","type":8,"code":0,"sourcePort":0,"destinationPort":0,"sequence":0,"ackSequence":0,"windowSize":0,"res":0,"urgent":false,"ack":false,"push":false,"reset":false,"syn":false,"fin":false,"urgp":0,"tcpOption":""}
```

### With buffering and JSON array wrapping

https://user-images.githubusercontent.com/1422834/180379103-01e83293-def4-4128-9b54-f2f8bd07e52e.mov

## Author

moznion (<moznion@mail.moznion.net>)

