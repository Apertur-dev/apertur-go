package apertur

import "context"

// StatsResource provides access to dashboard-level statistics.
type StatsResource struct {
	http *httpClient
}

// Get retrieves the current dashboard statistics.
func (s *StatsResource) Get(ctx context.Context) (*Stats, error) {
	var result Stats
	if err := s.http.request(ctx, "GET", "/api/v1/stats", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
