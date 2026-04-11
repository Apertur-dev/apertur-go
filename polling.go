package apertur

import (
	"context"
	"fmt"
	"time"
)

// PollingResource provides long-polling methods to retrieve uploaded images.
// Sessions must be created with LongPolling enabled.
type PollingResource struct {
	http *httpClient
}

// List returns images that are available for download via polling.
func (p *PollingResource) List(ctx context.Context, uuid string) (*PollResult, error) {
	var result PollResult
	if err := p.http.request(ctx, "GET", fmt.Sprintf("/api/v1/upload-sessions/%s/poll", uuid), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Download retrieves the raw image bytes for a polled image.
func (p *PollingResource) Download(ctx context.Context, uuid, imageID string) ([]byte, error) {
	return p.http.requestRaw(ctx, "GET", fmt.Sprintf("/api/v1/upload-sessions/%s/images/%s", uuid, imageID))
}

// Ack acknowledges that a polled image has been processed, removing it from
// subsequent poll results.
func (p *PollingResource) Ack(ctx context.Context, uuid, imageID string) (*AckResult, error) {
	var result AckResult
	path := fmt.Sprintf("/api/v1/upload-sessions/%s/images/%s/ack", uuid, imageID)
	if err := p.http.request(ctx, "POST", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PollAndProcess runs a continuous poll-download-ack loop until the context is
// cancelled. For each image found, handler is called with the image metadata
// and raw bytes. If handler returns an error, the loop stops and the error is
// returned. The default poll interval is 3 seconds.
func (p *PollingResource) PollAndProcess(ctx context.Context, uuid string, handler func(image PollImage, data []byte) error, opts PollProcessOptions) error {
	interval := time.Duration(opts.Interval) * time.Second
	if interval <= 0 {
		interval = 3 * time.Second
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result, err := p.List(ctx, uuid)
		if err != nil {
			return err
		}

		for _, image := range result.Images {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			data, err := p.Download(ctx, uuid, image.ID)
			if err != nil {
				return err
			}

			if err := handler(image, data); err != nil {
				return err
			}

			if _, err := p.Ack(ctx, uuid, image.ID); err != nil {
				return err
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}
