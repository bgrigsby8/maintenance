# Model brad-grigsby:maintenance:time-allowed

A [sensor](https://docs.viam.com/components/sensor/) component that reports whether the current local time falls within a configured maintenance window. Call `Readings` to get a single boolean value — `true` when the machine is inside the allowed window, `false` otherwise.

Overnight windows (e.g. `22:00`–`06:00`) are supported automatically.

## Configuration

### Attributes

| Name         | Type   | Inclusion | Description                                              |
|--------------|--------|-----------|----------------------------------------------------------|
| `start_time` | string | Required  | Start of the allowed window in `HH:MM` (24-hour) format |
| `end_time`   | string | Required  | End of the allowed window in `HH:MM` (24-hour) format   |

### Example Configuration

```json
{
  "start_time": "08:00",
  "end_time": "17:00"
}
```

## Readings

Returns a single key indicating whether the current time is within the configured window.

| Key          | Type | Description                                            |
|--------------|------|--------------------------------------------------------|
| `is_allowed` | bool | `true` if current time is within the maintenance window |

### Example Response

```json
{
  "is_allowed": true
}
```

