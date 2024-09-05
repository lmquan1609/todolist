package subscriber

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"log"
	"todolist/common"
	"todolist/common/asyncjob"
	"todolist/pubsub"
)

type subJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type pbEngine struct {
	serviceCtx goservice.ServiceContext
}

func NewEngine(serviceCtx goservice.ServiceContext) *pbEngine {
	return &pbEngine{serviceCtx}
}

func (engine *pbEngine) Start() error {
	engine.startSubTopic(common.TopicUserLikedItem, true,
		IncreaseLikeCountAfterUserLikeItem(engine.serviceCtx),
		PushNotificationAfterUserLikeItem(engine.serviceCtx),
	)

	engine.startSubTopic(common.TopicUserUnlikedItem, true,
		DecreaseLikeCountWhenUserUnlikeItem(engine.serviceCtx),
	)
	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}

func (engine *pbEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...subJob) error {
	ps := engine.serviceCtx.MustGet(common.PluginPubsub).(pubsub.PubSub)

	c, _ := ps.Subscribe(context.Background(), topic)

	for _, item := range jobs {
		log.Println("Setup subscriber for:", item.Title)
	}

	getJobHandler := func(job *subJob, message *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Println("running job for", job.Title, ", value:", message.Data())
			return job.Hld(ctx, message)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHldArr := make([]asyncjob.Job, len(jobs))
			for i := range jobs {
				jobHdl := getJobHandler(&jobs[i], msg)
				jobHldArr[i] = asyncjob.NewJob(jobHdl, asyncjob.WithName(jobs[i].Title))
			}

			group := asyncjob.NewGroup(isConcurrent, jobHldArr...)
			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}
