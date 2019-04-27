package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc93/echo-sdk/sdk"
)

func setup() {
	lstUser = []*User{
		&User{
			ID:       1,
			FullName: "John Nguyen",
			Coin:     1000,
			Age:      40,
			Address:  "80 Paster",
		},
		&User{
			ID:       2,
			FullName: "Ana Tran",
			Coin:     3000,
			Age:      21,
			Address:  "80 Lusic",
		},
	}

	total = len(lstUser)
}

// TestHandlerGetUsers is func testing for handler get users
func TestHandlerGetUsers(t *testing.T) {
	setup()
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	req := httptest.NewRequest(sdk.GET, "/api/v1/users", nil)
	req.Header.Set(sdk.HeaderContentType, "application/json")
	res := httptest.NewRecorder()
	c := app.NewContext(req, res)
	c.SetPath("/api/v1/users")
	if assert.NoError(t, handlerGetUsers(c), "API exec error") {
		var es []*User
		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err, "Parse body error")
		assert.NoError(t, json.Unmarshal(body, &es))
		assert.Equal(t, http.StatusOK, res.Code, "Http status code # 200")
		assert.Equal(t, 2, len(es))
	}
}

// TestHandlerGetUsersClear is func get empty list users
func TestHandlerGetUsersClear(t *testing.T) {
	t.Log("setupTest()")
	setup()
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	time.Sleep(1 * time.Second)
	req := httptest.NewRequest(sdk.GET, "/api/v1/users", nil)
	req.Header.Set(sdk.HeaderContentType, "application/json")
	res := httptest.NewRecorder()
	c := app.NewContext(req, res)
	c.SetPath("/api/v1/users")
	lstUser = nil
	if assert.NoError(t, handlerGetUsers(c), "API exec error") {
		assert.Equal(t, http.StatusNotFound, res.Code, "Http status code # 404")
		var es sdk.Error
		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err, "Parse body error")
		assert.NoError(t, json.Unmarshal(body, &es))
	}
}

