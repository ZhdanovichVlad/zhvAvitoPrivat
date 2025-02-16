package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/ZhdanovichVlad/go_final_project/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// BuyMerch mocks base method.
func (m *Mockrepository) BuyMerch(ctx context.Context, userUUID *string, merchInfo *entity.Merch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyMerch", ctx, userUUID, merchInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyMerch indicates an expected call of BuyMerch.
func (mr *MockrepositoryMockRecorder) BuyMerch(ctx, userUUID, merchInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyMerch", reflect.TypeOf((*Mockrepository)(nil).BuyMerch), ctx, userUUID, merchInfo)
}

// ExistsUser mocks base method.
func (m *Mockrepository) ExistsUser(ctx context.Context, userUUID *string) (*bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsUser", ctx, userUUID)
	ret0, _ := ret[0].(*bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistsUser indicates an expected call of ExistsUser.
func (mr *MockrepositoryMockRecorder) ExistsUser(ctx, userUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsUser", reflect.TypeOf((*Mockrepository)(nil).ExistsUser), ctx, userUUID)
}

// FindUser mocks base method.
func (m *Mockrepository) FindUser(ctx context.Context, userInfo *entity.User) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", ctx, userInfo)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockrepositoryMockRecorder) FindUser(ctx, userInfo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*Mockrepository)(nil).FindUser), ctx, userInfo)
}

// GetUserInfo mocks base method.
func (m *Mockrepository) GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", ctx, userUUID)
	ret0, _ := ret[0].(*entity.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserInfo indicates an expected call of GetUserInfo.
func (mr *MockrepositoryMockRecorder) GetUserInfo(ctx, userUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*Mockrepository)(nil).GetUserInfo), ctx, userUUID)
}

// SaveUser mocks base method.
func (m *Mockrepository) SaveUser(ctx context.Context, user *entity.User) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", ctx, user)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockrepositoryMockRecorder) SaveUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*Mockrepository)(nil).SaveUser), ctx, user)
}

// TransferCoins mocks base method.
func (m *Mockrepository) TransferCoins(ctx context.Context, userUUID, receiverUUID *string, amount *int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferCoins", ctx, userUUID, receiverUUID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferCoins indicates an expected call of TransferCoins.
func (mr *MockrepositoryMockRecorder) TransferCoins(ctx, userUUID, receiverUUID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferCoins", reflect.TypeOf((*Mockrepository)(nil).TransferCoins), ctx, userUUID, receiverUUID, amount)
}

// MocktokenGenerator is a mock of tokenGenerator interface.
type MocktokenGenerator struct {
	ctrl     *gomock.Controller
	recorder *MocktokenGeneratorMockRecorder
}

// MocktokenGeneratorMockRecorder is the mock recorder for MocktokenGenerator.
type MocktokenGeneratorMockRecorder struct {
	mock *MocktokenGenerator
}

// NewMocktokenGenerator creates a new mock instance.
func NewMocktokenGenerator(ctrl *gomock.Controller) *MocktokenGenerator {
	mock := &MocktokenGenerator{ctrl: ctrl}
	mock.recorder = &MocktokenGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktokenGenerator) EXPECT() *MocktokenGeneratorMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MocktokenGenerator) GenerateToken(userUUID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userUUID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MocktokenGeneratorMockRecorder) GenerateToken(userUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MocktokenGenerator)(nil).GenerateToken), userUUID)
}
