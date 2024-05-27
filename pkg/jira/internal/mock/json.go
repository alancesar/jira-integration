package mock

const SearchResponse = `
{
    "expand": "names,schema",
    "startAt": 0,
    "maxResults": 15,
    "total": 1,
    "issues": [
        {
            "id": "1234",
            "self": "https://bexs.atlassian.net/rest/api/3/issue/1234",
            "key": "MAQ-1234",
            "fields": {
                "summary": "Some story",
                "issuetype": {
                    "self": "https://bexs.atlassian.net/rest/api/3/issuetype/1234",
                    "id": "1234",
                    "description": "Some issue type description",
                    "iconUrl": "https://bexs.atlassian.net/rest/api/2/universal_avatar/view/type/issuetype/avatar/1234?size=medium",
                    "name": "Story",
                    "subtask": false,
                    "avatarId": 1234,
                    "hierarchyLevel": 0
                },
                "created": "2023-04-12T14:00:00.0-0300",
                "updated": "2023-04-13T16:00:00.0-0300",
                "parent": {
                    "id": "123",
                    "key": "MAQ-123",
                    "self": "https://bexs.atlassian.net/rest/api/3/issue/123",
                    "fields": {
                        "summary": "Some epic",
                        "status": {
                            "self": "https://bexs.atlassian.net/rest/api/3/status/123",
                            "description": "This status is managed internally by Jira Software",
                            "iconUrl": "https://bexs.atlassian.net/",
                            "name": "In Progress",
                            "id": "123",
                            "statusCategory": {
                                "self": "https://bexs.atlassian.net/rest/api/3/statuscategory/4",
                                "id": 4,
                                "key": "indeterminate",
                                "colorName": "yellow",
                                "name": "In Progress"
                            }
                        },
                        "priority": {
                            "self": "https://bexs.atlassian.net/rest/api/3/priority/3",
                            "iconUrl": "https://bexs.atlassian.net/images/icons/priorities/medium.svg",
                            "name": "Medium",
                            "id": "3"
                        },
                        "issuetype": {
                            "self": "https://bexs.atlassian.net/rest/api/3/issuetype/10000",
                            "id": "10000",
                            "description": "A big user story that needs to be broken down. Created by Jira Software - do not edit or delete.",
                            "iconUrl": "https://bexs.atlassian.net/images/icons/issuetypes/epic.svg",
                            "name": "Epic",
                            "subtask": false,
                            "hierarchyLevel": 1
                        }
                    }
                },
				"customfield_10025": 3.0,
                "customfield_10070": {
                    "self": "https://bexs.atlassian.net/rest/api/3/customFieldOption/10070",
                    "value": "Maquininha",
                    "id": "10070"
                },
                "customfield_10020": [
                    {
                        "id": 1,
                        "name": "Sprint 1",
                        "state": "closed",
                        "boardId": 66,
                        "goal": "Finish some big item",
                        "startDate": "2023-03-20T15:00:00.0Z",
                        "endDate": "2023-03-30T21:00:00.0Z",
                        "completeDate": "2023-03-30T21:00:00.0Z"
                    },
                    {
                        "id": 2,
                        "name": "Sprint 2",
                        "state": "active",
                        "boardId": 66,
                        "goal": "Finish another big item",
                        "startDate": "2023-04-03T15:00:00.0Z",
                        "endDate": "2023-04-13T21:00:00.0Z"
                    }
                ],
                "fixVersions": [
                    {
                        "self": "https://bexs.atlassian.net/rest/api/3/version/123",
                        "id": "123",
                        "description": "Release 2023/Q1",
                        "name": "2023 - Q1",
                        "archived": false,
                        "released": false,
                        "releaseDate": "2023-03-31"
                    }
                ],
                "customfield_10441": {
                    "self": "https://bexs.atlassian.net/rest/api/3/customFieldOption/10441",
                    "value": "Digital FX",
                    "id": "10441"
                },
                "priority": {
                    "self": "https://bexs.atlassian.net/rest/api/3/priority/3",
                    "iconUrl": "https://bexs.atlassian.net/images/icons/priorities/medium.svg",
                    "name": "Medium",
                    "id": "3"
                },
                "labels": [
                    "Some Label",
					"Another Label"
                ],
                "customfield_10106": {
                    "self": "https://bexs.atlassian.net/rest/api/3/customFieldOption/10106",
                    "value": "DigitalFx",
                    "id": "10106"
                },
                "assignee": {
                    "self": "https://bexs.atlassian.net/rest/api/3/user?accountId=abc123",
                    "accountId": "abc123",
                    "emailAddress": "some.user@bexsbanco.com.br",
                    "avatarUrls": {
                        "16x16": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/abc123/efg456/16",
                        "24x24": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/abc123/efg456/24",
                        "32x32": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/abc123/efg456/32",
                        "48x48": "https://avatar-management--avatars.us-west-2.prod.public.atl-paas.net/abc123/efg456/48"
                    },
                    "displayName": "Some User",
                    "active": true,
                    "timeZone": "America/Sao_Paulo",
                    "accountType": "atlassian"
                },
                "status": {
                    "self": "https://bexs.atlassian.net/rest/api/3/status/10052",
                    "description": "Done",
                    "iconUrl": "https://bexs.atlassian.net/",
                    "name": "Done",
                    "id": "10052",
                    "statusCategory": {
                        "self": "https://bexs.atlassian.net/rest/api/3/statuscategory/3",
                        "id": 3,
                        "key": "done",
                        "colorName": "green",
                        "name": "Done"
                    }
                }
            }
        }
    ]
}
`