// TestHandlerGetUser is func test get user by id
func TestHandlerGetUser(t *testing.T) {
	setup()
	t.Log("setupTest()")
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	tests := []struct {
		name     string
		id       int
		expected *User
		isError  bool
	}{
		{
			name: "Test 1",
			id:   1,
			expected: &User{
				ID:       1,
				FullName: "John Nguyen",
				Coin:     1000,
				Age:      40,
				Address:  "80 Paster",
			},
			isError: false,
		},
		{
			name: "Test 2",
			id:   2,
			expected: &User{
				ID:       2,
				FullName: "Ana Tran",
				Coin:     3000,
				Age:      21,
				Address:  "80 Lusic",
			},
			isError: false,
		},
		{
			name:     "Test 3",
			id:       0,
			expected: nil,
			isError:  true,
		},
		{
			name:     "Test 4",
			id:       -1,
			expected: nil,
			isError:  true,
		},
		{
			name:     "Test 5",
			id:       3,
			expected: nil,
			isError:  true,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(sdk.GET, "/api/v1/users", nil)
		req.Header.Set(sdk.HeaderContentType, "application/json")
		res := httptest.NewRecorder()
		c := app.NewContext(req, res)
		c.SetPath("/api/v1/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(test.id))
		if assert.NoError(t, handlerGetUser(c), "API exec error") {
			if !test.isError {
				assert.Equal(t, http.StatusOK, res.Code, "Http status code # 200")
				var es *User
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
				assert.NotNil(t, es, test.name+" Data nil")
				assert.NotNil(t, test.expected, test.name+" Test data panic")
				assert.Equal(t, test.expected.ID, es.ID, test.name)
				assert.Equal(t, test.expected.FullName, es.FullName, test.name)
				assert.Equal(t, test.expected.Coin, es.Coin, test.name)
				assert.Equal(t, test.expected.Age, es.Age, test.name)
				assert.Equal(t, test.expected.Address, es.Address, test.name)
			} else {
				assert.Equal(t, http.StatusNotFound, res.Code, test.name+" Http status code # 404")
				var es sdk.Error
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" , Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
			}
		}
	}
}

// TestHandlerAddUser is func create new user
func TestHandlerAddUser(t *testing.T) {
	setup()
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	app.ValidatorDefault()
	tests := []struct {
		name     string
		input    User
		expected *User
		isError  bool
	}{
		{
			name: "Test 1",
			input: User{
				FullName: "Test 1",
				Coin:     333,
				Age:      21,
				Address:  "123 cao van",
			},
			expected: &User{
				FullName: "Test 1",
				Coin:     333,
				Age:      21,
				Address:  "123 cao van",
			},
			isError: false,
		},
		{
			name: "Test 2",
			input: User{
				FullName: "Test 2",
				Coin:     5600,
				Age:      43,
				Address:  "27/8 thanh thai",
			},
			expected: &User{
				FullName: "Test 2",
				Coin:     5600,
				Age:      43,
				Address:  "27/8 thanh thai",
			},
			isError: false,
		},
	}

	for _, test := range tests {
		obj, err := json.Marshal(test.input)
		assert.NoError(t, err, test.name+" Parse body error")
		req := httptest.NewRequest(sdk.POST, "/api/v1/users", strings.NewReader(string(obj)))
		req.Header.Set(sdk.HeaderContentType, "application/json")
		res := httptest.NewRecorder()
		c := app.NewContext(req, res)
		c.SetPath("/api/v1/users")
		if assert.NoError(t, handlerAddUser(c), "API exec error") {
			if !test.isError {
				assert.Equal(t, http.StatusCreated, res.Code, "Http status code # 201")
				var es *User
				body, err := ioutil.ReadAll(res.Body)
				t.Log(string(body))
				assert.NoError(t, err, test.name+" Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
				assert.NotNil(t, es, test.name+" Data nil")
				assert.NotNil(t, test.expected, test.name+" Test data panic")
				assert.NotZero(t, es.ID, test.name)
				assert.Equal(t, test.expected.FullName, es.FullName, test.name)
				assert.Equal(t, test.expected.Coin, es.Coin, test.name)
				assert.Equal(t, test.expected.Age, es.Age, test.name)
				assert.Equal(t, test.expected.Address, es.Address, test.name)
			} else {
				assert.Equal(t, http.StatusBadRequest, res.Code, test.name+" Http status code # 400")
				var es sdk.Error
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" , Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
			}
		}
	}
}

// TestHandlerUpdateUser is func update  user
func TestHandlerUpdateUser(t *testing.T) {
	setup()
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	app.ValidatorDefault()
	tests := []struct {
		name     string
		input    User
		expected *User
		isError  bool
	}{
		{
			name: "Test 1",
			input: User{
				ID:       1,
				FullName: "Test 1",
				Coin:     333,
				Age:      21,
				Address:  "123 cao van",
			},
			expected: &User{
				ID:       1,
				FullName: "Test 1",
				Coin:     333,
				Age:      21,
				Address:  "123 cao van",
			},
			isError: false,
		},
		{
			name: "Test 2",
			input: User{
				ID:       3,
				FullName: "Test 2",
				Coin:     5600,
				Age:      43,
				Address:  "27/8 thanh thai",
			},
			expected: nil,
			isError:  true,
		},
	}

	for _, test := range tests {
		obj, err := json.Marshal(test.input)
		assert.NoError(t, err, test.name+" Parse body error")
		req := httptest.NewRequest(sdk.PUT, "/api/v1/users", strings.NewReader(string(obj)))
		req.Header.Set(sdk.HeaderContentType, "application/json")
		res := httptest.NewRecorder()
		c := app.NewContext(req, res)
		c.SetPath("/api/v1/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(test.input.ID))
		if assert.NoError(t, handlerUpdateUser(c), "API exec error") {
			if !test.isError {
				assert.Equal(t, http.StatusAccepted, res.Code, "Http status code # 202")
				var es *User
				body, err := ioutil.ReadAll(res.Body)
				t.Log(string(body))
				assert.NoError(t, err, test.name+" Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
				assert.NotNil(t, es, test.name+" Data nil")
				assert.NotNil(t, test.expected, test.name+" Test data panic")
				assert.NotZero(t, es.ID, test.name)
				assert.Equal(t, es.ID, test.expected.ID)
				assert.Equal(t, test.expected.FullName, es.FullName, test.name)
				assert.Equal(t, test.expected.Coin, es.Coin, test.name)
				assert.Equal(t, test.expected.Age, es.Age, test.name)
				assert.Equal(t, test.expected.Address, es.Address, test.name)
			} else {
				assert.Equal(t, http.StatusNotFound, res.Code, test.name+" Http status code # 404")
				var es sdk.Error
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" , Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
			}
		}
	}
}

// TestHandlerDeleteUser is func del  user
func TestHandlerDeleteUser(t *testing.T) {
	setup()
	app := sdk.New("test", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	app.ValidatorDefault()
	tests := []struct {
		name     string
		input    int
		expected *User
		isError  bool
	}{
		{
			name:     "Test 1",
			input:    1,
			expected: nil,
			isError:  false,
		},
		{
			name:     "Test 2",
			input:    3,
			expected: nil,
			isError:  true,
		},
		{
			name:     "Test 3",
			input:    1,
			expected: nil,
			isError:  true,
		},
	}

	for _, test := range tests {
		obj, err := json.Marshal(test.input)
		assert.NoError(t, err, test.name+" Parse body error")
		req := httptest.NewRequest(sdk.DELETE, "/api/v1/users", strings.NewReader(string(obj)))
		req.Header.Set(sdk.HeaderContentType, "application/json")
		res := httptest.NewRecorder()
		c := app.NewContext(req, res)
		c.SetPath("/api/v1/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(test.input))
		if assert.NoError(t, handlerDeleteUser(c), "API exec error") {
			if !test.isError {
				assert.Equal(t, http.StatusAccepted, res.Code, "Http status code # 202")
				var es sdk.Error
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" , Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
			} else {
				assert.Equal(t, http.StatusNotAcceptable, res.Code, test.name+" Http status code # 406")
				var es sdk.Error
				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err, test.name+" , Parse body error")
				assert.NoError(t, json.Unmarshal(body, &es), test.name)
			}
		}
	}
}
