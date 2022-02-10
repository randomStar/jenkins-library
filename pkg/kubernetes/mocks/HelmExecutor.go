// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HelmExecutor is an autogenerated mock type for the HelmExecutor type
type HelmExecutor struct {
	mock.Mock
}

// RunHelmAdd provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmAdd() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmInstall provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmInstall() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmLint provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmLint() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmPackage provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmPackage() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmPush provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmPush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmRegistryLogin provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmRegistryLogin() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmRegistryLogout provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmRegistryLogout() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmTest provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmTest() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmUninstall provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmUninstall() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunHelmUpgrade provides a mock function with given fields:
func (_m *HelmExecutor) RunHelmUpgrade() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
