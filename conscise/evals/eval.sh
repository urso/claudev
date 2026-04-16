#!/usr/bin/env bash
# Eval script for conscise plugin
# Compares baseline vs hook vs skill-invoked responses

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_DIR="$(dirname "$SCRIPT_DIR")"
SKILL_ONLY_DIR="/tmp/conscise-skill-only-$$"

# Prompts to test
PROMPTS=(
  "What are the pros and cons of microservices?"
  "Explain dependency injection"
  "Why use TypeScript over JavaScript?"
)

# Runs per config
RUNS=3

# Filler words to detect (case insensitive)
FILLER_WORDS="just|really|basically|actually|simply|essentially|certainly|of course|happy to|let me explain|here's what"

# Pleasantry patterns
PLEASANTRIES="sure|great question|glad you asked|absolutely|definitely"

# Setup skill-only plugin (no hooks)
setup_skill_only() {
  mkdir -p "$SKILL_ONLY_DIR/.claude-plugin"
  cp "$PLUGIN_DIR/.claude-plugin/plugin.json" "$SKILL_ONLY_DIR/.claude-plugin/"
  cp -r "$PLUGIN_DIR/skills" "$SKILL_ONLY_DIR/"
}

# cleanup handled in main after TMPDIR is created

run_baseline() {
  local prompt="$1"
  claude -p "$prompt" --setting-sources "" 2>/dev/null
}

run_hook() {
  local prompt="$1"
  claude -p "$prompt" --setting-sources "" --plugin-dir "$PLUGIN_DIR" 2>/dev/null
}

run_skill() {
  local prompt="$1"
  claude -p "/conscise:conscise then answer: $prompt" --setting-sources "" --plugin-dir "$SKILL_ONLY_DIR" 2>/dev/null
}

count_filler() {
  local text="$1"
  echo "$text" | grep -ioE "\b($FILLER_WORDS)\b" | wc -l | tr -d ' '
}

count_pleasantries() {
  local text="$1"
  echo "$text" | grep -ioE "\b($PLEASANTRIES)\b" | wc -l | tr -d ' '
}

avg_sentence_length() {
  local text="$1"
  # Split on . ! ? and count words per sentence
  local sentences words avg
  # Count sentences (periods, exclamations, questions not in code blocks)
  sentences=$(echo "$text" | grep -oE '[.!?]' | wc -l | tr -d ' ')
  words=$(echo "$text" | wc -w | tr -d ' ')
  if [[ $sentences -eq 0 ]]; then
    echo "0"
  else
    avg=$((words / sentences))
    echo "$avg"
  fi
}

count_arrows() {
  local text="$1"
  echo "$text" | grep -oE '→|->|=>' | wc -l | tr -d ' '
}

analyze() {
  local text="$1"
  local lines words chars filler pleasantries avg_sent arrows
  lines=$(echo "$text" | wc -l | tr -d ' ')
  words=$(echo "$text" | wc -w | tr -d ' ')
  chars=$(echo "$text" | wc -c | tr -d ' ')
  filler=$(count_filler "$text")
  pleasantries=$(count_pleasantries "$text")
  avg_sent=$(avg_sentence_length "$text")
  arrows=$(count_arrows "$text")
  echo "$lines,$words,$chars,$filler,$pleasantries,$avg_sent,$arrows"
}

# Aggregate metrics across runs
aggregate() {
  local -a values=("$@")
  local sum=0 count=${#values[@]}
  for v in "${values[@]}"; do
    sum=$((sum + v))
  done
  echo $((sum / count))
}

# Run single eval and write results to file
run_single_eval() {
  local mode="$1" prompt="$2" outfile="$3"
  local result
  case "$mode" in
    baseline) result=$(run_baseline "$prompt") ;;
    hook)     result=$(run_hook "$prompt") ;;
    skill)    result=$(run_skill "$prompt") ;;
  esac
  analyze "$result" > "$outfile"
}

# Main
setup_skill_only

TMPDIR=$(mktemp -d)
trap "rm -rf '$TMPDIR' '$SKILL_ONLY_DIR'" EXIT

echo "=== Conscise Plugin Eval (${RUNS} runs per config, parallel) ==="
echo ""

# Launch all runs in parallel
pids=()
for pi in "${!PROMPTS[@]}"; do
  prompt="${PROMPTS[$pi]}"
  for mode in baseline hook skill; do
    for ((r=1; r<=RUNS; r++)); do
      outfile="$TMPDIR/p${pi}_${mode}_${r}.csv"
      run_single_eval "$mode" "$prompt" "$outfile" &
      pids+=($!)
    done
  done
done

# Wait for all (ignore individual failures)
for pid in "${pids[@]}"; do
  wait "$pid" || true
done

echo "| Prompt | Mode | Words | Filler | Pleasantries | Avg Sent | Arrows |"
echo "|--------|------|-------|--------|--------------|----------|--------|"

# Collect results
for pi in "${!PROMPTS[@]}"; do
  prompt="${PROMPTS[$pi]}"
  short_prompt="${prompt:0:35}..."

  for mode in baseline hook skill; do
    words_arr=() filler_arr=() pleas_arr=() avg_arr=() arrows_arr=()

    for ((r=1; r<=RUNS; r++)); do
      outfile="$TMPDIR/p${pi}_${mode}_${r}.csv"
      if [[ -f "$outfile" ]]; then
        IFS=',' read -r _ words _ filler pleas avg arrows < "$outfile"
        words_arr+=("$words") filler_arr+=("$filler") pleas_arr+=("$pleas") avg_arr+=("$avg") arrows_arr+=("$arrows")
      fi
    done

    if [[ "$mode" == "baseline" ]]; then
      echo "| $short_prompt | $mode | $(aggregate "${words_arr[@]}") | $(aggregate "${filler_arr[@]}") | $(aggregate "${pleas_arr[@]}") | $(aggregate "${avg_arr[@]}") | $(aggregate "${arrows_arr[@]}") |"
    else
      echo "| | $mode | $(aggregate "${words_arr[@]}") | $(aggregate "${filler_arr[@]}") | $(aggregate "${pleas_arr[@]}") | $(aggregate "${avg_arr[@]}") | $(aggregate "${arrows_arr[@]}") |"
    fi
  done
done

echo ""
echo "Values are averages across $RUNS runs"
echo "Good concise output: fewer words, zero filler/pleasantries, shorter sentences, more arrows"
