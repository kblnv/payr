Simple reminder system

Example config file (config/registry.json):
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
        "name": "text"
      }
    }
  ],
  "transports": [
    {
      "name": "telegram",
      "sender": "",
      "channel_id": ""
    }
  ]
}
```