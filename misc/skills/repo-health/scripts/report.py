#!/usr/bin/env python3
"""Render metrics.json (+ optional reading.md) into report.html and report.md.

Usage: report.py <run-dir>
"""
import html, json, os, re, sys

# Reference categorical palette (dataviz skill), fixed order. Light / dark.
SERIES = [('#2a78d6', '#3987e5'), ('#eb6834', '#d95926'), ('#1baf7a', '#199e70'), ('#eda100', '#c98500'),
          ('#e87ba4', '#d55181'), ('#008300', '#008300'), ('#6250d6', '#9085e9'), ('#e34948', '#e66767')]
INTENTS = ['feature', 'fix', 'refactor', 'perf', 'test', 'docs', 'build', 'chore']
W = 760

CSS = '''
.viz-root{color-scheme:light;--surface-1:#fcfcfb;--surface-2:#f2f1ee;--grid:#e4e3df;--text-primary:#0b0b0b;
--text-secondary:#52514e;--text-muted:#7a7975;--other:#a3a29c;%s}
@media (prefers-color-scheme:dark){:root:where(:not([data-theme="light"])) .viz-root{color-scheme:dark;
--surface-1:#1a1a19;--surface-2:#242423;--grid:#33332f;--text-primary:#f0efec;--text-secondary:#c3c2b7;
--text-muted:#8f8e86;--other:#6b6a64;%s}}
body{margin:0;background:var(--surface-1);color:var(--text-primary);font:14px/1.5 system-ui,-apple-system,sans-serif}
main{max-width:980px;margin:0 auto;padding:24px}
h1{font-size:24px;margin:0 0 4px}h2{font-size:18px;margin:40px 0 4px;border-top:1px solid var(--grid);padding-top:24px}
h3{font-size:15px;margin:20px 0 4px}.sub{color:var(--text-secondary);margin:0 0 12px}
.tiles{display:grid;grid-template-columns:repeat(auto-fill,minmax(150px,1fr));gap:12px;margin:16px 0}
.tile{background:var(--surface-2);border-radius:8px;padding:12px}.tile .l{color:var(--text-secondary);font-size:12px}
.tile .v{font-size:24px;font-weight:600}
table{border-collapse:collapse;width:100%%;font-size:13px;font-variant-numeric:tabular-nums;margin:8px 0}
th,td{text-align:left;padding:4px 8px;border-bottom:1px solid var(--grid)}th{color:var(--text-secondary);font-weight:500}
td.n,th.n{text-align:right}details{margin:8px 0}summary{cursor:pointer;color:var(--text-secondary)}
.legend{display:flex;flex-wrap:wrap;gap:12px;font-size:12px;color:var(--text-secondary);margin:4px 0}
.legend i{display:inline-block;width:10px;height:10px;border-radius:2px;margin-right:4px;vertical-align:-1px}
svg text{fill:var(--text-secondary);font-size:11px}svg .val{fill:var(--text-primary)}
svg .grid{stroke:var(--grid);stroke-width:1}svg .mark:hover{opacity:.8}
.reading{background:var(--surface-2);border-radius:8px;padding:4px 20px}code{font-size:12px}
.flag{color:var(--text-primary);font-weight:600}
'''


def css():
    light = ''.join(f'--s{i + 1}:{l};' for i, (l, _) in enumerate(SERIES))
    dark = ''.join(f'--s{i + 1}:{d};' for i, (_, d) in enumerate(SERIES))
    return CSS % (light, dark)


def esc(s):
    return html.escape(str(s))


def fmt(v):
    if v is None:
        return '–'
    if isinstance(v, float):
        return f'{v:,.1f}'
    if isinstance(v, int) and abs(v) >= 10000:
        return f'{v / 1000:,.1f}K'
    return f'{v:,}' if isinstance(v, int) else str(v)


def nice_max(v):
    if v <= 0:
        return 1
    import math
    e = 10 ** math.floor(math.log10(v))
    for m in (1, 2, 2.5, 5, 10):
        if v <= m * e:
            return m * e
    return 10 * e


def ticks(vmax):
    # prefer steps like 5/10/20 over 12.5/18.75
    n = 5 if (vmax / 5) == int(vmax / 5) and vmax >= 5 else 4
    return [vmax * i / n for i in range(n + 1)]


