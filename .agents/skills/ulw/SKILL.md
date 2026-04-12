---
name: ulw
description: Ultrawork Worktree Launch — derives a conventional-commit branch name from current session context, creates a git worktree inside ./worktrees/<branch>, symlinks node_modules, copies .env, and continues implementation in the same terminal. Project-specific command for TwirChat.
license: MIT
compatibility: opencode
---

# ULW — Ultrawork Worktree Launch

**Act immediately on these instructions. Do not ask for confirmation before executing.**

---

## What this command does

1. Reads the current task from this session's conversation history
2. Generates a conventional-commit-style branch name
3. Creates a git worktree at `./worktrees/<branch>` (inside the project, gitignored)
4. Symlinks `node_modules`, copies `.env`, and desktop `.env`
5. Continues ALL implementation work inside the worktree — **same terminal, no new window**

---

## Step 1 — Understand Current Task

Analyze the conversation in this session. Identify:

- The type of change (feature / bug fix / refactor / etc.)
- The area or module being changed (e.g., `twitch-auth`, `chat-aggregator`, `youtube-adapter`, `overlay`)
- A concise, specific 2–5 word description of the change

If there is no prior conversation context, use `chore/new-session`.

---

## Step 2 — Generate Branch Name

Format: `<type>/<short-description>` or `<type>/<scope>/<short-description>`

| Type       | When to use                           |
| ---------- | ------------------------------------- |
| `feat`     | New feature or capability             |
| `fix`      | Bug fix                               |
| `refactor` | Restructuring without behavior change |
| `docs`     | Documentation only                    |
| `chore`    | Dependencies, configs, maintenance    |
| `test`     | Tests only                            |
| `style`    | CSS / formatting changes              |
| `ci`       | CI/CD pipeline                        |

**Rules:** lowercase, kebab-case, no special chars except `-` and `/`, max 60 chars.

**Examples for TwirChat:**

- `feat/youtube-live-chat-reconnect`
- `fix/kick-oauth-token-expiry`
- `refactor/chat-aggregator-dedup`
- `feat/overlay/custom-font-support`

---

## Step 3 — Create Worktree

Run these bash commands (replace `<branch>` with the generated name):

```bash
REPO_ROOT=$(git rev-parse --show-toplevel)
BRANCH="<branch>"
WORKTREE="$REPO_ROOT/worktrees/$BRANCH"

# Create worktree directory and branch
mkdir -p "$(dirname "$WORKTREE")"
git worktree add "$WORKTREE" -b "$BRANCH"

# Symlink node_modules (saves disk space)
ln -s "$REPO_ROOT/node_modules" "$WORKTREE/node_modules"

# Copy .env if it exists
[ -f "$REPO_ROOT/.env" ] && cp "$REPO_ROOT/.env" "$WORKTREE/.env"
```

---

## Step 4 — Continue Implementation

After the worktree is created:

- Use `workdir="<absolute-worktree-path>"` for **every** bash command
- Use absolute paths under `./worktrees/<branch>/` for **every** file read/write/edit
- Do NOT touch files in the original repo directory

Then proceed with `/ulw-loop` to implement the task.

---

## Step 5 — Cleanup (when done)

```bash
git worktree remove "$REPO_ROOT/worktrees/<branch>" --force
git branch -d "<branch>"
```

---

## Constraints

- **Never ask for confirmation** — act immediately
- **Never ask for the branch name** — derive it from context
- **Always use the worktree path** for all file operations after Step 3
