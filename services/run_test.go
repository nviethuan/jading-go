package services

import (
	"testing"
)

func TestRun(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run() bị panic: %v", r)
		}
	}()
	Run()
}
