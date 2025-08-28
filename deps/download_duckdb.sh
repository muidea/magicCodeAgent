#!/usr/bin/bash

DUCKDB_VERSION=${1}
TARGET="$(go env GOOS)-$(go env GOARCH)"
WORK_DIR="$(cd "$(dirname "$0")" && pwd)"
TARGET_DIR=${WORK_DIR}/${TARGET}
PRE_DIR=$(pwd)

if [ -d ${TARGET_DIR} ]; then
    if [ -f ${TARGET_DIR}/libduckdb_bundle.a ]; then
        echo "The '${TARGET}' directory already exists. Skipping download." 1>&2
        exit 0
    fi    
    rm -rf ${TARGET_DIR}
fi

cd ${WORK_DIR}

version_less_equal() {
  # Strip leading 'v' if present
  local ver1=${1#v}
  local ver2=${2#v}

  IFS='.' read -r -a v1 <<< "$ver1"
  IFS='.' read -r -a v2 <<< "$ver2"

  for i in 0 1 2; do
    local a=${v1[i]:-0}
    local b=${v2[i]:-0}
    if ((a < b)); then
      return 0
    elif ((a > b)); then
      return 1
    fi
  done
  return 0
}

# Until version 1.2.1, the bundle was provided as a single static library
if version_less_equal $DUCKDB_VERSION "1.2.1"; then
    wget https://github.com/duckdb/duckdb/releases/download/${DUCKDB_VERSION}/static-lib-${TARGET}.zip \
        && unzip -o static-lib-${TARGET}.zip -d ${TARGET_DIR} \
        && rm static-lib-${TARGET}.zip

    echo "The 'libduckdb_bundle.a' file is now available in the current directory." 1>&2
    cd ${PRE_DIR}
    exit 0
fi

wget https://github.com/duckdb/duckdb/releases/download/${DUCKDB_VERSION}/static-libs-${TARGET}.zip \
        && unzip -o static-libs-${TARGET}.zip -d ${TARGET_DIR} \
        && rm static-libs-${TARGET}.zip

cd ${TARGET_DIR}
find . -name '*.a' -exec mkdir -p {}.objects \; -exec mv {} {}.objects \;
find . -name '*.a' -execdir ar -x {} \;
echo ./*/*.o | xargs ar cr libduckdb_bundle.a

echo "The 'libduckdb_bundle.a' file is now available in the current directory." 1>&2
cd ${PRE_DIR}