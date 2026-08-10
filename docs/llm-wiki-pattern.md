---
type: concept
title: LLM Wiki Pattern
description: Compounding knowledge base architecture where the LLM acts as the programmer and a markdown wiki acts as the codebase.
resource: https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f
tags: [llm-wiki, concept, architecture, knowledge-base]
timestamp: 2026-08-09T17:20:00Z
---

# LLM Wiki Pattern

The **LLM Wiki Pattern** (inspired by Andrej Karpathy's design) models knowledge management analogously to software development:
> *"Obsidian is the IDE; the LLM is the programmer; the wiki is the codebase."*

Rather than relying purely on ad-hoc RAG or unindexed conversational memory, an LLM maintains a structured, compounding directory of interlinked Markdown files.

## Three-Layer Architecture

1. **Raw Sources (Immutable)**: Raw inputs (articles, transcripts, PDFs, notes) living in `raw/` or fetched via URLs. The LLM reads these sources but never writes to them.
2. **Wiki (Synthesized Codebase)**: Flat directory of interlinked `.md` files using OKF frontmatter and `[[slug]]` wikilinks. Maintained and refactored by the LLM.
3. **Schema (`AGENTS.md`)**: Defines the domain, page-type conventions, and operations (`ingest`, `query`, `lint`).

## Operations

- **Ingest**: Read raw input, write/update synthesis pages, and aggressively propagate wikilinks across related pages.
- **Query**: Read index/pages, follow `[[wikilinks]]`, and produce cited answers with optional write-back for novel insights.
- **Lint**: Run structural health audits (orphan detection, contradiction resolution, stale claim updates) using tools like [[analyze]].

## Integration with Wikigraph

`wikigraph` models this wiki structure as a Markov chain random walk. Running [[analyze]] or [[graph]] on an LLM-maintained wiki reveals topological health, hub centrality, and missing link suggestions that feed back into maintenance prompts.

## Sources

- [Karpathy's LLM Wiki Gist](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f)
- [[quickstart]]
- [[AGENTS]]
