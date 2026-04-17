---
title: Getting Started
description: Installation and first run
---

# Getting Started

Get up and running with `agrepl` in minutes.

## Installation

`agrepl` is written in Go. You can install it using `go install`:

```bash
go install github.com/taiwrash/agrepl@latest
```

Alternatively, download the pre-compiled binaries from the [Releases](https://github.com/taiwrash/agrepl/releases) page.

## First Run: Recording a Session

To start recording your agent's interactions, use the `record` command:

```bash
agrepl record --port 8080
```

This starts a local proxy server. Point your AI agent's LLM client to `http://localhost:8080`. For example, if you are using OpenAI's client:

```python
import openai

client = openai.OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="sk-..." # Your real key
)
```

As your agent runs, `agrepl` will record all requests and responses to the `.agent-replay/` directory.

## Replaying a Session

Once you have recorded a run (e.g., `run-001`), you can replay it:

```bash
agrepl replay run-001 --port 8080
```

Now, when your agent makes the same requests to the proxy, `agrepl` will serve the recorded responses instead of calling the live API.
