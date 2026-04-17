---
title: Why Local-First?
description: Privacy, speed, and determinism
---

# Why Local-First?

`agrepl` is built on a **local-first** philosophy. We believe that developer tools for AI agents should be fast, private, and deterministic.

## 1. Privacy & Security

AI agents often handle sensitive data: system prompts, proprietary business logic, and PII (Personally Identifiable Information). Cloud-based observability platforms require you to send this data to their servers.

With `agrepl`, **nothing leaves your machine**.
- Traces are stored in local JSON files.
- Root CAs are generated locally.
- You have total control over your data.

## 2. Speed and Developer Loop

Recording and replaying agents should be instantaneous. By running everything locally, `agrepl` eliminates:
- Network latency to external APIs during replay.
- Slow cloud dashboards.
- Dependency on a stable internet connection for debugging.

## 3. Determinism

Cloud environments are noisy. By replaying agent sessions locally from a static trace, you ensure that:
- Every tool call returns the exact same result.
- LLM non-determinism is eliminated by caching responses.
- You can debug edge cases in a "frozen" environment.

## 4. Portability

Because runs are just JSON files, you can:
- Commit them to your git repository for regression testing.
- Send a single file to a teammate to reproduce a bug.
- Build your own tools on top of the open JSON format.
