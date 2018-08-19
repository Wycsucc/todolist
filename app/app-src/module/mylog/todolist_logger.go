package mylog

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogFormat Nginx日志
type LogFormat struct {
	Time      string `json:"time"`
	Clientip  string `json:"clientip,omitempty"`
	RemoteIP  string `json:"remote_ip"`
	Method    string `json:"method"`
	URI       string `json:"uri"`
	Query     string `json:"query"`
	Proto     string `json:"proto"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	Host      string `json:"host"`
	Status    int    `json:"status"`
	Latency   int64  `json:"latency"`
	BytesIn   int64  `json:"bytes_in"`
	BytesOut  int64  `json:"bytes_out"`
	Origin    string `json:"origin"`
}

// Marshal return encoding string of format
func (format *LogFormat) Marshal() (string, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(format); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Log format which can be constructed using the following tags:
//
// - time_rfc3339
// - remote_ip
// - uri
// - host
// - method
// - path
// - referer
// - user_agent
// - status
// - latency (In microseconds)
// - latency_human (Human readable)
// - bytes_in (Bytes received)
// - bytes_out (Bytes sent)
//
// Example "${remote_ip} ${status}"

// Logger func
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		begin := time.Now()

		proto := c.Request.Proto
		host := c.Request.Host
		method := c.Request.Method
		uri := c.Request.URL.Path
		remoteIP := c.Request.RemoteAddr
		sli := strings.Split(remoteIP, ":")
		remoteIP = sli[0]

		userIP := c.Request.Header.Get("X-Forwarded-For")
		if userIP == "" {
			userIP = c.Request.Header.Get("X-Real-Ip")
		}

		bytesIn := c.Request.ContentLength
		query := c.Request.URL.RawPath

		userAgent := ""
		if val, ok := c.Request.Header["User-Agent"]; ok {
			userAgent = strings.Join(val, "")
		}
		referer := ""
		if val, ok := c.Request.Header["Referer"]; ok {
			referer = strings.Join(val, "")
		}
		// before request
		c.Next()
		// after request
		latency := time.Now().Sub(begin).Nanoseconds() / 1000000
		// access the status we are sending
		status := c.Writer.Status()
		bytesOut := int64(c.Writer.Size())

		format := &LogFormat{}
		format.Time = begin.Format("2018-08-19 19:52:55")
		format.RemoteIP = remoteIP
		format.Clientip = userIP
		format.UserAgent = userAgent
		format.Method = method
		format.URI = uri
		format.Query = query
		format.Proto = proto
		format.Referer = referer
		format.Host = host
		format.Status = status
		format.Latency = latency
		format.BytesIn = bytesIn
		format.BytesOut = bytesOut

		if result, err := format.Marshal(); err != nil {
			mylog := GetMylog().Access()
			mylog.Info(result)
		}
	}
}
