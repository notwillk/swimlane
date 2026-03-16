# swimlane

Use this skill when working with a repository that uses **swimlane** for task management.

## What is swimlane?

swimlane is a CLI kanban-style task system. Tickets are **Markdown files with YAML frontmatter**. They live in the repo (e.g. under `tickets/`) and are version-controlled with Git.

## Commands

- **`swimlane ls`** — List tickets. Use `--priority`, `--tag`, `--status`, `--ready` to filter. Prefix with `!` to negate (e.g. `--tag !infra`). Use `--csv` or `--json` for machine-readable output.
- **`swimlane create [title]`** — Create a new ticket (reads description from stdin unless `--no-description`). Use `--claim`, `--start`, etc. to also apply lifecycle actions.
- **`swimlane next`** — Print the path of the next ticket to implement (ready, todo, not blocked; respects priority and filters). Use the same filter flags as `ls` for sharding (e.g. `--tag backend`).
- **`swimlane done <ulid>`** — Mark a ticket complete by deleting its file. Use the ULID from the filename or from `swimlane ls`.
- **`swimlane schema-json config`** / **`swimlane schema-json ticket`** — Print JSON schema for config or ticket frontmatter.
- **`swimlane completion bash|zsh|fish`** — Generate shell completions.

## Ticket format

- **Filename**: `[ULID]-[slug].md` (e.g. `01J9T8ZK1BC5A9JH56T9Y9M1DX-implement-login-api.md`).
- **Frontmatter** (required): `priority` (p0–p4), `status` (todo | in-progress | done), `ready` (boolean).
- **Frontmatter** (optional): `title`, `blocked_by` (list of ULIDs), `tags` (list of strings).
- **Body**: Markdown implementation instructions.

## Workflow for agents

1. Run **`swimlane next`** (optionally with `--tag` or other filters) to get the path of the next ticket.
2. Read that file, implement the instructions in the body, and update the repo.
3. Run **`swimlane done <ulid>`** to mark it complete (delete the file).
4. Repeat.

## Config

Config is optional. Lookup order: `.swimlane.yaml`, `swimlane.yaml`, `~/.config/swimlane/config.yaml`. Override with `swimlane --config path/to/config.yaml`. Key fields: `tickets` (glob, e.g. `tickets/**/*.md`), `default_path`, `default.priority`, `default.ready`, `default.tags`.
