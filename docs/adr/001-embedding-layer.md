# ADR-001 — Embedding layer for semantic wiki search

**Date:** 2026-08  
**Status:** Superseded — decision revised; see *Revised Decision* below  
**Issue:** [#6](https://github.com/stephen-mcelhose/wikigraph/issues/6), [#7](https://github.com/stephen-mcelhose/wikigraph/issues/7)

---

## Revised Decision

**In the context of** adding semantic goal resolution to `wikigraph goal` (issue #6), where we need to embed wiki pages into a vector space so a natural-language query can be resolved to relevant pages,

**facing the concern** that the initially-preferred `vectorize` CLI lives in a private Bayer org repository (`bayer-int/csgdaa-code`) and cannot be treated as an open-source dependency, and that any embedding approach adds an external runtime dependency,

**we decided** to implement `wikigraph vectorize` as a proper subcommand using **[chromem-go](https://github.com/philippgille/chromem-go)** as the in-process vector store, with a pluggable embedding backend (Ollama as default, interface-driven so other backends can be added),

**to achieve** fully open-source, self-contained semantic search with no private dependencies and no mandatory running daemon,

**accepting** that users who want semantic goal resolution must have an embedding backend available (Ollama by default), and that the vector index is local to the machine.

---

## What we looked at

### Option A — `vectorize` CLI (csgdaa-code)

The `vectorize` binary (v0.2.21, Rust) ships with the csgdaa-code Homebrew tap. It does exactly what we need: chunk + embed Git repo files into a local SQLite vector DB, with semantic search via `vectorize local-search`.

**Rejected** because `bayer-int/csgdaa-code` is a **private Bayer org repo**. It cannot be declared as a dependency in an open-source project and cannot be assumed present in arbitrary developer environments.

| Capability | Notes |
| --- | --- |
| Chunk + embed from Git repo | ✅ `vectorize local-index` |
| Local SQLite vector store | ✅ `~/.local/share/csgdaa-code/vectorize/vectorize.db` |
| Semantic search | ✅ `vectorize local-search` |
| Open source / public | ❌ private repo |
| Installable without org access | ❌ requires Bayer Homebrew tap |

---

### Option B — chromem-go + Ollama

**[chromem-go](https://github.com/philippgille/chromem-go)** (MPL-2.0) is an embeddable in-process vector database for Go — "SQLite for vector search". It compiles directly into the binary with zero third-party runtime dependencies. Embedding generation is pluggable via an `EmbeddingFunc` interface; the library ships adapters for Ollama and LocalAI out of the box.

Ollama requires a local daemon (`ollama serve`), but it is widely installed and the `--host` flag lets users point at any compatible endpoint.

| Capability | Notes |
| --- | --- |
| In-process, no server | ✅ library embedded in binary |
| Open source | ✅ MPL-2.0 |
| Go-native | ✅ no CGo |
| Pluggable embedding backend | ✅ `EmbeddingFunc` interface |
| Local inference (no API key) | ✅ via Ollama or LocalAI |
| Embedding daemon required | ⚠️ Ollama must be running to embed (not to search) |

**Selected.**

---

### Option C — sqlite-lembed

**[sqlite-lembed](https://github.com/asg017/sqlite-lembed)** is a SQLite extension that runs llama.cpp inference in-process, generating embeddings via a SQL function against GGUF model files. No daemon at all — model runs inside the SQLite process.

**Not selected** because: (1) it is alpha (`v0.0.1-alpha`), with no license file committed; (2) it has no Go bindings — integration requires CGo + `go-sqlite3` + compiling the C extension; (3) GGUF model files must be separately downloaded and managed.

| Capability | Notes |
| --- | --- |
| No daemon at all | ✅ in-process via llama.cpp |
| Open source | ⚠️ no license file yet (sister projects MIT/Apache) |
| Go-native | ❌ CGo required |
| Maturity | ❌ alpha, API unstable |
| Model management | ❌ user downloads GGUF files manually |

---

## Consequences

- **`wikigraph vectorize .`** will use chromem-go to chunk and embed all wiki pages, persisting the index to `.vectors/wiki.db` (gob format, alongside the wiki).
- **`wikigraph goal --semantic`** will load the chromem-go index and perform cosine-similarity lookup to resolve a natural-language query to slugs, then feed those into the existing MFPT machinery.
- **Embedding backend** defaults to Ollama (`--host http://localhost:11434`, `--model nomic-embed-text`). The interface allows future backends (LocalAI, OpenAI-compatible, custom).
- **Searching does not require Ollama** — the stored vectors are used directly. Only `wikigraph vectorize` needs the daemon running.
- **Plain slug-based `goal`** remains fully self-contained with no new dependencies.
