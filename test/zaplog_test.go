package test_test

import (
	"gostart/test"
	"testing"

	"go.uber.org/zap"
)

func TestTestZapLogInit(t *testing.T) {
	test.TestZapLogInit()
	test.Zap.Debug("TestZapLogInit", zap.String("test", "test"))
	test.Zap.Info("TestZapLogInit", zap.String("test", "test"))
	test.Zap.Warn("TestZapLogInit", zap.String("test", "test"))
	test.Zap.Error("TestZapLogInit", zap.String("test", "test"))
}
