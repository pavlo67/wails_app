package logadapter

import (
	"fmt"
	"log"

	"github.com/pavlo67/base_go/lib/data"
	"github.com/pavlo67/base_go/lib/logger"
)

type Logger struct {
	key  string
	path string
}

func New(key string) *Logger { return &Logger{key: key} }

func (l *Logger) Debug(args ...interface{})                 { log.Print(args...) }
func (l *Logger) Debugf(t string, args ...interface{})       { log.Printf(t, args...) }
func (l *Logger) Info(args ...interface{})                  { log.Print(args...) }
func (l *Logger) Infof(t string, args ...interface{})        { log.Printf(t, args...) }
func (l *Logger) Warn(args ...interface{})                  { log.Print(args...) }
func (l *Logger) Warnf(t string, args ...interface{})        { log.Printf(t, args...) }
func (l *Logger) Error(args ...interface{})                 { log.Print(args...) }
func (l *Logger) Errorf(t string, args ...interface{})       { log.Printf(t, args...) }
func (l *Logger) Fatal(args ...interface{})                 { log.Fatal(args...) }
func (l *Logger) Fatalf(t string, args ...interface{})       { log.Fatalf(t, args...) }
func (l *Logger) Comment(text string)                       { log.Print(text) }
func (l *Logger) SetKey(key string)                         { l.key = key }
func (l *Logger) Key() string                               { return l.key }
func (l *Logger) SetPath(path string)                       { l.path = path }
func (l *Logger) Path() string                              { return l.path }
func (l *Logger) File(path string, append bool, data []byte) { log.Printf("logger.File(%s, append=%v, len=%d)", path, append, len(data)) }
func (l *Logger) Image(path string, getImage logger.GetImage, opts data.Map) {
	log.Printf("logger.Image(%s, bounds=%v)", path, getImage.Bounds())
}

func (l *Logger) String() string { return fmt.Sprintf("Logger{%s}", l.key) }
