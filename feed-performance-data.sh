#!/usr/bin/env bash

./bin/fetch --jql="project IN (\"Digital FX\", \"FX Core\") AND issuetype NOT IN (subTaskIssueTypes(), Theme, \"Sprint Config\") AND updated >= -10d"
