package mssql

//获取所有的表格
const GetAllTableName = `
select * from dbo.sysobjects where  OBJECTPROPERTY(id, N'IsUserTable') = 1
`
