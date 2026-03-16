# swimlane task management

Use this skill when you are a coding agent working in a repository that manages work with **swimlane**.

This skill tells you how to:
- Discover and pick the right ticket to work on.
- Apply changes safely based on ticket instructions.
- Keep tickets and config valid so other tools and agents can keep using them.

---

## Concepts

- **Tickets are files**: Each ticket is a Markdown file with YAML frontmatter.
- **Filenames encode identity**: `ULID-slug.md`, for example:
  - `01J9T8ZK1BC5A9JH56T9Y9M1DX-implement-login-api.md`
  - The ULID is the stable identifier; the slug is human-readable.
- **Frontmatter fields** (YAML, first `---` block):
  - Required:
    - `priority`: one of `p0`, `p1`, `p2`, `p3`, `p4` (lower = higher priority)
    - `status`: one of `todo`, `in-progress`, `done`
    - `ready`: boolean
  - Optional:
    - `$schema`: JSON Schema URI or path for the frontmatter (if present, new tickets often reuse it)
    - `title`: short description
    - `assignee`: user assigned to the ticket (e.g. for `swimlane ls --mine`)
    - `blocked_by`: ULID list of tickets that must finish first
    - `subtasks`: ULID list of tickets that break down this ticket’s work (when all subtasks are done, parent may be auto-closed per config)
    - `tags`: string list (e.g. `backend`, `frontend`, `infra`, `auth`)
- **Config file** (usually `.swimlane.yaml`, `swimlane.yaml`, or `~/.config/swimlane/config.yaml`):
  - `tickets`: glob for ticket discovery (e.g. `tickets/**/*.md`, `todo/**/*.md`)
  - `default_path`: where new tickets are written
  - `default`: default frontmatter (including optional `$schema`)

You should treat the ULID plus path as the canonical identifier, and treat frontmatter as the “task metadata” you must keep valid.

---

## CLI commands you should know

Always run these from the repo root unless a config or fixture instructs you otherwise.

- **List tickets**
  - `swimlane ls`
  - Flags:
    - `--priority` (e.g. `--priority p1`, `--priority !p0`)
    - `--tag` (e.g. `--tag backend`, `--tag !infra`)
    - `--status` (e.g. `--status todo`)
    - `--ready` (e.g. `--ready true`, `--ready !true`)
  - Output formats:
    - Default: human table, e.g. `P1 todo   01J9T8... implement-login-api`
    - `--csv`: machine-friendly CSV
    - `--json`: JSON array of tickets

- **Create a new ticket**
  - `swimlane create "implement login api"`
  - Prints the created path, e.g. `todo/01KKW52FZYXA77KTF1JFZSJJY0-life-cycle-actions.md`
  - Frontmatter is seeded from config `default` (including `$schema` if present).

- **Pick the next ticket to implement**
  - `swimlane next`
  - Returns a single path:
    - Only considers tickets where `ready == true` and `status == todo`.
    - Respects priority order `p0`..`p4`.
    - Skips tickets blocked by others; may return a blocking dependency if everything at a priority is blocked.
  - Accepts the same filters as `ls` (e.g. `swimlane next --tag backend`).

- **Mark a ticket done**
  - `swimlane done <ulid>`
  - Locates the ticket by ULID and **deletes the file**.
  - Use ULID from:
    - `swimlane ls` output, or
    - The filename itself.

- **Validate structures**
  - `swimlane static`
    - Validates:
      - Config YAML (if present) and its metadata.
      - Every ticket file’s frontmatter (YAML + required fields).
    - Exit code:
      - `0` if all checks pass.
      - Non-zero if any failure. All failures are printed; it does not stop after the first.
  - `swimlane schema-json config` / `swimlane schema-json ticket`
    - Print JSON Schema definitions for config and ticket frontmatter.

- **Other**
  - `swimlane completion bash|zsh|fish` — shell completions.
  - `swimlane skill` — built-in, minimal skill description (project-local; this `.agents` skill is for broader use across repos).

---

## How to pick and execute work

When you are asked to “work from swimlane tickets” in a repo:

1. **Ensure you’re in the right repo and config**
   - If there is a repo-local config:
     - Look for `.swimlane.yaml` / `swimlane.yaml` at the root.
     - Or use a provided config path: `swimlane --config path/to/config.yaml <command>`.
   - If the repo has **fixtures** (e.g. `fixtures/`), only use them when explicitly asked or for tests/demos.

2. **Discover candidate tickets**
   - Run:
     - `swimlane ls`
     - Or for specific slices:
       - `swimlane ls --tag backend`
       - `swimlane ls --priority p1`
       - `swimlane ls --status todo --ready true`
   - Use `--csv` or `--json` if you need structured inspection.

3. **Pick the next ticket**
   - Preferred:
     - `swimlane next`
     - Or shard by filters, e.g.:
       - `swimlane next --tag backend`
       - `swimlane next --priority !p0 --tag !infra`
   - Open the returned file path and read:
     - YAML frontmatter (metadata, dependencies).
     - Markdown body (implementation instructions).

