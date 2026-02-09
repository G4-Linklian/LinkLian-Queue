#!/bin/bash
set -euo pipefail

MODE="${1:-all}" # help | secret | sca | sast | security | format | vet | deps | build | test | all | full

# ---------- change to worker directory if not already there ----------
if [[ ! -f "go.mod" ]]; then
  if [[ -f "worker/go.mod" ]]; then
    cd worker
  else
    echo "❌ Could not find go.mod. Make sure you're in the worker directory or project root."
    exit 1
  fi
fi

# ---------- color codes ----------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ---------- status tracking (bash 3.2 compatible) ----------
declare -a test_order
declare -a test_passed
declare -a test_failed

add_test() { test_order+=("$1"); }

mark_passed() {
  test_passed+=("$1")
  echo -e "${GREEN}    ✅ PASSED${NC}"
}

mark_failed() {
  test_failed+=("$1")
  echo -e "${RED}    ❌ FAILED${NC}"
}

test_status() {
  local name="$1" i

  for i in ${test_passed[@]+"${test_passed[@]}"}; do
    [[ "$i" == "$name" ]] && echo "passed" && return 0
  done

  for i in ${test_failed[@]+"${test_failed[@]}"}; do
    [[ "$i" == "$name" ]] && echo "failed" && return 0
  done

  echo "pending"
}

