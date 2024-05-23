drop view if exists issues_by_sprint;
create view issues_by_sprint as
select issue.*,
       case
           when ((select 1 from json_each(labels) where value = 'Front-End') or lower(summary) like '[front]%')
               then 'Front-End'
           when (select 1 from json_each(labels) where value = 'Back-End') then 'Back-End'
           when (select 1 from json_each(labels) where value = 'QA') then 'QA'
           when (select 1 from json_each(labels) where value in ('DevOps', 'SRE', 'SRE-DevOps')) then 'DevOps'
           end             as kind,
       case
           when issue_type.name in ('Support', 'Bug', 'Sub-bug') then 'Operation'
           else 'Investment'
           end             as allocation,
       issue_type.id       as issue_type_id,
       issue_type.name     as issue_type_name,
       issue_type.subtask  as is_subtask,
       sprint.id           as sprint_id,
       sprint.name         as sprint_name,
       sprint.state        as sprint_state,
       sprint.completed_at as sprint_completed_at
from issues issue
         inner join (select issue_id, max(sprint_id) as sprint_id
                     from issue_sprints
                     group by issue_id) last_sprint on issue.id = last_sprint.issue_id
         inner join issue_types issue_type on issue.issue_type_id = issue_type.id
         inner join sprints sprint on last_sprint.sprint_id = sprint.id
where issue.status_id = 10291
  and issue.issue_type_id <> 10304;
