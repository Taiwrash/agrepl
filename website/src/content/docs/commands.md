---
title: Commands
description: Reference for all CLI commands
---

# Command Reference

`agrepl` provides a comprehensive set of commands to manage your agent runs.

## `record`

Starts a recording proxy.

```bash
agrepl record [flags]
```

- `--port`: The port to run the proxy on (default: 8080).
- `--ca-cert`: Path to a custom CA certificate for HTTPS interception.

## `replay`

Replays a previously recorded run.

```bash
agrepl replay [run-id] [flags]
```

- `--port`: The port to run the replay server on (default: 8080).

## `list`

Lists all recorded runs in the local workspace.

```bash
agrepl list
```

## `diff`

Compares two recorded runs.

```bash
agrepl diff [run-a] [run-b]
```

## `push`

Pushes a run to a remote storage server.

```bash
agrepl push [run-id]
```

## `pull`

Pulls a run from a remote storage server.

```bash
agrepl pull [run-id]
```
