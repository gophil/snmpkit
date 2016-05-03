package main

import (
	"flag"
	"fmt"
	"github.com/gophil/npd"
	"github.com/gophil/snmpkit/config"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//测试作业
type TestJob struct {
	Id     int
	Name   string
	Host   string
	Result string
}

func NewTestJob(id int, name string, host string) *TestJob {
	return &TestJob{
		Id:   id,
		Name: name,
		Host: host,
	}
}

//作业执行方法
func (t *TestJob) HandleJob() {
	time.Sleep(10 * time.Second)
	fmt.Println(t.Name, " is working ")
	t.Result = t.Name + "_result"
}

var (
	work_num = flag.String("w", "100", "num of worker num") //执行的协程数量
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	num, err := strconv.Atoi(*work_num)
	if err != nil {
		num = 100 //默认数
	}

	var wg sync.WaitGroup
	var mpwg sync.WaitGroup

	d := npd.NewDispatcherWithMQ(num, num, &wg, &mpwg)

	//设置消息发送函数
	d.SetMF(func(task npd.Task) {
		tj := (*task.TargetObj).(*TestJob)
		fmt.Printf("task: %s 正在往host: %s, 推送结果: %s \n ", tj.Name, tj.Host, tj.Result)
		time.Sleep(200 * time.Millisecond)
	})

	d.Run()
	defer d.Stop()

	wg.Add(1)
	mpwg.Add(1)

	go func() {
		items, err := config.LoadSwitchFromFile("/Users/lihaoquan/Desktop/switch..")
		if err != nil {
			panic(err)
		}

		for i, item := range items {
			job := npd.CreateTask(NewTestJob(i, "TEST JOB "+fmt.Sprintf("%d", i), item.Host), "HandleJob")
			d.SubmitTask(job)
		}

		fmt.Println("jobs are submit")
		wg.Done()
		mpwg.Done()
	}()

	wg.Wait()
	mpwg.Wait()
	fmt.Println("all jobs are finished")

}
