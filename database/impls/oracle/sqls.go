package oracle

// GetAllTableNameInOracle .
const GetAllTableNameInOracle = `
select
  ub.table_name
from
  user_tables ub
order by 
  ub.table_name
`

// GetSingleTableInfoInOracle .
const GetSingleTableInfoInOracle = `SELECT a.table_name,
a.column_name,
a.data_type,
CASE upper(a.data_type)
  WHEN upper('number') THEN
   decode(a.data_scale,
		  0,
		  to_char(a.data_precision),
		  a.data_precision || ',' || a.data_scale)
  WHEN upper('date') THEN
   ''
  ELSE
   to_char(a.data_length)
END data_length,
a.data_precision,
a.data_scale,
a.nullable,
a.data_default data_default,
b.comments column_comments,
c.comments table_comments,
nvl(d.constraint_type,'') constraint_type,
nvl(e.index_name,'') index_name
FROM user_tab_columns a
LEFT JOIN user_col_comments b ON (a.table_name = b.table_name AND
							a.column_name = b.column_name)
LEFT JOIN user_tab_comments c ON a.table_name = c.table_name
LEFT JOIN (SELECT cu.table_name,
			 cu.column_name,
			 wm_concat(au.constraint_type || '(' ||
					   au.constraint_name || '|' ||
					   nvl(cu.position, 0) || ')') constraint_type
		FROM user_cons_columns cu
		LEFT JOIN user_constraints au ON (cu.constraint_name =
										 au.constraint_name)
	   WHERE cu.table_name = @table_name
	   GROUP BY cu.table_name, cu.column_name) d ON d.column_name = a.column_name
LEFT JOIN (select t.column_name,
			 wm_concat('IDX(' || t.index_name || '|' ||
					   nvl(t.column_position, 0) || ')') index_name
		from user_ind_columns t, user_indexes i
	   where t.table_name = @table_name
		 and t.index_name = i.index_name
		 and i.uniqueness = 'NONUNIQUE'
	   group by t.column_name) e on e.column_name = a.column_name
WHERE a.table_name = @table_name
ORDER BY a.column_id
`
