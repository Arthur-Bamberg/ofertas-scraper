# Versioned Git hooks run `go test ./...` before commit

There is no Husky equivalent in the Go toolchain. This repo uses a versioned `.githooks/` directory and `git config core.hooksPath .githooks` (via `scripts/install-git-hooks.sh`) so every commit runs `go test ./...` with zero Node/Python/Lefthook dependencies. Hooks are opt-in per clone after install; CI remains the backstop if someone skips hooks.
