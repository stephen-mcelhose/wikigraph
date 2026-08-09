---
type: decision
title: "ADR-004: Use quantum-go as the canonical example wiki in how-to docs"
description: Use quantum-go as the canonical example wiki in all how-to guides to make verification reproducible and realistic.
status: accepted
timestamp: 2026-08-09T09:00:00Z
tags: [adr, documentation, examples]
---

# ADR-004: Use quantum-go as the canonical example wiki in how-to docs

## Status

Accepted

## Context

The how-to guides for `wikigraph` [[graph]], [[analyze]], [[goal]], and [[export]] need
a concrete wiki to use as the running example throughout steps and verification
sections. The alternative approaches were:

**`~/notes` (generic placeholder)** — The original choice. Pros: feels
universal. Cons: cannot be cloned and run; page counts and verification
output are invented, making it impossible to confirm examples are not stale.
The rubric requires verification sections to contain real, non-stale expected
values.

**`./docs` (this repo's own wiki)** — Always available. Cons: makes the guides
feel like documentation about wikigraph documenting itself, which could
confuse readers about whether the tool is generally applicable.

**A third-party wiki** — Demonstrates that wikigraph works on any Markdown
wiki, not just its own docs.

## Decision

Use `~/quantum-go/wiki` — the wiki directory inside
[github.com/stephen-mcelhose/quantum-go](https://github.com/stephen-mcelhose/quantum-go),
a real Go project with a well-linked knowledge base — as the canonical example
wiki in all four how-to guides.

All how-to guides include this prerequisite step:

```bash
git clone https://github.com/stephen-mcelhose/quantum-go ~/quantum-go
```

Verification sections show real output produced by running each subcommand
against `~/quantum-go/wiki` (36 pages, 295 edges as of this writing).

## Consequences

- **Verifiable**: anyone can clone quantum-go and reproduce the exact output
  shown in the docs.
- **Tool-agnostic signal**: the guides demonstrate wikigraph working on a real,
  independently-maintained Go project, not on its own documentation.
- **Maintenance**: if quantum-go's wiki grows significantly, verification
  numbers (Pages, Edges, Entropy rate) may drift. They should be re-checked
  when the how-to guides are updated.
- **Prerequisite cost**: readers must clone an additional repo before following
  the guide. This is offset by the benefit of having commands that actually run.
