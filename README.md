# vaultdiff

A CLI tool to diff and audit changes between HashiCorp Vault secret versions.

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git && cd vaultdiff && go build -o vaultdiff .
```

## Usage

Compare two versions of a Vault secret:

```bash
vaultdiff --path secret/data/myapp --v1 3 --v2 4
```

Show a full audit trail of changes across all versions:

```bash
vaultdiff --path secret/data/myapp --audit
```

**Example output:**

```
--- version 3 (2024-01-10 12:00:00)
+++ version 4 (2024-01-15 09:30:00)

~ DB_PASSWORD  [changed]
+ NEW_API_KEY  [added]
- OLD_TOKEN    [removed]
```

## Configuration

`vaultdiff` respects standard Vault environment variables:

| Variable | Description |
|---|---|
| `VAULT_ADDR` | Vault server address |
| `VAULT_TOKEN` | Authentication token |
| `VAULT_NAMESPACE` | Vault namespace (Enterprise) |

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

## License

MIT © [yourusername](https://github.com/yourusername)