def tick_label(v, unit=''):
    if v >= 1000:
        return f'{v / 1000:g}K{unit}'
    return f'{v:g}{unit}'


def bar_path(x, y, w, h, r=4, horizontal=False):
    """Bar with rounded data end, square at the baseline."""
    r = min(r, (w if not horizontal else h) / 2, (h if not horizontal else w))
    if r <= 0 or w <= 0 or h <= 0:
        return f'M{x},{y}h{max(w, 0)}v{max(h, 0)}h{-max(w, 0)}z'
    if horizontal:  # rounded right end
        return (f'M{x},{y}h{w - r}a{r},{r} 0 0 1 {r},{r}v{h - 2 * r}a{r},{r} 0 0 1 {-r},{r}h{-(w - r)}z')
    return (f'M{x},{y + h}v{-(h - r)}a{r},{r} 0 0 1 {r},{-r}h{w - 2 * r}a{r},{r} 0 0 1 {r},{r}v{h - r}z')


def legend(names, colors):
    return '<div class="legend">' + ''.join(
        f'<span><i style="background:{c}"></i>{esc(n)}</span>' for n, c in zip(names, colors)) + '</div>'


def columns(cats, series, height=240, unit='', title=''):
    """Stacked (or single) columns. series: [(name, color, values)]."""
    ml, mb, mt = 48, 28, 16
    pw, ph = W - ml - 8, height - mb - mt
    totals = [sum(s[2][i] or 0 for s in series) for i in range(len(cats))]
    vmax = nice_max(max(totals, default=0))
    band = pw / max(len(cats), 1)
    bw = min(24, band * 0.6)
    out = [f'<svg viewBox="0 0 {W} {height}" width="100%" role="img" aria-label="{esc(title)}">']
    for t in ticks(vmax):
        y = mt + ph - ph * t / vmax
        out.append(f'<line class="grid" x1="{ml}" x2="{W - 8}" y1="{y:.1f}" y2="{y:.1f}"/>'
                   f'<text x="{ml - 6}" y="{y + 4:.1f}" text-anchor="end">{tick_label(t, unit)}</text>')
    for i, c in enumerate(cats):
        x = ml + band * i + (band - bw) / 2
        base = mt + ph
        visible = [(n, col, v[i]) for n, col, v in series if v[i]]
        for k, (n, col, v) in enumerate(visible):
            h = ph * v / vmax
            top = k == len(visible) - 1
            gap = 2 if k > 0 else 0
            hh = max(h - gap, 0.5)
            y = base - h
            tip = f'<title>{esc(c)} · {esc(n)}: {fmt(v)}{unit}</title>'
            if top:
                out.append(f'<path class="mark" d="{bar_path(x, y, bw, hh)}" fill="{col}">{tip}</path>')
            else:
                out.append(f'<rect class="mark" x="{x:.1f}" y="{y:.1f}" width="{bw:.1f}" height="{hh:.1f}" fill="{col}">{tip}</rect>')
            base = y
        out.append(f'<text x="{x + bw / 2:.1f}" y="{height - 8}" text-anchor="middle">{esc(c)}</text>')
    out.append('</svg>')
    leg = legend([s[0] for s in series], [s[1] for s in series]) if len(series) > 1 else ''
    return leg + ''.join(out)


def hbars(rows, unit='', height_per=22, label_w=300, color='var(--s1)', title=''):
    """rows: [(label, value, tooltip)]"""
    h = height_per * len(rows) + 8
    vmax = nice_max(max((r[1] for r in rows), default=0))
    pw = W - label_w - 60
    out = [f'<svg viewBox="0 0 {W} {h}" width="100%" role="img" aria-label="{esc(title)}">']
    for i, (lab, v, tip) in enumerate(rows):
        y = 4 + i * height_per
        bh = min(14, height_per - 6)
        w = pw * v / vmax
        short = lab if len(lab) <= 46 else '…' + lab[-45:]
        out.append(f'<text x="{label_w - 8}" y="{y + bh - 3}" text-anchor="end">{esc(short)}</text>')
        out.append(f'<path class="mark" d="{bar_path(label_w, y, max(w, 1), bh, horizontal=True)}" fill="{color}">'
                   f'<title>{esc(tip)}</title></path>')
        out.append(f'<text class="val" x="{label_w + w + 6:.1f}" y="{y + bh - 3}">{fmt(v)}{unit}</text>')
    out.append('</svg>')
    return ''.join(out)


