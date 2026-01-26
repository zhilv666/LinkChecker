package netdisk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string // 测试名称
		input    int64  // 输入
		expected string // 期望输出
	}{
		{
			name:     "测试 0 字节",
			input:    0,
			expected: "0 B",
		},
		{
			name:     "测试 1KB",
			input:    1024,
			expected: "1.00 KB",
		},
		{
			name:     "测试 1MB",
			input:    1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "测试 1GB",
			input:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "测试 1TB",
			input:    1024 * 1024 * 1024 * 1024,
			expected: "1.00 TB",
		},
		{
			name:     "测试 1PB",
			input:    1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.00 PB",
		},
		{
			name:     "测试 1EB",
			input:    1024 * 1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.00 EB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := FormatSize(tt.input)
			assert.Equal(t, tt.expected, value, "Expected value does not match.")
		})
	}
}
