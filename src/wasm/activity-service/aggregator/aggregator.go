// Copyright 2024 The FIT SDK for Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Modifications made by Hikmatulloh Hari Mukti on 8 Sep 2024:
// - Allow Aggregate function to receive different struct types
//   and able to aggregate value of the same field name.
// - Add Replace and Fill functions.

package aggregator

import (
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/muktihari/fit/profile/basetype"
)

// Aggregate aggregates src and dst into dst using reflection where it
// should be a pointer to a struct, otherwise, it panics.
// The logic depends on the prefix of the Field's Name; If it starts with:
//   - "Total": sum of the two values.
//   - "Max" or "EnhancedMax": max value between the two values.
//   - "Min" or "EnhancedMin": min value between the two values.
//   - "Avg" or "EnhancedAvg": avg of the two values.
//   - "Num" and ends with "s": sum of the two values.
//     (e.g. NumSplits, NumLaps, NumSessions)
//   - Otherwise: fill with src's Value only if dst's Value is invalid.
func Aggregate(dst, src interface{}) {
	dv := reflect.Indirect(reflect.ValueOf(dst))
	sv := reflect.Indirect(reflect.ValueOf(src))

	for i := 0; i < dv.NumField(); i++ {
		f := dv.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		sField := sv.FieldByName(f.Name)
		if !sField.IsValid() {
			continue
		}

		if f.Type != reflect.TypeOf(time.Time{}) && reflect.Indirect(dv.Field(i)).Kind() == reflect.Struct {
			Aggregate(dv.Field(i).Interface(), sField.Interface())
			continue
		}

		switch {
		case strings.HasPrefix(f.Name, "Total"):
			sum(dv.Field(i), sField) // TotalElapsedTime, TotalCycles, etc.
		case strings.HasPrefix(f.Name, "Num") && strings.HasSuffix(f.Name, "s"):
			sum(dv.Field(i), sField) // NumSessions, NumLaps, NumSplits, etc.
		case strings.HasPrefix(f.Name, "Max") || strings.HasPrefix(f.Name, "EnhancedMax"):
			max(dv.Field(i), sField) // MaxHeartRate, MaxCadence, EnhancedMaxRespirationRate, etc.
		case strings.HasPrefix(f.Name, "Min") || strings.HasPrefix(f.Name, "EnhancedMin"):
			min(dv.Field(i), sField) // MinHeartRate, MinCadence, EnhancedMinAltitude, etc.
		case strings.HasPrefix(f.Name, "Avg") || strings.HasPrefix(f.Name, "EnhancedAvg"):
			avg(dv.Field(i), sField) // AvgHeartRate, AvgCadence, EnhancedAvgSpeed, etc.
		default:
			fill(dv.Field(i), sField) // Timestamp, Sport, Event, etc.
		}
	}
}

// Replace replaces dst fields with src when it has the same field name.
// It panics if dst and src are not pointer to a struct.
func Replace(dst, src interface{}) {
	dv := reflect.Indirect(reflect.ValueOf(dst))
	sv := reflect.Indirect(reflect.ValueOf(src))

	for i := 0; i < dv.NumField(); i++ {
		f := dv.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		sField := sv.FieldByName(f.Name)
		if !sField.IsValid() {
			continue
		}

		if f.Type != reflect.TypeOf(time.Time{}) && reflect.Indirect(dv.Field(i)).Kind() == reflect.Struct {
			Replace(dv.Field(i).Interface(), sField.Interface())
			continue
		}

		replace(dv.Field(i), sField)
	}
}

// Fill fill invalid value of dst fields with src when it has the same field name.
// It panics if dst and src are not pointer to a struct.
func Fill(dst, src interface{}) {
	dv := reflect.Indirect(reflect.ValueOf(dst))
	sv := reflect.Indirect(reflect.ValueOf(src))

	for i := 0; i < dv.NumField(); i++ {
		f := dv.Type().Field(i)
		if !f.IsExported() {
			continue
		}

		sField := sv.FieldByName(f.Name)
		if !sField.IsValid() {
			continue
		}

		if f.Type != reflect.TypeOf(time.Time{}) && reflect.Indirect(dv.Field(i)).Kind() == reflect.Struct {
			Fill(dv.Field(i).Interface(), sField.Interface())
			continue
		}

		fill(dv.Field(i), sField)
	}
}

var sInvalids = [...]int64{
	reflect.Int8:  int64(basetype.Sint8Invalid),
	reflect.Int16: int64(basetype.Sint16Invalid),
	reflect.Int32: int64(basetype.Sint32Invalid),
	reflect.Int64: int64(basetype.Sint64Invalid),
}

var uInvalids = [...]uint64{
	reflect.Uint8:  uint64(basetype.Uint8Invalid),
	reflect.Uint16: uint64(basetype.Uint16Invalid),
	reflect.Uint32: uint64(basetype.Uint32Invalid),
	reflect.Uint64: uint64(basetype.Uint64Invalid),
}

