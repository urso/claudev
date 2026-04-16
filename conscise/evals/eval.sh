#!/usr/bin/env bash
# Eval script for conscise plugin
# Compares baseline vs hook vs skill-invoked responses

set -euo pipefail

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

cleanup() {
  rm -rf "$SKILL_ONLY_DIR"
}
trap cleanup EXIT

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

# Main
setup_skill_only

echo "=== Conscise Plugin Eval (${RUNS} runs per config) ==="
echo ""
echo "| Prompt | Mode | Words | Filler | Pleasantries | Avg Sent | Arrows |"
echo "|--------|------|-------|--------|--------------|----------|--------|"

for prompt in "${PROMPTS[@]}"; do
  short_prompt="${prompt:0:35}..."

  # Arrays for aggregation
  b_words_arr=() b_filler_arr=() b_pleas_arr=() b_avg_arr=() b_arrows_arr=()
  h_words_arr=() h_filler_arr=() h_pleas_arr=() h_avg_arr=() h_arrows_arr=()
  s_words_arr=() s_filler_arr=() s_pleas_arr=() s_avg_arr=() s_arrows_arr=()

  for ((i=1; i<=RUNS; i++)); do
    baseline=$(run_baseline "$prompt")
    hook=$(run_hook "$prompt")
    skill=$(run_skill "$prompt")

    IFS=',' read -r _ b_words _ b_filler b_pleas b_avg b_arrows <<< "$(analyze "$baseline")"
    IFS=',' read -r _ h_words _ h_filler h_pleas h_avg h_arrows <<< "$(analyze "$hook")"
    IFS=',' read -r _ s_words _ s_filler s_pleas s_avg s_arrows <<< "$(analyze "$skill")"

    b_words_arr+=("$b_words") b_filler_arr+=("$b_filler") b_pleas_arr+=("$b_pleas") b_avg_arr+=("$b_avg") b_arrows_arr+=("$b_arrows")
    h_words_arr+=("$h_words") h_filler_arr+=("$h_filler") h_pleas_arr+=("$h_pleas") h_avg_arr+=("$h_avg") h_arrows_arr+=("$h_arrows")
    s_words_arr+=("$s_words") s_filler_arr+=("$s_filler") s_pleas_arr+=("$s_pleas") s_avg_arr+=("$s_avg") s_arrows_arr+=("$s_arrows")
  done

  echo "| $short_prompt | baseline | $(aggregate "${b_words_arr[@]}") | $(aggregate "${b_filler_arr[@]}") | $(aggregate "${b_pleas_arr[@]}") | $(aggregate "${b_avg_arr[@]}") | $(aggregate "${b_arrows_arr[@]}") |"
  echo "| | hook | $(aggregate "${h_words_arr[@]}") | $(aggregate "${h_filler_arr[@]}") | $(aggregate "${h_pleas_arr[@]}") | $(aggregate "${h_avg_arr[@]}") | $(aggregate "${h_arrows_arr[@]}") |"
  echo "| | skill | $(aggregate "${s_words_arr[@]}") | $(aggregate "${s_filler_arr[@]}") | $(aggregate "${s_pleas_arr[@]}") | $(aggregate "${s_avg_arr[@]}") | $(aggregate "${s_arrows_arr[@]}") |"
done

echo ""
echo "Values are averages across $RUNS runs"
echo "Good concise output: fewer words, zero filler/pleasantries, shorter sentences, more arrows"
