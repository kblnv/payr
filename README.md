# Simple Notification System

A lightweight notification system built with a plugin-based architecture. Example plugin implementations can be found in the `plugins/` directory.

The system accepts HTTP requests on the `/event` endpoint. When an event is received, the system determines which plugin should handle it, passes necessary data to that plugin, and then delivers the generated notification using the configured transport layer.

A plugin is a simple function that receives event data and plugin-specific configuration, and returns a notification message as a string. Plugins are loaded separately using Go's `plugin` package. Each plugin defines its own contract, including the structure of the event data it expects and the configuration parameters it requires.

A transport is responsible for delivering the generated message to its destination (for example, Telegram).

## Tech Stack

At the moment, the system is built without any external dependencies. The main design goal is to keep it simple.

## Planned Improvements

- **Scheduler**  
  A simple scheduler that allows configuring event execution time and recurrence.

- **More Transports**  
  Add additional transport implementations besides Telegram.

## How to Use

### Build and Run
Build core and plugins using `make`, then initialize config and run the server.

```sh
$ make core
$ make plugins
$ cd build
$ ./payr init
$ ./payr run
```

### Triggering an Event

Send an HTTP request to trigger an event manually:

```sh
curl -X POST 127.0.0.1:8080/event \
  -H "Content-Type: application/json" \
  -d '{"event": "hello", "transports": [{"transport": "telegram", "to": "<chat_id>"}], "meta": {"Name": "Guest"}}'
```
