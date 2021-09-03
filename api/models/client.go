package models

const (
	protocol = "http"
	host     = "localhost"
	port     = "3000"
)

func GetFrontEndUrl() string {
	return protocol + "://" + host + ":" + port
}
