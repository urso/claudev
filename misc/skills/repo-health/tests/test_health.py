import os, sys

import health


def commit_of(repo, sha):
    parents = repo.git('rev-list', '--parents', '-n1', sha).split()[1:]
    ts = int(repo.git('show', '-s', '--format=%at', sha))
    return {'sha': sha, 'parents': parents, 'ts': ts}


class TestLineAge:
    def test_root_commit_is_empty(self, repo):
        sha = repo.commit('init', {'a.txt': 'x\n'})
        assert health.line_age(commit_of(repo, sha)) == {}

    def test_counts_added_and_replaced(self, repo):
        repo.commit('init', {'a.txt': 'one\ntwo\n'})
        sha = repo.commit('edit', {'a.txt': 'one\nTWO\nthree\n'})
        # 2 added lines, 1 replaced line written just now -> recent
        assert health.line_age(commit_of(repo, sha)) == {'a.txt': [2, 1, 0]}

    def test_skips_blank_lines(self, repo):
        repo.commit('init', {'a.txt': 'one\n'})
        sha = repo.commit('edit', {'a.txt': 'one\n\n  \ntwo\n'})
        assert health.line_age(commit_of(repo, sha)) == {'a.txt': [1, 0, 0]}

    def test_content_that_looks_like_headers(self, repo):
        repo.commit('init', {'s.sql': '-- drop old index\nselect 1;\n', 'b.txt': 'keep\nold\n'})
        sha = repo.commit('edit', {'s.sql': 'select 1;\n++ added\n', 'b.txt': 'keep\nnew\n'})
        assert health.line_age(commit_of(repo, sha)) == {'s.sql': [1, 1, 0], 'b.txt': [1, 1, 0]}

    def test_ignores_noprefix_config(self, repo):
        repo.git('config', 'diff.noprefix', 'true')
        repo.commit('init', {'dir/a.txt': 'one\n'})
        sha = repo.commit('edit', {'dir/a.txt': 'ONE\n'})
        assert health.line_age(commit_of(repo, sha)) == {'dir/a.txt': [1, 1, 0]}

    def test_rename(self, repo):
        repo.commit('init', {'old.txt': 'a\nb\nc\nd\ne\n'})
        repo.git('mv', 'old.txt', 'new.txt')
        sha = repo.commit('rename', {'new.txt': 'a\nb\nc\nd\nE\n'})
        assert health.line_age(commit_of(repo, sha)) == {'new.txt': [1, 1, 0]}


class TestLookup:
    def test_prefix_lengths(self):
        m ={'abcdef1': 'fix', '1234567890': 'feature', 'f' * 40: 'docs'}
        intents = (m, sorted({len(k) for k in m}, reverse=True))
        assert health.lookup(intents, 'abcdef1' + '0' * 33) == 'fix'
        assert health.lookup(intents, '1234567890' + '0' * 30) == 'feature'
        assert health.lookup(intents, 'f' * 40) == 'docs'
        assert health.lookup(intents, '9' * 40) is None

    def test_load_intents(self, tmp_path):
        (tmp_path / 'intent.tsv').write_text('1234567890\tfix\nabcdef1 bogus\n\n')
        m, lens = health.load_intents(str(tmp_path))
        assert m == {'1234567890': 'fix'}
        assert lens == [10]

    def test_load_intents_missing(self, tmp_path):
        assert health.load_intents(str(tmp_path)) == ({}, [])


class TestExpand:
    shas = sorted(['aaa111', 'aab222', 'bbb333'])

    def test_unique_prefix(self):
        assert health.expand(self.shas, 'aab') == 'aab222'
        assert health.expand(self.shas, 'bbb333') == 'bbb333'

    def test_no_match(self):
        assert health.expand(self.shas, 'ccc') is None
        assert health.expand(self.shas, 'zzz') is None

    def test_empty_prefix(self):
        assert health.expand(self.shas, '') is None


class TestRevertChain:
    def test_drops_revert_and_target(self):
        commits = [
            {'sha': 'a' * 40, 'subject': 'feat: x', 'reverts': None},
            {'sha': 'b' * 40, 'subject': 'Revert "feat: x"', 'reverts': 'a' * 7},
            {'sha': 'c' * 40, 'subject': 'fix: y', 'reverts': None},
        ]
        assert health.revert_chain(commits) == {'a' * 40, 'b' * 40}


def test_main_runs_from_repo_root(repo, monkeypatch, capsys):
    repo.commit('init', {'sub/f.txt': 'x\n'})
    monkeypatch.chdir(repo.root / 'sub')
    monkeypatch.setattr(sys, 'argv', ['health.py', 'dir'])
    health.main()
    assert os.path.realpath(os.getcwd()) == os.path.realpath(repo.root)
    assert capsys.readouterr().out.strip().endswith(os.path.join('.git', 'repo-health'))
