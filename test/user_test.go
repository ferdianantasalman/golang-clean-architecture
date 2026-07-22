package test

import (
	"encoding/json"
	"golang-clean-architecture/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestLogin(t *testing.T) {
	TestRegister(t)

	user := GetFirstUser(t)

	requestBody := model.LoginUserRequest{
		Name:     user.Name,
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)
	assert.NotNil(t, responseBody.Data.RefreshToken)

	// ponytail: store JWT in DB so dependent tests can use user.Token for auth
	user.Token = responseBody.Data.Token
	err = db.Save(user).Error
	assert.Nil(t, err)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		Name:     "nonexistent",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	TestRegister(t)

	user := GetFirstUser(t)

	requestBody := model.LoginUserRequest{
		Name:     user.Name,
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestLogout(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := GetFirstUser(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestLogoutWrongAuthorization(t *testing.T) {
	ClearAll()
	TestLogin(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := GetFirstUser(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Name, responseBody.Data.Name)
	assert.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()
	TestLogin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := GetFirstUser(t)

	requestBody := model.UpdateUserRequest{
		Name: "Eko Kurniawan Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := GetFirstUser(t)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	user = GetFirstUser(t)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	ClearAll()
	TestLogin(t)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "wrong")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}

func TestRefresh(t *testing.T) {
	ClearAll()
	TestRegister(t)

	user := GetFirstUser(t)

	loginBody := model.LoginUserRequest{
		Name:     user.Name,
		Password: "rahasia",
	}
	bodyJson, err := json.Marshal(loginBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	loginResponse := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, loginResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, loginResponse.Data.RefreshToken)

	refreshBody := model.RefreshTokenRequest{
		RefreshToken: loginResponse.Data.RefreshToken,
	}
	bodyJson, err = json.Marshal(refreshBody)
	assert.Nil(t, err)

	request = httptest.NewRequest(http.MethodPost, "/api/users/_refresh", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	refreshResponse := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, refreshResponse)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, refreshResponse.Data.Token)
	assert.NotEmpty(t, refreshResponse.Data.RefreshToken)
	assert.NotEqual(t, loginResponse.Data.Token, refreshResponse.Data.Token)
	assert.NotEqual(t, loginResponse.Data.RefreshToken, refreshResponse.Data.RefreshToken)
}

func TestRefreshInvalid(t *testing.T) {
	ClearAll()

	refreshBody := model.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}
	bodyJson, err := json.Marshal(refreshBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_refresh", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "error", responseBody.Status)
}
