#!/bin/sh

set -euo pipefail

/scripts/download-tools.sh

exec /kitops s
