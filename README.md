# swimlane

CLI kanban-style task system built on Markdown files with YAML frontmatter. Local-first, Git-friendly, and automation-friendly.

## Install

```bash
curl -sSL https://raw.githubusercontent.com/notwillk/swimlane/main/scripts/install.sh | sh
```

Or with a specific version:

```bash
curl -sSL https://raw.githubusercontent.com/notwillk/swimlane/main/scripts/install.sh | sh -s -- v1.0.0
```

Binary is installed to `/usr/local/bin/swimlane` by default. Set `INSTALL_DIR` to override.

## Config

Optional. Lookup order: `.swimlane.yaml`, `swimlane.yaml`, `~/.config/swimlane/config.yaml`. Override with:

```bash
swimlane --config path/to/config.yaml
```

Example `.swimlane.yaml`:

```yaml
tickets: "tickets/**/*.md"
default_path: tickets
default:
  priority: p2
  ready: true
  tags: []
```

## Commands

| Command | Description |
|---------|-------------|
| `swimlane ls` | List tickets (use `--csv` or `--json`) |
| `swimlane new "title"` | Create a new ticket |
| `swimlane next` | Print path of next ticket to implement |
| `swimlane done <ulid>` | Mark ticket complete (deletes file) |
| `swimlane schema-json config` | Print JSON schema for config |
| `swimlane schema-json ticket` | Print JSON schema for ticket frontmatter |
| `swimlane completion bash\|zsh\|fish` | Generate shell completions |
| `swimlane skill` | Output coding agent skill (for Cursor etc.) |

### Filters (ls and next)

Use `--priority`, `--tag`, `--status`, `--ready`. Prefix with `!` to negate (e.g. `--tag !infra`).

## Ticket format

- **Filename**: `[ULID]-[slug].md` (e.g. `01J9T8ZK1BC5A9JH56T9Y9M1DX-implement-login-api.md`)
- **Frontmatter** (required): `priority` (p0–p4), `status` (todo | in-progress | done), `ready` (boolean)
- **Frontmatter** (optional): `title`, `blocked_by` (list of ULIDs), `tags`
- **Body**: Markdown implementation instructions

## Build

```bash
go build -o bin/swimlane ./cmd/swimlane
```

Or with [just](https://github.com/casey/just): `just build`.

## License

See [LICENSE](LICENSE).
