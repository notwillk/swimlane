---
title: Add "subtasks" field to tickets
priority: p2
status: todo
ready: true
blocked_by:
  - 01KKW52FZYXA77KTF1JFZSJJY0
tags: []
---

# Add "subtasks" field to tickets.

This should be an array of identifiers (e.g. ULIDs) that correspond to tickets that further break down a ticket into additional tasks.  The concept is that the work defined in the sub tasks entirely accounts for all the work necessary for the parent task (e.g. "make cookies" will have sub tasks "combine ingredients" and "bake").

There should be a new boolean property in the config for if when all sub-tasks are closed then the corresponding parent should also be closed.
