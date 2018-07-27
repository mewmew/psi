# PSI

[![Build Status](https://travis-ci.org/mewmew/psi.svg?branch=master)](https://travis-ci.org/mewmew/psi)
[![Coverage Status](https://coveralls.io/repos/github/mewmew/psi/badge.svg?branch=master)](https://coveralls.io/github/mewmew/psi?branch=master)
[![GoDoc](https://godoc.org/github.com/mewmew/psi?status.svg)](https://godoc.org/github.com/mewmew/psi)

PlayStation 1 emulator, symbolic execution engine and binary lifter.

* The preliminary aim of this repository is to implement a PlayStation I emulator in Go.
* The secondary aim of this repository is to develop a symbolic execution engine for PSX.
* The tertiary aim of this repository is to leverage the inferred CPU state of the symbolic execution engine to lift R3059 MIPS assembly to LLVM IR.

## Credits

Thanks to [Lionel Flandrin](https://github.com/simias) for sharing the [Playstation Emulation Guide](https://github.com/simias/psx-guide); a truly wonderful resource documenting the implementation details of a PSX emulator in Rust.
