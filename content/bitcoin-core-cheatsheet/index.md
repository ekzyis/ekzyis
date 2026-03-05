---
title: Bitcoin Core Cheatsheet
date: 2026-03-05
tags: bitcoin dev nixos
---

Build [Bitcoin Core](https://github.com/bitcoin/bitcoin) with
[bitcoin-nix-tools](https://github.com/uncomputable/bitcoin-nix-tools):

```
$ git clone git@github.com:bitcoin/bitcoin
$ git clone git@github.com:uncomputable/bitcoin-nix-tools
$ nix-shell bitcoin-nix-tools/
$ cd bitcoin
$ cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=1 -B build
$ cmake --build build -j8
```

<sub>With `CMAKE_EXPORT_COMPILE_COMMANDS`, `cmake` will create
compile_commands.json in build/. It is needed to enable code navigation with
`clangd`. It will look for this file in build/ by default.</sub>

Functional tests:

```
# all functional tests
$ python build/test/functional/test_runner.py
# specific functional test
$ python build/test/functional/test_runner.py feature_signet.py
```

Unit tests:

```
# all unit tests
$ build/test/test_bitcoin
# specific unit test
$ build/bin/test_bitcoin --run_test=txospenderindex_tests
```

Clean rebuild in one command:

```
$ rm -rf build && cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=1 -B build && cmake --build build -j8
```

<sub>Usually, only recompiling changed files with `cmake --build` is enough,
though.</sub>
