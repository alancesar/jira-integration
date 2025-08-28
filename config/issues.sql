create or replace view tasks as
(
select id,
       key,
       summary,
       status,
       issue_type,
       project,
       sprint_id,
       story_points,
       parent_id,
       assignee_id,
       reporter_id
from issues
where issue_type in ('Task', 'Technical debt', 'Refinement', 'Story', 'Support', 'Spike', 'Bug'));

---

create or replace view epics as
(
select id, key, summary, status, project, parent_id
from issues
where issue_type = 'Epic');

---

create or replace view themes as
(
select id, key, summary, status, project
from issues
where issue_type = 'Theme');
