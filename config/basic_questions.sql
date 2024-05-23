--- velocity
select sum(issues_by_sprint.story_points)   as "Story Points",
       issues_by_sprint.sprint_name         as "Sprint",
       issues_by_sprint.sprint_completed_at as "Completed At"
from issues_by_sprint
         inner join
     (select id from sprints where state = 'closed' order by completed_at desc limit 7) last_sprints
     on last_sprints.id = issues_by_sprint.sprint_id
         left join issue_products issue_product on issue_product.issue_id = issues_by_sprint.id
         left join products product on product.id = issue_product.product_id
where issues_by_sprint.is_subtask = false
group by 2, 3
order by 3;

--- story points by task type
select sum(issues_by_sprint.story_points)   as "Story Points",
       issues_by_sprint.issue_type_name     as "Issue Type",
       issues_by_sprint.sprint_name         as "Sprint",
       issues_by_sprint.sprint_completed_at as "Completed At"
from issues_by_sprint
         inner join
     (select id from sprints where state = 'closed' order by completed_at desc limit 7) last_sprints
     on last_sprints.id = issues_by_sprint.sprint_id
where issues_by_sprint.is_subtask = false
group by 2, 3, 4
order by 4;

--- average story points
select avg(story_points) as "Story Points",
       allocation        as "Allocation"
from (select sum(issues_by_sprint.story_points) as story_points,
             case
                 when allocation = 'Operation' then 'Operation'
                 when kind is null then 'Undefined'
                 else kind end                  as allocation
      from issues_by_sprint
               inner join
           (select id from sprints where state = 'closed' order by completed_at desc limit 7) last_sprints
           on last_sprints.id = issues_by_sprint.sprint_id
      where issues_by_sprint.is_subtask = false
      group by issues_by_sprint.sprint_name, 2
      order by issues_by_sprint.sprint_name, 2)
group by 2;

--- todos
select key,
       summary,
       case when kind is not null then true else false end              as has_label,
       case when products.issue_id is not null then true else false end as has_product
from issues_by_sprint
         inner join
     (select id from sprints where state = 'closed' order by completed_at desc limit 7) last_sprints
     on last_sprints.id = issues_by_sprint.sprint_id
         left join (select issue_id from issue_products group by issue_id) products
                   on products.issue_id = issues_by_sprint.id
where (kind is null or products.issue_id is null)
order by issues_by_sprint.id desc;

--- issues
select i.key                      as "Key",
       i.summary                  as "Summary",
       i.story_points             as "Story Points",
       i.issue_type_name          as "Issue Type",
       i.kind                     as "Kind",
       i.allocation               as "Allocation",
       group_concat(p.name, ', ') as "Products"
from issues_by_sprint i
         inner join
     (select id from sprints where state = 'closed' order by completed_at desc limit 7) last_sprints
     on last_sprints.id = i.sprint_id
         left join issue_products ip on i.id = ip.issue_id
         left join products p on ip.product_id = p.id
where i.is_subtask = false
group by i.key, i.created_at
order by i.created_at desc;