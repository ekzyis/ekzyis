---
title: Bitcoin Core Cheatsheet
date: 2026-03-05
frontpage: true
---

## Build

**Build [Bitcoin Core](https://github.com/bitcoin/bitcoin) with
[bix](https://github.com/bitcoin-dev-tools/bix):**

```
$ git clone git@github.com:bitcoin/bitcoin
$ git clone git@github.com:bitcoin-dev-tools/bix
$ nix develop bix/
$ cd bitcoin
$ cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=1 -B build
$ cmake --build build -j$(nproc)
```

<sub>With `CMAKE_EXPORT_COMPILE_COMMANDS`, `cmake` will create
compile_commands.json in build/. It is needed to enable code navigation with
`clangd`. It will look for this file in build/ by default.</sub>

**Clean rebuild in one command:**

```
$ rm -rf build && cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=1 -B build && cmake --build build -j$(nproc)
```

<sub>Usually, only recompiling changed files with `cmake --build` is enough,
though.</sub>

## Test

**Functional tests:**

```
# all functional tests
$ python build/test/functional/test_runner.py

# specific functional test
$ python build/test/functional/test_runner.py feature_signet.py
```

**Unit tests:**

```
# all unit tests
$ build/bin/test_bitcoin

# specific unit test
$ build/bin/test_bitcoin --run_test=txospenderindex_tests
```

## Fuzz

**Build for fuzzing:**

```
$ cmake --preset libfuzzer
$ cmake --build build_fuzz -j$(nproc)
```

**Fuzz specific target:**

```
$ FUZZ_TARGET=process_message
$ FUZZ=$FUZZ_TARGET build_fuzz/bin/fuzz
# fuzz target for 1 minute and generate corpus
$ mkdir -p fuzz_corpora/$FUZZ_TARGET
$ FUZZ=$FUZZ_TARGET build_fuzz/bin/fuzz -max_total_time=60 fuzz_corpora/$FUZZ_TARGET
```

**Build for fuzz coverage:**

```
$ cmake --preset libfuzzer -B build_cov \
    -DCMAKE_C_FLAGS="-fprofile-instr-generate -fcoverage-mapping" \
    -DCMAKE_CXX_FLAGS="-fprofile-instr-generate -fcoverage-mapping"
$ cmake --build build_cov -j$(nproc)
```

**Fuzz for coverage data:**

```
$ FUZZ=$FUZZ_TARGET build_cov/bin/fuzz -max_total_time=5 fuzz_corpora/$FUZZ_TARGET
$ llvm-profdata merge build_cov/raw_profile_data/*.profraw -o build_cov/coverage.profdata
```

**Generate HTML coverage report:**

```
$ llvm-cov show \
    --object=build_cov/bin/fuzz \
    -Xdemangler=llvm-cxxfilt \
    --instr-profile=build_cov/coverage.profdata \
    --ignore-filename-regex="src/crc32c/|src/leveldb/|src/minisketch/|src/secp256k1/|src/test/" \
    --format=html \
    --show-instantiation-summary \
    --show-line-counts-or-regions \
    --show-expansions \
    --output-dir=build_cov/coverage_report \
    --project-title="Bitcoin Core Fuzz Coverage Report"
```

<sub>See doc/fuzzing.md in the Bitcoin Core repository for more info.</sub>
