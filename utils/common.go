package utils

import (
	"math"
	"sync"
	"time"
)

func Debounce(f func(), delay time.Duration) func() {
	var timer *time.Timer
	return func() {
		if timer == nil {
			timer = time.AfterFunc(delay, f)
		}
		timer.Reset(delay)
	}
}

// Throttle thực thi function ngay lập tức và sau đó chỉ thực thi lại sau khoảng thời gian delay
// Hoạt động giống với throttle của RxJS
func Throttle(f func(), delay time.Duration) func() {
	var lastExecuted time.Time
	var mu sync.Mutex

	return func() {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		// Nếu chưa bao giờ thực thi hoặc đã qua đủ thời gian delay
		if lastExecuted.IsZero() || now.Sub(lastExecuted) >= delay {
			lastExecuted = now
			f()
		}
	}
}

func FloorTo(x float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Floor(x*pow) / pow
}

func Sum(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum
}


func Average(values []float64) float64 {
	sum := Sum(values)
	return sum / float64(len(values))
}
