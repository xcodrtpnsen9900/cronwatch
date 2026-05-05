# cronwatch

Monitors cron job execution and sends alerts on missed or failed runs via webhook.

## Installation

```bash
go install github.com/yourname/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronwatch.git && cd cronwatch && go build ./...
```

## Usage

Wrap your cron command with `cronwatch` to start monitoring it:

```bash
cronwatch --name "daily-backup" --schedule "0 2 * * *" --webhook "https://hooks.example.com/alert" -- /usr/local/bin/backup.sh
```

Configure via a YAML file for multiple jobs:

```yaml
# cronwatch.yaml
webhook: "https://hooks.example.com/alert"
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    command: /usr/local/bin/backup.sh
    timeout: 30m
  - name: hourly-sync
    schedule: "0 * * * *"
    command: /usr/local/bin/sync.sh
    timeout: 5m
```

```bash
cronwatch --config cronwatch.yaml
```

An alert is sent to the configured webhook if a job:
- Exits with a non-zero status code
- Exceeds its defined timeout
- Fails to run within the expected schedule window

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `cronwatch.yaml` |
| `--name` | Job name | required |
| `--schedule` | Cron expression | required |
| `--webhook` | Webhook URL for alerts | required |
| `--timeout` | Max allowed runtime | `1h` |

## License

MIT © yourname