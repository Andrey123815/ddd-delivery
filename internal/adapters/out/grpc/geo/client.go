package geo

import (
	"context"
	"fmt"
	"time"

	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conn *grpc.ClientConn
	pbGeoClient geopb.GeoClient
	timeout time.Duration
}

func NewClient(host string) (ports.GeoClient, error) {
	if host == "" {
		return nil, fmt.Errorf("geo: host is required")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	pbGeoClient := geopb.NewGeoClient(conn)

	return &client{
		conn:        conn,
		pbGeoClient: pbGeoClient,
		timeout:     5 * time.Second,
	}, nil
}

func (c *client) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	req := &geopb.GetGeolocationRequest{Street: street}

	resp, err := c.pbGeoClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}
	if resp == nil || resp.GetLocation() == nil {
		return kernel.Location{}, fmt.Errorf("geo: пустой ответ GetGeolocation")
	}

	location := resp.GetLocation()

	return kernel.NewLocation(uint8(location.GetX()), uint8(location.GetY()))
}

func (c *client) Close() error {
	return c.conn.Close()
}
