package testifyhelper

import (
	"testing"

	"github.com/lumiluminousai/testify/assert"
	"github.com/lumiluminousai/testify/mock"
	"github.com/lumiluminousai/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) DoSomething(arg string) error {
	args := m.Called(arg)
	return args.Error(0)
}

type Handler struct {
	Service *MockService
}

func Test_AssertExpectationsForMocks_WithAssertNumberOfCalls_Success(t *testing.T) {

	t.Run("Only_AssertNumberOfCalls", func(t *testing.T) {
		// Arrange
		mockService := new(MockService)
		mockService.On("DoSomething", "test").Return(nil).Once()

		handler := &Handler{
			Service: mockService,
		}

		// Act
		err := handler.Service.DoSomething("test")

		// Assert
		require.NoError(t, err)
		mockService.AssertNumberOfCalls(t, "DoSomething", 1)

		// This should pass without errors
		err = AssertExpectationsForMocks(t, handler)
		require.NoError(t, err)
	})

	t.Run("AssertNumberOfCalls_And_AssertCalled_Correctly Order", func(t *testing.T) {
		// Arrange
		mockService := new(MockService)
		mockService.On("DoSomething", "test").Return(nil).Once()

		handler := &Handler{
			Service: mockService,
		}

		// Act
		err := handler.Service.DoSomething("test")

		// Assert
		require.NoError(t, err)
		mockService.AssertNumberOfCalls(t, "DoSomething", 1)
		mockService.AssertCalled(t, "DoSomething", "test")

		// This should pass without errors
		err = AssertExpectationsForMocks(t, handler)
		require.NoError(t, err)
	})

	t.Run("AssertNumberOfCalls_And_AssertNotCalled_Correctly Order", func(t *testing.T) {
		// Arrange
		mockService := new(MockService)
		mockService.On("DoSomething", "test").Return(nil).Once()

		handler := &Handler{
			Service: mockService,
		}

		// Act
		err := handler.Service.DoSomething("test")

		// Assert
		require.NoError(t, err)
		mockService.AssertNumberOfCalls(t, "DoSomething", 1)

		mockService.AssertNumberOfCalls(t, "DoSomething2", 0)
		mockService.AssertNotCalled(t, "DoSomething2", "test")

		// This should pass without errors
		err = AssertExpectationsForMocks(t, handler)
		require.NoError(t, err)
	})
}

func Test_AssertExpectationsForMocks_WithoutAssertNumberOfCalls_Failure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	mockService.On("DoSomething", "test").Return(nil)

	handler := &Handler{
		Service: mockService,
	}

	// Act
	err := handler.Service.DoSomething("test")

	// Assert
	require.NoError(t, err)

	// Not calling AssertNumberOfCalls

	// This should fail and report missing AssertNumberOfCalls
	err = AssertExpectationsForMocks(t, handler)

	// Ensure err is not nil
	if assert.Error(t, err) {
		// Check the error message contains the expected message
		expectedErrorMessage := "assert expectations failed for mock field 'Service.Mock':\nMissing AssertNumberOfCalls for method 'DoSomething'.\nFAIL: 0 out of 1 expectation(s) were met.\n\tThe code you are testing needs to make 1 more call(s).\n\tat"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	}
}

func Test_AssertExpectationsForMocks_WithoutRealFunctionCalling_Failure(t *testing.T) {
	// Arrange
	mockService := new(MockService)
	mockService.On("DoSomething", "test").Return(nil)

	handler := &Handler{
		Service: mockService,
	}

	// Act

	// Assert

	// This should fail and report missing AssertNumberOfCalls
	err := AssertExpectationsForMocks(t, handler)

	// Ensure err is not nil
	if assert.Error(t, err) {
		// Check the error message contains the updated message
		expectedErrorMessage := "assert expectations failed for mock field 'Service.Mock':\nFAIL: 0 out of 1 expectation(s) were met.\n\tThe code you are testing needs to make 1 more call(s).\n\tat:"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	}
}

// func Test_MockMethodAssertions_OrderEnforcement(t *testing.T) {
// 	t.Run("AssertCalled_Before_AssertNumberOfCalls", func(t *testing.T) {
// 		// Arrange
// 		mockService := new(MockService)
// 		mockService.On("DoSomething", "test").Return(nil)

// 		// Act
// 		mockService.DoSomething("test")

// 		// Assert
// 		// Attempt to call AssertCalled before AssertNumberOfCalls
// 		called := mockService.AssertCalled(t, "DoSomething", "test")
// 		assert.False(t, called, "AssertCalled should return false if AssertNumberOfCalls was not called")

// 		// Check if the test has failed due to t.Errorf in AssertCalled
// 		if t.Failed() {
// 			t.Log("AssertCalled failed as expected because AssertNumberOfCalls was not called.")
// 		} else {
// 			t.Error("Expected AssertCalled to fail, but it did not.")
// 		}
// 	})

// 	t.Run("AssertNotCalled_Before_AssertNumberOfCalls", func(t *testing.T) {
// 		// Arrange
// 		mockService := new(MockService)
// 		mockService.On("DoSomethingElse", "test").Return(nil)

// 		// Act
// 		//mockService.DoSomethingElse("test")

// 		// Assert
// 		// Attempt to call AssertNotCalled before AssertNumberOfCalls
// 		notCalled := mockService.AssertNotCalled(t, "DoSomethingElse", "test")
// 		assert.False(t, notCalled, "AssertNotCalled should return false if AssertNumberOfCalls was not called")

// 		// Check if the test has failed due to t.Errorf in AssertNotCalled
// 		if t.Failed() {
// 			t.Log("AssertNotCalled failed as expected because AssertNumberOfCalls was not called.")
// 		} else {
// 			t.Error("Expected AssertNotCalled to fail, but it did not.")
// 		}
// 	})
// }
