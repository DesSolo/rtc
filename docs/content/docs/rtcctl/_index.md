---
date: '2025-09-14T21:33:12+03:00'
draft: false
title: 'rtcctl'
---

`rtcctl` is a command-line interface (CLI) tool designed to manage your RTC server installation via its API. It's perfect for automation scripts, CI/CD pipelines, and quick administrative tasks.

## Launch Options / Flags

The following flags are available for all `rtcctl` commands:

| Flag          | Shorthand | Default Value                  | Description                                                                                                            |
|:--------------|:----------|:-------------------------------|:-----------------------------------------------------------------------------------------------------------------------|
| `--url`       | `-u`      | `http://localhost:8080/api/v1` | The URL of the RTC server to connect to.                                                                               |
| `--token`     | `-t`      | (none)                         | The API access token. If not provided as a flag, the value will be taken from the `RTCCTL_TOKEN` environment variable. |
| `--log-level` | `-l`      | `0`                            | The logging level, using the `slog` format.                                                                            |
