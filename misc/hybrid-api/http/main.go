package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// arn:aws:kinesis:eu-central-1:627211582084:stream/productevents

var (
	stream = "productevents"
	region = "eu-central-1"
	sid    = "shardId-000000000000"
)

var kc *kinesis.Kinesis = nil
var streamName = aws.String("")
var start = ""

func init() {
	s := session.New(&aws.Config{Region: aws.String(region)})
	kc = kinesis.New(s)
	streamName = aws.String(stream)

	streams, err := kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		panic(err)
	}

	start = string(*(streams.StreamDescription.Shards[0].SequenceNumberRange.StartingSequenceNumber))

}

type Event struct {
	id string
}

type EventsResponse struct {
	Events             []Event
	LastSequenceNumber string
}

func MakeHttp() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		since := request.QueryStringParameters["since"]
		if since == "" {
			since = start
		}

		iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
			ShardId:                &sid,
			ShardIteratorType:      aws.String("AFTER_SEQUENCE_NUMBER"),
			StartingSequenceNumber: aws.String(since),
			StreamName:             streamName,
		})
		if err != nil {
			panic(err)
		}

		eventsResponse := EventsResponse{
			Events:             []Event{},
			LastSequenceNumber: since,
		}

		si := iteratorOutput.ShardIterator
		b := ""

		for {
			// get records use shard iterator for making request
			records, err := kc.GetRecords(&kinesis.GetRecordsInput{
				ShardIterator: si,
			})
			if err != nil {
				panic(err)
			}
			for _, r := range records.Records {
				ev := Event{}
				eventsResponse.Events = append(eventsResponse.Events, ev)
				b += "--MixedBoundaryString\r\n"
				b += fmt.Sprintf("%s\n", string(r.Data))
			}
			if len(records.Records) > 0 {
				eventsResponse.LastSequenceNumber = string(*((records.Records[len(records.Records)-1]).SequenceNumber))
			}
			if *records.MillisBehindLatest == 0 {
				break

			}
			si = records.NextShardIterator
		}
		b += "--MixedBoundaryString--\r\n"

		//headers:= map[string]string{"Content-Type": "application/soap...xml"}
		headers := map[string]string{
			"Link":         fmt.Sprintf("</proto/http/latest?since=%s>; rel=next", eventsResponse.LastSequenceNumber),
			"Content-Type": `multipart/mixed; boundary="MixedBoundaryString"`,
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       b,
		}, nil
	}
}

func main() {

	lambda.Start(MakeHttp())
}
