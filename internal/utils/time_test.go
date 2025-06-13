package utils

import (
	"testing"
	"time"
)

func TestSmartDurationFormat(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		// 微秒级别测试
		{
			name:     "微秒级别_基本",
			duration: 456*time.Microsecond + 789*time.Nanosecond,
			expected: "456μs", // 应该忽略纳秒
		},
		{
			name:     "微秒级别_边界",
			duration: 999*time.Microsecond + 999*time.Nanosecond,
			expected: "999μs",
		},

		// 毫秒级别测试 - 关键测试点
		{
			name:     "毫秒级别_忽略微秒",
			duration: 123*time.Millisecond + 456*time.Microsecond + 789*time.Nanosecond,
			expected: "123ms", // 关键：应该忽略微秒和纳秒
		},
		{
			name:     "毫秒级别_整数",
			duration: 500 * time.Millisecond,
			expected: "500ms",
		},

		// 秒级别测试
		{
			name:     "秒+毫秒_忽略微秒",
			duration: 2*time.Second + 123*time.Millisecond + 456*time.Microsecond,
			expected: "2s123ms", // 应该忽略微秒
		},
		{
			name:     "整秒",
			duration: 5 * time.Second,
			expected: "5s",
		},
		{
			name:     "秒级别_边界",
			duration: 1*time.Second + 999*time.Millisecond + 999*time.Microsecond,
			expected: "1s999ms",
		},

		// 分钟级别测试
		{
			name:     "分钟+秒+毫秒_忽略微秒",
			duration: 2*time.Minute + 30*time.Second + 123*time.Millisecond + 456*time.Microsecond,
			expected: "2m30s",
		},
		{
			name:     "整分钟",
			duration: 3 * time.Minute,
			expected: "3m",
		},
		{
			name:     "分钟+毫秒_无秒",
			duration: 1*time.Minute + 500*time.Millisecond,
			expected: "1m500ms",
		},

		// 边界情况测试
		{
			name:     "边界_1毫秒",
			duration: 1 * time.Millisecond,
			expected: "1ms",
		},
		{
			name:     "边界_1微秒",
			duration: 1 * time.Microsecond,
			expected: "1μs",
		},

		// 性能测试相关的真实场景
		{
			name:     "典型API响应时间",
			duration: 45*time.Millisecond + 234*time.Microsecond,
			expected: "45ms",
		},
		{
			name:     "数据库查询时间",
			duration: 2*time.Second + 345*time.Millisecond + 123*time.Microsecond,
			expected: "2s345ms",
		},
		{
			name:     "内存分配时间",
			duration: 123*time.Microsecond + 456*time.Nanosecond,
			expected: "123μs",
		},
	}

	// 执行所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SmartDurationFormat(tt.duration)
			result = SmartDurationFormatDeperacted(tt.duration)

			// 断言
			if result != tt.expected {
				t.Errorf("测试 %s 失败:\n输入: %v\n期望: %s\n实际: %s",
					tt.name, tt.duration, tt.expected, result)
			}
		})
	}
}

func BenchmarkSmartDurationFormat(b *testing.B) {
	durations := []time.Duration{
		789 * time.Nanosecond,
		456 * time.Microsecond,
		123 * time.Millisecond,
		2*time.Second + 123*time.Millisecond,
		2*time.Minute + 30*time.Second + 123*time.Millisecond,
	}

	b.ResetTimer()
	b.Run("format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, d := range durations {
				_ = SmartDurationFormat(d)
			}
		}
	})
	b.Run("Deprecated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, d := range durations {
				_ = SmartDurationFormatDeperacted(d)
			}
		}
	})
}
