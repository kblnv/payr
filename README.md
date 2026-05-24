# Simple Notification / Reminder System

A lightweight notification and reminder system built with a plugin-based architecture.

The system processes incoming events through a plugin mechanism and then delivers them using a transport layer. The main design goal is to remain small, efficient in CPU/RAM usage, and easily extensible for future improvements.

## Tech Stack

At the moment, the system is built without any external dependencies.  
All functionality is implemented using only the Go standard library.

## Planned Improvements

- **Scheduler**  
  A simple scheduler that allows configuring event execution time and recurrence.

- **External Plugins**  
  Support for third-party plugins without recompiling the core system, allowing plugins to be registered and used dynamically.

- **More Transports**  
  Add additional transport implementations besides Telegram.

## How to Use

### Build and Run

Build the project using `make` and start the server with a configuration file.

```sh
$ make build && ./build/payr --config <config_path>
```

## Configuration Example

```json
{
  "events": [
    {
      "name": "event 1",
      "transports": [
        "telegram"
      ],
      "plugin": {
        "type": "builtin",
        "name": "printer"
      }
    }
  ],
  "transports": [
    {
      "name": "telegram",
      "settings": {
        "sender": "",
        "channel_id": ""
      }
    }
  ]
}
```

### Configuration Overview

- `events` - list of available events
- `transports` - delivery mechanisms used to send notifications
- `plugin` - event processing logic

## Triggering an Event

Send an HTTP request to trigger an event manually:

```sh
curl -X POST 127.0.0.1:8080/event \
  -H "Content-Type: application/json" \
  -d '{"event":"event 1"}'
```

If the event exists in the configuration, it will be processed by the selected plugin and delivered through the configured transports.