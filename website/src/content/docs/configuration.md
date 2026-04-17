---
title: Configuration
description: Env vars and settings
---

# Configuration

`agrepl` can be configured via command-line flags or environment variables.

## Environment Variables

- `AGREPL_STORAGE_DIR`: Path to the directory where runs are stored (default: `./.agent-replay`).
- `AGREPL_REMOTE_URL`: URL of the remote storage server.
- `AGREPL_LOG_LEVEL`: Logging verbosity (`debug`, `info`, `warn`, `error`).

## SSL/TLS Interception

For HTTPS requests, `agrepl` generates a self-signed CA certificate on the first run. You may need to trust this certificate in your system or application to intercept encrypted traffic.

The CA certificate is located at:
`./.agent-replay/ca/ca.crt`

In Python (with `requests` or `openai`):

```python
import os
os.environ["REQUESTS_CA_BUNDLE"] = "./.agent-replay/ca/ca.crt"
```
