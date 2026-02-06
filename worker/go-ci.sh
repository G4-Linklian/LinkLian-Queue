#!/bin/bash
set -euo pipefail

MODE="${1:-all}" # unit | integration | all

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

add_test() {
  local name="$1"
  test_order+=("$name")
}

mark_passed() {
  local name="$1"
  test_passed+=("$name")
  echo -e "${GREEN}    ✅ PASSED${NC}"
}

mark_failed() {
  local name="$1"
  test_failed+=("$name")
  echo -e "${RED}    ❌ FAILED${NC}"
}

test_status() {
  local test_name="$1"
  local i
  
  # Check if passed
  for i in "${test_passed[@]}"; do
    [[ "$i" == "$test_name" ]] && echo "passed" && return 0
  done
  
  # Check if failed
  for i in "${test_failed[@]}"; do
    [[ "$i" == "$test_name" ]] && echo "failed" && return 0
  done
  
  echo "pending"
}

print_summary() {
  echo
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo -e "${BLUE}  TEST SUMMARY${NC}"
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  
  local test
  for test in "${test_order[@]}"; do
    local status
    status=$(test_status "$test")
    case "$status" in
      passed)
        echo -e "${GREEN}✅ PASSED${NC}  | $test"
        ;;
      failed)
        echo -e "${RED}❌ FAILED${NC}  | $test"
        ;;
      pending)
        echo -e "${YELLOW}⏭️  SKIPPED${NC} | $test"
        ;;
    esac
  done
  
  local passed_count=${#test_passed[@]}
  local failed_count=${#test_failed[@]}
  local total=${#test_order[@]}
  
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo -e "Total: $total | ${GREEN}Passed: $passed_count${NC} | ${RED}Failed: $failed_count${NC}"
  echo -e "${BLUE}════════════════════════════════════════${NC}"
  echo
  
  if [[ $failed_count -gt 0 ]]; then
    return 1
  fi
  return 0
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

need_cmd() { command -v "$1" >/dev/null 2>&1 || die "$1" "Tool required" "" "Install '$1' then retry"; }

need_env() {
  local var="$1" example="${2:-}"
  [[ -n "${!var:-}" ]] || die "integration tests" "Environment check" "" \
"export $var='$example'
# then rerun:
./go-ci.sh integration"
}

# ---------- docker compose helpers ----------
compose_file_check() {
  [[ -f "docker-compose.yml" || -f "compose.yml" || -f "compose.yaml" ]] || die "docker compose" "Compose file check" \
"No docker compose file found in project root." \
"Ensure docker-compose.yml (or compose.yml/compose.yaml) exists, then rerun:
./go-ci.sh integration"
}

compose_services() {
  docker compose config --services 2>/dev/null | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g' | sed 's/^ //; s/ $//'
}

pick_service() {
  local re="$1"
  local s
  for s in $(docker compose config --services 2>/dev/null); do
    if [[ "$s" =~ $re ]]; then
      echo "$s"
      return 0
    fi
  done
  echo ""
}

integration_fix_hints() {
  local svcs="$1"
  local pg="$2"
  local mq="$3"

  local target=""
  if [[ -n "$pg" || -n "$mq" ]]; then
    target="${pg} ${mq}"
    target="$(echo "$target" | sed 's/[[:space:]]\+/ /g' | sed 's/^ //; s/ $//')"
  else
    target="$svcs"
  fi

  cat <<EOF
docker compose ps
docker compose logs --tail=200 ${target}

# If a service is unhealthy / crash-looping:
docker compose restart ${target}
docker compose ps
docker compose logs --tail=200 ${target}

# Full reset (drops volumes):
docker compose down -v
docker compose up -d
EOF
}

# ---------- preflight ----------
need_cmd go

# ---------- register tests ----------
add_test "gofmt (style)"
add_test "go vet (static analysis)"
add_test "go mod download"
add_test "go build"
add_test "go test (unit)"

if [[ "$MODE" == "integration" || "$MODE" == "all" ]]; then
  add_test "docker compose up"
  add_test "docker compose status"
  add_test "go test (integration)"
  add_test "docker compose down"
fi

# ---------- steps ----------
step "gofmt (style)" "Checks Go code formatting (spacing, imports, indentation)."
files="$(gofmt -l . || true)"
if [[ -n "$files" ]]; then
  die "gofmt" "Code style / formatting" \
"These files are not formatted:
$files" \
"gofmt -w .
# then rerun:
./go-ci.sh $MODE"
fi
mark_passed "gofmt (style)"

step "go vet (static analysis)" "Catches common bugs (printf mismatch, suspicious constructs)."
run "go vet ./..." "go vet" "Static analysis" \
"go vet ./...
# Fix the issues shown, then rerun:
./go-ci.sh $MODE"
mark_passed "go vet (static analysis)"

step "go mod download" "Downloads and verifies module dependencies."
run "go mod download" "go mod download" "Dependency resolution" \
"go mod download
# If mod cache looks corrupted:
go clean -modcache && go mod download
# then rerun:
./go-ci.sh $MODE"
mark_passed "go mod download"

step "go build" "Compiles packages to ensure the project builds."
run "go build -v ./..." "go build" "Compilation" \
"go build ./...
# Fix compile errors above, then rerun:
./go-ci.sh $MODE"
mark_passed "go build"

step "go test (unit)" "Runs unit tests (*_test.go)."
run "go test -v ./... -count=1" "unit tests" "Unit tests" \
"go test -v ./... -count=1
# Fix failing tests, then rerun:
./go-ci.sh $MODE"
mark_passed "go test (unit)"

# ---------- integration (docker compose + tags=integration) ----------
if [[ "$MODE" == "integration" || "$MODE" == "all" ]]; then
  need_cmd docker
  compose_file_check

  # auto-detect service names (best-effort)
  svcs="$(compose_services)"
  [[ -n "$svcs" ]] || die "docker compose" "Compose config" \
"Could not read services from compose." \
"Try:
docker compose config
# fix compose errors, then rerun:
./go-ci.sh integration"

  postgres_svc="$(pick_service '(?i)(postgres|postg|pg|db)')"
  rabbit_svc="$(pick_service '(?i)(rabbit|mq)')"

  fix_compose="$(integration_fix_hints "$svcs" "$postgres_svc" "$rabbit_svc")"

  step "docker compose up" "Starts services from docker-compose.yml for integration tests."
  run "docker compose up -d" "docker compose up" "Start services" "$fix_compose"
  mark_passed "docker compose up"

  step "docker compose status" "Shows service status (useful when tests fail)."
  run "docker compose ps" "docker compose ps" "Service status" "$fix_compose"
  mark_passed "docker compose status"

  step "go test -tags=integration" "Runs integration tests against REAL services."
  need_env DATABASE_URL "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
  need_env RABBITMQ_URL "amqp://guest:guest@localhost:5672/"

  run "go test -v -tags=integration ./... -count=1" \
      "integration tests" "DB / RabbitMQ integration" \
"$fix_compose
# Rerun only integration:
./go-ci.sh integration"
  mark_passed "go test (integration)"

  step "docker compose down" "Stops services after tests (cleanup)."
  if ! docker compose down >/dev/null 2>&1; then
    echo "    ⚠️ cleanup failed; run manually:"
    echo "       docker compose down -v"
    mark_failed "docker compose down"
  else
    mark_passed "docker compose down"
  fi
fi

# ---------- print summary ----------
if print_summary; then
  exit 0
else
  exit 1
fi
