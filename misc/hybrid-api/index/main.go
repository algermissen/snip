package main

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-lambda-go/events"
)


func MakeIndex(body string) func(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context) (events.APIGatewayProxyResponse, error) {
        headers:= map[string]string{"Content-Type": "text/html"}

        return events.APIGatewayProxyResponse{
            StatusCode:200,
            Headers         :headers,
            Body : body,
       },nil
	}
}


var body = `<html>
<head>
  <title>Demo Hybrid API</title>
</head>
<body>
<h1>Demo Hybrid API</h1>
<p>
This is an example of an API that combines HTTP mechanisms with alternative feed
connector approaches (AWS Kinesis, SOAP)
</p>
<p>
The API provides a <a href="https://mnot.github.io/I-D/json-home/">JSON-Home</a> machine reaadable <a href="/prod/home">
discovery document</a>.
</p>
<a name="doc"></a>
<h2>Primary Use Case</h2>
<p>
This API provides the ability to replicate a set of data items from the source system. The API is event based,
meaning that each item provided by the API represents the occurrence of a change to an item managed by the
source system. By consuming the events, the client system is able to locally create a copy of the items
managed in the source system. Item meta data can be used to handle only those data items, the client system
is interested in.
</p>
<p>
</p>
</p>
<h2>Ordering Guarantees</h2>
<p>
<b>1. Per-Item Order:</b> Managed items have unique identifiers and item events reference the changed item using this identifier.
Events that pertain to the same data item are guranteed to be sent in the orderthey occurred, however,
the overall ordering of events may be different from the order in which they happened.
</p>
<p>
<b>2. Referential Integrity:</b>: If data items contain a reference to other data items, the referenced data items
will always have been placed into the queue. (Note: another approach to solbve this would be to make individual items
available using a simple GET-item API, so that the consumer can resolve an integrity problem itself by obtaining the
item directly. This would also remove the problem that the referenced item may in fact not have been changed and thus
should not be in the queue at all.
</p>
<h2>Message Format and Event Metadata</h2>
<p>
Data item events are represented as MIME messages to allow for rich meta data with each item. 
<br/>
<br/>
MIME Message headers are used to convey item meta data. TBD: provide hint at headers to be
expected.

</p>

<h2>The REST API</h2>
<p>
For the most loosely coupled way of integrating with the API is using the REST feed variant of the API.
It provides a page oriented way of iterating over batches of events from a given starting point to the
latest event. Pages are connected using the "next" link relation based Link headers.
<br/>
<br/>
Pages are represented using the multipart/mixed media type so that items can be processed iteratively,
without the need to apply a JSON for XML parser to the overall batch.
</p>
<h2>The SOAP API</h2>
<p>
For consumers that work natively with SOAP, the events are also provided as a SOAP Web Service. The link
relation "soap" can be used to discover from the home document the WSDL that describes the SOAP API.
</p>
<h2>The Kinesis-Stream API</h2>
<p>
Using the link relation "kinesis-stream" clients can discover from the home document the URI of the AWS Kinesis
stream that contains all the events. This stream canb be processed with the usinal Kinesis client SDKs.
</p>

</body>
</html>
`

func main() {

        lambda.Start(MakeIndex(body))
}
