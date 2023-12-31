// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "auth-service/models"

	repositories "auth-service/repositories"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, db, user
func (_m *UserRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	ret := _m.Called(ctx, db, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *models.User) error); ok {
		r0 = rf(ctx, db, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, db, args
func (_m *UserRepository) Get(ctx context.Context, db *gorm.DB, args *repositories.GetUserArgs) (*models.User, error) {
	ret := _m.Called(ctx, db, args)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repositories.GetUserArgs) (*models.User, error)); ok {
		return rf(ctx, db, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repositories.GetUserArgs) *models.User); ok {
		r0 = rf(ctx, db, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repositories.GetUserArgs) error); ok {
		r1 = rf(ctx, db, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, db, args
func (_m *UserRepository) List(ctx context.Context, db *gorm.DB, args *repositories.ListUsersArgs) ([]*models.User, error) {
	ret := _m.Called(ctx, db, args)

	var r0 []*models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repositories.ListUsersArgs) ([]*models.User, error)); ok {
		return rf(ctx, db, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repositories.ListUsersArgs) []*models.User); ok {
		r0 = rf(ctx, db, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repositories.ListUsersArgs) error); ok {
		r1 = rf(ctx, db, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRepository(t mockConstructorTestingTNewUserRepository) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
