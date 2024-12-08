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
type LoggerConfig struct {
	DiscordWebhookURL string // URL веб хука Discord
	TelegramChatID    string // ID чата Telegram, куда отправляются логи
	TelegramToken     string // Токен Telegram бота
	LogFilePath       string // Путь к файлу для записи логов
	ServiceName       string
}


// NewLogger создает новый логгер на основе zap с учетом настроек для Discord и Telegram.
func NewLogger(config LoggerConfig) *Logger {

	var cores []zapcore.Core

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoder.TimeKey = "time"

	// Логирование в Discord, если URL вебхука указан
	if config.DiscordWebhookURL != "" {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder),
			zapcore.AddSync(NewDiscordWriter(config.DiscordWebhookURL)),
			zap.InfoLevel, // Для Discord
		))
	}

	// Логирование в Telegram, если указаны токен и chat_id
	if config.TelegramChatID != "" && config.TelegramToken != "" {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoder),
			zapcore.AddSync(&telegramWriter{
				botToken:  config.TelegramToken,
				chatID: config.TelegramChatID,
			}),
			zap.WarnLevel, // Для Telegram логирование с WarnLevel и выше
		))
	}

	// Логирование в файл
	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(createLogFile(config.LogFilePath)),
		zap.InfoLevel,
	))

	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Логирование в консоль
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	))

	// Создаем логгер с добавленными компонентами
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{ZapLogger: logger, sName: config.ServiceName}
}

func (l *Logger) Shutdown() {
	err := l.ZapLogger.Sync()
	if err != nil {
		l.ZapLogger.Error(err.Error())
	}
}


func LoggerZap(botToken string, chatID int64, webhookDS string, serviceName string) *Logger {
	telegramWriter := NewTelegramWriter(botToken, chatID)
	discordWriter := NewDiscordWriter(webhookDS)

	// Определяем имя файла с логами, включающее "log", дату и время
	logFileName := fmt.Sprintf("docker/log/log_%s_%s.log", serviceName, time.Now().Format("2006-01-02_15-04-05"))

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
func (l *Logger) DebugStruct(s string, i interface{}, fields ...zap.Field) {
	l.ZapLogger.Debug(fmt.Sprintf("[%s] %s: %+v \n", l.sName, s, i), fields...)
}