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
SELECT
c.id,
c.name,
c.isnullable ,
b.name,
c.length,
COLUMNPROPERTY(c.id,c.name,'PRECISION') 整数位,
isnull(COLUMNPROPERTY(c.id,c.name,'Scale'),0) 小数位,
c.cdefault,
isnull(e.text,'') default_val,
case when exists(SELECT 1 FROM sysobjects where xtype='PK' and parent_obj=c.id and name in (
                     SELECT name FROM sysindexes WHERE indid in( SELECT indid FROM sysindexkeys WHERE id = c.id AND colid=c.colid))) then  1 else 0 end PK,
f.value

fROM syscolumns c
left join  sys.extended_properties f
on c.id=f.major_id and c.colid=f.minor_id  
left join syscomments e 
on c.cdefault = e.id
left join systypes b 
on c.xusertype = b.xusertype

where c.id =@{obj_id}
order by c.colorder
`

const QueryTableIdxs = `
select t.object_id, t.name ,f.value  from sys.tables t
left join  sys.extended_properties f
on t.object_id = f.major_id
and f.minor_id = 0 
where 1= 1 
&{t.name}
`
