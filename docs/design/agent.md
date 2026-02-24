# `tk agent` — Design Document

## Overview

`tk agent` is a conversational AI assistant for Tanka. Users working with Tanka spend a lot of time on repetitive configuration tasks — creating environments, instantiating Helm charts, updating job specs, scaling deployments. The agent acts as a senior SRE: it edits Jsonnet files, validates them, and shows diffs for human review and deployment.

The agent deliberately never runs `tanka apply`. All cluster deploys remain human-controlled. This is the central design constraint.

---

## Key Design Decisions

| Decision | Rationale |
|---|---|
| **Prepare-only** (no `apply`) | Keeps humans in the loop for all changes that affect running workloads. The agent is a preparation tool, not a deployment tool. |
| **Sandboxed to git repo root** | All file operations are anchored to the repo root discovered via go-git. The agent cannot touch files outside the repo regardless of which subdirectory it was invoked from. |
| **Ephemeral sessions** | Each invocation starts fresh. No persisted conversation history between runs. Simplifies state management and avoids stale context surprises. |
| **Git-gated startup** | Fails immediately with a clear error if not invoked inside a git repository. Prevents confusing partial-tool failures later. |
| **In-process Tanka API** | Tanka operations call Go functions directly rather than shelling out to the `tk` binary. This avoids PATH issues, is faster, and enables clean error propagation. |
| **Pure-Go git** | Uses `go-git` rather than the `git` binary. No external binary dependency; works anywhere the binary runs. |

---

## Architecture

### Code Layout

```
cmd/tk/
  agent.go                   # CLI command: flag parsing, startup checks, mode dispatch

pkg/agent/
  agent.go                   # ADK agent setup, system prompt, Run() loop
  config.go                  # Config loading: ~/.config/tanka/agent.yaml + env vars
  provider.go                # NewModel() factory: routes to Gemini/Anthropic/OpenAI
  repl.go                    # Interactive REPL (chzyer/readline) with ASCII art banner
  oneshot.go                 # Single-request mode
  verbose.go                 # Verbose logging: lipgloss borders, colordiff for diffs
  gemini.go                  # Gemini model (ADK native)
  anthropic.go               # Anthropic model (custom model.LLM adapter)
  openai.go                  # OpenAI model (custom model.LLM adapter)
  tools/
    tool.go                  # Shared helpers: mustSchema(), bind()
    filesystem.go            # Sandboxed file ops
    git.go                   # Read-only git operations (go-git)
    tanka.go                 # Tanka + Jsonnet in-process ops
    jb.go                    # jsonnet-bundler dependency management
    helm.go                  # Helm dependency management
```

### Request Flow

```
User input
    │
    ▼
cmd/tk/agent.go
  • find git root (go-git)
  • load config (yaml + env vars + flags)
  • validate API key
  • construct model.LLM
  • construct Agent (registers all tools)
    │
    ├─ one-shot arg? ──► RunOneShot() ──► Agent.Run() ──► print response
    │
    └─ interactive? ──► RunREPL() (readline loop)
                              │
                              └─► Agent.Run() on each input
                                        │
                                        ▼
                              ADK runner.Run() (iter)
                                        │
                          ┌─────────────┴──────────────┐
                          │                            │
                    LLM API call              Tool call dispatch
                          │                            │
                    model.LLM.GenerateContent()    tools.NewXxxTools()
                    (Gemini / Anthropic / OpenAI)      │
                                                functiontool.New() handlers
```

---

## Agent Framework: Google ADK Go

The agent loop is built on `google.golang.org/adk` (`v0.5.0`). ADK manages the multi-turn conversation loop: it calls the LLM, detects tool call requests, invokes the registered tools, appends results to the conversation, and loops until the LLM produces a final text response.

**Core ADK types used:**

| Type | Role |
|---|---|
| `llmagent.New(Config)` | Creates the agent with a model, system prompt, and tool list |
| `runner.New(Config)` | Orchestrates the run loop, session, and event dispatch |
| `session.InMemoryService()` | Per-invocation in-memory session store |
| `functiontool.New(Config, handler)` | Wraps a Go function as an LLM-callable tool |
| `model.LLM` | Interface implemented by each provider adapter |

**Provider integration:**

| Provider | How |
|---|---|
| Gemini | ADK native via `google.golang.org/adk/model/gemini` |
| Anthropic | Custom `model.LLM` adapter in `anthropic.go`, backed by `anthropic-sdk-go` |
| OpenAI | Custom `model.LLM` adapter in `openai.go`, backed by `openai-go` |

