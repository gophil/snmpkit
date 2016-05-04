# snmpkit
a tool for executing snmp request


###如何使用 ?

`Example: `

1.安装npd
```go
    $ go get github.com/gophil/npd
    
    $ go github.com/cdevr/WapSNMP
```

2.运行编译脚本
```go
    $ make build
```

3.执行静态链接文件
```
$ cd builds 
$ ./snmpdemo -datafile=/your/datafile/path -w 10000 -i 10 -timeout 500 -oids=1.3.6.1.2.1.31.1.1.1.6
```

4.参数说明

> w : 最大工作并发任务数

> i : 任务执行间隔

> timeout: snmp采集超时时间

> oids: snmp oid, 多个以逗号分隔开

> datafile : 数据文件: 格式为 `[{"host": "1.1.1.1", "community": "public"} ...]`