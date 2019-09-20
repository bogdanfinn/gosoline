// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import time "time"

// Viper is an autogenerated mock type for the Viper type
type Viper struct {
	mock.Mock
}

// AddConfigPath provides a mock function with given fields: _a0
func (_m *Viper) AddConfigPath(_a0 string) {
	_m.Called(_a0)
}

// AllKeys provides a mock function with given fields:
func (_m *Viper) AllKeys() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// AllSettings provides a mock function with given fields:
func (_m *Viper) AllSettings() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// AutomaticEnv provides a mock function with given fields:
func (_m *Viper) AutomaticEnv() {
	_m.Called()
}

// Get provides a mock function with given fields: _a0
func (_m *Viper) Get(_a0 string) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetBool provides a mock function with given fields: key
func (_m *Viper) GetBool(key string) bool {
	ret := _m.Called(key)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetDuration provides a mock function with given fields: _a0
func (_m *Viper) GetDuration(_a0 string) time.Duration {
	ret := _m.Called(_a0)

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(string) time.Duration); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// GetFloat64 provides a mock function with given fields: _a0
func (_m *Viper) GetFloat64(_a0 string) float64 {
	ret := _m.Called(_a0)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}

// GetInt provides a mock function with given fields: _a0
func (_m *Viper) GetInt(_a0 string) int {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetString provides a mock function with given fields: _a0
func (_m *Viper) GetString(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetStringMapString provides a mock function with given fields: key
func (_m *Viper) GetStringMapString(key string) map[string]string {
	ret := _m.Called(key)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string) map[string]string); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// GetStringSlice provides a mock function with given fields: key
func (_m *Viper) GetStringSlice(key string) []string {
	ret := _m.Called(key)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// IsSet provides a mock function with given fields: _a0
func (_m *Viper) IsSet(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MergeConfig provides a mock function with given fields: in
func (_m *Viper) MergeConfig(in io.Reader) error {
	ret := _m.Called(in)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader) error); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: key, value
func (_m *Viper) Set(key string, value interface{}) {
	_m.Called(key, value)
}

// SetConfigType provides a mock function with given fields: in
func (_m *Viper) SetConfigType(in string) {
	_m.Called(in)
}

// SetDefault provides a mock function with given fields: _a0, _a1
func (_m *Viper) SetDefault(_a0 string, _a1 interface{}) {
	_m.Called(_a0, _a1)
}

// SetEnvPrefix provides a mock function with given fields: _a0
func (_m *Viper) SetEnvPrefix(_a0 string) {
	_m.Called(_a0)
}