The Anthropic and OpenAI adapters translate ADK's `model.LLMRequest` / `model.LLMResponse` (which use `google.golang.org/genai` types internally) into the wire formats each SDK expects, then translate responses back.

**Tool registration:**

Tools are registered at agent construction time. Each tool group's constructor (`NewFileTools`, `NewGitTools`, etc.) returns `[]adktool.Tool`, which are passed directly to `llmagent.Config.Tools`. The `functiontool.New` signature is:

```go
functiontool.New(
    functiontool.Config{Name, Description, InputSchema *jsonschema.Schema},
    func(ctx adktool.Context, input map[string]any) (map[string]any, error),
)
```

The `mustSchema(raw string)` helper parses inline JSON Schema strings at startup (panics on invalid JSON — caught during development, not at runtime). The `bind(input, &params)` helper JSON-round-trips `map[string]any` into a typed params struct, handling arrays and nested types correctly.

---

## Configuration

**Config file:** `~/.config/tanka/agent.yaml`

```yaml
provider: gemini          # gemini | anthropic | openai
model: gemini-2.0-flash   # provider-specific model ID
# api_key: ""             # prefer environment variables
```

**Environment variables** (take priority over config file):

| Variable | Purpose |
|---|---|
| `TANKA_AGENT_PROVIDER` | Provider selection |
| `TANKA_AGENT_MODEL` | Model selection |
| `GEMINI_API_KEY` or `GOOGLE_API_KEY` | Gemini credentials |
| `ANTHROPIC_API_KEY` | Anthropic credentials |
| `OPENAI_API_KEY` | OpenAI credentials |

**CLI flags** (take priority over everything):

```
--provider   gemini | anthropic | openai
--model      model ID for the chosen provider
```

**Load order:** CLI flags → env vars → config file → defaults (`provider=gemini`, `model=gemini-2.0-flash`)

**Default models by provider:**

| Provider | Default model |
|---|---|
| `gemini` | `gemini-2.0-flash` |
| `anthropic` | `claude-opus-4-6` |
| `openai` | `gpt-4o` |

---

## Tools

### `file` — Sandboxed filesystem (`tools/filesystem.go`)

All paths are validated against the repo root via `safePath()`, which resolves the path and confirms it doesn't escape with `..` traversal.

| Tool | Parameters | Returns |
|---|---|---|
| `file_read` | `path`, `offset` (lines, default 0), `limit` (lines, default/max 500) | File contents (paginated) |
| `file_write` | `path`, `content` | Success message |
| `file_list` | `glob_pattern`, `offset` (default 0), `limit` (default 200) | Matching relative paths (paginated) |
| `file_search` | `glob_pattern`, `text_query`, `offset` (default 0), `limit` (default 200) | Matching lines with file:line: prefix (paginated) |
| `file_delete` | `path` | Success message |

All three read/list/search tools include an informational pagination header when results are truncated, e.g.:

```
[file.go: lines 1–500 of 1200] (700 more lines, use offset=500 to continue)
[1–200 of 350 files matching "**/*.jsonnet"] (150 more, use offset=200 to continue)
[1–200 of 412 matches for "apiVersion"] (212 more, use offset=200 to continue)
```

The double-glob `**` pattern is handled by `matchDoubleGlob()` since Go's `filepath.Match` doesn't natively support it.

### `git` — Read-only repository operations (`tools/git.go`)

Uses `go-git` (pure Go; no `git` binary required). The git tools are intentionally read-only — the agent inspects history but does not stage, commit, or manage branches. Those actions remain the user's responsibility.

| Tool | Parameters | Returns |
|---|---|---|
| `git_log` | `limit` (default 10) | Recent commit history |
| `git_show` | `ref` (default `HEAD`; accepts hash, `HEAD~1`, branch name) | Full commit details: author, date, message, and unified diff of all changed files |

### `tanka` — Tanka + Jsonnet (`tools/tanka.go`)

Calls the Tanka Go API in-process. No subprocess.

| Tool | Parameters | Returns |
|---|---|---|
| `tanka_find_environments` | `path` (use `"."` for all) | Environments with name, namespace, API server, labels |
| `tanka_show` | `env_path` | Rendered YAML manifests |
| `tanka_diff` | `env_path` | Cluster diff (requires kubectl connectivity) |
| `tanka_validate_jsonnet` | `file_path` | Lint errors/warnings or `✓ valid` |
| `tanka_format_jsonnet` | `file_path` | Formatted content (does not modify file) |

