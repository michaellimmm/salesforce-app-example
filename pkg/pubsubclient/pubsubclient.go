package pubsubclient

import (
	"context"
	"crypto/x509"
	"fmt"
	"github/michaellimmm/salesforce-app-example/gen/pubsubapi"
	"io"
	"log"
	"os"

	"github.com/linkedin/goavro/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	appetite int32 = 5

	tokenHeader    = "accesstoken"
	instanceHeader = "instanceurl"
	tenantHeader   = "tenantid"
)

type (
	PubSubClient struct {
		logger       *zap.Logger
		conn         *grpc.ClientConn
		pubSubClient pubsubapi.PubSubClient
		schemaCache  map[string]*goavro.Codec
	}

	Auth struct {
		AccessToken string
		InstanceUrl string
		OrgID       string
	}
)

func NewPubSubClient(logger *zap.Logger) *PubSubClient {
	dialOpts := []grpc.DialOption{}

	grpcEndpoint := os.Getenv("SALESFORCE_GRPC_ENDPOINT")
	certs := getCerts()
	creds := credentials.NewClientTLSFromCert(certs, "")
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.DialContext(context.Background(), grpcEndpoint, dialOpts...)
	if err != nil {
		logger.Fatal("failed to connect salesforce", zap.Error(err))
	}

	return &PubSubClient{
		logger:       logger,
		conn:         conn,
		pubSubClient: pubsubapi.NewPubSubClient(conn),
	}
}

func getCerts() *x509.CertPool {
	if certs, err := x509.SystemCertPool(); err == nil {
		return certs
	}

	return x509.NewCertPool()
}

func (c *PubSubClient) Close() {
	c.conn.Close()
}

func (p *PubSubClient) GetTopic(
	ctx context.Context,
	auth Auth,
	topic string) (*pubsubapi.TopicInfo, error) {
	var trailer metadata.MD
	req := &pubsubapi.TopicRequest{
		TopicName: topic,
	}

	newCtx := p.getAuthContext(ctx, auth)

	resp, err := p.pubSubClient.GetTopic(newCtx, req, grpc.Trailer(&trailer))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *PubSubClient) GetSchema(
	ctx context.Context,
	auth Auth,
	schemaID string) (*pubsubapi.SchemaInfo, error) {
	var trailer metadata.MD

	newCtx := p.getAuthContext(ctx, auth)

	req := &pubsubapi.SchemaRequest{
		SchemaId: schemaID,
	}

	res, err := p.pubSubClient.GetSchema(newCtx, req, grpc.Trailer(&trailer))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PubSubClient) Subscribe(
	ctx context.Context,
	auth Auth,
	topicName string,
	replayPreset pubsubapi.ReplayPreset,
	replayID []byte) ([]byte, error) {
	subscribeClient, err := p.pubSubClient.Subscribe(ctx)
	if err != nil {
		return replayID, err
	}

	initialFetchRequest := &pubsubapi.FetchRequest{
		TopicName:    topicName,
		ReplayPreset: replayPreset,
		NumRequested: appetite,
	}
	if replayPreset == pubsubapi.ReplayPreset_CUSTOM && replayID != nil {
		initialFetchRequest.ReplayId = replayID
	}

	err = subscribeClient.Send(initialFetchRequest)
	if err == io.EOF {
		p.logger.Info("WARNING - EOF error returned from initial Send call, proceeding anyway")
	} else if err != nil {
		return replayID, err
	}

	requestedEvents := initialFetchRequest.NumRequested

	// NOTE: the replayID should ve stored in a persistent dadta store rather than being stored in a variable
	log.Printf("replayID: %s", string(replayID))
	curReplayID := replayID
	for {
		resp, err := subscribeClient.Recv()
		if err == io.EOF {
			return curReplayID, fmt.Errorf("stream closed")
		} else if err != nil {
			return curReplayID, err
		}

		for _, event := range resp.Events {
			codec, err := p.fetchCodec(ctx, auth, event.GetEvent().GetSchemaId())
			if err != nil {
				return curReplayID, err
			}

			parsed, _, err := codec.NativeFromBinary(event.GetEvent().GetPayload())
			if err != nil {
				return curReplayID, err
			}

			body, ok := parsed.(map[string]interface{})
			if !ok {
				return curReplayID, fmt.Errorf("error casting parsed event: %v", body)
			}

			curReplayID = event.GetReplayId()

			log.Printf("event body: %+v\n", body)

			requestedEvents--
			if requestedEvents < appetite {
				fetchRequest := &pubsubapi.FetchRequest{
					TopicName:    topicName,
					NumRequested: appetite,
				}

				err = subscribeClient.Send(fetchRequest)
				if err == io.EOF {
					// If the Send call returns an EOF error then print log
					p.logger.Info("Warning - EOF error returned for subsequent Send call, proceeding anyway")
				} else if err != nil {
					return curReplayID, err
				}

				requestedEvents += fetchRequest.NumRequested
			}
		}
	}
}

// Unexported helper function to retrieve the cached codec from the PubSubClient's schema cache. If the schema ID is not found in the cache
// then a GetSchema call is made and the corresponding codec is cached for future use
func (p *PubSubClient) fetchCodec(ctx context.Context, auth Auth, schemaId string) (*goavro.Codec, error) {
	codec, ok := p.schemaCache[schemaId]
	if ok {
		return codec, nil
	}

	schema, err := p.GetSchema(ctx, auth, schemaId)
	if err != nil {
		return nil, err
	}

	codec, err = goavro.NewCodec(schema.GetSchemaJson())
	if err != nil {
		return nil, err
	}

	p.schemaCache[schemaId] = codec

	return codec, nil
}

func (p *PubSubClient) getAuthContext(ctx context.Context, auth Auth) context.Context {
	return metadata.NewOutgoingContext(
		ctx, metadata.Pairs(
			tokenHeader, auth.AccessToken,
			instanceHeader, auth.InstanceUrl,
			tenantHeader, auth.OrgID))
}
