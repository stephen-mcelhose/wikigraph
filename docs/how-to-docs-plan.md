---
type: proposal
title: Create How-To Documentation for Subcommands
description: Plan to implement targeted how-to documentation guides for the four wikigraph subcommands.
resource: https://github.com/stephen-mcelhose/wikigraph/issues/3
tags: [documentation, cli, help]
timestamp: 2026-08-09T06:39:32Z
---

## Problem

Users of `wikigraph` lack practical, hands-on tutorials ("how-to" guides) for the CLI's four subcommands (`graph`, `analyze`, `goal`, `export`), which prevents them from fully leveraging the Markov chain analysis and health reporting capabilities on their markdown wikis.

## Proposed Solution

Create detailed, step-by-step how-to guides for each of the four subcommands. Each guide will provide clear examples, expected outputs, and practical use cases.

### Scope of Guides

1. **`graph` How-To**:
   - Focus: Generating and customising interactive HTML force-directed graphs.
   - Key steps: pointing the tool to a wiki, configuring node/link parameters, viewing the HTML output.

2. **`analyze` How-To**:
   - Focus: High-quality wiki health analysis.
   - Key steps: locating orphans/sinks, identifying central pages, and utilizing link suggestions to improve wiki structure.
   - **Note**: This will be the largest and most complex how-to doc given the analytical depth of the subcommand.

3. **`goal` How-To**:
   - Focus: Interactive path-finding/learning paths toward a target goal page using Mean First Passage Time (MFPT).
   - Key steps: specifying goal nodes, analyzing paths, outputting HTML subgraphs.
   - **Note**: Since MFPT pathfinding is a more aspirational/inspirational use case, this guide may start as a draft.

4. **`export` How-To**:
   - Focus: Exporting graph structure data.
   - Key steps: converting markdown wiki links to programmatic outputs (JSON, CSV, or DOT format) for external processing.

## Alternatives Considered

- **Inline CLI help only**: Keeping all instructions in `--help`. *Rejected* because command line usage output lacks context, screenshots/examples, and explanatory narrative.
- **Single monolithic README**: Keeping all instructions in a single file. *Rejected* because subcommand instructions will be highly detailed and deserve modular, searchable files.

## Open Questions

- Should we host these guides in a `docs/how-to/` directory within the repository, or publish them as a wiki/external documentation site?
- For the interactive elements in `graph` and `goal`, should we include pre-rendered screenshots or interactive preview links in the guides?

## Implementation Notes

- Use the OKF `how-to` format for each individual guide once they are written.
- Start with the simpler ones (`export`, `graph`), tackle the complex `analyze` guide, and draft `goal`.

## Completed Guides

All four guides were delivered on 2026-08-09:

- [[export]]
- [[graph]]
- [[analyze]]
- [[goal]] *(draft)*
