# Daedalus — Codebase Assessment

## Testing (highest priority gap)

- [ ] Add unit tests for `git/git.go` (`Diff`, `Commit`, `ValidateCommit`) — code runs `git commit`, no safety net
- [ ] Add unit tests for `config/loader.go` — Viper-backed config with dotenv at `$HOME/.config/daedalus/config`
- [ ] Add unit tests for `gemini/gemini.go` — mock the `genai.Client`, cover success / API error / `instantiate()` failure paths
- [ ] Add unit tests for `application/stylus.go` and `application/oracle.go` orchestration logic
- [ ] Add a `make test` target (and wire `go test ./...` into CI)
- [ ] Set up coverage reporting (`go test -cover` / `go test -coverprofile`)

## CI/CD

- [ ] Add GitHub Actions workflow: `lint` (`go vet`, `staticcheck`, `golangci-lint`) + `test` on push/PR
- [ ] Add a release workflow that runs on tag push and uploads binaries (the `make release` target is ready but not wired)
- [ ] Add a `make lint` target

## Documentation

- [ ] Expand `README.md` — usage examples for `stylus` and `oracle`, command reference, configuration docs
- [ ] Add `CHANGELOG.md` (start from the 33 commits on `main`)
- [ ] Add `CONTRIBUTING.md` if this will be shared
- [ ] Add a short architecture note (the package layout — `cmd/`, `application/`, `daedalus/`, `gemini/`, `git/`, `config/` — is non-obvious)
- [ ] Replace placeholder Cobra `Long` descriptions in `cmd/root.go`, `cmd/oracle.go`, `cmd/stylus.go`

## Bugs / Broken Behavior

- [ ] **Fix Makefile version resolution** — `Makefile:3` calls `node -p "require('./package.json').version"` but no `package.json` exists; version always reports `0.0.0`. Either add a `package.json`, switch to `git describe --tags`, or read from a Go `var Version = "..."`
- [ ] **Create at least one git tag** (e.g. `v0.1.0`) so `git describe` works and releases are meaningful
- [ ] **Finish the `sync.Once` refactor in `gemini/gemini.go`** — `once` and `initErr` fields are declared (lines 25-27) but `instantiate()` still uses the explicit mutex; `initErr` is never read. Either complete it or revert the half-done change
- [ ] **Remove or implement `daedalus/util.go` and `daedalus/weaver.go`** — both are 1-line empty stubs (likely abandoned)

## Code Quality

- [ ] **Make the Gemini model configurable** — `gemini/gemini.go:17` hardcodes `"gemini-3.1-flash-lite"`; add a config key or `--model` flag
- [ ] **Implement `oracle config` without `--key`** — `cmd/oracle.go:46-48` currently does `os.Exit(0)` with no output (stub)
- [ ] **Clean up the empty-case fallthrough in `application/oracle.go:28-31`** — the `""` case has an empty body and falls into `default`; rewrite as an explicit default
- [ ] **Address the `TODO` in `gemini/gemini.go:58`** — _"Add validation if model already reached the limit so we can switch to another model"_
- [ ] **Tighten error messages in `git/git.go:37` and `git/git.go:56`** — inconsistent punctuation inside `fmt.Errorf` strings
- [ ] **Run `goimports` / `gofmt`** across the repo and clean up any flagged files (e.g. `main.go`)
- [ ] **Review uncommitted working-tree changes** in `application/oracle.go` and `gemini/gemini.go` — commit or stash

## Configuration & Hygiene

- [ ] Add `.editorconfig` for Go formatting consistency
- [ ] Move the built binary out of the source tree — `daedalus/daedalus` (16MB arm64 binary) sits inside the package dir, even if gitignored
- [ ] Add `.dockerignore` if/when a Dockerfile is added
- [ ] Audit `.gitignore` — currently only 3 entries; consider adding `.idea/`, `.vscode/`, `coverage.out`, `*.log`

## Architecture / Structure (optional, future)

- [ ] Define an `LLMProvider` interface in `daedalus/` so other providers (OpenAI, Anthropic, local) can plug in beside Gemini
- [ ] Consider an interface for the `git` wrapper so it can be mocked in tests
- [ ] Split `application/stylus.go` if it grows — diff parsing, prompt building, and commit writing are distinct concerns

## Security (low priority for personal tool, but worth a pass)

- [ ] Audit how Gemini API key is loaded from the dotenv file — confirm it's never logged or echoed
- [ ] Sanitize the LLM-generated commit message before passing to `git commit` (e.g. strip backticks, newlines, control chars) — current code appears to pass the raw model output through
