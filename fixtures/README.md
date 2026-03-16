# Fixtures

Representative swimlane “boards” for testing and demos. Each subdirectory is a self-contained example (config + tickets). Run commands from inside a fixture directory, e.g.:

```bash
cd fixtures/simple && swimlane ls
cd fixtures/blocked && swimlane next
cd fixtures/with-tags && swimlane ls --tag backend
```

| Fixture       | Description |
|---------------|-------------|
| **empty**     | No tickets; fresh board. |
| **simple**    | A few tickets, mixed statuses (todo, in-progress, done), no dependencies. |
| **blocked**   | Tickets with `blocked_by`; tests “next” when work is blocked. |
| **with-tags** | Tickets with various tags; tests `--tag` / `--tag !value` filtering. |
