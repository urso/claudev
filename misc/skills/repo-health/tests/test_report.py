import report


def test_table_escapes_cells():
    out = report.table(['Title'], [['<img src=x onerror=alert(1)>'], ['a & b']])
    assert '<img' not in out
    assert '&lt;img src=x onerror=alert(1)&gt;' in out
    assert 'a &amp; b' in out


def test_table_escapes_headers():
    assert '&lt;h&gt;' in report.table(['<h>'], [])


def test_table_numeric_class():
    out = report.table(['Name', 'Count'], [['x', 12345]], numeric=(1,))
    assert '<td class="n">12.3K</td>' in out
    assert '<td class="">x</td>' in out
