package logiclog

import (
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as Json instead of the default ASCII formatter.
	formatter := Formatter{
		ChildFormatter: &log.JSONFormatter{},
		Line:           true,
		Package:        false,
		File:           true,
		BaseNameOnly:   false,
	}
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var ServerName string
var Environment string
var Instancekey string

func InitConfig(serverName, env, instanceKey, level string) {
	ServerName = serverName
	Environment = env
	Instancekey = instanceKey

	formatLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	}
	log.SetLevel(formatLevel)
}

func CtxLogger(ctx iris.Context) *log.Entry {
	user := "-"
	if len(ctx.GetHeader("X-User-Name")) > 0 {
		user = ctx.GetHeader("X-User-Name")
	}

	reqId := ctx.GetID()

	contextLogger := log.WithFields(log.Fields{
		"user":         user,
		"req_id":       reqId,
		"server_name":  ServerName,
		"environment":  Environment,
		"instance_key": Instancekey,
		"logic_type":   "logic",
	})
	return contextLogger
}

func Logger() *log.Entry {
	contextLogger := log.WithFields(log.Fields{
		"server_name":  ServerName,
		"environment":  Environment,
		"instance_key": Instancekey,
		"logic_type":   "logic",
	})

	return contextLogger
}
