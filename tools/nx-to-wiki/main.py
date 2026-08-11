#!/usr/bin/env python3
"""
nx-to-wiki — Convert a NetworkX named graph into a flat directory of .md files
with [[slug]] wikilinks, for wikigraph experiments.

See docs/adr-009-wiki-gen-make-vs-buy.md for the make-vs-buy rationale and
https://github.com/stephen-mcelhose/wikigraph/issues/32 for the naming spec.

Usage:
    python3 tools/nx-to-wiki/main.py --graph karate --out /tmp/karate-wiki
    python3 tools/nx-to-wiki/main.py --graph florentine --out /tmp/florentine-wiki
    python3 tools/nx-to-wiki/main.py --graph barbell --m1 5 --m2 3 --out /tmp/barbell-wiki
    python3 tools/nx-to-wiki/main.py --graph caveman --l 3 --k 4 --p 0.2 --out /tmp/caveman-wiki
    python3 tools/nx-to-wiki/main.py --graph planted-partition --l 3 --k 10 \\
        --p-in 0.5 --p-out 0.05 --out /tmp/pp-wiki
    python3 tools/nx-to-wiki/main.py --graph watts-strogatz --n 20 --k 4 --p 0.1 \\
        --role-names --out /tmp/ws-wiki
    python3 tools/nx-to-wiki/main.py --graph lfr --n 50 --tau1 2.5 --tau2 1.5 \\
        --mu 0.3 --out /tmp/lfr-wiki

Requires: networkx (pip install networkx)
"""

import argparse
import os
import statistics
import sys

try:
    import networkx as nx
except ImportError:
    sys.exit("ERROR: networkx is required. Install with: pip install networkx")


GRAPH_CHOICES = [
    "karate",
    "krackhardt-kite",
    "florentine",
    "barbell",
    "caveman",
    "planted-partition",
    "watts-strogatz",
    "lfr",
]


# ---------------------------------------------------------------------------
# Slug helpers
# ---------------------------------------------------------------------------

def slugify(label: str) -> str:
    """Return a wikilink-safe slug (letters, digits, hyphens; starts with letter)."""
    s = str(label).strip().lower().replace("_", "-").replace(" ", "-")
    s = "".join(ch for ch in s if ch.isalnum() or ch == "-")
    if not s or not s[0].isalpha():
        s = "n-" + s
    return s


def pad(i: int, count: int) -> str:
    """Zero-pad i (0-indexed) to fit count items, minimum width 2."""
    width = max(2, len(str(max(count - 1, 0))))
    return str(i).zfill(width)


# ---------------------------------------------------------------------------
# Tier 1 — graph builders that also return a node -> slug map
# ---------------------------------------------------------------------------

def build_karate():
    G = nx.karate_club_graph()
    groups: dict[str, list] = {}
    for n, data in G.nodes(data=True):
        club = data.get("club", "")
        prefix = "mhi" if club == "Mr. Hi" else "off"
        groups.setdefault(prefix, []).append(n)
    slugs = {}
    for prefix, nodes in groups.items():
        for i, n in enumerate(sorted(nodes)):
            slugs[n] = f"{prefix}-{pad(i, len(nodes))}"
    return G, slugs, "tier1-karate-club-label"


def build_krackhardt_kite():
    G = nx.krackhardt_kite_graph()
    # No built-in labels; caller falls through to Tier 2/3.
    return G, {}, None


def build_florentine():
    G = nx.florentine_families_graph()
    slugs = {n: slugify(n) for n in G.nodes()}
    return G, slugs, "tier1-florentine-family-name"


def build_barbell(m1: int, m2: int):
    G = nx.barbell_graph(m1, m2)
    # Structural position is known from the generator's own node numbering:
    #   [0, m1)             -> left clique
    #   [m1, m1 + m2)       -> bridging path
    #   [m1 + m2, 2*m1+m2)  -> right clique
    slugs = {}
    for n in G.nodes():
        if n < m1:
            slugs[n] = f"left-{pad(n, m1)}"
        elif n < m1 + m2:
            slugs[n] = f"bridge-{pad(n - m1, m2)}"
        else:
            idx = n - (m1 + m2)
            slugs[n] = f"right-{pad(idx, m1)}"
    return G, slugs, "tier1-barbell-structural-position"


def build_caveman(l: int, k: int, p: float, seed):
    G = nx.relaxed_caveman_graph(l, k, p, seed=seed)
    # Clique membership is known from the generator: clique i occupies
    # nodes [i*k, (i+1)*k).
    slugs = {}
    for n in G.nodes():
        clique_idx = n // k
        within = n % k
        slugs[n] = f"cave{clique_idx}-{pad(within, k)}"
    return G, slugs, "tier1-caveman-clique-membership"


