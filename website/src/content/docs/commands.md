---
title: Commands
description: Reference for all CLI commands
---

# Command Reference

`agrepl` provides a comprehensive set of commands to manage your agent runs.

## `record`

Starts a recording proxy and executes your agent command.

```bash
agrepl record -- [command] [args...]
```

Example:
```bash
agrepl record -- python agent.py --task "research"
```

## `replay`

Replays a previously recorded run offline.

```bash
agrepl replay [run-id] -- [command] [args...]
```

Example:
```bash
agrepl replay run-001 -- python agent.py
```

## `list`

Lists all recorded runs in the local index.

```bash
agrepl list
```

## `diff`

Compares two recorded runs to see how behavior changed.

```bash
agrepl diff [run-a] [run-b]
```

## `auth`

Manage your agrepl Cloud identity.

```bash
agrepl auth login
agrepl auth logout
agrepl auth status
```

## `share`

Register and upload a run for team collaboration.

```bash
agrepl share [run-id]
```

## `pull`

Download a shared run from your team workspace.

```bash
agrepl pull [share-id]
```
