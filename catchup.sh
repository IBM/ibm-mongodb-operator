#!/usr/bin/env bash

git checkout master
git fetch upstream
git merge upstream/master
git reset --hard HEAD
git push origin master
