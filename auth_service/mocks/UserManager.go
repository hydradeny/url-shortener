// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	user "github.com/hydradeny/url-shortener/auth_service/internal/service/user"
	mock "github.com/stretchr/testify/mock"
)

// UserManager is an autogenerated mock type for the UserManager type
type UserManager struct {
	mock.Mock
}

type UserManager_Expecter struct {
	mock *mock.Mock
}

func (_m *UserManager) EXPECT() *UserManager_Expecter {
	return &UserManager_Expecter{mock: &_m.Mock}
}

// CheckPasswordByEmail provides a mock function with given fields: ctx, in
func (_m *UserManager) CheckPasswordByEmail(ctx context.Context, in *user.CheckPassword) (*user.User, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for CheckPasswordByEmail")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.CheckPassword) (*user.User, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.CheckPassword) *user.User); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.CheckPassword) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserManager_CheckPasswordByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckPasswordByEmail'
type UserManager_CheckPasswordByEmail_Call struct {
	*mock.Call
}

// CheckPasswordByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - in *user.CheckPassword
func (_e *UserManager_Expecter) CheckPasswordByEmail(ctx interface{}, in interface{}) *UserManager_CheckPasswordByEmail_Call {
	return &UserManager_CheckPasswordByEmail_Call{Call: _e.mock.On("CheckPasswordByEmail", ctx, in)}
}

func (_c *UserManager_CheckPasswordByEmail_Call) Run(run func(ctx context.Context, in *user.CheckPassword)) *UserManager_CheckPasswordByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*user.CheckPassword))
	})
	return _c
}

func (_c *UserManager_CheckPasswordByEmail_Call) Return(_a0 *user.User, _a1 error) *UserManager_CheckPasswordByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserManager_CheckPasswordByEmail_Call) RunAndReturn(run func(context.Context, *user.CheckPassword) (*user.User, error)) *UserManager_CheckPasswordByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, in
func (_m *UserManager) Create(ctx context.Context, in *user.CreateUser) (*user.User, error) {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.CreateUser) (*user.User, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.CreateUser) *user.User); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.CreateUser) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserManager_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserManager_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - in *user.CreateUser
func (_e *UserManager_Expecter) Create(ctx interface{}, in interface{}) *UserManager_Create_Call {
	return &UserManager_Create_Call{Call: _e.mock.On("Create", ctx, in)}
}

func (_c *UserManager_Create_Call) Run(run func(ctx context.Context, in *user.CreateUser)) *UserManager_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*user.CreateUser))
	})
	return _c
}

func (_c *UserManager_Create_Call) Return(_a0 *user.User, _a1 error) *UserManager_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserManager_Create_Call) RunAndReturn(run func(context.Context, *user.CreateUser) (*user.User, error)) *UserManager_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserManager) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserManager_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type UserManager_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserManager_Expecter) GetByEmail(ctx interface{}, email interface{}) *UserManager_GetByEmail_Call {
	return &UserManager_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *UserManager_GetByEmail_Call) Run(run func(ctx context.Context, email string)) *UserManager_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserManager_GetByEmail_Call) Return(_a0 *user.User, _a1 error) *UserManager_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserManager_GetByEmail_Call) RunAndReturn(run func(context.Context, string) (*user.User, error)) *UserManager_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserManager creates a new instance of UserManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserManager {
	mock := &UserManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}