def build_planted_partition(l: int, k: int, p_in: float, p_out: float, seed):
    G = nx.planted_partition_graph(l, k, p_in, p_out, seed=seed)
    slugs = _slugs_from_community_attr(G, attr="block")
    return G, slugs, "tier1-planted-partition-block-attr"


def build_watts_strogatz(n: int, k: int, p: float, seed):
    G = nx.watts_strogatz_graph(n, k, p, seed=seed)
    # Anonymous graph: no built-in labels.
    return G, {}, None


def build_lfr(n: int, tau1: float, tau2: float, mu: float, seed, extra: dict):
    kwargs = dict(seed=seed)
    kwargs.update({k: v for k, v in extra.items() if v is not None})
    G = nx.LFR_benchmark_graph(n, tau1, tau2, mu, **kwargs)
    slugs = _slugs_from_community_attr(G, attr="community")
    return G, slugs, "tier1-lfr-community-attr"


def _slugs_from_community_attr(G, attr: str):
    """Map each node to a c{community}-{n} slug from a community/block attribute.

    Handles both scalar community ids (e.g. `block`, an int) and set-valued
    community assignments (e.g. `community`, a set of co-members produced by
    LFR_benchmark_graph).
    """
    raw = {n: G.nodes[n].get(attr) for n in G.nodes()}
    # Canonicalize community identity so equal sets/ints map to the same id.
    canon_keys = {}
    for n, v in raw.items():
        key = frozenset(v) if isinstance(v, (set, frozenset)) else v
        canon_keys[n] = key
    unique = sorted(set(canon_keys.values()), key=lambda k: sorted(k) if isinstance(k, frozenset) else (k,))
    community_id = {key: i for i, key in enumerate(unique)}

    groups: dict[int, list] = {}
    for n in G.nodes():
        cid = community_id[canon_keys[n]]
        groups.setdefault(cid, []).append(n)

    slugs = {}
    for cid, nodes in groups.items():
        for i, n in enumerate(sorted(nodes)):
            slugs[n] = f"c{cid}-{pad(i, len(nodes))}"
    return slugs


# ---------------------------------------------------------------------------
# Tier 2 — structural role walk for anonymous graphs
# ---------------------------------------------------------------------------

def percentile(values: list[float], pct: float) -> float:
    if not values:
        return 0.0
    if len(values) == 1:
        return values[0]
    qs = statistics.quantiles(values, n=100, method="inclusive")
    idx = min(max(int(round(pct)) - 1, 0), len(qs) - 1)
    return qs[idx]


def assign_role_slugs(G) -> dict:
    """Classify each node into hub/connector/leaf/node via degree, betweenness,
    and clustering coefficient, then assign {role}-{rank} slugs."""
    degree = dict(G.degree())
    betweenness = nx.betweenness_centrality(G)
    # Clustering is computed for completeness / future thresholds; not used
    # in the classification cutoffs below beyond being available per spec.
    nx.clustering(G)

    deg_p75 = percentile(list(degree.values()), 75)
    bet_p75 = percentile(list(betweenness.values()), 75)

    roles: dict = {}
    for n in G.nodes():
        d, b = degree[n], betweenness[n]
        if d == 1:
            roles[n] = "leaf"
        elif d >= deg_p75 and b >= bet_p75:
            roles[n] = "hub"
        elif b >= bet_p75:
            roles[n] = "connector"
        else:
            roles[n] = "node"

    groups: dict[str, list] = {}
    for n, role in roles.items():
        groups.setdefault(role, []).append(n)

    slugs = {}
    for role, nodes in groups.items():
        for i, n in enumerate(sorted(nodes)):
            slugs[n] = f"{role}-{pad(i, len(nodes))}"
    return slugs


# ---------------------------------------------------------------------------
# Tier 3 — fallback
# ---------------------------------------------------------------------------

def fallback_slugs(G) -> dict:
    nodes = list(G.nodes())
    return {n: f"page-{pad(i, len(nodes))}" for i, n in enumerate(sorted(nodes))}


# ---------------------------------------------------------------------------
# Writer
# ---------------------------------------------------------------------------

