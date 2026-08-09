---
type: decision
title: Embedding Layer for Semantic Wiki Search
description: Chose chromem-go + Ollama for in-process semantic goal resolution, rejecting a private Bayer dependency and an alpha CGo library.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/6
tags: [adr, embedding, semantic-search, chromem-go, ollama]
timestamp: 2026-08-09T06:54:46Z
status: superseded
---

# ADR-001 — Embedding layer for semantic wiki search

## Context

When adding semantic goal resolution to `wikigraph goal` (issue #6), we needed to embed wiki pages into a vector space so a natural-language query can be resolved to relevant pages. We faced concerns that the initially preferred `vectorize` CLI lives in a private repository (`bayer-int/csgdaa-code`) and cannot be treated as an open-source dependency, and that any embedding approach adds an external runtime dependency.

Option A (`vectorize` CLI in private repo) was rejected due to org access requirements. Option C (`sqlite-lembed`) was rejected because it requires CGo and manual GGUF model management. Option B (`chromem-go` + Ollama) was selected as an in-process, open-source Go vector store.

## Decision

In the context of adding semantic goal resolution to `wikigraph goal`, facing the need for open-source self-contained embeddings without private dependencies or mandatory server daemons, we decided to use chromem-go as the in-process vector store with a pluggable embedding backend (Ollama default), to achieve zero-CGo open-source semantic search, accepting that users generating embeddings need an available embedding backend.

## Consequences

- **`wikigraph vectorize .`** uses chromem-go to chunk and embed all wiki pages, persisting the index to `.vectors/wiki.db`.
- **`wikigraph goal --semantic`** loads the chromem-go index and performs cosine-similarity lookup to resolve natural-language queries (see [[goal]]).
- **Embedding backend** defaults to Ollama (`--host http://localhost:11434`, `--model nomic-embed-text`).
- **Searching does not require Ollama** — stored vectors are used directly.
- **Plain slug-based `goal`** remains fully self-contained.
- **Slug naming convention** for all wiki pages is documented in [[adr-002-slug-resolution]].

## Sources

- [GitHub issue #6](https://github.com/stephen-mcelhose/wikigraph/issues/6)
- [chromem-go](https://github.com/philippgille/chromem-go)
- [sqlite-lembed](https://github.com/asg017/sqlite-lembed)
