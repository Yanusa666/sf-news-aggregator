package logger

import (
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	DefaultTimestampFieldName   = "time"
	TimestampFieldName          = "timestamp"
	DefaultLevelFieldName       = "level"
	DefaultMessageFieldName     = "message"
	DefaultErrorStackFieldName  = "stacktrace"
	DefaultCallerSkipFrameCount = 2
	Debug                       = "debug"
	debugMessage                = "debug log is enabled"
)

func init() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = DefaultTimestampFieldName
	zerolog.LevelFieldName = DefaultLevelFieldName
	zerolog.MessageFieldName = DefaultMessageFieldName
	zerolog.ErrorStackFieldName = DefaultErrorStackFieldName
	zerolog.CallerSkipFrameCount = DefaultCallerSkipFrameCount
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func NewLogger(w io.Writer, logLevel string) (zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(w).
		Level(lvl).With().
		Timestamp().
		Logger().
		Hook(UnixTimestampHook{})

	if lvl == zerolog.DebugLevel {
		logger.Debug().Bool(Debug, true).Msg(debugMessage)
	}

	return logger, nil
}

type UnixTimestampHook struct{}

func (h UnixTimestampHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Int64(TimestampFieldName, time.Now().Unix())
}
