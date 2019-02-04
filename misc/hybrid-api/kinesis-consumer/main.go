package main

import (
	"flag"
	"fmt"
    "time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// arn:aws:kinesis:eu-central-1:627211582084:stream/productevents

var (
	stream = flag.String("stream", "product-events", "your stream name")
	region = flag.String("region", "eu-west-1", "your AWS region")
)

func read_shard(kc *kinesis.Kinesis, i int,shard *kinesis.Shard) {
    sid := shard.ShardId;
	streamName := aws.String(*stream)

	iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
		ShardId:           sid,
		//ShardIteratorType: aws.String("TRIM_HORIZON"),
		//ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		//ShardIteratorType: aws.String("LATEST"),
		ShardIteratorType: aws.String("AFTER_SEQUENCE_NUMBER"),
        StartingSequenceNumber: shard.SequenceNumberRange.StartingSequenceNumber,
		StreamName: streamName,
	})
	if err != nil {
		panic(err)
	}
	si:= iteratorOutput.ShardIterator

    for {
	  fmt.Printf("[ %d ] Getting next records\n", i)

	  records, err := kc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: si,
	  })
	  fmt.Printf("[ %d ] Read done, got %d records\n", i, len(records.Records))
	  if err != nil {
	 	panic(err)
	  }
      for _,r := range records.Records {
	    fmt.Printf("[ %d ] ================================================================\n", i)
	    fmt.Printf("%s\n\n", string(r.Data))
      }
      if(*records.MillisBehindLatest == 0) {
	    fmt.Printf("[ %d ] At end, sleeping\n", i)
        time.Sleep(10 * time.Second)
      }
      si = records.NextShardIterator
    }
}

func main() {
	flag.Parse()

	s := session.New(&aws.Config{Region: aws.String(*region)})
	kc := kinesis.New(s)
	streamName := aws.String(*stream)

	streams, err := kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s\n", sid)
	//fmt.Printf("%v\n", streams)

    for i,shard := range streams.StreamDescription.Shards {
	  //fmt.Printf("%v\n", shard)
      go read_shard(kc,i,shard);
    }

    fmt.Print("Press ENTER to quit..")
    fmt.Scanln()
    fmt.Println("done")


}

