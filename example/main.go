package main

import (
	"fmt"
	"github.com/changsongl/delay-queue-client/client"
	"github.com/changsongl/delay-queue-client/consumer"
	"github.com/changsongl/delay-queue-client/job"
	"time"
)

func main() {
	// job object
	myTopic, myID := "my-topic", "my-id"
	j, err := job.New(myTopic, myID, job.DelayOption(2*time.Second), job.TTROption(30*time.Second))
	if err != nil {
		panic(err)
	}

	// client
	cli := client.NewClient("http://127.0.0.1:8000")
	// adding job to delay queue, if job is exist will be failed
	if err = cli.AddJob(j); err != nil {
		panic(err)
	}

	// replace the job, even if the job is exists
	if err = cli.ReplaceJob(j); err != nil {
		panic(err)
	}

	// delete the job
	if err = cli.DeleteJob(myTopic, myID); err != nil {
		panic(err)
	}

	// pop the job from queue, no recommended. please use consumer.
	topic, id, body, delay, ttr, err := cli.PopJob(myTopic, 3*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(topic, id, body, delay, ttr)

	// finish the job, after having processed the job
	if err = cli.FinishJob(myTopic, myID); err != nil {
		panic(err)
	}

	// consumer jobs
	c := consumer.New(
		cli,
		topic,
		consumer.WorkerNumOption(1),
		consumer.PopTimeoutOption(3*time.Second),
	)
	ch := c.Consume()
	for jobMsg := range ch {
		id := jobMsg.GetID()
		body := jobMsg.GetBody()

		// do your job
		fmt.Println(id, body)

		if id == "xxx" {
			// job is not valid anymore
			if err = jobMsg.Finish(); err != nil {
				// do something
			}
			continue
		}

		if err = jobMsg.Finish(); err != nil {
			// do something
		}
	}
}
