#!/usr/bin/env bash

set -euo pipefail

VERSION_FILE="VERSION"
BUMP_TYPE="${1:-}"

if [[ "$BUMP_TYPE" != "patch" && "$BUMP_TYPE" != "minor" && "$BUMP_TYPE" != "major" ]]; then
  echo "Usage: $0 {patch|minor|major}" >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Erreur : des changements ne sont pas commités. Commit ou stash d'abord." >&2
  git status --short
  exit 1
fi

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$CURRENT_BRANCH" != "main" ]]; then
  read -r -p "⚠️  Tu es sur '$CURRENT_BRANCH', pas 'main'. Continuer quand même ? [y/N] " CONFIRM
  if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 1
  fi
fi

if [[ ! -f "$VERSION_FILE" ]]; then
  echo "0.0.0" > "$VERSION_FILE"
fi

CURRENT_VERSION="$(cat "$VERSION_FILE" | tr -d '[:space:]')"

if [[ ! "$CURRENT_VERSION" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  echo "Erreur : contenu de VERSION invalide ('$CURRENT_VERSION'), attendu MAJOR.MINOR.PATCH" >&2
  exit 1
fi

MAJOR="${BASH_REMATCH[1]}"
MINOR="${BASH_REMATCH[2]}"
PATCH="${BASH_REMATCH[3]}"

case "$BUMP_TYPE" in
  patch)
    PATCH=$((PATCH + 1))
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
NEW_TAG="v${NEW_VERSION}"

if git rev-parse "$NEW_TAG" >/dev/null 2>&1; then
  echo "Erreur : le tag $NEW_TAG existe déjà." >&2
  exit 1
fi

echo "$NEW_VERSION" > "$VERSION_FILE"

git add "$VERSION_FILE"
git commit -m "chore(release): $NEW_TAG"
git tag -a "$NEW_TAG" -m "$NEW_TAG"

echo ""
echo "✅ Version bump : $CURRENT_VERSION -> $NEW_VERSION"
echo "   Commit créé et tag $NEW_TAG posé localement."
echo ""
echo "📤 Push de '$CURRENT_BRANCH' et du tag $NEW_TAG..."

git push origin "$CURRENT_BRANCH" --follow-tags

echo "✅ Publié sur origin."