func sum(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Int() != sInvalids[k] && src.Int() != sInvalids[k] {
			dst.SetInt(dst.Int() + src.Int())
		} else if src.Int() != sInvalids[k] {
			dst.Set(src)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Uint() != uInvalids[k] && src.Uint() != uInvalids[k] {
			dst.SetUint(dst.Uint() + src.Uint())
		} else if src.Uint() != uInvalids[k] {
			dst.Set(src)
		}
	case reflect.Float32, reflect.Float64:
		if !math.IsNaN(dst.Float()) && !math.IsNaN(src.Float()) {
			dst.SetFloat(dst.Float() + src.Float())
		} else if !math.IsNaN(src.Float()) {
			dst.Set(src)
		}
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				sum(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				sum(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			sum(dst.Index(i), src.Index(i))
		}
	}
}

func max(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Int() != sInvalids[k] && src.Int() != sInvalids[k] {
			if dst.Int() < src.Int() {
				dst.Set(src)
			}
		} else if src.Int() != sInvalids[k] {
			dst.Set(src)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Uint() != uInvalids[k] && src.Uint() != uInvalids[k] {
			if dst.Uint() < src.Uint() {
				dst.Set(src)
			}
		} else if src.Uint() != uInvalids[k] {
			dst.Set(src)
		}
	case reflect.Float32, reflect.Float64:
		if !math.IsNaN(dst.Float()) && !math.IsNaN(src.Float()) {
			if dst.Float() < src.Float() {
				dst.Set(src)
			}
		} else if !math.IsNaN(src.Float()) {
			dst.Set(src)
		}
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				max(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				max(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			max(dst.Index(i), src.Index(i))
		}
	}
}

func min(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Int() != sInvalids[k] && src.Int() != sInvalids[k] {
			if dst.Int() > src.Int() {
				dst.Set(src)
			}
		} else if src.Int() != sInvalids[k] {
			dst.Set(src)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Uint() != uInvalids[k] && src.Uint() != uInvalids[k] {
			if dst.Uint() > src.Uint() {
				dst.Set(src)
			}
		} else if src.Uint() != uInvalids[k] {
			dst.Set(src)
		}
	case reflect.Float32, reflect.Float64:
		if !math.IsNaN(dst.Float()) && !math.IsNaN(src.Float()) {
			if dst.Float() > src.Float() {
				dst.Set(src)
			}
		} else if !math.IsNaN(src.Float()) {
			dst.Set(src)
		}
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				min(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				min(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			min(dst.Index(i), src.Index(i))
		}
	}
}

func avg(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Int() != sInvalids[k] && src.Int() != sInvalids[k] {
			dst.SetInt((dst.Int() + src.Int()) / 2)
		} else if src.Int() != sInvalids[k] {
			dst.Set(src)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Uint() != uInvalids[k] && src.Uint() != uInvalids[k] {
			dst.SetUint((dst.Uint() + src.Uint()) / 2)
		} else if src.Uint() != uInvalids[k] {
			dst.Set(src)
		}
	case reflect.Float32, reflect.Float64:
		if !math.IsNaN(dst.Float()) && !math.IsNaN(src.Float()) {
			dst.SetFloat((dst.Float() + src.Float()) / 2)
		} else if !math.IsNaN(src.Float()) {
			dst.Set(src)
		}
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				avg(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				avg(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			avg(dst.Index(i), src.Index(i))
		}
	}
}

func fill(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Int() == sInvalids[k] {
			dst.Set(src)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Uint() == uInvalids[k] {
			dst.Set(src)
		}
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(dst.Float()) {
			dst.Set(src)
		}
	case reflect.String:
		if dst.String() == basetype.StringInvalid || dst.String() == "\x00" {
			dst.Set(src)
		}
	case reflect.Bool:
		if !dst.Bool() {
			dst.Set(src)
		}
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				fill(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				fill(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			fill(dst.Index(i), src.Index(i))
		}
	case reflect.Struct:
		if dst.IsZero() && dst.Type() == reflect.TypeOf(time.Time{}) {
			dst.Set(src)
		}
	}
}

func replace(dst, src reflect.Value) {
	switch k := dst.Kind(); k {
	case reflect.Slice:
		if dst.Len() >= src.Len() {
			for i := 0; i < dst.Len(); i++ {
				if i >= src.Len() {
					break
				}
				replace(dst.Index(i), src.Index(i))
			}
		} else {
			for i := 0; i < src.Len(); i++ {
				if i >= dst.Len() {
					dst.Set(reflect.AppendSlice(dst, src.Slice(i, src.Len())))
					break
				}
				replace(dst.Index(i), src.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < dst.Len(); i++ {
			replace(dst.Index(i), src.Index(i))
		}
	default:
		dst.Set(src)
	}
}
