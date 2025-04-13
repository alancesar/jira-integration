#!/usr/bin/env bash

./bin/fetch --jql="project IN (\"Digital FX\", \"FX Core\") AND issuetype IN (Epic) AND updated >= -10d"
./bin/fetch --jql="project IN (\"Digital FX\", \"FX Core\") AND issuetype NOT IN (subTaskIssueTypes(), Theme, Epic, \"Sprint Config\") AND updated >= -10d"
