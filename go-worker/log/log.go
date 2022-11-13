package log

import "github.com/sirupsen/logrus"

func HttpResponseLogger(url, method, body string, statusCode int) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"url": url, "method": method, "status": statusCode})
}
