package gcppubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"pkg.agungdp.dev/candi/candihelper"
	"pkg.agungdp.dev/candi/candishared"
	"pkg.agungdp.dev/candi/codebase/interfaces"
	"pkg.agungdp.dev/candi/tracer"
)

type publisher struct {
	client *pubsub.Client
}

// NewPublisher gcp
func NewPublisher(client *pubsub.Client) interfaces.Publisher {
	return &publisher{
		client: client,
	}
}

func (p *publisher) PublishMessage(ctx context.Context, args *candishared.PublisherArgument) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "GCPPubSub:PublishMessage")
	defer trace.Finish()

	message := &pubsub.Message{
		Data:        candihelper.ToBytes(args.Data),
		PublishTime: time.Now(),
	}
	message.Attributes = make(map[string]string, len(args.Header))
	for k, v := range args.Header {
		if val, ok := v.(string); ok {
			message.Attributes[k] = val
		}
	}

	result := p.client.Topic(args.Topic).Publish(ctx, message)
	serverID, err := result.Get(ctx)
	trace.Log("server_id", serverID)
	return err
}