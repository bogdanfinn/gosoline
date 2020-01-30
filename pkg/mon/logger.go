package mon

import (
	"context"
	"fmt"
	"github.com/applike/gosoline/pkg/tracing"
	"github.com/getsentry/raven-go"
	"github.com/jonboulle/clockwork"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	Trace = "trace"
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
	Fatal = "fatal"
	Panic = "panic"
)

type LoggerSettings struct {
	Level           string                 `cfg:"level" default:"info" validate:"required"`
	Format          string                 `cfg:"format" default:"console" validate:"required"`
	TimestampFormat string                 `cfg:"timestamp_format" default:"15:04:05.000" validate:"required"`
	Tags            map[string]interface{} `cfg:"tags"`
}

var levels = map[string]int{
	Trace: 0,
	Debug: 1,
	Info:  2,
	Warn:  3,
	Error: 4,
	Fatal: 5,
	Panic: 6,
}

func levelPriority(level string) int {
	return levels[level]
}

const (
	ChannelDefault   = "default"
	FormatConsole    = "console"
	FormatGelf       = "gelf"
	FormatGelfFields = "gelf_fields"
	FormatJson       = "json"
)

type Tags map[string]interface{}
type ConfigValues map[string]interface{}
type Fields map[string]interface{}
type EcsMetadata map[string]interface{}

type Metadata struct {
	channel       string
	context       context.Context
	contextFields Fields
	fields        Fields
	tags          Tags
}

type formatter func(timestamp string, level string, msg string, err error, data *Metadata) ([]byte, error)

var formatters = map[string]formatter{
	FormatConsole:    formatterConsole,
	FormatGelf:       formatterGelf,
	FormatGelfFields: formatterGelfFields,
	FormatJson:       formatterJson,
}

//go:generate mockery -name Sentry
type Sentry interface {
	Capture(packet *raven.Packet, captureTags map[string]string) (eventID string, ch chan error)
}

type GosoLog interface {
	Logger
	Option(options ...LoggerOption) error
}

//go:generate mockery -name Logger
type Logger interface {
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})
	Error(err error, msg string)
	Errorf(err error, msg string, args ...interface{})
	Fatal(err error, msg string)
	Fatalf(err error, msg string, args ...interface{})
	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	Panic(err error, msg string)
	Panicf(err error, msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	WithChannel(channel string) Logger
	WithContext(ctx context.Context) Logger
	WithFields(fields map[string]interface{}) Logger
}

type logger struct {
	clock       clockwork.Clock
	output      io.Writer
	outputLck   *sync.Mutex
	ctxResolver []ContextFieldsResolver
	hooks       []LoggerHook

	level           int
	format          string
	timestampFormat string

	data Metadata
}

func NewLogger() *logger {
	return NewLoggerWithInterfaces(clockwork.NewRealClock(), os.Stdout)
}

func NewLoggerWithInterfaces(clock clockwork.Clock, out io.Writer) *logger {
	logger := &logger{
		clock:       clock,
		output:      out,
		outputLck:   &sync.Mutex{},
		ctxResolver: make([]ContextFieldsResolver, 0),
		hooks:       make([]LoggerHook, 0),
		level:       levelPriority(Info),
		format:      FormatConsole,
		data: Metadata{
			channel:       ChannelDefault,
			contextFields: make(Fields),
			fields:        make(Fields),
			tags:          make(Tags),
		},
	}

	return logger
}

func (l *logger) copy() *logger {
	return &logger{
		clock:           l.clock,
		outputLck:       l.outputLck,
		output:          l.output,
		ctxResolver:     l.ctxResolver,
		hooks:           l.hooks,
		level:           l.level,
		format:          l.format,
		timestampFormat: l.timestampFormat,
		data:            l.data,
	}
}

func (l *logger) Option(options ...LoggerOption) error {
	for _, opt := range options {
		if err := opt(l); err != nil {
			return err
		}
	}

	return nil
}

func (l *logger) WithChannel(channel string) Logger {
	cpy := l.copy()
	cpy.data.channel = channel

	return cpy
}

func (l *logger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	cpy := l.copy()
	cpy.data.context = ctx

	for _, r := range l.ctxResolver {
		newContextFields := r(ctx)
		cpy.data.contextFields = mergeMapStringInterface(cpy.data.contextFields, newContextFields)
	}

	return cpy
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	cpy := l.copy()
	cpy.data.fields = mergeMapStringInterface(l.data.fields, fields)

	return cpy
}

func (l *logger) Info(args ...interface{}) {
	l.log(Info, fmt.Sprint(args...), nil, Fields{})
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.log(Info, fmt.Sprintf(msg, args...), nil, Fields{})
}

func (l *logger) Debug(args ...interface{}) {
	if l.level > levels[Debug] {
		return
	}

	l.log(Debug, fmt.Sprint(args...), nil, Fields{})
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	if l.level > levels[Debug] {
		return
	}

	l.log(Debug, fmt.Sprintf(msg, args...), nil, Fields{})
}

func (l *logger) Warn(args ...interface{}) {
	l.log(Warn, fmt.Sprint(args...), nil, Fields{})
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.log(Warn, fmt.Sprintf(msg, args...), nil, Fields{})
}