No `tanka_apply`.

### `jb` — jsonnet-bundler (`tools/jb.go`)

Calls the jsonnet-bundler library directly. No `jb` binary required. Panic recovery is wrapped around `jbpkg.Ensure()` since the jb library uses `panic()` in some error paths.

| Tool | Parameters | Returns |
|---|---|---|
| `jb_init` | `path` | Creates `jsonnetfile.json` in directory |
| `jb_install` | `path`, `packages []string` (optional) | Installs deps or adds + vendors new ones |
| `jb_update` | `path`, `packages []string` (optional) | Updates all or specific deps |
| `jb_list` | `path` | Dependencies with declared and locked versions |

---

## System Prompt

The agent is given a fixed system prompt at construction time in `agent.go`. It establishes an SRE persona with explicit behavioral constraints:

```
You are a senior Site Reliability Engineer specializing in Tanka and Jsonnet-based
Kubernetes configuration management.

Principles you always follow:
- Minimize blast radius: prefer the smallest possible scope of change. Target one
  environment at a time; don't touch shared libraries unless essential.
- Staged rollout: change one environment first, validate it, then propagate to others.
- Minimal, reviewable changes: keep diffs small and easy for humans to review.
  Avoid reformatting unrelated code.
- Always validate: after every file change, call tanka_validate_jsonnet on modified
  files, run tanka_diff for each affected environment (if cluster is reachable), then
  summarise the changes for the user before suggesting next steps.
- You cannot and will not apply to any cluster. Your job is to prepare correct,
  validated configuration changes for human review and deployment.
- You do not manage git: no staging, committing, or branch creation. That is the
  user's responsibility. You may read history with git_log and git_show.

Tanka workflow reminders:
- Environments live in subdirectories (often environments/) and have a spec.json
- Shared libraries live in lib/ or vendor/
- Always run tanka_find_environments to discover the repo structure before making changes
- After making changes: validate jsonnet → (optional) tanka_diff, then present a clear
  summary of every file changed and what was changed, so the user can review and commit
```

---

## CLI Interface

### Startup sequence (`cmd/tk/agent.go`)

1. `findGitRoot()` — open git repo with `DetectDotGit: true`; fail with clear error if not in a repo
2. `agent.LoadConfig()` — load `~/.config/tanka/agent.yaml`, apply env var overrides
3. Apply `--provider` / `--model` flag overrides
4. `cfg.Validate()` — verify API key is present; error message names the specific missing env var
5. `agent.NewModel()` — construct the `model.LLM` implementation
6. `agent.NewAgent()` — register all tools, create ADK runner + in-memory session
7. Dispatch to `RunOneShot()` or `RunREPL()` based on whether args were provided

### REPL (`pkg/agent/repl.go`)

Uses `chzyer/readline` for line editing, persistent history, and Ctrl+C handling.

- **Prompt:** `❯ `
- **History file:** `~/.config/tanka/agent_history`
- **`/exit` or Ctrl+D:** exit cleanly
- **`/clear`:** reset conversation (`Agent.Reset()` creates a new ADK session)
- **`/context`:** dump the raw ADK session event log (commit hashes, token counts, full tool I/O) for debugging
- **`/help`:** print available commands and keyboard shortcuts
- **Ctrl+C:** cancels nothing, returns a fresh prompt (the ADK run loop is synchronous; true async interrupt would require goroutine + context cancellation)