def lines(cats, series, height=220, unit='', title=''):
    """series: [(name, color, values)] — values may contain None."""
    ml, mb, mt, mr = 48, 28, 16, 60
    pw, ph = W - ml - mr, height - mb - mt
    vmax = nice_max(max((v for s in series for v in s[2] if v is not None), default=0))
    step = pw / max(len(cats) - 1, 1)
    out = [f'<svg viewBox="0 0 {W} {height}" width="100%" role="img" aria-label="{esc(title)}">']
    for t in ticks(vmax):
        y = mt + ph - ph * t / vmax
        out.append(f'<line class="grid" x1="{ml}" x2="{W - mr}" y1="{y:.1f}" y2="{y:.1f}"/>'
                   f'<text x="{ml - 6}" y="{y + 4:.1f}" text-anchor="end">{tick_label(t, unit)}</text>')
    for i, c in enumerate(cats):
        out.append(f'<text x="{ml + step * i:.1f}" y="{height - 8}" text-anchor="middle">{esc(c)}</text>')
    for name, col, vals in series:
        pts = [(ml + step * i, mt + ph - ph * v / vmax, v, cats[i]) for i, v in enumerate(vals) if v is not None]
        if not pts:
            continue
        out.append(f'<polyline fill="none" stroke="{col}" stroke-width="2" stroke-linejoin="round" '
                   f'stroke-linecap="round" points="{" ".join(f"{x:.1f},{y:.1f}" for x, y, _, _ in pts)}"/>')
        for x, y, v, c in pts:
            out.append(f'<circle class="mark" cx="{x:.1f}" cy="{y:.1f}" r="4" fill="{col}" stroke="var(--surface-1)" '
                       f'stroke-width="2"><title>{esc(c)} · {esc(name)}: {fmt(v)}{unit}</title></circle>')
        x, y, v, _ = pts[-1]
        out.append(f'<text class="val" x="{x + 8:.1f}" y="{y + 4:.1f}">{fmt(v)}{unit}</text>')
    out.append('</svg>')
    leg = legend([s[0] for s in series], [s[1] for s in series]) if len(series) > 1 else ''
    return leg + ''.join(out)


def scatter(points, xlab, ylab, height=300, label_top=6, title=''):
    """points: [(x, y, label, tooltip)]"""
    ml, mb, mt, mr = 56, 40, 16, 16
    pw, ph = W - ml - mr, height - mb - mt
    xmax = nice_max(max((p[0] for p in points), default=0))
    ymax = nice_max(max((p[1] for p in points), default=0))
    out = [f'<svg viewBox="0 0 {W} {height}" width="100%" role="img" aria-label="{esc(title)}">']
    for t in ticks(ymax):
        y = mt + ph - ph * t / ymax
        out.append(f'<line class="grid" x1="{ml}" x2="{W - mr}" y1="{y:.1f}" y2="{y:.1f}"/>'
                   f'<text x="{ml - 6}" y="{y + 4:.1f}" text-anchor="end">{tick_label(t)}</text>')
    for t in ticks(xmax):
        x = ml + pw * t / xmax
        out.append(f'<text x="{x:.1f}" y="{mt + ph + 16}" text-anchor="middle">{tick_label(t)}</text>')
    out.append(f'<text x="{ml + pw / 2}" y="{height - 4}" text-anchor="middle">{esc(xlab)}</text>'
               f'<text x="12" y="{mt + ph / 2}" transform="rotate(-90 12 {mt + ph / 2})" text-anchor="middle">{esc(ylab)}</text>')
    ranked = sorted(points, key=lambda p: -(p[0] / xmax) * (p[1] / ymax))
    for x, y, lab, tip in points:
        cx, cy = ml + pw * x / xmax, mt + ph - ph * y / ymax
        out.append(f'<circle class="mark" cx="{cx:.1f}" cy="{cy:.1f}" r="5" fill="var(--s1)" fill-opacity=".85" '
                   f'stroke="var(--surface-1)" stroke-width="2"><title>{esc(tip)}</title></circle>')
    for x, y, lab, _ in ranked[:label_top]:
        cx, cy = ml + pw * x / xmax, mt + ph - ph * y / ymax
        anchor = 'end' if cx > W * 0.7 else 'start'
        dx = -8 if anchor == 'end' else 8
        out.append(f'<text class="val" x="{cx + dx:.1f}" y="{cy - 6:.1f}" text-anchor="{anchor}">{esc(lab)}</text>')
    out.append('</svg>')
    return ''.join(out)


