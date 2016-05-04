package main

import (
	"flag"
	"fmt"
	wapsnmp "github.com/cdevr/WapSNMP"
	"github.com/gophil/npd"
	"github.com/gophil/snmpkit/config"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TestJob struct {
	Id          string
	Host        string
	Community   string
	Oid         string
	Result      string
	Timeout     int
	fail        bool
	failMessage string
}

func NewTestJob(id string, host string, community string, oid string, timeout int) *TestJob {
	t := &TestJob{
		Id:        id,
		Oid:       oid,
		Host:      host,
		Timeout:   timeout,
		Community: community,
	}

	return t
}

func (t *TestJob) SetFailure(message string) {
	t.fail = true
	t.failMessage = message
}

func (t *TestJob) GetSnmpValue() string {
	version := wapsnmp.SNMPv2c
	target := t.Host
	id := t.Oid

	oid := wapsnmp.MustParseOid(id)
	fmt.Printf("Contacting %v %v %v\n", target, t.Community, version)

	wsnmp, err := wapsnmp.NewWapSNMP(target, t.Community, version, time.Duration(t.Timeout)*time.Millisecond, 0)
	defer wsnmp.Close()
	if err != nil {

		t.SetFailure("ECW")
		//fmt.Printf("Error creating wsnmp => %v\n", wsnmp)
		return ""
	}
	table, err := wsnmp.GetTable(oid)
	if err != nil {

		t.SetFailure("ECT")
		//fmt.Printf("Error getting table => %v\n", wsnmp)
		return ""
	}

	return strconv.Itoa(len(table))
}

//作业执行方法
func (t *TestJob) HandleJob() {
	v := t.GetSnmpValue()
	t.Result = v

}

var (
	work_num = flag.Int("w", 100, "num of worker num")              //执行的协程数量
	interval = flag.Int("i", 10, "interval of worker execute")      //任务执行间隔
	timeout  = flag.Int("timeout", 500, "timeout of smmp get data") //SNMP调用超时
	oids     = flag.String("oids", "", "oids for snmp")             //oids 数据
	datafile = flag.String("datafile", "", "oids for snmp")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	num := *work_num

	oidstr := *oids
	filePath := *datafile

	if filePath == "" {
		println("no data file found")
		return
	}

	if oidstr == "" {
		println("no oid found")
		return
	}

	oidlist := strings.Split(oidstr, ",")
	if len(oidlist) == 0 {
		println("no oids found")
		return
	}

	var wg sync.WaitGroup
	var mpwg sync.WaitGroup

	d := npd.NewDispatcherWithMQ(num, num, &wg, &mpwg)

	//设置消息发送函数
	d.SetMF(func(task npd.Task) {
		tj := (*task.TargetObj).(*TestJob)
		time.Sleep(100 * time.Millisecond)
		if !tj.fail {
			fmt.Printf("task: %s(%s) 正在往host: %s, 推送结果: %s \n ", tj.Id, tj.Oid, tj.Host, tj.Result)
		} else {
			fmt.Printf("正在发生错误信息 %s : %s \n", tj.failMessage, tj.Oid)
		}

	})

	itvl := time.Duration(*interval)
	d.RunWithLimiter(itvl * time.Millisecond)
	defer d.Stop()

	start := time.Now()

	wg.Add(1)
	mpwg.Add(1)

	go func() {
		items, err := config.LoadSwitchFromFile(filePath)
		if err != nil {
			panic(err)
		}

		for i, item := range items {
			i_index := strconv.Itoa(i)
			for j, oid := range oidlist {
				j_index := strconv.Itoa(j)
				job := npd.CreateTask(NewTestJob(i_index+"-"+j_index, item.Host, item.Community, oid, *timeout), "HandleJob")
				d.SubmitTask(job)
			}

		}

		fmt.Println("jobs are submit")
		wg.Done()
		mpwg.Done()
	}()

	wg.Wait()
	mpwg.Wait()

	cost := fmt.Sprintln(time.Now().Sub(start).Seconds())

	fmt.Println("all jobs are finished, cost (seconds): ", cost)

}