4. **Implement the ticket**
   - Follow instructions in the body.
   - Make code/config changes as described.
   - Keep frontmatter fields valid:
     - Do not change `priority` / `status` / `ready` semantics unless asked.
     - Keep `blocked_by` and `tags` structurally correct if you touch them.
   - If the repo defines additional conventions for tickets (e.g. custom sections), follow those as well.

5. **Run validation / tests**
   - At minimum:
     - `swimlane static`
   - If the repo has additional commands (examples; may vary by repo):
     - `just static`
     - `just test`
   - Fix any issues and re-run until they pass.

6. **Mark the ticket done (if appropriate)**
   - When the work is complete and tests/validation pass:
     - Extract ULID from the ticket filename or `ls`:
       - Example filename: `01KKW52FZYXA77KTF1JFZSJJY0-life-cycle-actions.md`
       - ULID: `01KKW52FZYXA77KTF1JFZSJJY0`
     - Run:
       - `swimlane done 01KKW52FZYXA77KTF1JFZSJJY0`
   - This deletes the ticket file; the history remains in Git.

---

## Filters and sharding strategies

Filters work the same for `ls` and `next` and support negation with `!`.

- **Priority-based selection**
  - Examples:
    - `swimlane next --priority p0` — most critical work.
    - `swimlane next --priority !p0` — anything except p0.

- **Tag-based sharding**
  - Common patterns:
    - Backend agent: `swimlane next --tag backend`
    - Frontend agent: `swimlane next --tag frontend`
    - Infra agent: `swimlane next --tag infra`
  - Exclusions:
    - Skip infra: `swimlane next --tag !infra`

- **Status / readiness**
  - Prefer to keep `next` focused on `status=todo` and `ready=true` (which `next` enforces).
  - Use `ls` to inspect other states:
    - `swimlane ls --status in-progress`
    - `swimlane ls --ready !true`

When coordinating with other agents, choose disjoint filter sets (e.g. one agent per tag or priority range) to minimize conflicts.

---

## Config, schemas, and fixtures

- **Config discovery**
  - Default lookup:
    - `.swimlane.yaml`
    - `swimlane.yaml`
    - `~/.config/swimlane/config.yaml`
  - Override:
    - `swimlane --config path/to/config.yaml <command>`

- **`default.$schema`**
  - Some configs set `default.$schema` (e.g. a path like `.schemas/ticket.json`).
  - New tickets created with `swimlane create` will copy this `$schema` into ticket frontmatter.
  - When editing frontmatter, keep `$schema` pointing at the intended schema unless asked to change it.

- **Fixtures**
  - Repos may provide under `fixtures/`:
    - `fixtures/empty/`, `fixtures/simple/`, `fixtures/blocked/`, `fixtures/with-tags/`, etc.
  - Use these when:
    - Writing tests.
    - Demonstrating behavior of `ls`, `next`, or dependency resolution.
  - Do **not** treat fixture tickets as real work unless the instructions explicitly say so.

---

## Safety and best practices for agents

- **Do**
  - Use `swimlane next` to select work instead of scanning files manually whenever possible.
  - Keep ticket frontmatter valid YAML and conforming to the known schema:
    - Required: `priority`, `status`, `ready`.
    - Only use allowed values for `priority` and `status`.
  - Run `swimlane static` before and after large batches of changes.
  - Use filters to avoid stepping on other agents’ work (e.g. shard by tag).

- **Do not**
  - Do **not** rename ULIDs, change ticket filenames arbitrarily, or reuse ULIDs.
  - Do **not** manually delete ticket files; prefer `swimlane done <ulid>` so tooling stays in sync.
  - Do **not** introduce new frontmatter keys outside the schema unless repo instructions explicitly allow it.
  - Do **not** modify tickets outside your scope (e.g. changing unrelated priorities/tags) unless the user asks for that refactor.

By following this skill, you help keep the swimlane board consistent so that humans and other agents can reliably coordinate on the same set of tickets.

---

## Using swimlane during planning (agentic coders)

When you are in a **planning mode** (for example, a user asks you to “create a plan” or you have an explicit plan stage) and the repository uses swimlane:

1. **Represent plan items as tickets, not prose todos**
   - Instead of inventing your own ad-hoc todo list in the plan, you should:
     - Identify concrete, actionable tasks that belong on the swimlane board.
     - For each such task that does not already exist, create a new ticket with:
       - `swimlane create "short, action-oriented title"` (optionally with `--config` if needed).
     - When a task naturally breaks into sub-items, create a **parent ticket** for the overall work and **subtask tickets** for each part; then set the parent’s `subtasks` frontmatter to the list of those subtask ULIDs (in order). That way the board reflects the hierarchy and parents can auto-close when all subtasks are done (per repo config).
   - Use the config’s `default` block to pre-fill `priority`, `ready`, `tags`, and optional `$schema`.

