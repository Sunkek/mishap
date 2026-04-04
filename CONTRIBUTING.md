# Contributing to mishap

Thanks for your interest. mishap is a small, focused package and contributions
should stay in that spirit — precise, well-tested, and zero new dependencies.

## Before you open a PR

1. **Open an issue first** for anything beyond a trivial fix. A brief discussion
   saves everyone time if the direction isn't right.
2. **Check the existing tests** — the test suite is the specification. If a
   scenario isn't tested, that's a gap worth filling before changing behaviour.

## Development setup

```
git clone https://github.com/sunkek/mishap
cd mishap
make test   # race-enabled tests
make lint   # go vet
```

`staticcheck` is optional but appreciated:

```
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## What good contributions look like

**Fixes** — a failing test that demonstrates the bug, then the minimal fix.
Avoid changing unrelated code in the same PR.

**Features** — a clear rationale (what production problem does this solve?),
a test that would fail without the change, and updated documentation if the
public API changes.

**Documentation** — plain language, concrete examples. The README's notes on
`Code` implementing `error` and the `Is` matching semantics are the most
important things to keep accurate.

## API stability

mishap follows semantic versioning. The public API is stable at v1 — breaking
changes require a major version bump and a clear entry in the CHANGELOG.

## Commit messages

Short imperative subject line, 72 characters max. Reference an issue number
when relevant: `fix wrap panic on nil option (#12)`.
