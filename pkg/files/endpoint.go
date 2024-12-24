package files

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"net/url"
)

type endpoint struct {
}

func (e *endpoint) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	if params.Endpoint == nil {
		// fallback to default
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	u, err := url.Parse(*params.Endpoint)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}

	if params.Bucket != nil {
		u = u.JoinPath(*params.Bucket)
	}

	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