2. **Name and tag tickets appropriately**
   - Titles should be short and specific (e.g. “implement static command for YAML validation”, “add devcontainer feature for swimlane”).
   - When choosing `tags` (via config defaults or manual edits), prefer:
     - Functional domains: `backend`, `frontend`, `infra`, `docs`, `cli`, etc.
     - Any project-specific tagging conventions if present.

3. **Plan output should list tickets, not prose descriptions**
   - When returning a plan to the user, **do not** list free-form bullet points describing what you will do.
   - Instead, list the newly created tickets, for example:
     - Their **paths**: `todo/01KKW6C11RGE9XY4P29Z0ASZFC-add-assignees.md`
     - Or ULID + title pairs: `01KKW6C11RGE9XY4P29Z0ASZFC — add assignees`
   - Clearly distinguish:
     - Tickets you just created in this planning step.
     - Existing tickets you intend to work on (refer to them by ULID/path, do not duplicate them).

4. **Keep the plan and board in sync**
   - The plan you present should be derivable from the ticket set:
     - If the user agrees to the plan, subsequent implementation steps should operate **through the tickets** (using `next`, `ls`, `done`), not through a separate internal todo list.
   - If the user asks to adjust the plan (e.g. “drop this item” or “add one more step”), you should:
     - Update tickets accordingly (create, remove, or edit tickets as requested), then
     - Present the updated ticket list instead of editing a prose plan only.

5. **Coordination across agents**
   - If multiple agentic coders (e.g. different tools/LLMs) are active:
     - Each agent should prefer to plan and track work via swimlane tickets.
     - Use tags and priorities to partition responsibilities instead of local todo lists.
     - Plans should reference shared tickets, not private task names.

The high-level rule: **during planning, treat “creating or updating swimlane tickets” as the primary way to express your plan**. The plan you show to the user should primarily be a view of those tickets, not a separate prose checklist.

---

## Batch implementation flows (e.g. “implement all p2s”)

Sometimes you will be asked to implement **many tickets in one session**, for example:
- “implement all p2s”
- “implement all todos”

When following these instructions, you should:

### Implement all tickets at or above a priority (e.g. “implement all p2s”)

1. **Interpret the request**
   - “implement all p2s” means:
     - Implement all tickets with `priority == p2` (and possibly higher priorities, if present).
     - Stop when there are no more eligible `p2` (or higher) tickets that are ready and unblocked.

2. **Loop using `swimlane next`**
   - Use a loop like:
     - Call `swimlane next --priority p0 --priority p1 --priority p2` (or repo-specific equivalent).
     - If the command fails with “no next ticket”, stop.
     - If it prints a path:
       - Open the ticket, implement it, validate (`swimlane static`, tests), and mark it done (`swimlane done <ulid>`).
       - Repeat.
   - Be aware:
     - Because of dependency resolution, `swimlane next` may sometimes return a **lower priority dependency** needed to unblock a higher-priority ticket.
     - After several iterations, check remaining candidates explicitly with:
       - `swimlane ls --priority p0 --priority p1 --priority p2 --status todo --ready true`
       - If any tickets remain, continue the loop or report back to the user why they are still blocked.

3. **Reporting progress**
   - When you summarize what you did, report:
     - Which tickets (ULIDs and titles) you completed.
     - Which `p2` (or above) tickets remain and whether they are blocked (`blocked_by`) or not ready.

### Implement “all todos” with user confirmation

When the request does **not** specify a priority (e.g. “implement all todos”), proceed more cautiously:

1. **Estimate the scope**
   - List candidates first, for example:
     - `swimlane ls --status todo --ready true --csv` (or `--json`) and count them.
   - Determine:
     - Total number of ready `todo` tickets (N).

2. **Ask for confirmation and constraints**
   - Before starting a large batch, ask the user something equivalent to:
     - “There are **N ready todo tickets**.  
        Do you want me to:
        - implement **all** of them,  
        - implement **only the first M**, or  
        - implement **only tickets at or above a given priority (e.g. p2 and above)**?”
   - Respect the user’s answer:
     - If they specify a maximum count M:
       - Loop through `swimlane next` at most M times, implementing each returned ticket.
     - If they specify a minimum priority:
       - Use the “implement all tickets at or above a priority” approach above.

3. **Execution loop**
   - For each ticket you decide to implement as part of this batch:
     - Use `swimlane next` (optionally filtered).
     - Implement, validate (`swimlane static`, tests), and call `swimlane done <ulid>`.
   - If `swimlane next` returns a ticket that does **not** satisfy the agreed constraints (e.g. lower priority than requested), you should:
     - Explain to the user why (likely a dependency) and either:
       - Obtain consent to implement that dependency as part of the batch, or
       - Skip it and adjust the plan.

4. **End-of-batch behavior**
   - Stop when one of these is true:
     - You have implemented the agreed number of tickets (M).
     - There are no more `todo` tickets satisfying the agreed priority/filters.
     - `swimlane next` reports “no next ticket”.
   - Summarize:
     - Tickets completed.
     - Tickets remaining (with ULID, title, and priority).
     - Any tickets you intentionally skipped and why (e.g. non-ready, deeply blocked).

