---
title: lnfuzz
date: 2026-09-20
sn_id: 1577486
frontpage: true
---

Lightning fuzzers (the people, not the software) now have their own website!

> # About
>
> We are a group of security researchers and developers working to improve the security and
> robustness of the Bitcoin Lightning Network.
>
> We use fuzzing and other techniques to find bugs, [publishing](https://lnfuzz.org/bugs/) them on
> this site along with [advisories](https://lnfuzz.org/advisories/) for security issues. We also
> write [technical posts](https://lnfuzz.org/posts/) about our fuzzing tools and techniques.
>
> ## Why
>
> The Lightning Network’s four main implementations (LND, Core Lightning, LDK, and Eclair) all speak
> the same peer protocol. They have to interoperate with each other, and they have to hold up
> against a peer that is not playing along. Both of those are hard to get right, and neither is well
> covered by the tests a project writes for itself, which tend to assume a well-behaved peer.
>
> That is the gap our tooling helps to address. Our fuzzing tools can play the role of a misbehaving
> peer, triggering unusual code paths and edge cases. This approach uncovers many kinds of issues,
> ranging from protocol violations to crashes and assertion failures, memory safety errors, and
> more.

For now, the only project is [smite](https://lnfuzz.org/projects/smite/), a coverage-guided snapshot
fuzzer for all four major lightning implementations. We use snapshot fuzzing because lightning is
stateful, and nondeterminism makes findings difficult to reproduce. To ensure deterministic
execution, we reset the node between fuzz runs to a snapshot using [Nyx](https://nyx-fuzz.com/), a
specialized virtual machine.
