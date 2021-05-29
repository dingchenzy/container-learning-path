## 开发mysql_exporter

### 监控指标

```go
mysql 可用性
    操作失败
        select 1;
        ping
慢查询次数
	// 需要开启配置，在 my.cnf 中添加
	// slow_query_log    slow_query_log_file=/var/log/mariadb/mariadb-slow.log    long_query_time=2
    show global status where variable_name='low_queries'
容量:
    qps:
    // 每秒请求数，也就是在一秒之内请求的数量
        show global status where variable_name='Queries'
    tps:
    // 服务器每秒处理的事务数，因为 mysql 分为多个动作所以处理的事务分为多种
        insert, update, delete *
        com_insert
        com_update
        com_delete
        com_select
        com_replace
    连接:
        // 查看当前线程运行
        show global status where variable_name='Threads_connected'
        // 最大连接数量
        show global variables where variable_name= 'max_connections';

    流量：
        // 查看接收字节数
        show global status where variable_name='Bytes_received'
        // 查看发送字节数
        show global status where variable_name='Bytes_sent'

prometheus.Collector interface


// mysql连接信息 => mysql host, port, user
```

## 需求分析

代码中实现了 `Slow_querys` 慢查询，可用性检测以及连接的线程数量。

### 代码结构

```go
.
├── controllers		// 代码主要的处理器
│   ├── base.go	// 基础查询功能，为其他处理器提供
│   ├── longtimeselect.go		// slow_querys 慢查询
│   ├── threads.go	// 线程数量
│   └── up.go		// 可用性检测
├── go.mod
├── go.sum
├── main.go			// 整个程序的入口
└── vendor
```

### 代码内容

#### main.go

```go
package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/exporter/mysql_exporter/controllers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn    = "localuser:123.com@tcp(localhost:3306)/mysql?charset=utf8mb4&loc=PRC&parseTime=true"
	driver = "mysql"
	addr   = ":9999"
)

/*
	项目的 controllers 部分是负责提取指标并生成 collector 接口中的方法
	然后通过 registrer 注册
	最后生成
*/

func main() {
	// 打开数据库
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("db open：", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("db ping：", err)
	}

	// 注册指标
	prometheus.Register(controllers.NewUpCollertor(db))
	prometheus.Register(controllers.NewThreadCoThreadsCollector(db))
	prometheus.Register(controllers.NewLongtimeLongtimeCollector(db))
	// 开启服务
	http.Handle("/metrics/", promhttp.Handler())

	// 启动 web 服务
	http.ListenAndServe(addr, nil)
}
```

#### longtimeselect.go

```go
package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type LongtimeCollector struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

func NewLongtimeLongtimeCollector(sqldb *sql.DB) *LongtimeCollector {
	return &LongtimeCollector{
		mysqlSelcetor: &mysqlSelcetor{sqldb},
		desc:          prometheus.NewDesc("mysql_select_long_time", "this is select timeout", nil, nil),
	}
}

func (l *LongtimeCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- l.desc
}

func (l *LongtimeCollector) Collect(metric chan<- prometheus.Metric) {
	metric <- prometheus.MustNewConstMetric(l.desc, prometheus.CounterValue, l.mysqlSelcetor.showStatus("Slow_queries"))
}
```

#### thread.go

```go
package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type ThreadsCollector struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

func NewThreadCoThreadsCollector(db *sql.DB) *ThreadsCollector {
	return &ThreadsCollector{
		mysqlSelcetor: newMysqlSelcector(db),
		desc:          prometheus.NewDesc("mysql_thread_count", "current thread count", nil, nil),
	}
}

func (t *ThreadsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- t.desc
}

func (t *ThreadsCollector) Collect(metrics chan<- prometheus.Metric) {
	value := t.mysqlSelcetor.showStatus("Threads_connected")
	metrics <- prometheus.MustNewConstMetric(t.desc, prometheus.GaugeValue, value)
}
```

#### up.go

```go
package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type UpController struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

// 构建函数，并封装 newMysqlSelector 中的查询方法
func NewUpCollertor(sqldb *sql.DB) *UpController {
	return &UpController{
		newMysqlSelcector(sqldb),
		prometheus.NewDesc("mysql_up", "mysql up or down", nil, nil),
	}
}

// desc 中集成着 label 的名称，指标的名称，help 的内容等
func (Up *UpController) Describe(descs chan<- *prometheus.Desc) {
	descs <- Up.desc
}

func (Up *UpController) Collect(metrics chan<- prometheus.Metric) {
	var up float64 = 1
	if err := Up.mysqlSelcetor.mysqlconnection.Ping(); err != nil {
		up = 0
	}
	// 生成 metrics
	// metrics 内部集成着 desc 以及指标的类型，label 的 value 等内容
	metrics <- prometheus.MustNewConstMetric(Up.desc, prometheus.GaugeValue, up)
}
```







