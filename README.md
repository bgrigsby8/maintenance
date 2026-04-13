# Module maintenance

This module provides a sensor that reports whether the current time falls within a configured maintenance window. It is useful for enabling or disabling robot behaviors based on a scheduled time range.

## Models

This module provides the following model(s):

- [`brad-grigsby:maintenance:time-allowed`](brad-grigsby_maintenance_time-allowed.md) - A sensor that returns `true` when the current local time is within a configured start/end window, and `false` otherwise.
