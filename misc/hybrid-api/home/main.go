package main

import (
	"encoding/json"
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-lambda-go/events"
	"fmt"
)

// JsonHome is an implementation of
// https://tools.ietf.org/html/draft-nottingham-json-home-04.
type JsonHome struct {
	Resources map[string]Resource `json:"resources"`
}

// https://tools.ietf.org/html/draft-nottingham-json-home-04#section-3
type Resource struct {
	Rel          string       `json:"-"`
	Href         string       `json:"href,omitempty"`
	HrefTemplate string       `json:"href-template,omitempty"`
	HrefVars     HrefVarsType `json:"href-vars,omitempty"`
	Hints        Hints        `json:"hints,omitempty"`
}

//
type HrefVarsType map[string]string

// https://tools.ietf.org/html/draft-nottingham-json-home-04#section-4
type Hints struct {
	Allow []string `json:"allow,omitempty"`
	//Formats TBD https://tools.ietf.org/html/draft-nottingham-json-home-04#section-5
	AcceptPatch     []string  `json:"accept-patch,omitempty"`
	AcceptPost      []string  `json:"accept-patch,omitempty"`
	AcceptRanges    []string  `json:"accept-ranges,omitempty"`
	AcceptPrefer    []string  `json:"accept-prefer,omitempty"`
	Docs            string    `json:"docs,omitempty"`
	PreconditionReq []string  `json:"precondition-req,omitempty"`
	AuthReq         []AuthReq `json:"auth-req,omitempty"`
	Status          []string  `json:"status,omitempty"`
}

// https://tools.ietf.org/html/draft-nottingham-json-home-04#section-4.9
type AuthReq struct {
	Scheme string   `json:"scheme,omitempty"`
	Realms []string `json:"realms,omitempty"`
}

// NewResource creates a resource with a fixed href (and not an href-template)
func NewResource(rel string, href string) Resource {
	return Resource{Rel: rel, Href: href}
}

// NewResource creates a resource with an href-template.
func NewTemplatedResource(rel string, hrefTemplate string, vars HrefVarsType) Resource {
	return Resource{Rel: rel, HrefTemplate: hrefTemplate, HrefVars: vars}
}

// NewJsonHome creates a new JSON home document structure from the provided
// set of resources.
// Resources will be stored in a map with the link relations used as map key.
// This means that, for resources with duplicate link relations, only one
// resource will be stored in the map.
func NewJsonHome(resources ...Resource) JsonHome {

	var jh = JsonHome{make(map[string]Resource)}
	for _, r := range resources {
		jh.Resources[r.Rel] = r
	}

	return jh
}

// GetResource is used for looking up a resource based on a provided
// link relation. If there is no resource for that link relation nil
// is returned.
func (jh *JsonHome) GetResource(rel string) *Resource {
	r, ok := jh.Resources[rel]
	if ok {
		return &r
	} else {
		return nil
	}
}

// MakeHome takes a list of resources and produces a lambda Handler
// that responds with a JSON home document containing the resources.
// MakeHome panics if the JSON home document cannot be created as a
// string.
func MakeHome(resources ...Resource) func(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	jh := NewJsonHome(resources...)
	bytes, err := json.Marshal(&jh)
	if err != nil {
		panic(fmt.Sprintf("Cannot serialize JSON home %v, %v", jh, err))
	}
	jsonHome := string(bytes)

	return func(ctx context.Context) (events.APIGatewayProxyResponse, error) {
        headers:= map[string]string{"Content-Type": "application/json-home"}

        return events.APIGatewayProxyResponse{
            StatusCode:200,
            Headers         :headers,
            Body : jsonHome,
       },nil
	}
}


// String returns string representation of JSON Home object. This
// is intended to make mocking for tests simpler.
// Will panic if the JSON home cannot be serialized.
func (jh *JsonHome) String() string {
	bytes, err := json.Marshal(jh)
	if err != nil {
		panic(fmt.Sprintf("Cannot serialize JSON home %v, %v", jh, err))
	}
	return string(bytes)
}

