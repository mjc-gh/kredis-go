package kredis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func (s *KredisTestSuite) TestScalarBoolWithUnknownConnection() {
	_, e := NewBool("foo", WithConfigName("unknown"))

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarBool() {
	k, e := NewBool("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(false)

	s.NoError(e)
	s.False(k.Value())

	e = k.SetValue(true)

	s.NoError(e)
	s.True(k.Value())

	ttl, e := k.TTL()
	s.NoError(e)
	s.Equal(time.Duration(0), ttl)
}

func (s *KredisTestSuite) TestScalarBoolWithDefaultValue() {
	k, e := NewBoolWithDefault("foo", true)

	s.NoError(e)
	s.True(k.Value())

	k, e = NewBoolWithDefault("bar", false)
	s.NoError(e)
	s.False(k.Value())

	k, e = NewBoolWithDefault("bar", true)
	s.NoError(e)
	s.False(k.Value())
}

func (s *KredisTestSuite) TestScalarBoolBadConnection() {
	k, _ := NewBool("foo", WithConfigName("badconn"))
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())

	k, e = NewBoolWithDefault("foo", true, WithConfigName("badconn"))
	s.Error(e)
	s.Nil(k)
}

func (s *KredisTestSuite) TestScalarBoolExpiresIn() {
	k, _ := NewBool("foo", WithExpiry("1ms"))

	e := k.SetValue(true)
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)

	k, _ = NewBool("foo", WithExpiry("1s"))
	k.SetValue(true)
	ttl, _ := k.TTL()
	s.Greater(ttl, time.Duration(0))
}

func (s *KredisTestSuite) TestScalarBoolWithContext() {
	k, e := NewBool("b", WithContext(context.TODO()))
	s.NoError(e)

	s.Empty(k.Value())
	s.Equal("context.TODO", fmt.Sprintf("%v", k.ctx))
}

func (s *KredisTestSuite) TestNewString() {
	k, e := NewString("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue("bar")

	s.NoError(e)
	s.Equal("bar", k.Value())
}

func (s *KredisTestSuite) TestNewStringWithDefaultValue() {
	k, e := NewStringWithDefault("foo", "bar")

	s.NoError(e)
	s.Equal("bar", k.Value())

	k, e = NewStringWithDefault("bar", "baz")
	s.NoError(e)
	s.Equal("baz", k.Value())

	k, e = NewStringWithDefault("bar", "quix")
	s.NoError(e)
	s.Equal("baz", k.Value())
}

func (s *KredisTestSuite) TestScalarStringExpiresIn() {
	k, _ := NewString("foo", WithExpiry("1ms"))

	e := k.SetValue("bar")
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func (s *KredisTestSuite) TestScalarIntegerWithUnknownConnection() {
	_, e := NewInteger("foo", WithConfigName("unknown"))

	s.Error(e)
}

func (s *KredisTestSuite) TestScalarInteger() {
	k, e := NewInteger("foo")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(1234)

	s.NoError(e)
	s.Equal(1234, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerWithDefaultValue() {
	k, e := NewIntegerWithDefault("foo", 5678)

	s.NoError(e)
	s.Equal(5678, k.Value())
}

func (s *KredisTestSuite) TestScalarIntegerBadConnection() {
	k, _ := NewInteger("foo", WithConfigName("badconn"))
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())

	k, e = NewIntegerWithDefault("foo", -1, WithConfigName("badconn"))
	s.Error(e)
	s.Nil(k)
}

func (s *KredisTestSuite) TestScalarIntegerExpiresIn() {
	k, _ := NewInteger("foo", WithExpiry("1ms"))

	e := k.SetValue(828)
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func (s *KredisTestSuite) TestScalarFloatWithUnknownConnection() {
	_, e := NewFloat("foo", WithConfigName("unknown"))
	s.Error(e)
}

func (s *KredisTestSuite) TestScalarFloat() {
	k, e := NewFloat("foo")
	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(12.34)
	s.NoError(e)
	s.Equal(float64(12.34), k.Value())
}

func (s *KredisTestSuite) TestScalarFloatWithDefaultValue() {
	k, e := NewFloatWithDefault("foo", 56.78)

	s.NoError(e)
	s.Equal(float64(56.78), k.Value())
}

func (s *KredisTestSuite) TestScalarFloatBadConnection() {
	k, _ := NewFloat("foo", WithConfigName("badconn"))
	_, e := k.ValueResult()

	s.Error(e)
	s.Empty(k.Value())

	k, e = NewFloatWithDefault("foo", -1, WithConfigName("badconn"))
	s.Error(e)
	s.Nil(k)
}

func (s *KredisTestSuite) TestScalarFloatExpiresIn() {
	k, _ := NewFloat("foo", WithExpiry("1ms"))

	e := k.SetValue(828)
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func (s *KredisTestSuite) TestNewTime() {
	k, e := NewTime("t")

	s.NoError(e)
	s.Empty(k.Value())

	now := time.Now()

	e = k.SetValue(now)
	s.NoError(e)

	// The canonical way to strip a monotonic clock reading is to use t =
	// t.Round(0).
	s.Equal(now.Round(0), k.Value().Local())
}

func (s *KredisTestSuite) TestNewTimeWithDefaultValue() {
	yesterday := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	k, e := NewTimeWithDefault("t1", yesterday)

	s.NoError(e)
	s.Equal(yesterday.Round(0), k.Value().Local())

	k, e = NewTimeWithDefault("t2", now)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value().Local())

	k, e = NewTimeWithDefault("t2", yesterday)
	s.NoError(e)
	s.Equal(now.Round(0), k.Value().Local())
}

func (s *KredisTestSuite) TestScalarTimeExpiresIn() {
	k, _ := NewTime("foo", WithExpiry("1ms"))

	e := k.SetValue(time.Now())
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func (s *KredisTestSuite) TestNewJSON() {
	kjson := NewKredisJSON(`{"key":"value"}`)
	k, e := NewJSON("json")

	s.NoError(e)
	s.Empty(k.Value())

	e = k.SetValue(kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())
}

func (s *KredisTestSuite) TestNewJSONWithDefaultValue() {
	kjson := NewKredisJSON(`{"key":"default"}`)
	kjson2 := NewKredisJSON(`{"k2":"v2"}`)
	k, e := NewJSONWithDefault("foo", kjson)

	s.NoError(e)
	s.Equal(*kjson, k.Value())

	k, e = NewJSONWithDefault("bar", kjson2)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())

	k, e = NewJSONWithDefault("bar", kjson)
	s.NoError(e)
	s.Equal(*kjson2, k.Value())
}

func (s *KredisTestSuite) TestScalarJSONExpiresIn() {
	k, _ := NewJSON("foo", WithExpiry("1ms"))

	e := k.SetValue(NewKredisJSON(`[1]`))
	s.NoError(e)

	time.Sleep(5 * time.Millisecond)

	v, e := k.ValueResult()
	s.Equal(redis.Nil, e)
	s.Nil(v)
}

func TestScalarJSONUnmarshal(t *testing.T) {
	var data interface{}

	kj := NewKredisJSON("[1]")
	err := kj.Unmarshal(&data)

	nums := data.([]interface{})

	assert.NoError(t, err)
	assert.Contains(t, nums, 1.0)
}