func (l *logger) Error(err error, msg string) {
	l.logError(Error, err, msg)
}

func (l *logger) Errorf(err error, msg string, args ...interface{}) {
	l.logError(Error, err, fmt.Sprintf(msg, args...))
}

func (l *logger) Fatal(err error, msg string) {
	l.logError(Fatal, err, msg)
	os.Exit(1)
}

func (l *logger) Fatalf(err error, msg string, args ...interface{}) {
	l.logError(Fatal, err, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (l *logger) Panic(err error, msg string) {
	l.logError(Panic, err, msg)
	panic(err)
}

func (l *logger) Panicf(err error, msg string, args ...interface{}) {
	l.logError(Panic, err, fmt.Sprintf(msg, args...))
	panic(err)
}

func (l *logger) logError(level string, err error, msg string) {
	if l.data.context != nil {
		span := tracing.GetSpan(l.data.context)

		if span != nil {
			span.AddError(err)
		}
	}

	l.log(level, msg, err, Fields{
		"stacktrace": getStackTrace(1),
	})
}

func (l *logger) log(level string, msg string, logErr error, fields Fields) {
	levelNo := levels[level]

	if levelNo < l.level {
		return
	}

	cpyData := l.data
	cpyData.fields = mergeMapStringInterface(cpyData.fields, fields)

	for _, h := range l.hooks {
		if err := h.Fire(level, msg, logErr, &cpyData); err != nil {
			l.err(err)
		}
	}

	timestamp := FormatTime(l.clock.Now(), l.timestampFormat)
	buffer, err := formatters[l.format](timestamp, level, msg, logErr, &cpyData)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	l.write(buffer)
}

func (l *logger) err(err error) {
	timestamp := FormatTime(l.clock.Now(), l.timestampFormat)
	buffer, err := formatters[l.format](timestamp, Error, err.Error(), err, &l.data)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	l.write(buffer)
}

func (l *logger) write(buffer []byte) {
	l.outputLck.Lock()
	defer l.outputLck.Unlock()

	_, err := l.output.Write(buffer)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

// getStackTrace constructs the current stacktrace. depthSkip defines how many steps of the
// stacktrace should be skipped. This is useful to not clutter the stacktrace with logging
// function calls.
func getStackTrace(depthSkip int) string {
	depthSkip = depthSkip + 1 // Skip this function in stacktrace
	maxDepth := 50
	traces := make([]string, maxDepth)

	// Get traces
	var depth int
	for depth = 0; depth < maxDepth; depth++ {
		function, _, line, ok := runtime.Caller(depth)

		if !ok {
			break
		}

		var traceStrBuilder strings.Builder
		traceStrBuilder.WriteString("\t")
		traceStrBuilder.WriteString(runtime.FuncForPC(function).Name())
		traceStrBuilder.WriteString(":")
		traceStrBuilder.WriteString(strconv.Itoa(line))
		traceStrBuilder.WriteString("\n")
		traces[depth] = traceStrBuilder.String()
	}

	// Assemble stacktrace in reverse order
	var strBuilder strings.Builder
	strBuilder.WriteString("\n")
	for i := depth; i > depthSkip; i-- {
		strBuilder.WriteString(traces[i])
	}
	return strBuilder.String()
}

func mergeMapStringInterface(receiver map[string]interface{}, input map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(receiver)+len(input))

	for k, v := range receiver {
		newMap[k] = prepareForLog(v)
	}

	for k, v := range input {
		newMap[k] = prepareForLog(v)
	}

	return newMap
}

func prepareForLog(v interface{}) interface{} {
	switch t := v.(type) {
	case error:
		// Otherwise errors are ignored by `encoding/json`
		return t.Error()

	case map[string]interface{}:
		// perform a deep copy of any maps contained in this map element
		// to ensure we own the object completely
		return mergeMapStringInterface(t, nil)

	default:
		// same as before, but handle the case of the map mapping to something
		// different than interface{}
		// should quite rarely get hit, otherwise you are using too complex objects for your logs
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map:
			iter := rv.MapRange()
			newMap := make(map[string]interface{}, rv.Len())

			for iter.Next() {
				keyValue := iter.Key()
				elemValue := iter.Value()
				newMap[fmt.Sprint(keyValue.Interface())] = prepareForLog(elemValue.Interface())
			}

			return newMap

		case reflect.Ptr, reflect.Interface:
			if rv.IsNil() {
				return nil
			}

			return prepareForLog(rv.Elem().Interface())

		case reflect.Struct:
			rvt := rv.Type()
			newMap := make(map[string]interface{}, rv.NumField())

			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				if !field.CanInterface() {
					continue
				}
				newMap[rvt.Field(i).Name] = prepareForLog(field.Interface())
			}

			return newMap

		case reflect.Slice, reflect.Array:
			if rv.Kind() == reflect.Slice && rv.IsNil() {
				return nil
			}

			newArray := make([]interface{}, rv.Len())

			for i := range newArray {
				newArray[i] = prepareForLog(rv.Index(i).Interface())
			}

			return newArray

		default:
			return v
		}
	}
}
