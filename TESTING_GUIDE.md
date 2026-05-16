# Testing Your Gemini Chat Agent with Agrepl

This guide shows you how to use `agrepl` to record and replay sessions with your Gemini chat agent. This allows you to create deterministic tests and debug your agent's behavior without making real LLM calls.

## 1. Setup

Ensure you have your `GEMINI_API_KEY` set in your environment:

```bash
export GEMINI_API_KEY=your_api_key_here
```

## 2. Record a Session

To record an interactive session, use the `agrepl record` command followed by the command to run your agent:

```bash
agrepl record go run examples/chat_agent.go
```

1.  Interact with the agent as you normally would.
2.  Type `exit` to finish the session.
3.  `agrepl` will save the execution trace and provide a **Run ID** (e.g., `run_12345`).

## 3. List Recorded Runs

You can see all your recorded runs using:

```bash
agrepl list
```

## 4. Replay and Verify

To replay a session and verify it still behaves correctly without hitting the real Gemini API:

```bash
agrepl replay <run-id> go run examples/chat_agent.go
```

During replay:
- `agrepl` intercepts the Gemini SDK calls.
- It returns the exact same responses that were recorded.
- No real network requests are made to Google's servers.

## 5. Summary and Comparison

You can see a summary of the replay:

```bash
agrepl replay <run-id> go run examples/chat_agent.go --summary
```

If you change your code (e.g., how you format history) and want to see if it still works with old responses, use the replay mode. If the agent makes a call that doesn't match the recording, `agrepl` will report a mismatch (unless `--fallback` is used).

## 6. Integrating into Tests

You can use `agrepl` in your `Makefile` or CI/CD pipeline to ensure your agent logic remains consistent:

```makefile
test-agent:
	agrepl replay run_latest go run examples/chat_agent.go --summary
```
