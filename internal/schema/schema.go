package schema

// Config is the JSON schema for swimlane config files.
const Config = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "swimlane config",
  "type": "object",
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
      "properties": {
        "priority": { "type": "string", "enum": ["p0", "p1", "p2", "p3", "p4"] },
        "ready": { "type": "boolean" },
        "tags": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    }
  }
}
`

// Ticket is the JSON schema for ticket frontmatter.
const Ticket = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "swimlane ticket frontmatter",
  "type": "object",
  "required": ["priority", "status", "ready"],
  "properties": {
    "title": { "type": "string" },
    "priority": { "type": "string", "enum": ["p0", "p1", "p2", "p3", "p4"] },
    "status": { "type": "string", "enum": ["todo", "in-progress", "done"] },
    "ready": { "type": "boolean" },
    "blocked_by": {
      "type": "array",
      "items": { "type": "string", "description": "ULID of blocking ticket" }
    },
    "tags": {
      "type": "array",
      "items": { "type": "string" }
    }
  }
}
`
