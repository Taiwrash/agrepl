# agrepl: Debug AI agents with deterministic replay

Record LLM and API interactions. Replay them offline with zero network calls. Reproduce any agent run instantly.

## The Problem
AI agents are non-deterministic. Debugging a failing run by re-executing your script is slow, expensive, and often impossible to reproduce because of LLM variability and changing API state.

## The Solution
`agrepl` is a local-first proxy that intercepts your agent's traffic. Once recorded, you can replay the entire session offline. **Zero network requests made.**

- **Reproduce** any agent run with 100% fidelity.
- **Debug** without burning API credits.
- **Works with any CLI tool** (Python, Node, Go, curl, etc.).
- **Model Agnostic**: Works with Gemini, OpenAI, Anthropic, or any custom API.

## Quick Start (30s)

### 1. Install
```bash
curl -sSL https://raw.githubusercontent.com/taiwrash/agrepl/main/scripts/install.sh | bash
```

### 2. Record
Run your agent through `agrepl`. It captures all HTTP(S) and LLM interactions.
```bash
agrepl record -- python agent.py
```

### 3. Replay (Offline)
Re-run with the same command. `agrepl` serves recorded responses from your local machine.
```bash
agrepl replay run-001 -- python agent.py
```
> **Killer Feature:** Replay works without an internet connection. No new API calls are made.

## Use Cases
- **Debug failing workflows**: Replay the exact session where your agent drifted.
- **CI/CD Testing**: Run integration tests with guaranteed, deterministic outputs.
- **Save Credits**: Iterate on tool logic or response handling without hitting paid APIs.
- **Bug Reports**: Attach a `run.json` to a ticket so teammates can reproduce the bug locally.

## Replay vs. Logging
Logging tells you what happened. **Replay lets you re-live it.** 
Traditional tracing tools show you a post-mortem. `agrepl` provides a live, interactive environment where your code *thinks* it's talking to the real world, but it's actually talking to a deterministic local cache.

## Roadmap

- [x] Basic CLI (record, replay, list, diff)
- [x] Local JSON storage
- [x] MITM Proxy for HTTP/HTTPS interception
- [x] Automatic Root CA trust injection
- [x] Structural JSON matching for robust replays
- [x] Binary data support (for gRPC and images)
- [ ] gRPC / HTTP/2 optimized matching
- [ ] Official SDK Wrappers (Python, Node.js)
- [ ] Fuzzy/AI-based request matching
- [ ] Enhanced Remote Storage (Push/Pull)
- [ ] CI/CD integration for regression testing
- [ ] Token cost estimation & analytics

## License

MIT
