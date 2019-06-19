package clients

type Client interface {
	GetAccessToken() error
	UpdateLastActivity(*string, *int64) error
}
