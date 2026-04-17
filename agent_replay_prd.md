# Product Requirements Document (PRD)

## Product Name
agent-replay

## Summary
agent-replay is a lightweight CLI tool that enables developers to record and deterministically replay AI agent executions. It focuses on reproducibility, debugging, and offline testing of agent workflows by capturing LLM interactions and external tool/API calls.

---

## Problem Statement

AI agents are inherently non-deterministic due to:
- LLM variability
- External API/tool dependencies
- Changing environments

This leads to:
- Difficult debugging
- Irreproducible failures
- High iteration cost
- Hidden regressions

There is currently no simple, local-first tool to reliably reproduce agent runs.

---

## Goals

1. Enable deterministic replay of agent executions
2. Provide a simple CLI-first developer experience
3. Support offline debugging of agent workflows
4. Capture both LLM and tool interactions

---

## Non-Goals (MVP)

- No UI/dashboard
- No multi-agent orchestration
- No deployment system
- No distributed runtime
- No deep framework integrations initially

---

## Target Users

- AI engineers building agents
- Backend/infra engineers working with LLM workflows
- Developers debugging tool-using agents

---

## Core Features (MVP)

### 1. Record Mode

Command:
```
agent-replay record <command>
```

Description:
- Executes a user command (e.g., Python script, notebook runner)
- Intercepts and records:
  - LLM calls (Gemini SDK)
  - HTTP requests/responses

Output:
- Stores execution trace locally

---

### 2. Replay Mode

Command:
```
agent-replay replay <run-id>
```

Description:
- Replays a previously recorded execution
- Returns recorded responses instead of making real calls

Properties:
- Deterministic
- Offline-capable

---

### 3. Run Storage

Directory structure:
```
.agent-replay/
  runs/
    run-001.json
```

Run format:
```
{
  "run_id": "run-001",
  "steps": [
    {
      "type": "llm",
      "input": "...",
      "output": "..."
    },
    {
      "type": "http",
      "request": {},
      "response": {}
    }
  ]
}
```

---

## System Architecture

### High-Level

```
CLI (Go)
  ↓
Interceptor Layer
  ↓
Recorder / Replay Engine
  ↓
Local Storage (JSON)
```

---

## Technical Design

### Language
- Go (CLI + core engine)

### LLM Integration
- Gemini SDK
- Wrap client calls to capture:
  - prompt
  - response
  - metadata

### Interception Strategy

#### HTTP Interception
- Use custom HTTP transport in Go
- Wrap requests and responses

#### Gemini Interception
- Build wrapper client around Gemini SDK

---

## Key Components

### 1. CLI Layer
Responsibilities:
- Parse commands
- Manage run lifecycle

Commands:
- record
- replay

---

### 2. Recorder
Responsibilities:
- Capture events (LLM, HTTP)
- Serialize to JSON

---

### 3. Replay Engine
Responsibilities:
- Match incoming requests
- Return stored responses

Matching Strategy:
- Exact match (MVP)
- Hash-based lookup (later)

---

### 4. Storage Layer
- File-based JSON storage
- Simple read/write operations

---

## Developer Experience

### Example Flow

1. Record run:
```
agent-replay record python test-ml.ipynb
```

2. Replay run:
```
agent-replay replay run-001
```

---

## Success Metrics

- Time to first successful replay < 5 minutes
- Deterministic replay success rate > 95%
- CLI usability (low friction setup)

---

## Risks & Mitigations

### 1. LLM Non-Determinism
Mitigation:
- Always cache LLM responses

### 2. Request Matching Complexity
Mitigation:
- Start with exact matching
- Improve later with fuzzy matching

### 3. Integration Friction
Mitigation:
- Keep setup minimal
- Provide examples

---

## Future Enhancements

- Diff between runs
- Cost tracking
- Tool-level mocking
- CI integration
- Pluggable runtimes (Daytona, etc.)

---

## Timeline (MVP)

Week 1:
- CLI setup
- basic recording

Week 2:
- replay engine
- HTTP interception

Week 3:
- Gemini integration
- polish

---

## Open Questions

- How to best inject interception into Python workflows?
- Should we support multiple languages early?
- What level of matching is needed for replay accuracy?

---

## Conclusion

agent-replay is a focused, developer-first tool that solves a critical gap in the AI agent ecosystem: reproducibility. By starting simple (record + replay), it creates a strong foundation for future expansion into a full agent infrastructure layer.

