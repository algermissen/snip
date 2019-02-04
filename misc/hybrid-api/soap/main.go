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

func init() {
	s := session.New(&aws.Config{Region: aws.String(region)})
	kc = kinesis.New(s)
	streamName = aws.String(stream)
}

type EventsRequest struct {
	AfterSequenceNumber string
}

type Event struct {
	id string
}

type EventsResponse struct {
	Events             []Event
	LastSequenceNumber string
}

func MakeSoap() func(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context) (events.APIGatewayProxyResponse, error) {

		request := EventsRequest{
			AfterSequenceNumber: "49585037527304342602018146810276658730793015177990111234",
		}

		iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
			ShardId:                &sid,
			ShardIteratorType:      aws.String("AFTER_SEQUENCE_NUMBER"),
			StartingSequenceNumber: aws.String(request.AfterSequenceNumber),
			StreamName:             streamName,
		})
		if err != nil {
			panic(err)
		}

        eventsResponse := EventsResponse{
			Events:             []Event{},
			LastSequenceNumber: "",
		}

	    si:= iteratorOutput.ShardIterator
		b := ""
        for {
		// get records use shard iterator for making request
		records, err := kc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: si,
		})
		if err != nil {
			panic(err)
		}
        fmt.Printf("XXXXXXXXXXXXXXXX %d",len(records.Records))
		for _, r := range records.Records {
			ev := Event{}
			eventsResponse.Events = append(eventsResponse.Events,ev)
			b += fmt.Sprintf("================================\n%s\n", string(r.Data))
		}
        if ( len(records.Records) > 0) {
			eventsResponse.LastSequenceNumber = string(*((records.Records[len(records.Records)-1]).SequenceNumber))
        }
    if(*records.MillisBehindLatest == 0) {
        break;

    }
    si = records.NextShardIterator
}

		//headers:= map[string]string{"Content-Type": "application/soap...xml"}
		headers := map[string]string{
			"Content-Type": "text/plain",
			"Link":         fmt.Sprintf("<http://..../?last=%s>; rel=next", eventsResponse.LastSequenceNumber),
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       b,
		}, nil
	}
}

func main() {

	lambda.Start(MakeSoap())
}
