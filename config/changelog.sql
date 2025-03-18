create or replace view issues_changelog as
(
select issues.id as issue_id,
       min(started_at.created_at) started_at,
       max(done_at.created_at)    done_at
from issues
         inner join changelogs started_at
                    on issues.id = started_at.issue_id and started_at."to" in ('In Progress', 'In Development')
         inner join changelogs done_at on issues.id = done_at.issue_id and done_at."to" in ('Done')
where story_points is not null
group by issues.id,
         issues.story_points
    );
