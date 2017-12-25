# slack-filesave

## install

require Golang environment and setup GOPATH.

```
$ go get github.com/hashibiroko/slack-filesave
```

## usage

set user slack api token.

```
$ slack-filesave -token=xxxxxx-xxxxxxxxx -type=image
```

### flags

| name | description | default | require |
| :--- | :---------- | :-----: | :-----: |
| token | your slack api token |  | true |
| type | filter file type |  |  |
| private | include private files | false |  |