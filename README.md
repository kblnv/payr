# Simple Notification / Reminder System

A lightweight notification and reminder system built with a plugin-based architecture. Example plugin implementations can be found in the `plugins/` directory.

The system accepts HTTP requests on the `/event` endpoint. When an event is received, the system determines which plugin should handle it, passes necessary data to that plugin, and then delivers the generated notification using the configured transport layer.

A plugin is a simple function that receives event data and plugin-specific configuration, and returns a notification message as a string. Plugins are loaded separately using Go's `plugin` package, allowing new functionality to be added without modifying the core application. Each plugin defines its own contract, including the structure of the event data it expects and the configuration parameters it requires.

A transport is responsible for delivering the generated message to its destination (for example, Telegram).

## Tech Stack

At the moment, the system is built without any external dependencies. All functionality is implemented using only the Go standard library. The main design goal is to remain small, efficient in CPU/RAM usage, and easily extensible for future improvements. 

## Planned Improvements

- **Scheduler**  
  A simple scheduler that allows configuring event execution time and recurrence.

- **More Transports**  
  Add additional transport implementations besides Telegram.

## How to Use

### Build and Run

Build core and plugins using `make` and start the server with a configuration file.

```sh
$ make core
$ make plugins
$ ./build/payr --config <config_path>
```

## Configuration Example

```json
{
  "server": {
    "address": "127.0.0.1:8080"
  },

  "plugins_dir": "./build/plugins",

  "handlers": {
    "hello": {
      "plugin": "template",
      "settings": {
        "template": "Hello, {{ .Name }}!"
      }
    }
  },

  "transports": {
    "telegram": {
      "bot_token": "<bot_token>",
      "chat_id": "<chat_id>"
    }
  },

  "events": [
    {
      "name": "hello",
      "transports": ["telegram"],
      "handler": "hello"
    }
  ]
}
```

## Triggering an Event

Send an HTTP request to trigger an event manually:

```sh
curl -X POST 127.0.0.1:8080/event \
  -H "Content-Type: application/json" \
  -d '{"event":"hello", "meta": {"Name": "Guest"}}'
```