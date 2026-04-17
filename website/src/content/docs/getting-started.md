---
title: Getting Started
description: Installation and first run
---

# Getting Started

Get up and running with `agrepl` in minutes.

## Installation

The easiest way to install `agrepl` is via our installation script:

```bash
curl -sSL https://raw.githubusercontent.com/taiwrash/agrepl/main/scripts/install.sh | bash
```

Alternatively, you can build from source:

```bash
git clone https://github.com/taiwrash/agrepl
cd agrepl
make install
```

## First Run: Recording a Session

To start recording your agent's interactions, simply wrap your agent's execution command with `agrepl record`:

```bash
agrepl record -- python agent.py
```

This starts a local MITM proxy, automatically trusts the local CA for the child process, and records all HTTP(S) and LLM interactions to the `.agent-replay/` directory.

## Replaying a Session

Once you have recorded a run (e.g., `run-001`), you can replay it offline:

```bash
agrepl replay run-001 -- python agent.py
```

`agrepl` will serve the recorded responses instead of calling live APIs. **Zero network requests will be made.**
