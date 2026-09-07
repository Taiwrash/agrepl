# agrepl: Debug AI agents with deterministic replay

<p align="center">
  <a href="https://arxiv.org/abs/2607.16200"><img src="https://img.shields.io/badge/arXiv-2607.16200-b31b1b.svg" alt="arXiv Paper"></a>
  <a href="https://agrepl.taiwrash.xyz"><img src="https://img.shields.io/badge/docs-agrepl.taiwrash.xyz-blue.svg" alt="Docs"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License: MIT"></a>
</p>

Record LLM and API interactions. Replay them offline with zero network calls. Reproduce any agent run instantly.

> Grounded in research: Based on the paper *"Deterministic Replay for AI Agent Systems"* ([arXiv:2607.16200](https://arxiv.org/abs/2607.16200)), `agrepl` introduces formal replay invariants and noise-aware differential analysis to make stochastic agent runs fully reproducible.

## The Problem
AI agents that couple LLMs with external tools and APIs are inherently non-deterministic. Debugging a failing run by re-executing your script is slow, expensive, and often impossible to reproduce due to LLM sampling variance, external API state mutation, and environmental noise.

## The Solution: Beyond Observability
Traditional tracing tools (LangSmith, W&B) let you **observe** bugs after they happen. `agrepl` lets you **reproduce** them with mathematical determinism.

`agrepl` is a deterministic execution layer for AI agents. It intercepts external interactions at the transport layer, serialises them into structured execution traces, and serves them back with 100% fidelity ($F = 1.0$).

- **Deterministic Invariant**: Replay is a canonical key lookup, not a re-execution. Zero network calls, zero token consumption.
- **Zero-Instrumentation**: Transport-level MITM proxy. Works out of the box with Python, Node, Go, curl, or any CLI without code modifications or SDKs.
- **Noise-Aware Diffing**: Distinguishes semantic agent drift from non-semantic infrastructure noise (CDN headers, dynamic timestamps).
- **Shareable Ground Truth**: `share` an execution trace → `pull` on another machine → Replay the exact bug locally.

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
Traditional tracing tools show you a post-mortem. `agrepl` provides a live, interactive environment where your code *thinks* it's talking to the real world, but it's actually talking to a deterministic local cache. Empirical evaluation across five benchmark agent workloads ($n = 250$ replay instances) demonstrates replay fidelity $F = 1.0$ and a median per-step latency reduction of 98.3%.

**See how we compare to LangSmith, VCR.py, and more: [agrepl vs. The Orbit](https://agrepl.taiwrash.xyz/compare)**

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

## Citation

`agrepl` is developed alongside academic research on deterministic execution systems for stochastic AI agents. If you use `agrepl` in your research, agent evaluations, or benchmark suites, please cite our paper:

```bibtex
@article{mudasiru2026deterministic,
  title={Deterministic Replay for AI Agent Systems},
  author={Mudasiru, Rasheed},
  journal={arXiv preprint arXiv:2607.16200},
  year={2026},
  eprint={2607.16200},
  archivePrefix={arXiv},
  primaryClass={cs.AI},
  url={https://arxiv.org/abs/2607.16200},
  doi={10.48550/arXiv.2607.16200}
}
```

The full paper is available on arXiv: [arXiv:2607.16200 [cs.AI]](https://arxiv.org/abs/2607.16200).

## License

MIT

