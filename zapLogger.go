package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Logger struct {
	ZapLogger *zap.Logger
	sName     string
}

func LoggerZap(botToken string, chatID int64, webhookDS string, serviceName string) *Logger {
	telegramWriter := NewTelegramWriter(botToken, chatID)
	discordWriter := NewDiscordWriter(webhookDS)

	// Определяем имя файла с логами, включающее "log", дату и время
	logFileName := fmt.Sprintf("docker\\log\\log_%s.log", time.Now().Format("2006-01-02_15-04-05"))

	// Определяем WriteSyncer для файла
	fileWriteSyncer := zapcore.AddSync(createLogFile(logFileName))

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout", logFileName},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	cfgNew := cfg.EncoderConfig
	cfgNew.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				zapcore.AddSync(telegramWriter),
				cfg.Level,
			), zapcore.NewTee(zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				zapcore.AddSync(discordWriter),
				cfg.Level,
			), zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				fileWriteSyncer,
				cfg.Level,
			)))
		}),
		zap.AddCallerSkip(1),
	)

	if err != nil {
		fmt.Printf("Ошибка при создании логгера: %v\n", err)
		return nil
	}

	defer logger.Sync()
	return &Logger{ZapLogger: logger, sName: serviceName}
}
func LoggerZapDEV() *Logger {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil
	}

	defer logger.Sync()

	logger.Info("Develop Running")

	return &Logger{ZapLogger: logger}
}

func LoggerZapTelegram(botToken string, chatID int64, name ...string) *Logger {
	telegramWriter := NewTelegramWriter(botToken, chatID)

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	if len(name) > 0 {
		cfg.InitialFields = map[string]interface{}{
			"zoneName": name[0],
		}
	}

	cfgNew := cfg.EncoderConfig
	cfgNew.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				core, zapcore.NewCore(zapcore.NewConsoleEncoder(cfgNew), zapcore.AddSync(telegramWriter), cfg.Level))
		}),
		zap.AddCallerSkip(1),
	)

	if err != nil {
		fmt.Printf("Ошибка при создании логгера: %v\n", err)
		return nil
	}

	defer logger.Sync()
	return &Logger{ZapLogger: logger}
}
func LoggerZapDiscord(webhookDS string, name ...string) *Logger {
	discordWriter := NewDiscordWriter(webhookDS)

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	if len(name) > 0 {
		cfg.InitialFields = map[string]interface{}{
			"zoneName": name[0],
		}
	}

	cfgNew := cfg.EncoderConfig
	cfgNew.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				zapcore.AddSync(discordWriter),
				cfg.Level,
			))
		}),
		zap.AddCallerSkip(1),
	)

	if err != nil {
		fmt.Printf("Ошибка при создании логгера: %v\n", err)
		return nil
	}

	defer logger.Sync()
	return &Logger{ZapLogger: logger}
}

func LoggerZapTelegram1(botToken string, chatID int64, name ...string) *Logger {
	// Создаем писатель для телеграма
	telegramWriter := NewTelegramWriter(botToken, chatID)

	// Настройки конфигурации логгера
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	// Добавляем имя, если оно передано в качестве аргумента
	if len(name) > 0 {
		cfg.InitialFields = map[string]interface{}{
			"zoneName": name[0],
		}
	}

	// Создаем конфигурацию кодировщика для логов в консоль
	cfgNew := cfg.EncoderConfig
	cfgNew.EncodeLevel = zapcore.CapitalLevelEncoder

	// Создаем мульти-синкер для вывода логов в несколько мест
	consoleOutput := zapcore.Lock(os.Stdout)
	telegramOutput := zapcore.AddSync(telegramWriter)
	multiOutput := zapcore.NewMultiWriteSyncer(consoleOutput, telegramOutput)

	// Создаем ядро логгера
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfgNew),
		multiOutput,
		cfg.Level,
	)

	// Создаем логгер
	logger := zap.New(core, zap.AddCaller())

	return &Logger{ZapLogger: logger}
}

func (l *Logger) ErrorErr(err error) {
	l.ZapLogger.Error(fmt.Sprintf("[%s] Произошла ошибка", l.sName), zap.Error(err))
}
func (l *Logger) Debug(s string, fields ...zap.Field) {
	l.ZapLogger.Debug(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}
func (l *Logger) Info(s string, fields ...zap.Field) {
	l.ZapLogger.Info(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}
func (l *Logger) Warn(s string, fields ...zap.Field) {
	l.ZapLogger.Warn(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}
func (l *Logger) Error(s string, fields ...zap.Field) {
	l.ZapLogger.Error(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}
func (l *Logger) Panic(s string, fields ...zap.Field) {
	l.ZapLogger.Panic(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}
func (l *Logger) Fatal(s string, fields ...zap.Field) {
	l.ZapLogger.Fatal(fmt.Sprintf("[%s] %s\n", l.sName, s), fields...)
}

func (l *Logger) InfoStruct(s string, i interface{}, fields ...zap.Field) {
	l.ZapLogger.Info(fmt.Sprintf("[%s] %s: %+v \n", l.sName, s, i), fields...)
}
