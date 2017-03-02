# BTS

Post batch execution to Slack with a memo.

## Usage

bts memo -- cmd...

```
$ bts "Changed C Param, It may improve accuracy" -- python some_job.py --input 001.json
```

## Install

```
$ go get github.com/maruware/bts
```

Set Environment

```
export SLACK_TOKEN=xoxb-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxx
export BTS_SLACK_CHANNEL_NAME=times_maruware
```
