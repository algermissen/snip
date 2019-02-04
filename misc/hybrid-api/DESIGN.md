Data Replication Message Design

Data Replication APIs at HSE24 should leverage existing Web standards as much as possible, even if they are not realized primarily as REST APIs. Web standards are naturally designed to emphasise loose coupling because it is imperative for the success of the Web that  its components (the Web sites, the browsers, etc.) can change it independent intervals. Thus, when aiming for loosely coupled systems, it is allways a good idea to look at what is in use on the Web today already.

HSE24 aims for using an event log approach for replication scenarios, where changes to the items of a set of resources are published as change events on an event log.

We can leverage a variety of existing Web standards and their existing open source implementation by applying te concept of a MIME message to representating the change events. Every change event message in HSE24 data replication APIs will then be known to have the general MIME message syntax, as defined in RFC 2045.  https://tools.ietf.org/html/rfc2045

This format should generally be familiar; it is te format used by Email or HTTP requests/responses, too. An example MIME message looks like this

Content-Type: text/html
Content-Length: 2043

<html>
  <title>...</title>
  <body>...</body>
</html>

We can build on the combination of meta data and body in a single message, and make use of a number of standard MIME headers to describe the change events and how to process them. Here is an extended example:


Message-ID: e1857a7c-9b78-4317-9f59-236d241013e3
Content-Type: application/vnd.hse24.product+json
Content-Length: 183
Slug: instant-ingwer-getraenk-extrastark-371116
Content-Location: http://api.mercator.hse24.de/products/371116
Etag: "1234567890"
Control: PUT 
Distribution: DE, FR
Link: <http://api.mercator.hse24.de/products/371116>; rel=self

{
  "id": "371116",
  "title" : "Instant Ingwer Getränk extrastark"
  "description": "Freuen Sie sich auf ein voluminöses raumfüllendes Aromaereignis! ..."
  ... more fields ...
}


... exlplain in table ..
Message-ID: e1857a7c-9b78-4317-9f59-236d241013e3
Content-Type: application/vnd.hse24.product+json
Content-Length: 183
Slug: instant-ingwer-getraenk-extrastark-371116
Content-Location: http://api.mercator.hse24.de/products/371116
Etag: "1234567890"
Control: PUT 
Distribution: DE, FR
Link: <http://api.mercator.hse24.de/products/371116>; rel=self


Using the Control header, we can also easily solve the problem of how to represent deleted items so that the consumer can react on the deletion and remove the item from its own collection:

Message-ID: 7704fb34-be3a-45cb-953a-af92477a0f83
Link: <http://api.mercator.hse24.de/products/371116>; rel=self
Control: DELETE 


The use of the Content-Type header leverages standard MIME mechanisms for allowing different kinds of represented items to be mixed in a single event log, because the content type information allows the consumer to distinguish on a per-item basis its interpretation. In addition MIME content type information allows the evolution of content formats over time if some additional HTTP mechanims are leveraged. (See below PROVIDE EXAMPLE)

For example, the following message references an alternative representation using a Link header with relation "alternate". Consumers can actively choose whether they can process the content as provided in the message or need to access the alternate representation (which might be an older format kept for a grace period).

Message-ID: e1857a7c-9b78-4317-9f59-236d241013e3
Content-Type: application/vnd.hse24.nextgen-product+json
Content-Length: 183
Slug: instant-ingwer-getraenk-extrastark-371116
Content-Location: http://api.mercator.hse24.de/products/371116
Etag: "1234567890"
Control: PUT 
Distribution: DE, FR
Link: <http://api.mercator.hse24.de/products/371116>; rel=self
Link: <http://api.mercator.hse24.de/products/371116>; rel=alternate; type="application/vnd.hse24.product+json"

{
  "id": "371116",
  "title" : "Instant Ingwer Getränk extrastark"
  "description": "Freuen Sie sich auf ein voluminöses raumfüllendes Aromaereignis! ..."
  ... more fields ...
}


Using in snapshots

In addition to the pure event log scenario (and the single HTTP request/response use case), the MIME Message based approach can also very naturally be combined with the mulipart media type family to combine multiple events into a single messages. This can be useful, for example, for snapshot scenarios and REST-based paged event feeds.

This is how the syntax of such a multipart message would look (note the boundary identifier --abcdefgh)
    

Content-Type: multipart/mixed; boundary=abcdefgh

--abcdefgh 
Message-ID: e1857a7c-9b78-4317-9f59-236d241013e3
Content-Type: application/vnd.hse24.product+json
Content-Length: 183
Slug: instant-ingwer-getraenk-extrastark-371116
Content-Location: http://api.mercator.hse24.de/products/371116
Etag: "1234567890"
Control: PUT 
Distribution: DE, FR
Link: <http://api.mercator.hse24.de/products/371116>; rel=self

{
  "id": "371116",
  "title" : "Instant Ingwer Getränk extrastark"
  "description": "Freuen Sie sich auf ein voluminöses raumfüllendes Aromaereignis! ..."
  ... more fields ...
}

--abcdefgh 
Message-ID: e1857a7c-9b78-4317-9f59-236d241013e3
Content-Type: application/vnd.hse24.product+json
Content-Length: 183
Slug: instant-ingwer-getraenk-extrastark-371116
Content-Location: http://api.mercator.hse24.de/products/371116
Etag: "1234567890"
Control: PUT 
Distribution: DE, FR
Link: <http://api.mercator.hse24.de/products/371116>; rel=self

{
  "id": "371116",
  "title" : "Instant Ingwer Getränk extrastark"
  "description": "Freuen Sie sich auf ein voluminöses raumfüllendes Aromaereignis! ..."
  ... more fields ...
}

--abcdefgh-- 

For periodic snapshots of the full item set, item representations could be saved in sets of 10000, using such multipart documents. One million items, for example, would yield 100 of such snapshot parts.

For REST-based event feeds, events can be packaged up in the same way to form a sequence of event sets that can be linked for page-oriented processing of the event history. Object stores like S3 make a perfect and simple storage mechanism for such historic (and hence immutable) sets of events.








