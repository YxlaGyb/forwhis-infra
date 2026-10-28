package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mattn/go-isatty"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

type LoggerConfig struct {
	Formatter LogFormatter
	Output    io.Writer
	SkipPaths []string
	Skip      Skipper
}

type Skipper func(c fiber.Ctx) bool

type LogFormatter func(params LogFormatterParams) string

type LogFormatterParams struct {
	Request      *http.Request
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
	isTerm       bool
	BodySize     int
	Keys         map[any]any
}

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode
	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return white
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func (p *LogFormatterParams) MethodColor() string {
	method := p.Method
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func (p *LogFormatterParams) ResetColor() string {
	return reset
}

func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[LM] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func DisableConsoleColor() {
	consoleColorMode = disableColor
}

func ForceConsoleColor() {
	consoleColorMode = forceColor
}

func Logger() fiber.Handler {
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithFormatter(f LogFormatter) fiber.Handler {
	return LoggerWithConfig(LoggerConfig{Formatter: f})
}

func LoggerWithWriter(out io.Writer, notlogged ...string) fiber.Handler {
	return LoggerWithConfig(LoggerConfig{Output: out, SkipPaths: notlogged})
}

func LoggerWithConfig(conf LoggerConfig) fiber.Handler {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}
	out := conf.Output
	if out == nil {
		out = os.Stdout
	}
	notlogged := conf.SkipPaths
	isTerm := true
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}
	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		raw := string(c.Request().URI().QueryString())
		nextErr := c.Next()
		if _, ok := skip[path]; ok || (conf.Skip != nil && conf.Skip(c)) {
			return nextErr
		}
		clientIP := c.IP()
		if forwardedIP, ok := c.Locals("clientIP").(string); ok && forwardedIP != "" {
			clientIP = forwardedIP
		}
		req := &http.Request{
			Method:     c.Method(),
			URL:        &url.URL{Path: path, RawQuery: raw},
			Header:     make(http.Header),
			Host:       c.Hostname(),
			RemoteAddr: clientIP,
		}
		for name, values := range c.GetReqHeaders() {
			req.Header[name] = append([]string(nil), values...)
		}
		param := LogFormatterParams{
			Request: req,
			isTerm:  isTerm,
			Keys:    make(map[any]any),
		}
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = clientIP
		param.Method = c.Method()
		param.StatusCode = c.Response().StatusCode()
		if nextErr != nil {
			param.ErrorMessage = nextErr.Error()
		}
		param.BodySize = len(c.Response().Body())
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		if _, err := fmt.Fprint(out, formatter(param)); err != nil {
			if nextErr != nil {
				return nextErr
			}
			return err
		}
		return nextErr
	}
}
