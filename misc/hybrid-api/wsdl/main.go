package main

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-lambda-go/events"
)


func MakeIndex(body string) func(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	return func(ctx context.Context) (events.APIGatewayProxyResponse, error) {
        headers:= map[string]string{"Content-Type": "application/wsdl+xml"}

        return events.APIGatewayProxyResponse{
            StatusCode:200,
            Headers         :headers,
            Body : body,
       },nil
	}
}


var body = `<?xml version="1.0"?>

<!-- root element wsdl:definitions defines set of related services -->
<wsdl:definitions name="StockQuote"
             targetNamespace="http://example.com/stockquote.wsdl"
             xmlns:tns="http://example.com/stockquote.wsdl"
             xmlns:xsd1="http://example.com/stockquote.xsd"
             xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
             xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">

    <!-- wsdl:types encapsulates schema definitions of communication types; here using xsd -->
    <wsdl:types>

       <!-- all type declarations are in a chunk of xsd -->
       <xsd:schema targetNamespace="http://example.com/stockquote.xsd"
              xmlns:xsd="http://www.w3.org/2000/10/XMLSchema">


           <!-- xsd definition: TradePriceRequest [... tickerSymbol string ...] -->
           <xsd:element name="TradePriceRequest">
              <xsd:complexType>
                  <xsd:all>
                      <xsd:element name="tickerSymbol" type="string"/>
                  </xsd:all>
              </xsd:complexType>
           </xsd:element>


           <!-- xsd definition: TradePrice [... price float ...] -->
           <xsd:element name="TradePrice">
              <xsd:complexType>
                  <xsd:all>
                      <xsd:element name="price" type="float"/>
                  </xsd:all>
              </xsd:complexType>
           </xsd:element>

       </xsd:schema>
    </wsdl:types>

    <!-- request GetLastTradePriceInput is of type TradePriceRequest -->
    <wsdl:message name="GetLastTradePriceInput">
        <wsdl:part name="body" element="xsd1:TradePriceRequest"/>
    </wsdl:message>

    <!-- request GetLastTradePriceOutput is of type TradePrice -->
    <wsdl:message name="GetLastTradePriceOutput">
        <wsdl:part name="body" element="xsd1:TradePrice"/>
    </wsdl:message>

    <!-- wsdl:portType describes messages in an operation -->
    <wsdl:portType name="StockQuotePortType">

    <!-- the value of wsdl:operation eludes me -->
        <wsdl:operation name="GetLastTradePrice">
           <wsdl:input message="tns:GetLastTradePriceInput"/>
           <wsdl:output message="tns:GetLastTradePriceOutput"/>
        </wsdl:operation>
    </wsdl:portType>

    <!-- wsdl:binding states a serialization protocol for this service -->
    <wsdl:binding name="StockQuoteSoapBinding"
                  type="tns:StockQuotePortType">

        <!-- leverage off soap:binding document style @@@(no wsdl:foo pointing at the soap binding) -->
        <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http"/>

        <!-- semi-opaque container of network transport details classed by soap:binding above @@@ -->
        <wsdl:operation name="GetLastTradePrice">

           <!-- again bind to SOAP? @@@ -->
           <soap:operation soapAction="http://example.com/GetLastTradePrice"/>
           <!-- furthur specify that the messages in the wsdl:operation "" use SOAP? @@@ -->
           <wsdl:input>
               <soap:body use="literal"/>
           </wsdl:input>
           <wsdl:output>
               <soap:body use="literal"/>
           </wsdl:output>
        </wsdl:operation>
    </wsdl:binding>

    <!-- wsdl:service names a new service "StockQuoteService" -->
    <wsdl:service name="StockQuoteService">
    <wsdl:documentation>My first service</wsdl:documentation>

        <!-- connect it to the binding "StockQuoteBinding" above -->
        <wsdl:port name="StockQuotePort"
                   binding="tns:StockQuoteBinding">

           <!-- give the binding an network address -->
           <soap:address location="http://example.com/stockquote"/>
        </wsdl:port>
    </wsdl:service>

</wsdl:definitions>
`

func main() {

        lambda.Start(MakeIndex(body))
}
