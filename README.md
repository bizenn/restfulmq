# restfulmq
RESTful Interface Message Queue Server

## Install

```
% go install github.com/bizenn/restfulmq
```

## Usage

### Run

```
% restfulmq [<config.json>]
time="2017-03-10T16:36:09+09:00" level=info msg=Start 
```

If restfulmq is invoked without config.json, It make only one unbuffered queue named "/"
and listen port 8888 on all network interfaces.

### Enqueue

```
% curl -X POST --data "foo" http://localhost:8888/
```

### Dequeue

```
% curl http://localhost:8888/
foo
```

## Configuration

```
{
    "host": "127.0.0.1",
    "port": 8086,
    "logpath": "/tmp/restrulmq.log",
    "queues": [
        {"path": "/jobstart", "capacity": 10},
        {"path": "/jobfinish", "capacity": 0}
    ]
}
```

## Build

```
% make build
```

## Cross Build

```
% make xbuild
```

Default building Linux/x86_64.

```
% make xbuild OS=freebsd
```

Building FreeBSD/x86_64.
