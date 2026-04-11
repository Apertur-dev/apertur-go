package apertur

import (
	"context"
	"net/url"
	"strconv"
)

// UploadsResource provides methods for listing uploaded images.
type UploadsResource struct {
	http *httpClient
}

// List returns a paginated list of uploaded images.
func (u *UploadsResource) List(ctx context.Context, params ListParams) (*UploadsListPage, error) {
	path := "/api/v1/uploads"
	qs := url.Values{}
	if params.Page != nil {
		qs.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PageSize != nil {
		qs.Set("pageSize", strconv.Itoa(*params.PageSize))
	}
	if encoded := qs.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var result UploadsListPage
	if err := u.http.request(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Recent returns the most recently uploaded images.
func (u *UploadsResource) Recent(ctx context.Context, params LimitParams) ([]UploadRow, error) {
	path := "/api/v1/uploads/recent"
	if params.Limit != nil {
		path += "?limit=" + strconv.Itoa(*params.Limit)
	}

	var result []UploadRow
	if err := u.http.request(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