def write_wiki(D, slugs: dict, out_dir: str) -> None:
    os.makedirs(out_dir, exist_ok=True)
    seen = set()
    for n in D.nodes():
        slug = slugs[n]
        if slug in seen:
            sys.exit(f"ERROR: duplicate slug {slug!r} for node {n!r} — naming collision")
        seen.add(slug)

    for n in D.nodes():
        slug = slugs[n]
        out_links = sorted({slugs[v] for v in D.successors(n) if v != n})
        body = " ".join(f"[[{s}]]" for s in out_links)
        content = f"# {slug}\n\n{body}\n"
        with open(os.path.join(out_dir, f"{slug}.md"), "w") as f:
            f.write(content)


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(
        description="Convert a NetworkX named graph into a wikigraph-compatible wiki."
    )
    parser.add_argument("--graph", required=True, choices=GRAPH_CHOICES)
    parser.add_argument("--out", required=True, help="Output directory for .md files")
    parser.add_argument(
        "--role-names",
        action="store_true",
        help="Opt-in Tier 2 structural role walk for graphs with no built-in labels",
    )
    parser.add_argument(
        "--directed",
        action="store_true",
        help="No-op: output is always symmetric-directed via to_directed() (default)",
    )
    parser.add_argument("--seed", type=int, default=None)

    # barbell
    parser.add_argument("--m1", type=int, default=5)
    parser.add_argument("--m2", type=int, default=3)
    # caveman / planted-partition (shared l, k, p)
    parser.add_argument("--l", type=int, default=3, help="Number of cliques/blocks")
    parser.add_argument("--k", type=int, default=4, help="Nodes per clique/block, or WS ring degree")
    parser.add_argument("--p", type=float, default=0.1, help="Rewiring/relaxation probability")
    parser.add_argument("--p-in", type=float, default=0.5, dest="p_in")
    parser.add_argument("--p-out", type=float, default=0.05, dest="p_out")
    # watts-strogatz / lfr
    parser.add_argument("--n", type=int, default=50)
    parser.add_argument("--tau1", type=float, default=2.5)
    parser.add_argument("--tau2", type=float, default=1.5)
    parser.add_argument("--mu", type=float, default=0.3)
    parser.add_argument("--average-degree", type=float, default=None, dest="average_degree")
    parser.add_argument("--min-degree", type=int, default=None, dest="min_degree")
    parser.add_argument("--max-degree", type=int, default=None, dest="max_degree")
    parser.add_argument("--min-community", type=int, default=None, dest="min_community")
    parser.add_argument("--max-community", type=int, default=None, dest="max_community")

    args = parser.parse_args()

    if args.graph == "karate":
        G, slugs, tier = build_karate()
    elif args.graph == "krackhardt-kite":
        G, slugs, tier = build_krackhardt_kite()
    elif args.graph == "florentine":
        G, slugs, tier = build_florentine()
    elif args.graph == "barbell":
        G, slugs, tier = build_barbell(args.m1, args.m2)
    elif args.graph == "caveman":
        G, slugs, tier = build_caveman(args.l, args.k, args.p, args.seed)
    elif args.graph == "planted-partition":
        G, slugs, tier = build_planted_partition(args.l, args.k, args.p_in, args.p_out, args.seed)
    elif args.graph == "watts-strogatz":
        G, slugs, tier = build_watts_strogatz(args.n, args.k, args.p, args.seed)
    elif args.graph == "lfr":
        extra = dict(
            average_degree=args.average_degree,
            min_degree=args.min_degree,
            max_degree=args.max_degree,
            min_community=args.min_community,
            max_community=args.max_community,
        )
        if extra["average_degree"] is None and extra["min_degree"] is None:
            extra["average_degree"] = 5.0
        if extra["min_community"] is None:
            extra["min_community"] = max(10, args.n // 10)
        G, slugs, tier = build_lfr(args.n, args.tau1, args.tau2, args.mu, args.seed, extra)
    else:
        sys.exit(f"unknown graph {args.graph!r}")

    if not slugs:
        if args.role_names:
            slugs = assign_role_slugs(G)
            tier = "tier2-structural-role-walk"
        else:
            slugs = fallback_slugs(G)
            tier = "tier3-fallback-page-n"
            print(
                f"WARNING: {args.graph!r} has no built-in labels; falling back to "
                f"page-N slugs. Pass --role-names for structural role slugs.",
                file=sys.stderr,
            )

    D = G.to_directed()
    write_wiki(D, slugs, args.out)

    print(f"Generated {D.number_of_nodes()} pages, {D.number_of_edges()} directed edges in {args.out}")
    print(f"  graph:  {args.graph}")
    print(f"  naming: {tier}")


if __name__ == "__main__":
    main()
