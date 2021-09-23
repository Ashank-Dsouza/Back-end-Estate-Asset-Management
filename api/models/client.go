package models

const (
	protocol = "http"
	host     = "3.108.41.85"
	port     = "3000"
)

func GetFrontEndUrl() string {
	return protocol + "://" + host + ":" + port
}
