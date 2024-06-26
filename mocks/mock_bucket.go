// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Kachyr/findyourpet/findyourpet-backend/pkg/awsS3 (interfaces: S3ServiceI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/Kachyr/findyourpet/findyourpet-backend/pkg/models"
	manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	gomock "github.com/golang/mock/gomock"
)

// MockS3ServiceI is a mock of S3ServiceI interface.
type MockS3ServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockS3ServiceIMockRecorder
}

// MockS3ServiceIMockRecorder is the mock recorder for MockS3ServiceI.
type MockS3ServiceIMockRecorder struct {
	mock *MockS3ServiceI
}

// NewMockS3ServiceI creates a new mock instance.
func NewMockS3ServiceI(ctrl *gomock.Controller) *MockS3ServiceI {
	mock := &MockS3ServiceI{ctrl: ctrl}
	mock.recorder = &MockS3ServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3ServiceI) EXPECT() *MockS3ServiceIMockRecorder {
	return m.recorder
}

// UploadPhotos mocks base method.
func (m *MockS3ServiceI) UploadPhotos(arg0 []string, arg1 string) ([]models.Photo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadPhotos", arg0, arg1)
	ret0, _ := ret[0].([]models.Photo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadPhotos indicates an expected call of UploadPhotos.
func (mr *MockS3ServiceIMockRecorder) UploadPhotos(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadPhotos", reflect.TypeOf((*MockS3ServiceI)(nil).UploadPhotos), arg0, arg1)
}

// UploadSinglePhoto mocks base method.
func (m *MockS3ServiceI) UploadSinglePhoto(arg0, arg1 string) (*manager.UploadOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadSinglePhoto", arg0, arg1)
	ret0, _ := ret[0].(*manager.UploadOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadSinglePhoto indicates an expected call of UploadSinglePhoto.
func (mr *MockS3ServiceIMockRecorder) UploadSinglePhoto(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadSinglePhoto", reflect.TypeOf((*MockS3ServiceI)(nil).UploadSinglePhoto), arg0, arg1)
}
