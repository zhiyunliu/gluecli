package mysql

// QueryMysql .
const QueryColumnInfo = `
select
c.table_name,
c.column_name,
c.column_default,
c.is_nullable,
c.column_type,
c.column_key,
c.column_comment,
c.ordinal_position,
t.table_comment,
GROUP_CONCAT(CONCAT(if(co.constraint_type is null,'IDX',co.constraint_type),'(',s.index_name,',',s.SEQ_IN_INDEX,')')) con
from
information_schema.columns c
left join information_schema.tables t on 
                (c.table_name = t.table_name and t.table_schema =c.table_schema)
left join information_schema.statistics s on 
                ( s.column_name = c.column_name and s.table_name = c.table_name and s.table_schema =c.table_schema  )
left join information_schema.key_column_usage kc on 
                ( kc.table_name = s.table_name and kc.column_name =s.column_name and kc.table_schema =s.table_schema and s.index_name=kc.constraint_name )
left join information_schema.table_constraints co on 
                ( kc.constraint_name = co.constraint_name and co.table_name = kc.table_name and co.table_schema =kc.table_schema ) 
where
c.table_schema =@schema
group by 
c.table_name,
c.column_name,
c.column_default,
c.is_nullable,
c.column_type,
c.column_key,
c.column_comment,
c.ordinal_position,
t.table_comment
order by c.table_name,c.ordinal_position
`

const GetTableColumns = `
SELECT
	a.table_name,
	GROUP_CONCAT( a.column_name ) column_name 
FROM
	(
	SELECT
		c.table_name,
		c.column_name 
	FROM
		information_schema.COLUMNS c 
	WHERE
		c.table_schema = @schema
	ORDER BY
		c.table_name,
		c.ordinal_position 
	) a 
GROUP BY
	a.table_name
`
