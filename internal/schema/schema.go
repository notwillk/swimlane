package schema

// Config is the JSON schema for swimlane config files.
const Config = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "swimlane config",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "tickets": {
      "type": "string",
      "description": "Glob pattern for discovering ticket files"
    },
    "default_path": {
      "type": "string",
      "description": "Directory where new tickets are created"
    },
    "default": {
      "type": "object",
      "description": "Default values applied when creating tickets",
      "additionalProperties": false,
      "properties": {
        "$schema": {
          "type": "string",
          "description": "JSON schema URI for ticket frontmatter (applied when creating new tickets)"
        },
        "priority": { "type": "string", "enum": ["p0", "p1", "p2", "p3", "p4"] },
        "ready": { "type": "boolean" },
        "tags": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    },
    "actions": {
      "type": "object",
      "description": "Optional per-action command overrides; keys are action names (create, assign, etc.), value has 'command' with {arg-name} placeholders",
      "additionalProperties": {
        "type": "object",
        "properties": { "command": { "type": "string" } }
      }
    },
    "close_parent_when_subtasks_done": {
      "type": "string",
      "enum": ["never", "always", "when-empty", "when-matches"],
      "description": "When a subtask is marked done, close parent if all subtasks are done: never, always, when-empty (only if parent body empty), when-matches (LLM/normalized comparison)"
    }
  }
}
`

// Ticket is the JSON schema for ticket frontmatter.
const Ticket = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "swimlane ticket frontmatter",
  "type": "object",
  "additionalProperties": false,
  "required": ["priority", "status", "ready"],
  "properties": {
    "$schema": { "type": "string", "description": "JSON schema URI for this frontmatter" },
    "title": { "type": "string" },
    "priority": { "type": "string", "enum": ["p0", "p1", "p2", "p3", "p4"] },
    "status": { "type": "string", "enum": ["todo", "in-progress", "done"] },
    "ready": { "type": "boolean" },
    "assignee": { "type": "string", "description": "User assigned to this ticket" },
    "blocked_by": {
      "type": "array",
      "items": { "type": "string", "description": "ULID of blocking ticket" }
    },
    "subtasks": {
      "type": "array",
      "items": { "type": "string", "description": "ULID of subtask ticket" }
    },
    "tags": {
      "type": "array",
      "items": { "type": "string" }
    }
  }
}
`
