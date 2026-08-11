#!/usr/bin/env python3
"""
wiki-gen.py — Generate synthetic wiki structures for wikigraph experiments.

Produces a flat directory of .md files with wikilinks matching wikigraph's
slug rules: filename without .md, globally unique, matching [A-Za-z][A-Za-z0-9-]*.

Topology: N sides, each with a hub index highly connected to all pages on
that side. A bridge node connects all hubs. Each side has S sections
(level-2 indexes), each with L leaf pages (level-3). Leaf pages only link
up to their section; section indexes link up to their hub; hubs link to
bridge + all section indexes + all leaves on their side.

Usage:
    python3 tools/wiki-gen.py --out /tmp/wiki-200
    python3 tools/wiki-gen.py --sides 2 --sections 9 --leaves 10 --out /tmp/wiki-200
    python3 tools/wiki-gen.py --sides 3 --sections 5 --leaves 8 --out /tmp/wiki-custom
"""

import argparse
import os
import sys


def slugify(label: str) -> str:
    """Return a wikilink-safe slug (letters, digits, hyphens; starts with letter)."""
    return label.lower().replace("_", "-").replace(" ", "-")


def write_page(out_dir: str, slug: str, title: str, body: str) -> None:
    path = os.path.join(out_dir, f"{slug}.md")
    with open(path, "w") as f:
        f.write(f"# {title}\n\n{body}\n")


def link(slug: str) -> str:
    return f"[[{slug}]]"


def links(slugs: list[str]) -> str:
    return " ".join(link(s) for s in slugs)


def generate(out_dir: str, sides: int, sections: int, leaves: int) -> int:
    os.makedirs(out_dir, exist_ok=True)

    side_labels = [chr(ord("a") + i) for i in range(sides)]
    hub_slugs = [f"hub-{s}" for s in side_labels]

    # Compute all slugs up front so we can detect clashes early.
    all_slugs: set[str] = set()

    def register(slug: str) -> str:
        if slug in all_slugs:
            sys.exit(f"ERROR: duplicate slug {slug!r} — adjust parameters")
        all_slugs.add(slug)
        return slug

    register("bridge")
    for h in hub_slugs:
        register(h)

    section_slugs: dict[str, list[str]] = {}
    leaf_slugs: dict[str, list[str]] = {}

    for s in side_labels:
        section_slugs[s] = []
        for sec_n in range(1, sections + 1):
            sec_slug = register(f"{s}-s{sec_n:02d}")
            section_slugs[s].append(sec_slug)
            leaf_slugs[sec_slug] = []
            for leaf_n in range(1, leaves + 1):
                leaf_slug = register(f"{s}-s{sec_n:02d}-p{leaf_n:02d}")
                leaf_slugs[sec_slug].append(leaf_slug)

    # --- bridge ---
    write_page(
        out_dir,
        "bridge",
        "Bridge",
        f"Central bridge connecting all community hubs.\n\n{links(hub_slugs)}",
    )

    # --- hubs ---
    for s in side_labels:
        hub_slug = f"hub-{s}"
        all_sec = section_slugs[s]
        all_leaves = [leaf for sec in all_sec for leaf in leaf_slugs[sec]]
        out_links = ["bridge"] + all_sec + all_leaves
        write_page(
            out_dir,
            hub_slug,
            f"Hub {s.upper()}",
            f"Primary hub for community {s.upper()}. "
            f"Links to bridge, all section indexes, and all pages on this side.\n\n"
            f"{links(out_links)}",
        )

    # --- section indexes (level 2) ---
    for s in side_labels:
        hub_slug = f"hub-{s}"
        for sec_slug in section_slugs[s]:
            sec_leaves = leaf_slugs[sec_slug]
            out_links = [hub_slug] + sec_leaves
            write_page(
                out_dir,
                sec_slug,
                f"Section {sec_slug}",
                f"Level-2 section index.\n\n{links(out_links)}",
            )

    # --- leaf pages (level 3) ---
    for s in side_labels:
        for sec_slug in section_slugs[s]:
            for leaf_slug in leaf_slugs[sec_slug]:
                write_page(
                    out_dir,
                    leaf_slug,
                    f"Page {leaf_slug}",
                    f"Leaf page.\n\n{link(sec_slug)}",
                )

    total = len(all_slugs)
    print(
        f"Generated {total} pages in {out_dir}\n"
        f"  bridge:   1\n"
        f"  hubs:     {sides}\n"
        f"  sections: {sides} × {sections} = {sides * sections}\n"
        f"  leaves:   {sides} × {sections} × {leaves} = {sides * sections * leaves}"
    )
    return total


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Generate a synthetic wiki for wikigraph experiments."
    )
    parser.add_argument("--out", required=True, help="Output directory for .md files")
    parser.add_argument("--sides", type=int, default=2, help="Number of communities (default 2)")
    parser.add_argument("--sections", type=int, default=9, help="Sections per hub (default 9)")
    parser.add_argument("--leaves", type=int, default=10, help="Leaf pages per section (default 10)")
    args = parser.parse_args()

    if args.sides < 1:
        sys.exit("--sides must be >= 1")
    if args.sides > 26:
        sys.exit("--sides max 26 (uses a-z labels)")
    if args.sections < 1:
        sys.exit("--sections must be >= 1")
    if args.leaves < 1:
        sys.exit("--leaves must be >= 1")

    generate(args.out, args.sides, args.sections, args.leaves)


if __name__ == "__main__":
    main()
