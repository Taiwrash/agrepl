# agrepl: Debug AI agents with deterministic replay

Record LLM and API interactions. Replay them offline with zero network calls. Reproduce any agent run instantly.

## The Problem
AI agents are non-deterministic. Debugging a failing run by re-executing your script is slow, expensive, and often impossible to reproduce because of LLM variability and changing API state.

## The Solution: Beyond Observability
Traditional tracing tools (LangSmith, W&B) let you **observe** bugs. `agrepl` lets you **reproduce** them.

`agrepl` is a deterministic execution layer for AI agents. It captures every LLM and API interaction and serves them back with 100% fidelity.

- **Deterministic**: Replay is a lookup, not a re-execution. Zero network calls.
- **Zero-Instrumentation**: No SDKs. Works with Python, Node, Go, curl, or any CLI.
- **Shareable Truth**: `share` a run → `pull` on another machine → Replay the same bug instantly.

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

### 3. Replay
Re-run instantly. `agrepl` finds the original command and re-executes it in a deterministic sandbox.
```bash
agrepl run run-010
```
> **The Edge:** Replay works without an internet connection. No new API calls are made. Zero latency. Zero cost.

### 4. Share (Team Collaboration)
Share your run with your team for collaborative debugging.
```bash
agrepl auth login
agrepl share run-001
```
`agrepl` will generate a unique ID. Your teammates can then use `agrepl pull [id]` to reproduce the exact failure on their machines.

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
