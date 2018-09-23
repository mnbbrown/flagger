package flagger

import (
	"reflect"
	"testing"
)

func TestSimpleFlag(t *testing.T) {
	flagger, err := NewRedisFlagger("localhost:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	f := &Flag{
		Name:          "test",
		Type:          PERCENT,
		InternalValue: 0,
	}
	if err := flagger.SaveFlag(f); err != nil {
		t.Error(err)
		return
	}

	flag, err := flagger.GetFlag("test")
	if err != nil {
		t.Error(err)
		return
	}

	if reflect.DeepEqual(flag, f) == false {
		t.Error("Flag incorrectly saved or retrieved")
		return
	}
	if flag.Value() != false {
		t.Error("Flag incorrectly saved or retrieved")
		return
	}
}

func TestFlagWithSingleTag(t *testing.T) {
	flagger, err := NewRedisFlagger("localhost:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	f := &Flag{
		Name:          "test",
		Type:          BOOL,
		InternalValue: 1,
		Tags:          []string{"foo"},
	}
	if err := flagger.SaveFlag(f); err != nil {
		t.Error(err)
		return
	}

	flag, err := flagger.GetFlagWithTags("test", []string{"foo"})
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.DeepEqual(flag, f) == false {
		t.Error("Flag incorrectly saved or retrieved")
		return
	}
	if flag.Value() != true {
		t.Error("Flag incorrectly saved or retrieved")
		return
	}
}

func TestFlagWithMultipleTags(t *testing.T) {
	flagger, err := NewRedisFlagger("localhost:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	f := &Flag{
		Name:          "test",
		Type:          BOOL,
		InternalValue: 1,
		Tags:          []string{"foo", "buzz"},
	}
	if err := flagger.SaveFlag(f); err != nil {
		t.Error(err)
		return
	}

	flag, err := flagger.GetFlagWithTags("test", []string{"foo", "buzz"})
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.DeepEqual(flag, f) == false {
		t.Error("Flag incorrectly saved or retrieved")
		return
	}
	if flag.Value() != true {
		t.Error("Falg incorrectly saved or retrieved")
		return
	}
}

func TestFlagNotFound(t *testing.T) {
	flagger, err := NewRedisFlagger("localhost:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = flagger.GetFlag("nofoundtest")
	if err != ErrFlagNotFound {
		t.Error("Should return flag not found")
		return
	}
}

func TestFlagsNotFound(t *testing.T) {
	flagger, err := NewRedisFlagger("localhost:6379", 0)
	if err != nil {
		t.Error(err)
		return
	}

	f := &Flag{
		Name:          "test",
		Type:          BOOL,
		InternalValue: 1,
		Tags:          []string{"foo", "buzz"},
	}
	if err := flagger.SaveFlag(f); err != nil {
		t.Error(err)
		return
	}

	_, err = flagger.GetFlagWithTags("test", []string{"tag"})
	if err != ErrFlagNotFound {
		t.Error("Should return flag not found")
		return
	}

	flag, err := flagger.GetFlagWithTags("test", []string{"foo", "buzz"})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(f, flag) {
		t.Error("Not the same flag")
		return
	}
}