// ResourceBuilder provides a fluent builder API for resources.
type ResourceBuilder struct {
	Rel          string
	Href         string
	HrefTemplate string
	HrefVars     HrefVarsType
	Hints        Hints
}

// NewResourceBuilder creates a new ResourceBuilder using the
// specified relation for the to be build resource.
func NewResourceBuilder(rel string) *ResourceBuilder {
	return &ResourceBuilder{
		Rel: rel,
	}
}

// WithHref sets the href-property for the built resource.
// It is the caller's responsibility not to call this together
// with WithHrefTemplate(); there is no internal checking for
// this.
func (rb *ResourceBuilder) WithHref(href string) *ResourceBuilder {
	rb.Href = href
	return rb
}

// WithHrefTemplate sets the href-template and vars properties for the built resource.
// It is the caller's responsibility not to call this together
// with WithHref(); there is no internal checking for this.
// The provided vars strings are interpreted as var,value,var,value,...
// If an odd number of variable arguments is provided, the last var will
// be silently ignored.
func (rb *ResourceBuilder) WithHrefTemplate(hrefTemplate string, vars ...string) *ResourceBuilder {

	rb.HrefTemplate = hrefTemplate
	rb.HrefVars = HrefVarsType{}
	for i := 0; i < len(vars)-1; i += 2 {
		rb.HrefVars[vars[i]] = vars[i+1]
	}

	return rb
}

// WithHintAuthReq sets a single auth-req hint for the built resource.
func (rb *ResourceBuilder) WithHintAuthReq(scheme string, realms ...string) *ResourceBuilder {
	rb.Hints.AuthReq = []AuthReq{
		AuthReq{
			Scheme: scheme,
			Realms: realms,
		},
	}
	return rb
}

// WithHintAllow sets allow property for the built resource
func (rb *ResourceBuilder) WithHintAllow(allow ...string) *ResourceBuilder {
	rb.Hints.Allow = allow
	return rb
}

// WithHintPreconditionEtag sets the precondition-req property
// of the build resource to ["etag"]
func (rb *ResourceBuilder) WithHintPreconditionEtag() *ResourceBuilder {
	rb.Hints.PreconditionReq = []string{"etag"}
	return rb
}

// WithHintPreconditionLastModified sets the precondition-req property
// of the build resource to ["last-modified"]
func (rb *ResourceBuilder) WithHintPreconditionLastModified() *ResourceBuilder {
	rb.Hints.PreconditionReq = []string{"last-modified"}
	return rb
}

// WithHintDocs sets the docs property of the built resource
func (rb *ResourceBuilder) WithHintDocs(docs string) *ResourceBuilder {
	rb.Hints.Docs = docs
	return rb
}

// Build the resource from this builder.
func (rb *ResourceBuilder) Build() Resource {
	return Resource{
		Rel:          rb.Rel,
		Href:         rb.Href,
		HrefTemplate: rb.HrefTemplate,
		HrefVars:     rb.HrefVars,
		Hints:        rb.Hints,
	}
}



//func HandleRequest(ctx context.Context) (string, error) {
//        return fmt.Sprintf("Hello"), nil
//}

func main() {
        resources := []Resource{
		NewResourceBuilder("kinesis-stream").
			WithHref("arn:aws:kinesis:eu-central-1:627211582084:stream/productevents").
			WithHintDocs("/prod/#doc").
			Build(),
		NewResourceBuilder("soap").
			WithHref("/prod/soap/wsdl").
			WithHintDocs("/prod/#doc").
			Build(),
		NewResourceBuilder("http").
        WithHrefTemplate("/prod/http/latest{?since}", "since", "http://registry.hse24.de/vars/since").
			WithHintAllow("GET").WithHintAuthReq("Basic").
			WithHintDocs("/prod/#doc").
			Build(),
	}

        lambda.Start(MakeHome(resources...))
}
