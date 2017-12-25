# slack-filesave

## Install

require Golang environment and setup GOPATH.

```
$ go get github.com/hashibiroko/slack-filesave
```

## Usage

set user slack api token.

```
$ slack-filesave -token=xxxxxx-xxxxxxxxx -type=image
```

### Flags

| name | description | default | require |
| :--- | :---------- | :-----: | :-----: |
| token | your slack api token |  | true |
| type | filter file type | all |  |
| private | include private files | false |  |
| delete | delete downloaded file from slack | false |  |