On startup the REPL prints an ASCII-art rendering of the Tanka logo in bright yellow (matching the brand's `#FACA10` gradient colour), followed by a one-line tagline.

### One-shot mode (`pkg/agent/oneshot.go`)

```
tk agent "create a new staging environment"
```

Calls `Agent.Run()` once, prints the final response, exits.

---

## Run Loop and Output

`Agent.Run()` iterates events from `runner.Run()`. Each event can carry:

- **`FunctionCall`** part — printed as `[tool: <name>] <args (truncated to 120 chars)>`
- **`FunctionResponse`** part — printed as `[tool: <name>] <output (truncated to 120 chars)>`
- **`Text`** part — if it precedes tool calls, printed immediately; if it's the final response (`event.IsFinalResponse()`), accumulated and returned

The `summarize()` helper strips newlines and truncates to a configurable length for tool call display.

### Output rendering

**Final LLM responses** are rendered as styled Markdown using `charmbracelet/glamour` (`glamour.WithAutoStyle()`, word-wrap disabled). A `*glamour.TermRenderer` is initialised in `NewAgent()` and stored on the struct; if initialisation fails the response is printed as plain text. Both `RunOneShot` and `RunREPL` pass the response through `Agent.renderMarkdown()`.

**Verbose mode** (`-v` flag) is handled in `verbose.go`:

- Tool call lines (`▶ name(args)`) are printed in bold cyan.
- Tool response blocks are wrapped in a left-border style via `lipgloss.NormalBorder()` with a dim-grey (`color "240"`) foreground.
- If the response output looks like a diff (starts with `diff `, `---`, or contains `\n+++ `), it is passed through `term.Colordiff()` before being wrapped in the border — producing coloured `+`/`-` lines.
- Intermediate LLM text is printed in yellow.

**`pkg/term/colordiff.go`** — the head-line regex was broadened from `^diff -u -N.*` to `^diff .*` so that `diff --git a/... b/...` headers are also rendered in bold blue.

---

## Provider Adapters

Both Anthropic and OpenAI adapters implement `model.LLM`:

```go
type LLM interface {
    Name() string
    GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error]
}
```

The adapters translate:
- `model.LLMRequest.Config.SystemInstruction` → provider-specific system format
- `model.LLMRequest.Contents []genai.Content` → provider message history (with role mapping)
- `model.LLMRequest.Config.Tools` → provider tool/function definitions
- Provider response → `model.LLMResponse` with `TurnComplete` flag

`genaiSchemaToMap()` (shared by both adapters, defined in `anthropic.go`) converts `*genai.Schema` (which uses uppercase Gemini type names like `"OBJECT"`, `"STRING"`) to lowercase JSON Schema maps that Anthropic and OpenAI expect.

**Known subtlety (Anthropic):** `ToolInputSchemaParam.Type` must be explicitly set to `"object"` — not the zero value. The `omitzero` tag on the field causes an empty `ToolInputSchemaParam{}` to be omitted from the serialized JSON entirely, resulting in a 400 `input_schema: Field required` error from the API.

---

## Dependencies

New dependencies added for this feature:

```
google.golang.org/adk v0.5.0           # Agent framework (loop, session, tool dispatch)
google.golang.org/genai                # Google GenAI SDK (shared types used by ADK)
github.com/anthropics/anthropic-sdk-go # Anthropic Messages API
github.com/openai/openai-go            # OpenAI Chat Completions API
github.com/go-git/go-git/v5            # Git operations (pure Go, no binary)
github.com/charmbracelet/glamour       # Markdown rendering for LLM responses
github.com/charmbracelet/lipgloss      # Terminal styling for verbose tool blocks
github.com/chzyer/readline             # REPL line editor with history
gopkg.in/yaml.v3                       # Config file parsing
```

Already present in the module (reused):

```
github.com/jsonnet-bundler/jsonnet-bundler  # jb tools use library directly
github.com/google/go-jsonnet               # Jsonnet linter (tanka_validate_jsonnet)
```

---

## Post-change Validation Workflow

Encoded in the system prompt and tool descriptions. After every file change the agent should:

1. Call `tanka_validate_jsonnet` on each modified `.jsonnet` / `.libsonnet` file
2. Call `tanka_find_environments` to identify affected environments
3. Call `tanka_diff` for each affected environment (skippable if no cluster connectivity)
4. Present a clear summary of every file changed and what was changed, so the user can review, stage, and commit themselves

The agent does not call `git_add`, `git_commit`, or create branches. Those steps are entirely the user's responsibility.

---

## What Was Explicitly Left Out

- **`tanka_apply`** — by design; cluster deploys are human-controlled
- **Streaming output** — ADK's `GenerateContent` receives `stream=false`; non-streaming is simpler and sufficient for a CLI tool
- **Persistent session history** — ephemeral by design; each invocation starts fresh
- **Git write operations** (`git_add`, `git_commit`, `git_branch_create`, `git_branch_checkout`, `git_status`, `git_diff`) — the agent prepares changes for human review; committing and branching are deliberately left to the user
- **GitHub tools** (`github_push`, `github_pr_create`, `github_pr_list`) — removed; the agent's scope is local file editing and validation, not remote git or PR management
- **Worktree / stash / merge / rebase** — not needed for the prepare-and-review workflow