def stacked_hbars(rows, names, colors, label_w=300, title='', group=1, ends=None, unit=''):
    """Stacked bars. rows: [(label, [values per name], total?)]. Without a total a row
    fills the width (100% stacked). With group > 1, every `group` rows form one block:
    the first row carries the bold block label, blocks are separated by extra space.
    ends: optional text per row, drawn after the bar."""
    hp, gap = 20, 12 if group > 1 else 0
    blocks = (len(rows) + group - 1) // group
    h = hp * len(rows) + gap * blocks + 8
    pw = W - label_w - (70 if ends else 16)
    out = [f'<svg viewBox="0 0 {W} {h}" width="100%" role="img" aria-label="{esc(title)}">']
    for i, row in enumerate(rows):
        lab, vals = row[0], row[1]
        y, bh = 4 + i * hp + gap * (i // group), 14
        # optional explicit total: bar fills only sum(vals)/total of the width
        total = (row[2] if len(row) > 2 else None) or sum(vals) or 1
        short = lab if len(lab) <= 46 else '…' + lab[-45:]
        cls = ' class="val"' if group > 1 and i % group == 0 else ''
        out.append(f'<text{cls} x="{label_w - 8}" y="{y + bh - 3}" text-anchor="end">{esc(short)}</text>')
        x = label_w
        vis = [(n, c, v) for n, c, v in zip(names, colors, vals) if v]
        for k, (n, c, v) in enumerate(vis):
            w = pw * v / total - (2 if k < len(vis) - 1 else 0)
            tip = f'<title>{esc(lab)} · {esc(n)}: {fmt(round(v, 1))}{unit}</title>'
            if k == len(vis) - 1:
                out.append(f'<path class="mark" d="{bar_path(x, y, max(w, 1), bh, horizontal=True)}" fill="{c}">{tip}</path>')
            else:
                out.append(f'<rect class="mark" x="{x:.1f}" y="{y}" width="{max(w, .5):.1f}" height="{bh}" fill="{c}">{tip}</rect>')
            x += w + 2
        if ends and ends[i]:
            out.append(f'<text class="val" x="{x + 4:.1f}" y="{y + bh - 3}">{esc(ends[i])}</text>')
    out.append('</svg>')
    return legend(names, colors) + ''.join(out)


def table(headers, rows, numeric=()):
    th = ''.join(f'<th class="{"n" if i in numeric else ""}">{esc(h)}</th>' for i, h in enumerate(headers))
    body = ''.join('<tr>' + ''.join(f'<td class="{"n" if i in numeric else ""}">{esc(fmt(v))}</td>'
                                    for i, v in enumerate(r)) + '</tr>' for r in rows)
    return f'<table><thead><tr>{th}</tr></thead><tbody>{body}</tbody></table>'


def details(summary, inner):
    return f'<details><summary>{esc(summary)}</summary>{inner}</details>'


def md_table(headers, rows):
    out = ['| ' + ' | '.join(headers) + ' |', '|' + '---|' * len(headers)]
    out += ['| ' + ' | '.join(fmt(v) for v in r) + ' |' for r in rows]
    return '\n'.join(out)


def md_to_html(text):
    """Minimal markdown: headings, bullets, tables, bold, code, paragraphs."""
    def inline(s):
        s = esc(s)
        s = re.sub(r'`([^`]+)`', r'<code>\1</code>', s)
        return re.sub(r'\*\*([^*]+)\*\*', r'<strong>\1</strong>', s)
    out, lst, tbl = [], False, []

    def flush_tbl():
        nonlocal tbl
        if tbl:
            rows = [[c.strip() for c in r.strip('|').split('|')] for r in tbl if not re.match(r'^\|?\s*-', r.replace('|', '').strip() or '-') or '---' not in r]
            rows = [r for r in rows if not all(set(c) <= set('-: ') for c in r)]
            out.append('<table>' + ''.join('<tr>' + ''.join(f'<{"th" if i == 0 else "td"}>{inline(c)}</{"th" if i == 0 else "td"}>' for c in r) + '</tr>'
                                           for i, r in enumerate(rows)) + '</table>')
            tbl = []
    for l in text.splitlines():
        if l.startswith('|'):
            tbl.append(l)
            continue
        flush_tbl()
        if re.match(r'^\s*[-*] ', l):
            if not lst:
                out.append('<ul>')
                lst = True
            out.append(f'<li>{inline(re.sub(r"^\s*[-*] ", "", l))}</li>')
            continue
        if lst:
            out.append('</ul>')
            lst = False
        m = re.match(r'^(#{1,4}) (.*)', l)
        if m:
            out.append(f'<h3>{inline(m.group(2))}</h3>')
        elif l.strip():
            out.append(f'<p>{inline(l)}</p>')
    flush_tbl()
    if lst:
        out.append('</ul>')
    return '\n'.join(out)


def render(run):
    m = json.load(open(os.path.join(run, 'metrics.json')))
    reading_path = os.path.join(run, 'reading.md')
    reading = open(reading_path).read() if os.path.exists(reading_path) else ''
    s = m['summary']
    H, M = [], []

    H.append(f'<h1>Repository health · {esc(m["repo"])}</h1>'
             f'<p class="sub">HEAD {esc(m["head"])} · {s["first"]} – {s["last"]} · generated {esc(m["generated"])}. '
             f'Revert chains ({s["revert_chain_commits"]} commits) and generated files excluded. '
             f'Rework = replaced lines younger than {m["params"]["recent_days"]} days.</p>')
    M.append(f'# Repository health · {m["repo"]}\n\nHEAD {m["head"]} · {s["first"]} – {s["last"]} · generated {m["generated"]}\n')
    tiles = [('Code size', f'{fmt(s["loc"])} lines'), ('Analyzed commits', s['analyzed_commits']),
             ('Churn', f'{fmt(s["churn"])} lines'), ('Rework rate', f'{s["rework_pct"]}%'),
             ('Authors (active 90d)', f'{s["authors"]} ({s["active_authors"]})'), ('Merged PRs', s['merged_prs'])]
    H.append('<div class="tiles">' + ''.join(f'<div class="tile"><div class="l">{esc(l)}</div><div class="v">{esc(v)}</div></div>'
                                             for l, v in tiles) + '</div>')
    M.append(md_table(['Metric', 'Value'], tiles))
    if s['intent_labelled'] < s['analyzed_commits']:
        note = f'{s["analyzed_commits"] - s["intent_labelled"]} commits use keyword-based intent.'
        H.append(f'<p class="sub">{esc(note)}</p>')
        M.append(f'\n_{note}_')

    if reading:
        H.append('<h2>Reading</h2><div class="reading">' + md_to_html(reading) + '</div>')
        M.append('\n## Reading\n\n' + reading.strip())

    # churn by intent
    months = [r['month'] for r in m['churn_by_month']]
    used = [i for i in INTENTS if any(r[i] for r in m['churn_by_month'])]
    series = [(i, f'var(--s{INTENTS.index(i) + 1})', [r[i] for r in m['churn_by_month']]) for i in used]
    H.append('<h2>Churn by intent</h2><p class="sub">Lines added + deleted per month, by commit intent.</p>')
    H.append(columns(months, series, title='Churn by intent'))
    rows = [[r['month']] + [r[i] for i in used] for r in m['churn_by_month']]
    H.append(details('Table', table(['Month'] + used, rows, numeric=range(1, len(used) + 1))))
    M.append('\n## Churn by intent\n\n' + md_table(['Month'] + used, rows))

    # rework
    rw = m['rework_by_month']
    H.append('<h2>Rework</h2><p class="sub">Share of changed lines that replace code younger than '
             f'{m["params"]["recent_days"]} days. Legacy = replacing older code (refactoring, maintenance).</p>')
    H.append('<h3>Rework rate by month</h3>' + lines([r['key'] for r in rw], [('rework', 'var(--s1)', [r['rework_pct'] for r in rw])], unit='%', title='Rework by month'))
    H.append('<h3>New vs replaced lines</h3>' + columns([r['key'] for r in rw], [
        ('new', 'var(--s1)', [r['new'] for r in rw]), ('replaced recent', 'var(--s2)', [r['recent'] for r in rw]),
        ('replaced legacy', 'var(--s3)', [r['legacy'] for r in rw])], title='New vs replaced'))
    rw_names = ['new', 'replaced recent', 'replaced legacy']
    rw_colors = ['var(--s1)', 'var(--s2)', 'var(--s3)']

    def rework_bars(rs):
        top = max((r['new'] + r['recent'] + r['legacy'] for r in rs), default=1)
        return stacked_hbars([(r['key'], [r['new'], r['recent'], r['legacy']], top) for r in rs], rw_names, rw_colors,
                             ends=[f'{r["rework_pct"]}%' for r in rs], unit=' lines')
    ra = m['rework_by_area'][:20]
    ri = m['rework_by_intent']
    H.append('<h3>Rework by area</h3><p class="sub">Sorted by rework rate (shown at the bar end). '
             'Bar length = changed lines.</p>' + rework_bars(ra))
    H.append(details('Table', table(['Area', 'New', 'Recent', 'Legacy', 'Rework %'],
                                    [[r['key'], r['new'], r['recent'], r['legacy'], r['rework_pct']] for r in ra], numeric=(1, 2, 3, 4))))
    H.append('<h3>Rework by intent</h3>' + rework_bars(ri))
    H.append(details('Table', table(['Intent', 'New', 'Recent', 'Legacy', 'Rework %'],
                                    [[r['key'], r['new'], r['recent'], r['legacy'], r['rework_pct']] for r in ri], numeric=(1, 2, 3, 4))))
    M.append('\n## Rework\n\n' + md_table(['Month', 'New', 'Recent', 'Legacy', 'Rework %'], [[r['key'], r['new'], r['recent'], r['legacy'], r['rework_pct']] for r in rw]))
    M.append('\n### By area\n\n' + md_table(['Area', 'New', 'Recent', 'Legacy', 'Rework %'], [[r['key'], r['new'], r['recent'], r['legacy'], r['rework_pct']] for r in ra]))
    M.append('\n### By intent\n\n' + md_table(['Intent', 'New', 'Recent', 'Legacy', 'Rework %'], [[r['key'], r['new'], r['recent'], r['legacy'], r['rework_pct']] for r in ri]))

    # fixes
    fx = [r for r in m['fix_by_area'] if r['lines'] >= 500][:20]
    H.append('<h2>Fix share by area</h2><p class="sub">Share of an area\'s churn coming from bug-fix commits. Areas with ≥500 changed lines.</p>')
    H.append(hbars([(r['area'], r['fix_pct'], f'{r["area"]}: {r["fix_commits"]}/{r["commits"]} commits are fixes') for r in fx], unit='%'))
    frows = [[r['area'], r['commits'], r['fix_commits'], r['lines'], r['fix_pct']] for r in fx]
    H.append(details('Table', table(['Area', 'Commits', 'Fix commits', 'Lines', 'Fix %'], frows, numeric=(1, 2, 3, 4))))
    M.append('\n## Fix share by area\n\n' + md_table(['Area', 'Commits', 'Fix commits', 'Lines', 'Fix %'], frows))

    # hotspots
    hs = m['hotspots']
    H.append('<h2>Hotspots</h2><p class="sub">Revisions (excluding chores) × whitespace complexity (sum of indentation levels). '
             'Top-right files change often and are hard to change.</p>')
    H.append(scatter([(h['revisions'], h['complexity'], h['path'].rsplit('/', 1)[-1],
                       f'{h["path"]}: {h["revisions"]} revisions, complexity {h["complexity"]}, {h["loc"]} lines') for h in hs],
                     'revisions', 'complexity', title='Hotspots'))
    hrows = [[h['path'], h['revisions'], h['recent_revisions'], h['loc'], h['complexity'], h['authors']] for h in hs[:15]]
    H.append(details('Table', table(['File', 'Revisions', 'Last 90d', 'Lines', 'Complexity', 'Authors'], hrows, numeric=(1, 2, 3, 4, 5))))
    M.append('\n## Hotspots\n\n' + md_table(['File', 'Revisions', 'Last 90d', 'Lines', 'Complexity', 'Authors'], hrows))

    # coupling
    cp = m['coupling'][:20]
    H.append('<h2>Change coupling</h2><p class="sub">File pairs that change together (degree = shared ÷ average revisions). '
             'Cross-area pairs that are not tests hint at hidden dependencies.</p>')
    crows = [[c['a'], c['b'], c['shared'], c['degree'], ('cross-area' if c['cross_area'] else '') + (' test' if c['test_pair'] else '')] for c in cp]
    def pair(c):
        a, b = c['a'].rsplit('/', 1)[-1], c['b'].rsplit('/', 1)[-1]
        return f'{a} ↔ {b}' + (' ⚑' if c['cross_area'] and not c['test_pair'] else '')
    H.append('<p class="sub">⚑ = cross-area, non-test pair. Hover for full paths and degree.</p>')
    H.append(hbars([(pair(c), c['shared'], f'{c["a"]} ↔ {c["b"]}: {c["shared"]} shared commits, degree {c["degree"]}') for c in cp],
                   label_w=360))
    H.append(details('Table', table(['File A', 'File B', 'Shared', 'Degree', 'Flags'], crows, numeric=(2, 3))))
    M.append('\n## Change coupling\n\n' + md_table(['File A', 'File B', 'Shared', 'Degree', 'Flags'], crows))

    # ownership
    own = m['ownership']
    # Fixed author colors across editors and reviewers: top 4 by combined share, rest "Other".
    weight = {}
    for o in own:
        for t in o['top'] + o.get('reviewers', []):
            weight[t['author']] = weight.get(t['author'], 0) + t['pct']
    top_auth = sorted(weight, key=lambda a: -weight[a])[:4]
    names = top_auth + ['Other']
    colors = [f'var(--s{i + 1})' for i in range(len(top_auth))] + ['var(--other)']

    def shares(lst):
        d = {t['author']: t['pct'] for t in lst}
        vals = [d.get(a, 0) for a in top_auth]
        return vals + [max(100 - sum(vals), 0)] if lst else [0] * len(names)
    orows = []
    for o in own:
        orows.append((o['area'], shares(o['top'])))
        rv, rp = o.get('reviewers', []), o.get('reviewed_pct', 0) or 0
        orows.append((f'reviewers ({rp:.0f}% reviewed)' if rv else 'reviewers (none)',
                      [v * rp / 100 for v in shares(rv)], 100))
    H.append('<h2>Knowledge distribution</h2><p class="sub">Per area: editors (surviving lines at HEAD, git blame) '
             'and reviewers (PR reviews weighted by changed lines). The reviewer bar\'s length is the reviewed share '
             'of the area\'s PR churn. '
             'Truck factor = editors needed to cover over half the lines.</p>')
    H.append(stacked_hbars(orows, names, colors, title='Editors and reviewers', group=2, unit='%'))
    def who(lst):
        return ', '.join(f'{x["author"]} {x["pct"]:.0f}%' for x in lst[:2])

    def dom(o):
        d = o.get('dominant_commit')
        return f'{d["pct"]:.0f}% {d["author"]}: {d["subject"][:50]} ({d["intent"]})' if d and d['pct'] >= 25 else ''
    otab = [[o['area'], o['lines'], o['truck_factor'], who(o['top']), who(o.get('knowledge', [])),
             who(o.get('reviewers', [])), o.get('reviewed_pct'), dom(o), o['inactive_pct']] for o in own]
    ohdr = ['Area', 'Lines', 'Truck factor', 'Blame owners', 'History owners', 'Reviewers', 'Reviewed %',
            'Dominant commit (≥25%)', 'Inactive %']
    H.append('<p class="sub">Editors: blame owners = who last touched surviving lines (chores and revert chains passed through). '
             'History owners = lines added over time, halving every 180 days. '
             'Reviewers = PR reviewers other than the author, weighted by changed lines with the same decay. '
             'Reviewed % = share of PR-merged churn reviewed by someone else. '
             'A dominant commit means one change (often a refactor) decides the blame owner.</p>')
    H.append(details('Table', table(ohdr, otab, numeric=(1, 2, 6, 8))))
    M.append('\n## Knowledge distribution\n\n' + md_table(ohdr, otab))

    # survival
    sv = m['survival']
    H.append('<h2>Code survival</h2><p class="sub">Share of lines added in a month that still exist at HEAD. '
             'Recent months are naturally high.</p>')
    H.append(columns([r['month'] for r in sv], [('survival', 'var(--s1)', [r['survival_pct'] for r in sv])], unit='%', title='Survival'))
    srows = [[r['month'], r['added'], r['alive'], r['survival_pct']] for r in sv]
    H.append(details('Table', table(['Month', 'Added', 'Alive', 'Survival %'], srows, numeric=(1, 2, 3))))
    M.append('\n## Code survival\n\n' + md_table(['Month', 'Added', 'Alive', 'Survival %'], srows))

    # PR flow
    pr = m['pr_by_month']
    if pr:
        pm = [r['month'] for r in pr]
        H.append('<h2>Pull request flow</h2><p class="sub">Merged PRs by merge month. Reviews exclude the PR author.</p>')
        H.append('<h3>Merged PRs</h3>' + columns(pm, [('merged', 'var(--s1)', [r['merged'] for r in pr])], title='Merged PRs'))
        H.append('<h3>Cycle time (open → merge), hours</h3>' + lines(pm, [
            ('median', 'var(--s1)', [r['cycle_p50_h'] for r in pr]), ('p90', 'var(--s2)', [r['cycle_p90_h'] for r in pr])], unit='h', title='Cycle time'))
        H.append('<h3>Time to first review, median hours</h3>' + lines(pm, [('first review', 'var(--s1)', [r['first_review_p50_h'] for r in pr])], unit='h'))
        H.append('<h3>Merged without review</h3>' + columns(pm, [('unreviewed', 'var(--s1)', [r['unreviewed_pct'] for r in pr])], unit='%'))
        H.append('<h3>PR size (lines changed)</h3>' + lines(pm, [
            ('median', 'var(--s1)', [r['size_p50'] for r in pr]), ('p90', 'var(--s2)', [r['size_p90'] for r in pr])], title='PR size'))
        prows = [[r['month'], r['merged'], r['cycle_p50_h'], r['cycle_p90_h'], r['first_review_p50_h'], r['size_p50'], r['size_p90'],
                  r['rounds_p50'], r['unreviewed_pct']] for r in pr]
        hdr = ['Month', 'Merged', 'Cycle p50 h', 'Cycle p90 h', 'First review p50 h', 'Size p50', 'Size p90', 'Reviews p50', 'Unreviewed %']
        H.append(details('Table', table(hdr, prows, numeric=range(1, 9))))
        sl = [[f'#{p["number"]}', p['title'], p['author'], p['cycle_h'], p['size'], p['rounds']] for p in m['pr_slowest']]
        H.append(details('Slowest PRs', table(['PR', 'Title', 'Author', 'Cycle h', 'Size', 'Reviews'], sl, numeric=(3, 4, 5))))
        M.append('\n## Pull request flow\n\n' + md_table(hdr, prows))
        M.append('\n### Slowest PRs\n\n' + md_table(['PR', 'Title', 'Author', 'Cycle h', 'Size', 'Reviews'], sl))

    doc = (f'<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width">'
           f'<title>Repository health · {esc(m["repo"])}</title><style>{css()}</style></head>'
           f'<body class="viz-root"><main>{"".join(H)}</main></body></html>')
    with open(os.path.join(run, 'report.html'), 'w') as f:
        f.write(doc)
    with open(os.path.join(run, 'report.md'), 'w') as f:
        f.write('\n'.join(M) + '\n')
    print(os.path.join(run, 'report.html'))


if __name__ == '__main__':
    if len(sys.argv) != 2:
        sys.exit(__doc__)
    render(sys.argv[1])
