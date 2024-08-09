import os
from typing import List
import subprocess


def load() -> str:
    return open(os.path.join(os.path.dirname(__file__), 'funcs.go')).read()


def save(content: str):
    with open(os.path.join(os.path.dirname(__file__), 'funcs.go'), 'w') as f:
        f.write(content)


def inject(content: str, key: str, lines: List[str]) -> str:
    key = key.upper()
    mark_beg = f'__BEG_GEN:{key}__'
    mark_end = f'__END_GEN:{key}__'

    all_lines = content.split('\n')

    idx_beg = -1
    idx_end = -1

    for i, line in enumerate(all_lines):
        if mark_beg in line:
            idx_beg = i
        if mark_end in line:
            idx_end = i

    if idx_beg == -1 or idx_end == -1:
        return content

    all_lines = all_lines[:idx_beg + 1] + lines + all_lines[idx_end:]

    return '\n'.join(all_lines)


with_neg = [
    'uint8',
    'uint16',
    'uint32',
    'uint64',
    'int8',
    'int16',
    'int32',
    'int64',
    'float32',
    'float64',
    'complex64',
    'complex128',
    'int',
    'uint',
]

with_add = with_neg + ['string', 'uintptr']


content = load()

content = inject(content, 'add', [
    f"""case {t}:
        return a.({t}) + b.({t}), nil""" for t in with_add
])

content = inject(content, 'neg', [
    f"""case {t}:
        return -a, nil""" for t in with_neg
])

save(content)

subprocess.run(['go', 'fmt', os.path.join(
    os.path.dirname(__file__), 'funcs.go')])
