import os, subprocess, sys

import pytest

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'scripts'))


class Repo:
    def __init__(self, root):
        self.root = root

    def git(self, *a):
        return subprocess.run(('git',) + a, cwd=self.root, check=True, capture_output=True, text=True).stdout

    def write(self, path, text):
        p = self.root / path
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(text)

    def commit(self, msg, files):
        for path, text in files.items():
            self.write(path, text)
        self.git('add', '-A')
        self.git('commit', '-qm', msg)
        return self.git('rev-parse', 'HEAD').strip()


@pytest.fixture
def repo(tmp_path, monkeypatch):
    for k in ('GIT_DIR', 'GIT_WORK_TREE', 'GIT_INDEX_FILE'):
        monkeypatch.delenv(k, raising=False)
    monkeypatch.setenv('GIT_CONFIG_GLOBAL', os.devnull)
    monkeypatch.setenv('GIT_CONFIG_NOSYSTEM', '1')
    r = Repo(tmp_path)
    r.git('init', '-q', '-b', 'main')
    r.git('config', 'user.name', 'Test')
    r.git('config', 'user.email', 'test@example.com')
    monkeypatch.chdir(tmp_path)
    return r
