package sqls

const QueryTables = `
select t.object_id, t.name ,f.value  from sys.tables t
left join  sys.extended_properties f
on t.object_id = f.major_id
and f.minor_id = 0 
where 1= 1 
&{t.name}
`

const QueryTableCols = ` 
SELECt
c.name col_name,
b.name col_type,
COLUMNPROPERTY(c.id,c.name,'PRECISION') length,
isnull(COLUMNPROPERTY(c.id,c.name,'Scale'),0) decimal,
isnull(e.text,'') default_val,
case when c.isnullable = 0 then '否' else '是' end isnullable ,
f.value comments,
COLUMNPROPERTY(c.id,c.name,'IsIdentity') isidentity

from syscolumns c
left join  sys.extended_properties f
on c.id=f.major_id and c.colid=f.minor_id  
left join syscomments e 
on c.cdefault = e.id
left join systypes b 
on c.xusertype = b.xusertype

where c.id =  @{obj_id}
order by c.colorder
 

`

const QueryTableIdxs = `
select 
row_number() over(partition by idx_name order by col_sort ) sort_val,
*
from (
SELECT 
case 
when ax.type = 1 then 'PK'
when ax.type = 2 and ax.is_unique = 1 then 'UNQ'
when ax.type = 2 and ax.is_unique = 0 then 'IDX'
end idx_type,

a.name idx_name,d.name col_name,d.colid col_sort

FROM  sysindexes  a 
join sys.indexes  ax on a.id = ax.object_id and a.name = ax.name
JOIN  sysindexkeys  b  ON  a.id=b.id  AND  a.indid=b.indid 
JOIN  sysobjects  c  ON  b.id=c.id 
JOIN  syscolumns  d  ON  b.id=d.id  AND  b.colid=d.colid 
WHERE  a.indid  NOT IN(0,255) 
AND   c.id= @{obj_id}
 ) t


`
