
# 4. 表设计

## 4.1. 系统表

### 4.1.1. 数据字典行政区域信息[dds_area_info]

| 字段名          | 类型          | 默认值 | 为空  | 约束  | 描述                                |
| --------------- | ------------- | :----: | :---: | :---: | :---------------------------------- |
| canton_code     | varchar2(128) |        |  是   |  PK   | 新政区编号                          |
| chinese_name    | varchar2(512) |        |  是   |       | 中文名称                            |
| parent_code     | varchar2(128) |        |  是   |       | 父级编号                            |
| grade           | number(2)     |        |  是   |       | 行政等(0-全国 1-省级 2-市级 3-县级) |
| full_spell      | varchar2(256) |        |  是   |       | 全拼                                |
| simple_spell    | varchar2(32)  |        |  是   |       | 简称                                |
| sort_id         | number(3)     |        |  是   |       | 排序 id                             |
| status          | number(2)     |   0    |  是   |       | 状态(0-启用 1-禁用)                 |
| longitude_value | varchar2(128) |        |  是   |       | 经度                                |
| latitude_value  | varchar2(128) |        |  是   |       | 纬度                                |


###  1. 商户信息[ots_merchant_info]
	
 | 字段名      | 类型         | 默认值  | 为空  |        约束        | 描述     |
 | ----------- | ------------ | :-----: | :---: | :----------------: | :------- |
 | mer_no      | varchar2(32) |         |  否   |         PK         | 编号     |
 | mer_name    | varchar2(64) |         |  否   | UNQ(merchant_name) | 名称     |
 | mer_crop    | varchar2(64) |         |  否   |                    | 公司     |
 | mer_type    | number(1)    |         |  否   |                    | 类型     |
 | bd_uid      | number(20)   |         |  否   |                    | 商务人员 |
 | status      | number(1)    |    0    |  否   |                    | 状态     |
 | create_time | date         | sysdate |  否   |                    | 创建时间 |