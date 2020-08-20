package util

import "os"

const (
	ENV_KEY_HOST = "JPAAS_HOST"
	//ENV_KEY_PORT          = "JPAAS_HTTP_PORT"
	//ENV_KEY_PORT_ORIGINAL = "JPAAS_HOST_PORT_8080"
)

func GetHostAndE() (string, int) {
	DOCKER_HOST := os.Getenv(ENV_KEY_HOST)
	//DOCKER_PORT := os.Getenv(ENV_KEY_PORT)
	//if DOCKER_PORT == "" {
	//	DOCKER_PORT = os.Getenv(ENV_KEY_PORT_ORIGINAL)
	//}
	if DOCKER_HOST != "" {
		return DOCKER_HOST, 1
	} else {
		hostname, err := os.Hostname()
		if err != nil {
			return hostname, 2
		}
		return "127.0.0.1", 3
	}
}
