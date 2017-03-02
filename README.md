# BTS

Post batch execution to Slack with a memo.

## Usage

Execute

```
$ bts "Changed C Param, It may improve accuracy" -- python some_job.py --input 001.json
```

Then

<img width="356" alt="ss_ 2017-03-02 at 22 26 33" src="https://cloud.githubusercontent.com/assets/1129887/23509054/cd8545f6-ff97-11e6-9824-b26556db6ade.png">

## Install

```
$ go get github.com/maruware/bts
```

Set Environment

```
export SLACK_TOKEN=xoxb-xxxxxxxxxxxx-xxxxxxxxxxxxxxxxxx
export BTS_SLACK_CHANNEL_NAME=times_maruware
```
