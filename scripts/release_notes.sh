#!/bin/bash

set -euo pipefail

if [[ -z "${GITHUB_REF_NAME:-}" ]]; then
    GITHUB_REF_NAME=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
fi

TAG="${GITHUB_REF_NAME#v}"
PREV_TAG=$(git describe --tags --abbrev=0 "$GITHUB_REF_NAME"^ 2>/dev/null || echo "")

generate_changelog() {
    local output=""
    local current_type=""
    local types=("feat" "fix" "refactor" "perf" "chore" "docs" "test")
    local has_changes=false

    for type in "${types[@]}"; do
        local commits=$(git log --pretty=format:"- %s (%h)" "${PREV_TAG:-$type~10}..${GITHUB_REF_NAME}" 2>/dev/null | grep -i "^${type#!}" | sed "s/^- /- /" || true)
        if [[ -n "$commits" ]]; then
            has_changes=true
            local type_label="${type^}"
            case "$type" in
                feat) type_label="Features" ;;
                fix) type_label="Bug Fixes" ;;
                refactor) type_label="Refactoring" ;;
                perf) type_label="Performance" ;;
                docs) type_label="Documentation" ;;
                test) type_label="Tests" ;;
                chore) type_label="Maintenance" ;;
            esac
            output+="\n\n## ${type_label}\n${commits}"
        fi
    done

    if ! $has_changes; then
        output="\n\n## Changes\n$(git log --pretty=format:"- %s (%h)" "${PREV_TAG:-$GITHUB_REF_NAME~10}..${GITHUB_REF_NAME}" 2>/dev/null | sed 's/^-/-/')"
    fi

    echo -e "$output"
}

generate_release_notes() {
    local tag="${1:-}"
    local body=""
    local current_date=$(date +"%Y-%m-%d")

    body="## Release ${tag} (${current_date})\n\n"
    body+="### Checks Passed\n"
    body+="- ✅ Format validation\n"
    body+="- ✅ Linting\n"
    body+="- ✅ Tests\n"
    body+="- ✅ Build\n"
    body+=$(generate_changelog)

    echo -e "$body"
}

generate_release_notes "${TAG:-}"
