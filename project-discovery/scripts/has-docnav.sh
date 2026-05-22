#!/bin/bash
# Check if docnav plugin is installed
# Returns: "found" or "not found"

claude plugin list 2>/dev/null | grep -q docnav && echo "found" || echo "not found"
