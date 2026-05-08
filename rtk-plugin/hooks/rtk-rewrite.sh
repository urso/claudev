#!/bin/bash
if ! command -v rtk &>/dev/null; then
  exit 0
fi

rtk hook claude
