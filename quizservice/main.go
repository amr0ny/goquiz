package quizservice

import (
	"fmt"
	"github.com/amr0ny/goquiz/common/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

const defaultServicePort = "80001"
const defaultPprofPort = "8081"

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func initPprof(pprofPort string) {
	log, err := common.GetLogger()
	if err != nil {
		fmt.Println(err)
	}
	addr := fmt.Sprintf("localhost: %v", pprofPort)
	go func() {
		log.Infoln("Starting pprof on %v", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Errorf("failed running pprof: %v", err)
		}
	}()
}

func main() {
	err := common.InitLoggerConfig(&common.Config{Filepath: getEnv("LOG_PATH", defaultLogPath)})
	if err != nil {
		fmt.Println(err)
	}
	log, err := common.GetLogger()
	if err != nil {
		fmt.Println(err)
	}
	if strings.ToLower(getEnv("DEBUG", "")) == "true" {
		log.Level = logrus.DebugLevel
	}

	if strings.ToLower(getEnv("PROFILING_DISABLED", "FALSE")) == "false" {
		pprofPort := getEnv("PPROF_PORT", defaultPprofPort)
		initPprof(pprofPort)
	}
	port := getEnv("SERVICE_PORT", defaultServicePort)

}
