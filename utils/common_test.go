package utils

import (
	"sync"
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	t.Run("should execute function after delay", func(t *testing.T) {
		var executed bool
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executed = true
			mu.Unlock()
		}

		debounced := Debounce(f, 100*time.Millisecond)

		// Gọi function nhiều lần
		debounced()
		debounced()
		debounced()

		// Đợi một chút để function được thực thi
		time.Sleep(150 * time.Millisecond)

		mu.Lock()
		if !executed {
			t.Error("Function should have been executed after delay")
		}
		mu.Unlock()
	})

	t.Run("should not execute function immediately", func(t *testing.T) {
		var executed bool
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executed = true
			mu.Unlock()
		}

		debounced := Debounce(f, 100*time.Millisecond)

		// Gọi function
		debounced()

		// Kiểm tra ngay lập tức - function chưa được thực thi
		mu.Lock()
		if executed {
			t.Error("Function should not be executed immediately")
		}
		mu.Unlock()
	})

	t.Run("should reset timer on subsequent calls", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		debounced := Debounce(f, 50*time.Millisecond)

		// Gọi function lần đầu
		debounced()

		// Đợi 25ms (chưa đủ delay)
		time.Sleep(25 * time.Millisecond)

		// Gọi function lần nữa - reset timer
		debounced()

		// Đợi thêm 75ms (đủ delay từ lần gọi cuối)
		time.Sleep(75 * time.Millisecond)

		mu.Lock()
		if executionCount != 1 {
			t.Errorf("Function should be executed exactly once, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should handle multiple rapid calls", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		debounced := Debounce(f, 50*time.Millisecond)

		// Gọi function nhiều lần liên tiếp
		for i := 0; i < 10; i++ {
			debounced()
			time.Sleep(10 * time.Millisecond)
		}

		// Đợi để function được thực thi
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		if executionCount != 1 {
			t.Errorf("Function should be executed exactly once after multiple rapid calls, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should work with zero delay", func(t *testing.T) {
		var executed bool
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executed = true
			mu.Unlock()
		}

		debounced := Debounce(f, 0)

		debounced()

		// Với delay = 0, function sẽ được thực thi ngay lập tức
		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		if !executed {
			t.Error("Function should be executed immediately with zero delay")
		}
		mu.Unlock()
	})
}

func TestThrottle(t *testing.T) {
	t.Run("should execute function immediately on first call", func(t *testing.T) {
		var executed bool
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executed = true
			mu.Unlock()
		}

		throttled := Throttle(f, 100*time.Millisecond)

		// Gọi function lần đầu
		throttled()

		// Kiểm tra ngay lập tức - function đã được thực thi
		mu.Lock()
		if !executed {
			t.Error("Function should be executed immediately on first call")
		}
		mu.Unlock()
	})

	t.Run("should not execute function within throttle period", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		throttled := Throttle(f, 100*time.Millisecond)

		// Gọi function lần đầu
		throttled()

		// Gọi function nhiều lần trong khoảng thời gian throttle
		for i := 0; i < 5; i++ {
			throttled()
			time.Sleep(10 * time.Millisecond)
		}

		mu.Lock()
		if executionCount != 1 {
			t.Errorf("Function should be executed exactly once, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should execute function again after throttle period", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		throttled := Throttle(f, 50*time.Millisecond)

		// Gọi function lần đầu
		throttled()

		// Đợi để qua thời gian throttle
		time.Sleep(60 * time.Millisecond)

		// Gọi function lần nữa
		throttled()

		mu.Lock()
		if executionCount != 2 {
			t.Errorf("Function should be executed twice, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should handle multiple calls with proper throttling", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		throttled := Throttle(f, 50*time.Millisecond)

		// Gọi function lần đầu
		throttled()

		// Gọi nhiều lần trong thời gian throttle
		for i := 0; i < 3; i++ {
			throttled()
			time.Sleep(10 * time.Millisecond)
		}

		// Đợi để qua thời gian throttle
		time.Sleep(60 * time.Millisecond)

		// Gọi function lần nữa
		throttled()

		// Gọi nhiều lần nữa trong thời gian throttle
		for i := 0; i < 3; i++ {
			throttled()
			time.Sleep(10 * time.Millisecond)
		}

		mu.Lock()
		if executionCount != 2 {
			t.Errorf("Function should be executed exactly twice, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should work with zero delay", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		throttled := Throttle(f, 0)

		// Gọi function nhiều lần
		for i := 0; i < 5; i++ {
			throttled()
		}

		mu.Lock()
		if executionCount != 5 {
			t.Errorf("Function should be executed 5 times with zero delay, got %d", executionCount)
		}
		mu.Unlock()
	})

	t.Run("should be thread safe", func(t *testing.T) {
		var executionCount int
		var mu sync.Mutex

		f := func() {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		throttled := Throttle(f, 50*time.Millisecond)

		// Gọi function từ nhiều goroutine
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				throttled()
			}()
		}

		wg.Wait()

		mu.Lock()
		if executionCount != 1 {
			t.Errorf("Function should be executed exactly once with concurrent calls, got %d", executionCount)
		}
		mu.Unlock()
	})
}
