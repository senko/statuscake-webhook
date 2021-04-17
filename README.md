# StatusCake Webhook handler

This is a simple HTTP service to handle host
status change (up/down) webhooks from StatusCake.

When a DOWN status is triggered, the service executes
a command. That's it.

## Installation

A recent version of Go is required. To build:

```
cd src/
go build webhook.go
```

## Configuration

Configuration is done via environment variables. Port
to listen on should be specified in `PORT` env variable.

With each webhook, StatusCake sends a token to verify
the sender (so someone else can't fake it). More info
is available at [StatusCake Webhook handler documentation](https://www.statuscake.com/kb/knowledge-base/how-to-use-the-web-hook-url/).

Token should be specified in the `TOKEN` environment
variable. HTTP requests without a matching token will
be ignored.

## Running

To run, specify the command (with optional parameters) to
execute once the host DOWN status is received:

```
webhook /usr/bin/some-command arg1 arg2
```

## License

This software is released to Public Domain.
Copyright is for chums.

