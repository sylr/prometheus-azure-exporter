#!/usr/bin/env bash

setup_git() {
  git config --global user.email "support@travis-ci.org"
  git config --global user.name "Travis"
}

vendor() {
    make vendor
}

git_checkout() {
  if [[ "$TRAVIS_EVENT_TYPE" != "pull_request" ]]; then
    exit 0
  fi

  if [[ "$TRAVIS_PULL_REQUEST_BRANCH" != dependabot* ]]; then
    exit 0
  fi

  git remote set-url origin git@github.com:sylr/prometheus-azure-exporter.git
  git fetch origin +refs/heads/$TRAVIS_PULL_REQUEST_BRANCH:refs/remotes/origin/$TRAVIS_PULL_REQUEST_BRANCH
  git checkout -b $TRAVIS_PULL_REQUEST_BRANCH origin/$TRAVIS_PULL_REQUEST_BRANCH
  git push --set-upstream origin $TRAVIS_PULL_REQUEST_BRANCH
}

push_back() {
  if [[ "$TRAVIS_EVENT_TYPE" != "pull_request" ]]; then
    exit 0
  fi

  if [[ "$TRAVIS_PULL_REQUEST_BRANCH" !=  dependabot* ]]; then
    exit 0
  fi

  git push origin
}

case "$1" in
    vendor)     vendor
        ;;
    push-back)  push_back
        ;;
    git-checkout) git_checkout
        ;;
    *)          echo >&2 "Wrong argument '$1'" && exit 1
        ;;
esac
