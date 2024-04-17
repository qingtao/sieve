package filter

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestSieve_1(t *testing.T) {
	type Case struct {
		name   string
		filter *Filter
		in     string
		want   string
		want2  string
	}

	filter := New()
	// 添加
	filter.Add([]string{"苹果", "西红柿", "葡萄"})
	// 移除
	filter.Remove([]string{"葡萄"})
	tests := []Case{
		{
			name:   "1",
			filter: filter,
			in:     "我想吃葡萄和西红柿，苹果也不错",
			want:   "西红柿",
			want2:  "我想吃葡萄和***，**也不错",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 搜索 (第一个关键词)
			got1, _ := tt.filter.Search(tt.in)
			if got1 != tt.want {
				t.Errorf("Search() = %v, want %v", got1, tt.want)
			}
			// 替换
			got2, _ := tt.filter.Replace(tt.in)
			if got2 != tt.want2 {
				t.Errorf("Replace() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestSieve_2(t *testing.T) {
	filter := New()
	filter.Add([]string{"fuck", "操你*"})
	type Case struct {
		name   string
		filter *Filter
		in     string
		want   string
	}
	tests := []Case{
		{
			name:   "1",
			filter: filter,
			in:     "fuck you!",
			want:   "**** you!",
		},
		{
			name:   "2",
			filter: filter,
			in:     "fuck you! 操你x",
			want:   "**** you! ***",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.filter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSieve_3(t *testing.T) {
	type Case struct {
		name   string
		filter *Filter
		file   string
		in     string
		want   string
		want1  map[uint8][]string
	}

	tests := []Case{
		{
			name:   "1",
			file:   "testdata/keywords.txt",
			filter: New(),
			in:     "你是傻b么？这么傻呢！",
			want:   "你是**么？这么傻呢！",
			want1: map[uint8][]string{
				1: {"傻b"},
			},
		},
		{
			name:   "2",
			file:   "testdata/keywords.txt",
			filter: New(),
			in:     "二手房怎么样",
			want:   "**房怎么样",
			want1: map[uint8][]string{
				1: {"二手"},
			},
		},
	}
	for _, tt := range tests {
		fails, err := tt.filter.AddByFile(tt.file, 1, true)
		if err != nil {
			t.Errorf("AddByFile error %s", err)
			return
		}
		if len(fails) > 0 {
			t.Errorf("添加单词失败: %+v", fails)
		}
		got, got1 := tt.filter.Replace(tt.in)
		if got != tt.want {
			t.Errorf("Replace() = %v, want %v", got, tt.want)
		}
		if !maps.EqualFunc(got1, tt.want1, func(v1, v2 []string) bool {
			return slices.Equal(v1, v2)
		}) {
			t.Errorf("Replace() = %+v, want %+v", got1, tt.want1)
		}
	}
}

func TestSieve_5(t *testing.T) {
	type Case struct {
		name   string
		filter *Filter
		in     string
		want   string
		want2  string
	}

	filter := New()
	// 添加
	filter.Add([]string{"苹果", "苹果**本"})
	// 移除
	tests := []Case{
		{
			name:   "1",
			filter: filter,
			in:     "我想吃葡萄和西红柿，苹果也不错",
			want:   "苹果",
			want2:  "我想吃葡萄和西红柿，**也不错",
		},
		{
			name:   "2",
			filter: filter,
			in:     "我想买一台苹果笔记本",
			want:   "苹果笔记本",
			want2:  "我想买一台*****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 搜索 (第一个关键词)
			got1, _ := tt.filter.Search(tt.in)
			if got1 != tt.want {
				t.Errorf("Search() = %v, want %v", got1, tt.want)
			}
			// 替换
			got2, _ := tt.filter.Replace(tt.in)
			if got2 != tt.want2 {
				t.Errorf("Replace() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func BenchmarkReplace(b *testing.B) {
	longText := strings.Repeat("哦😯哈HA", 20) // 100字符
	filter := New()
	filter.AddByFile("testdata/keywords.txt", 1, true)
	for i := 0; i < b.N; i++ {
		filter.Replace(longText)
	}
}

func BenchmarkSearch(b *testing.B) {
	longText := strings.Repeat("哦😯哈HA", 20) // 100字符
	filter := New()
	filter.AddByFile("testdata/keywords.txt", 1, true)
	for i := 0; i < b.N; i++ {
		filter.Search(longText)
	}
}
