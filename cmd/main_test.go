package main

import (
    "strings"
    "testing"
    "go.uber.org/zap/zapcore"
)

func TestParseLogLevel(t *testing.T) {
    cases := map[string]zapcore.Level{
        "debug":   zapcore.DebugLevel,
        "DEBUG":   zapcore.DebugLevel,
        "info":    zapcore.InfoLevel,
        "Info":    zapcore.InfoLevel,
        "warn":    zapcore.WarnLevel,
        "warning": zapcore.WarnLevel,
        "error":   zapcore.ErrorLevel,
        "err":     zapcore.ErrorLevel,
        "panic":   zapcore.PanicLevel,
        "fatal":   zapcore.FatalLevel,
        "":        zapcore.InfoLevel, // default fallback
        "garbage": zapcore.InfoLevel, // unknown fallback
    }
    for in, expected := range cases {
        lvl := parseLogLevel(strings.ToLower(in))
        if lvl != expected {
            t.Fatalf("input %q => %v, expected %v", in, lvl, expected)
        }
    }
}
// normalize mimics the lower-casing done in main when reading env
