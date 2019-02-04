package main

import (
	"flag"
	"fmt"
    "encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
    "github.com/google/uuid"
)

var (
	stream = flag.String("stream", "productevents", "your stream name")
	region = flag.String("region", "eu-central-1", "your AWS region")
)

// Product type structs
type Product struct {
    Id string `json:"id"`
    Description string `json:"description"`
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
	fmt.Printf("Writing to stream %s\n",*(streams.StreamDescription.StreamARN))


	// put 10 records using PutRecords API
	entries := make([]*kinesis.PutRecordsRequestEntry, 10)
	for i := 0; i < len(entries); i++ {
        id := fmt.Sprintf("%d",i)
        k,d := makeProduct(id)
		entries[i] = &kinesis.PutRecordsRequestEntry{
			Data:         k,
			PartitionKey: aws.String(d),
		}
	}
	fmt.Printf("%v\n", entries)
	putsOutput, err := kc.PutRecords(&kinesis.PutRecordsInput{
		Records:    entries,
		StreamName: streamName,
	})
	if err != nil {
		panic(err)
	}
	// putsOutput has Records, and its shard id and sequece enumber.
    fmt.Printf("Done. Records failed: %v\n", putsOutput) //.FailedRecordCount)

}


func makeProduct(product_id string) ([]byte,string) {

    message_id := uuid.New().String()
    product_title := "Some product title " + product_id

    product := Product{
        Id:product_id,
        Description:"lkhdlkq hflqeh feh wflweh flkwehfl",
    }
    b, err := json.MarshalIndent(&product, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

    //Control: PUT (https://tools.ietf.org/html/rfc5536#section-3.2.3)
    //Distribution: DE, FR  (https://tools.ietf.org/html/rfc5536#section-3.2.4)

    s := fmt.Sprintf(`Message-ID: %s
Content-Location: http://api.oppla.io/products/%s
Content-Type: application/vnd.hse24.product+json
Etag: "1234567890"
Content-Length: %d
Content-Title: %s
Control: PUT (https://tools.ietf.org/html/rfc5536#section-3.2.3)
Distribution: DE, FR
Link: <http://pdm.hse24.de/products/%s>; rel=self

%s
`, message_id, product_id, len(b),product_title, product_id,string(b))

	return []byte(s),product_id

}
