---
title: Life cycle actions
priority: p2
status: todo
ready: true
blocked_by: ['01KKW6C11RGE9XY4P29Z0ASZFC']
tags: []
---

# Add lifecycle actions to the CLI

| action     | description                                       |
|------------|---------------------------------------------------|
| create     | creates a new ticket (replaces `new`)             |
| assign     | assign a ticket to a given user                   |
| claim      | assign a ticket to the current user               |
| unclaim    | unassign a ticket                                 |
| start      | move a ticket to the in-progress status           |
| stop       | move a ticket to the todo status                  |
| complete   | move a ticket to the done status                  |
| delete     | delete the ticket                                 |
| activate   | Set the value of `ready` to true                  |
| deactivate | Set the value of `ready` to false                 |

## Additional "Create" description

Creates a ticket, with the optionally given title.  Unless `--no-description` is given, it takes stdin as the contents of the ticket (with helpful message around "ctrl+d to end" or similar)

### Dogpiling actions

If additional actions do not require extra information (e.g. claim but not assign), then there should be a `--[action-name]` argument that can be added to also perform that action.  (If multiple actions that change the same value, e.g. claim and unclaim, an error shoudl be thrown and no action should be taken)

## Custom invocatons

### command

In the config file, there should be an optional section for "actions" as an object with optional actions (see list).  each one should be an object with a property `command` that is the CLI command that will be executed when that action is invoked.  arguments (e.g. title when creating a new ticket) should be represented by `{arg-name}` (e.g. {title}) and substituted into the command before execution.  Note, it should not add any wrapping syntax (like quotes), so it's up to the defined command to implement that as necessary (e.g. `some-command --title="{title}"`)

### Future

Create P3 tickets for these features, but mark them as ready=false:
* Support for multiple override types, including "HTTP", "Linear", "Jira".  Add an Integration Key section to the config that maps necessary API keys for these to env var names that hold the value
* A github action(s) to handle `claim` `unclaim` `start` `stop` actions that update the repo while bypassing the typical PR flow.