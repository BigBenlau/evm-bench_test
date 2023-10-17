#!/usr/bin/env bash

git add .

read -p "请输入您的 GitHub Token: " github_token

git commit -m "dcyUpdate"

git push https://$github_token@github.com/BigBenlau/evm-bench_test.git  main:dcy
