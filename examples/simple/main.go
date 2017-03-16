package main

import (
	"log"
	"time"

	"github.com/dewey4iv/dawdle/drivers/inmem"
	"github.com/dewey4iv/dawdle/processor"
	"github.com/dewey4iv/dawdle/scheduler"
	"github.com/the-control-group/event-analytics-service/tasks"
	"github.com/the-control-group/event-analytics-service/tasks/linkAccountOrderSuccess"
)

func main() {
	store, err := inmem.New()
	if err != nil {
		panic(err)
	}

	registrar := processor.NewRegistrar()
	tasks.Setup(registrar)

	scheduler, err := scheduler.New(
		scheduler.WithStore(store),
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 500; i++ {
		firstTask := linkAccountOrderSuccess.New(linkAccountOrderSuccess.Args{
			Argument: "hey hey",
			Delay:    time.Millisecond * time.Duration(i*50),
		})

		if err = scheduler.Schedule(firstTask); err != nil {
			panic(err)
		}
	}

	processor, err := processor.New(
		processor.WithStore(store),
		processor.WithRegistrar(registrar),
	)
	if err != nil {
		panic(err)
	}

	if false {
		log.Println(processor)
	}

	select {}
}
