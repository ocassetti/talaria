// Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package presto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAppendBlocks(b *testing.B) {
	var block PrestoThriftBigint
	for i := 0; i < 25000; i++ {
		block.Append(int64(i))
	}

	var blocks []PrestoThriftBlock
	for i := 0; i < 10; i++ {
		blocks = append(blocks, PrestoThriftBlock{BigintData: &block})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var output PrestoThriftBigint
		output.AppendBlock(blocks...)
	}
}

func TestAppend_Bigint(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftBigint
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "int64 appended",
			input: int64(321),
			output: &PrestoThriftBigint{
				Nulls: []bool{false, false},
				Longs: []int64{123, 321},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftBigint{
				Nulls: []bool{false, true},
				Longs: []int64{123, 0},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftBigint{
				Nulls: []bool{false},
				Longs: []int64{456},
			},
			output: &PrestoThriftBigint{
				Nulls: []bool{false, false},
				Longs: []int64{123, 456},
			},
			size:  20,
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftBigint{
			Nulls: []bool{false},
			Longs: []int64{123},
		}

		if array, ok := td.input.(PrestoThriftBigint); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().BigintData)
		})
	}
}

func TestAppend_Varchar(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftVarchar
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "varchar appended",
			input: "hello",
			output: &PrestoThriftVarchar{
				Nulls: []bool{false, false},
				Sizes: []int32{2, 5},
				Bytes: []byte("hihello"),
			},
			outputRes: 11,
			size:      19,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftVarchar{
				Nulls: []bool{false, true},
				Sizes: []int32{2, 0},
				Bytes: []byte("hi"),
			},
			outputRes: 6,
			size:      14,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftVarchar{
				Nulls: []bool{false},
				Sizes: []int32{5},
				Bytes: []byte("world"),
			},
			output: &PrestoThriftVarchar{
				Nulls: []bool{false, false},
				Sizes: []int32{2, 5},
				Bytes: []byte("hiworld"),
			},
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftVarchar{
			Nulls: []bool{false},
			Sizes: []int32{2},
			Bytes: []byte("hi"),
		}

		if array, ok := td.input.(PrestoThriftVarchar); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().VarcharData)
		})
	}
}

func TestAppend_Double(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftDouble
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "int64 appended",
			input: float64(321),
			output: &PrestoThriftDouble{
				Nulls:   []bool{false, false},
				Doubles: []float64{123, 321},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftDouble{
				Nulls:   []bool{false, true},
				Doubles: []float64{123, 0},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftDouble{
				Nulls:   []bool{false},
				Doubles: []float64{456},
			},
			output: &PrestoThriftDouble{
				Nulls:   []bool{false, false},
				Doubles: []float64{123, 456},
			},
			size:  20,
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftDouble{
			Nulls:   []bool{false},
			Doubles: []float64{123},
		}

		if array, ok := td.input.(PrestoThriftDouble); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().DoubleData)
		})
	}
}

func TestAppend_Int32(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftInteger
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "int32 appended",
			input: int32(321),
			output: &PrestoThriftInteger{
				Nulls: []bool{false, false},
				Ints:  []int32{123, 321},
			},
			outputRes: 6,
			size:      12,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftInteger{
				Nulls: []bool{false, true},
				Ints:  []int32{123, 0},
			},
			outputRes: 6,
			size:      12,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftInteger{
				Nulls: []bool{false},
				Ints:  []int32{456},
			},
			output: &PrestoThriftInteger{
				Nulls: []bool{false, false},
				Ints:  []int32{123, 456},
			},
			size:  12,
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftInteger{
			Nulls: []bool{false},
			Ints:  []int32{123},
		}

		if array, ok := td.input.(PrestoThriftInteger); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().IntegerData)
		})
	}
}

