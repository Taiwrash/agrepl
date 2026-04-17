---
title: Introduction
description: Overview of agrepl
---

# Introduction to agrepl

`agrepl` (Agent Replay) is a powerful CLI tool and library designed for debugging, testing, and monitoring AI agents. It acts as a transparent proxy between your AI agent and LLM providers (like OpenAI, Anthropic, or local models), recording all interactions for later analysis and replay.

## Key Features

- **Transparent Proxying**: Intercepts HTTP/LLM calls with zero-config for most setups.
- **Recording & Replay**: Capture agent sessions and replay them exactly as they occurred.
- **Deterministic Testing**: Use recorded responses to test agent behavior without hitting live APIs.
- **Visual Diffing**: Compare runs to see how changes in your agent's code or prompts affect LLM interactions.
- **Remote Storage**: Push and pull recorded runs to a central server for team collaboration.

## Why use agrepl?

Testing AI agents is notoriously difficult due to the stochastic nature of LLMs. `agrepl` solves this by giving you a reliable way to capture "what happened" and "what changed". Whether you're debugging a complex multi-step reasoning loop or regression testing a new prompt, `agrepl` provides the visibility you need.
