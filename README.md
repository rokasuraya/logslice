# logslice

A fast command-line tool for extracting and filtering structured log ranges by timestamp or field value.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Extract logs between two timestamps
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log

# Filter by a specific field value
logslice --field level=error app.log

# Combine timestamp range with field filter
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --field service=api app.log

# Read from stdin
cat app.log | logslice --field level=warn
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of timestamp range (RFC3339) |
| `--to` | End of timestamp range (RFC3339) |
| `--field` | Filter by field in `key=value` format |
| `--format` | Log format: `json` (default) or `logfmt` |
| `--out` | Output file (defaults to stdout) |

---

## Requirements

- Go 1.21 or later

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)