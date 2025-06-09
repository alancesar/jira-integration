#!/usr/bin/env bash

range="-15d"

./bin/fetch --jql="issuetype IN (Theme) AND updated >= "${range}
./bin/fetch --jql="project IN (\"Digital FX\", \"FX Core\") AND issuetype IN (Epic)  AND updated >= "${range}
./bin/fetch --jql="project IN (\"Digital FX\", \"FX Core\") AND issuetype NOT IN (subTaskIssueTypes(), Theme, Epic, \"Sprint Config\") AND updated >= "${range}
./bin/sync active future