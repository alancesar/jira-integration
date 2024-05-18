create view if not exists issues_by_sprint as
select issue.*,
       issue_type.id as issue_type_id,
       issue_type.name as issue_type_name,
       issue_type.subtask as is_subtask,
       sprint.id as sprint_id,
       sprint.name as sprint_name,
       sprint.state as sprint_state,
       sprint.completed_at as sprint_completed_at
from issues issue
         inner join (select issue_id, MAX(sprint_id) as sprint_id
                     from issue_sprints
                     group by issue_id) last_sprint on issue.id = last_sprint.issue_id
         inner join issue_types issue_type on issue.issue_type_id = issue_type.id
         inner join sprints sprint on last_sprint.sprint_id = sprint.id