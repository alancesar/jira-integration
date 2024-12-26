#!/usr/bin/env bash

echo 'feeding data from jira api...'
./sync --full=false
echo 'copying database to metabase container'
docker cp sqlite.db metabase:/db