print_summary() {
  echo
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo -e "${BLUE}  TEST SUMMARY${NC}"
  echo -e "${BLUE}════════════════════════════════════════${NC}"

  local t status
  for t in "${test_order[@]}"; do
    status="$(test_status "$t")"
    case "$status" in
      passed)  echo -e "${GREEN}✅ PASSED${NC}  | $t" ;;
      failed)  echo -e "${RED}❌ FAILED${NC}  | $t" ;;
      pending) echo -e "${YELLOW}⏭️  SKIPPED${NC} | $t" ;;
    esac
  done

  local passed_count=${#test_passed[@]}
  local failed_count=${#test_failed[@]}
  local total=${#test_order[@]}

  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo -e "Total: $total | ${GREEN}Passed: $passed_count${NC} | ${RED}Failed: $failed_count${NC}"
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo

  [[ $failed_count -eq 0 ]]
}

# ---------- helpers ----------
step() {
  echo
  echo -e "${BLUE}==> $1${NC}"
  echo "    $2"
}

die() {
  echo
  echo -e "${RED}❌ FAIL at: $1${NC}"
  echo "    What it checks: $2"
  [[ -n "${3:-}" ]] && { echo "    Details:"; echo "$3"; }
  [[ -n "${4:-}" ]] && { echo; echo "    👉 Fix (run these):"; echo "$4"; }
  exit 1
}

run() {
  # run "<cmd>" "<step>" "<what>" "<fix>"
  local cmd="$1" name="$2" what="$3" fix="${4:-}"
  local out; out="$(mktemp)"
  bash -lc "$cmd" >"$out" 2>&1 || die "$name" "$what" "$(cat "$out")" "$fix"
  rm -f "$out"
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "$1" "Tool required" "" "Install '$1' then retry"
}

# Install a Go-based tool if missing (puts binary in GOPATH/bin)
ensure_go_tool() {
  # ensure_go_tool <binary> <module@version>
  local bin="$1" mod="$2"
  local gopath gobin

  gopath="$(go env GOPATH 2>/dev/null || true)"
  [[ -z "$gopath" ]] && die "go env GOPATH" "Go env" "GOPATH is empty" "Check Go install and env"

  gobin="${GOBIN:-$gopath/bin}"
  mkdir -p "$gobin"

  if command -v "$bin" >/dev/null 2>&1; then
    return 0
  fi

  if [[ -x "$gobin/$bin" ]]; then
    export PATH="$gobin:$PATH"
    return 0
  fi

  step "Install $bin" "Missing '$bin' — installing via: go install $mod"
  run "go install $mod" "go install $bin" "Tool installation" \
"go install $mod
# then rerun:
./go-ci.sh $MODE"

  export PATH="$gobin:$PATH"

  command -v "$bin" >/dev/null 2>&1 || die "$bin" "Tool installation" \
"Installed but still not found in PATH (expected in $gobin)" \
"Add to PATH:
export PATH=\"$gobin:\$PATH\""
}

usage() {
  cat <<'EOF'
Usage:
  ./go-ci.sh [mode]

Modes:
  help        Show this help
  secret      Secret scan (gitleaks)
  sca         Dependency vuln scan (govulncheck)
  sast        SAST scan (gosec)
  security    runs: secret -> sca -> sast

  format      gofmt check (fails if changes needed)
  vet         go vet ./...
  deps        go mod download
  build       go build ./...
  test        go test ./...
  all         runs: format -> vet -> deps -> build -> test
  full        runs: security -> all

Examples:
  chmod +x go-ci.sh
  ./go-ci.sh secret
  ./go-ci.sh sca
  ./go-ci.sh sast
  ./go-ci.sh security
  ./go-ci.sh all
  ./go-ci.sh full
EOF
}

if [[ "$MODE" == "help" || "$MODE" == "-h" || "$MODE" == "--help" ]]; then
  usage
  exit 0
fi

# ---------- preflight ----------
need_cmd go

# ---------- register tests ----------
add_test "gitleaks (secret scan)"
add_test "govulncheck (SCA)"
add_test "gosec (SAST)"
add_test "gofmt (style)"
add_test "go vet (static analysis)"
add_test "go mod download"
add_test "go build"
add_test "go test (unit)"

# ---------- step functions ----------
run_secret() {
  step "gitleaks (secret scan)" "Detects leaked secrets (API keys, tokens, passwords) in the working tree."
  ensure_go_tool "gitleaks" "github.com/gitleaks/gitleaks/v8@latest"
  run "gitleaks detect --redact --no-git" "gitleaks" "Secret scanning" \
"gitleaks detect --redact --no-git
# Remove secrets, rotate keys, then rerun:
./go-ci.sh secret"
  mark_passed "gitleaks (secret scan)"
}

run_sca() {
  step "govulncheck (SCA)" "Checks known vulnerabilities in Go module dependencies (go.mod/go.sum)."
  ensure_go_tool "govulncheck" "golang.org/x/vuln/cmd/govulncheck@latest"
  run "govulncheck ./..." "govulncheck" "Dependency vulnerability scanning" \
"govulncheck ./...
# Upgrade/replace vulnerable deps, then rerun:
./go-ci.sh sca"
  mark_passed "govulncheck (SCA)"
}

run_sast() {
  step "gosec (SAST)" "Analyzes Go code for common security issues (crypto, injection patterns, file perms, etc.)."
  ensure_go_tool "gosec" "github.com/securego/gosec/v2/cmd/gosec@latest"
  run "gosec ./..." "gosec" "Static application security testing" \
"gosec ./...
# Fix findings or configure rules, then rerun:
./go-ci.sh sast"
  mark_passed "gosec (SAST)"
}

run_format() {
  step "gofmt (style)" "Checks Go code formatting (spacing, imports, indentation)."
  local files
  files="$(gofmt -l . || true)"
  if [[ -n "$files" ]]; then
    die "gofmt" "Code style / formatting" \
"These files are not formatted:
$files" \
"gofmt -w .
# then rerun:
./go-ci.sh format"
  fi
  mark_passed "gofmt (style)"
}

run_vet() {
  step "go vet (static analysis)" "Catches common bugs (printf mismatch, suspicious constructs)."
  run "go vet ./..." "go vet" "Static analysis" \
"go vet ./...
# Fix the issues shown, then rerun:
./go-ci.sh vet"
  mark_passed "go vet (static analysis)"
}

run_deps() {
  step "go mod download" "Downloads and verifies module dependencies."
  run "go mod download" "go mod download" "Dependency resolution" \
"go mod download
# If mod cache looks corrupted:
go clean -modcache && go mod download
# then rerun:
./go-ci.sh deps"
  mark_passed "go mod download"
}

run_build() {
  step "go build" "Compiles packages to ensure the project builds."
  run "go build -v ./..." "go build" "Compilation" \
"go build ./...
# Fix compile errors above, then rerun:
./go-ci.sh build"
  mark_passed "go build"
}

run_test() {
  step "go test (unit)" "Runs unit tests (*_test.go)."
  run "go test -v ./... -count=1" "go test" "Unit tests" \
"go test -v ./... -count=1
# Fix failing tests, then rerun:
./go-ci.sh test"
  mark_passed "go test (unit)"
}

run_security() {
  run_secret
  run_sca
  run_sast
}

# ---------- run selected mode ----------
case "$MODE" in
  secret)
    run_secret
    ;;
  sca)
    run_sca
    ;;
  sast)
    run_sast
    ;;
  security)
    run_security
    ;;
  format)
    run_format
    ;;
  vet)
    run_vet
    ;;
  deps)
    run_deps
    ;;
  build)
    run_build
    ;;
  test)
    run_test
    ;;
  all)
    run_format
    run_vet
    run_deps
    run_build
    run_test
    ;;
  full)
    run_security
    run_format
    run_vet
    run_deps
    run_build
    run_test
    ;;
  *)
    echo "Unknown mode: $MODE"
    usage
    exit 2
    ;;
esac

# ---------- print summary ----------
if print_summary; then
  exit 0
else
  exit 1
fi
