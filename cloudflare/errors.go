package cloudflare

import "github.com/sirupsen/logrus"

func HandlerErrors(args ...interface{}) {
	logrus.WithFields(logrus.Fields{"error": args}).Error("Error occurred")
}
