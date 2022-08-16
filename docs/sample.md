
# 4. 表设计

## 4.1. 系统表



###  1. 商户信息[ots_merchant_info]
	
| 字段名      | 类型         | 默认值  | 为空  |        约束        | 描述     |
| ----------- | ------------ | :-----: | :---: | :----------------: | :------- |
| mer_no      | varchar2(32) |         |  否   |         PK         | 编号     |
| mer_name    | varchar2(64) |         |  否   | UNQ(merchant_name) | 名称     |
| mer_crop    | varchar2(64) |         |  是   |                    | 公司     |
| mer_type    | number(1)    |         |  否   |                    | 类型     |
| bd_uid      | number(20)   |         |  否   |                    | 商务人员 |
| status      | number(1)    |    0    |  否   |                    | 状态     |
| create_time | date         | sysdate |  否   |                    | 创建时间 |