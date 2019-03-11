package logger

import (
	"bytes"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var levelToColor = map[logrus.Level](func(string, ...interface{}) string){
	logrus.PanicLevel: color.RedString,
	logrus.FatalLevel: color.RedString,
	logrus.ErrorLevel: color.RedString,
	logrus.WarnLevel:  color.YellowString,
	logrus.InfoLevel:  color.GreenString,
	logrus.DebugLevel: color.BlueString,
	logrus.TraceLevel: color.CyanString,
}

//CustomFormatter is a custom logrus.Formatter implementation
type CustomFormatter struct {
	DisableTimestamp bool
	DisableSeverity  bool
}

//Format formats a logrus.Entry into a byte array to view
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	if !f.DisableTimestamp {
		b.WriteString(entry.Time.Format(time.RFC3339))
		b.WriteByte(' ')
	}

	if !f.DisableSeverity {
		b.WriteString("- ")
		colorFunc := levelToColor[entry.Level]
		b.WriteString(colorFunc(entry.Level.String()))
		b.WriteByte(' ')
	}

	for _, v := range entry.Data {
		b.WriteByte('[')
		f.appendValue(b, v)
		b.WriteByte(']')
		b.WriteByte(' ')
	}

	if len(entry.Message) > 0 {
		b.WriteString(entry.Message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *CustomFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		b.WriteString(value)
	case error:
		errmsg := value.Error()
		b.WriteString(errmsg)
	default:
		fmt.Fprint(b, value)
	}
}

// time hostname service[pid]: severity [file:line] msg
