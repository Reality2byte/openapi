// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/openapi/tools/cli/internal/openapi (interfaces: Parser,Merger)
//
// Generated by this command:
//
//	mockgen -destination=../openapi/mock_openapi.go -package=openapi github.com/mongodb/openapi/tools/cli/internal/openapi Parser,Merger
//

// Package openapi is a generated GoMock package.
package openapi

import (
	reflect "reflect"

	load "github.com/oasdiff/oasdiff/load"
	gomock "go.uber.org/mock/gomock"
)

// MockParser is a mock of Parser interface.
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance.
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// CreateOpenAPISpecFromPath mocks base method.
func (m *MockParser) CreateOpenAPISpecFromPath(arg0 string) (*load.SpecInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOpenAPISpecFromPath", arg0)
	ret0, _ := ret[0].(*load.SpecInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOpenAPISpecFromPath indicates an expected call of CreateOpenAPISpecFromPath.
func (mr *MockParserMockRecorder) CreateOpenAPISpecFromPath(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOpenAPISpecFromPath", reflect.TypeOf((*MockParser)(nil).CreateOpenAPISpecFromPath), arg0)
}

// MockMerger is a mock of Merger interface.
type MockMerger struct {
	ctrl     *gomock.Controller
	recorder *MockMergerMockRecorder
}

// MockMergerMockRecorder is the mock recorder for MockMerger.
type MockMergerMockRecorder struct {
	mock *MockMerger
}

// NewMockMerger creates a new mock instance.
func NewMockMerger(ctrl *gomock.Controller) *MockMerger {
	mock := &MockMerger{ctrl: ctrl}
	mock.recorder = &MockMergerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerger) EXPECT() *MockMergerMockRecorder {
	return m.recorder
}

// MergeOpenAPISpecs mocks base method.
func (m *MockMerger) MergeOpenAPISpecs(arg0 []string) (*Spec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeOpenAPISpecs", arg0)
	ret0, _ := ret[0].(*Spec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MergeOpenAPISpecs indicates an expected call of MergeOpenAPISpecs.
func (mr *MockMergerMockRecorder) MergeOpenAPISpecs(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeOpenAPISpecs", reflect.TypeOf((*MockMerger)(nil).MergeOpenAPISpecs), arg0)
}