func TestAppend_Boolean(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftBoolean
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "bool appended",
			input: false,
			output: &PrestoThriftBoolean{
				Nulls:    []bool{false, false},
				Booleans: []bool{true, false},
			},
			outputRes: 4,
			size:      8,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftBoolean{
				Nulls:    []bool{false, true},
				Booleans: []bool{true, false},
			},
			outputRes: 4,
			size:      8,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftBoolean{
				Nulls:    []bool{false},
				Booleans: []bool{false},
			},
			output: &PrestoThriftBoolean{
				Nulls:    []bool{false, false},
				Booleans: []bool{true, false},
			},
			size:  8,
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftBoolean{
			Nulls:    []bool{false},
			Booleans: []bool{true},
		}

		if array, ok := td.input.(PrestoThriftBoolean); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().BooleanData)
		})
	}
}

func TestAppend_Timestamp(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftTimestamp
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "timestamp appended",
			input: int64(321),
			output: &PrestoThriftTimestamp{
				Nulls:      []bool{false, false},
				Timestamps: []int64{123, 321},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc:  "timestamp appended",
			input: time.Unix(321, 0),
			output: &PrestoThriftTimestamp{
				Nulls:      []bool{false, false},
				Timestamps: []int64{123, 321},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftTimestamp{
				Nulls:      []bool{false, true},
				Timestamps: []int64{123, 0},
			},
			outputRes: 10,
			size:      20,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftTimestamp{
				Nulls:      []bool{false},
				Timestamps: []int64{456},
			},
			output: &PrestoThriftTimestamp{
				Nulls:      []bool{false, false},
				Timestamps: []int64{123, 456},
			},
			size:  20,
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftTimestamp{
			Nulls:      []bool{false},
			Timestamps: []int64{123},
		}

		if array, ok := td.input.(PrestoThriftTimestamp); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().TimestampData)
		})
	}
}

func TestAppend_Json(t *testing.T) {
	tests := []struct {
		desc      string
		input     interface{}
		output    *PrestoThriftJson
		outputRes int
		size      int
		count     int
	}{
		{
			desc:  "json appended",
			input: "hello",
			output: &PrestoThriftJson{
				Nulls: []bool{false, false},
				Sizes: []int32{2, 5},
				Bytes: []byte("hihello"),
			},
			outputRes: 11,
			size:      19,
			count:     2,
		},
		{
			desc:  "json appended",
			input: json.RawMessage("hello"),
			output: &PrestoThriftJson{
				Nulls: []bool{false, false},
				Sizes: []int32{2, 5},
				Bytes: []byte("hihello"),
			},
			outputRes: 11,
			size:      19,
			count:     2,
		},
		{
			desc:  "null value appended",
			input: nil,
			output: &PrestoThriftJson{
				Nulls: []bool{false, true},
				Sizes: []int32{2, 0},
				Bytes: []byte("hi"),
			},
			outputRes: 6,
			size:      14,
			count:     2,
		},
		{
			desc: "block appended",
			input: PrestoThriftJson{
				Nulls: []bool{false},
				Sizes: []int32{5},
				Bytes: []byte("world"),
			},
			output: &PrestoThriftJson{
				Nulls: []bool{false, false},
				Sizes: []int32{2, 5},
				Bytes: []byte("hiworld"),
			},
			count: 2,
		},
	}

	for _, td := range tests {
		output := &PrestoThriftJson{
			Nulls: []bool{false},
			Sizes: []int32{2},
			Bytes: []byte("hi"),
		}

		if array, ok := td.input.(PrestoThriftJson); ok {
			t.Run(td.desc, func(*testing.T) {
				output.AppendBlock(*array.AsBlock())
				assert.Equal(t, td.count, output.Count(), td.desc)
			})
			continue
		}

		// Append a single value
		t.Run(td.desc, func(*testing.T) {
			res := output.Append(td.input)
			assert.Equal(t, td.outputRes, res, td.desc)
			assert.Equal(t, td.output, output, td.desc)
			assert.Equal(t, td.size, output.Size(), td.desc)
			assert.Equal(t, td.count, output.Count(), td.desc)
			assert.NotNil(t, output.AsBlock().JsonData)
		})
	}
}