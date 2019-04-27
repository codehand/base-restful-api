package main

import (
	"net/http"
	"strconv"

	"gitlab.com/jc93/echo-sdk/sdk"
)

// User is struct
type User struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
	Age      int    `json:"age"`
	Coin     int    `json:"coin"`
}

func handlerTest(c sdk.Context) error {
	return c.JSON(http.StatusOK, sdk.OKStatus)
}

func handlerGetUsers(c sdk.Context) error {
	if len(lstUser) == 0 {
		return c.JSON(http.StatusNotFound, sdk.ErrNotFound)
	}
	return c.JSON(http.StatusOK, lstUser)
}

func handlerGetUser(c sdk.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, sdk.NewError("REQ_INVALID", "ID not found"))
	}
	for _, u := range lstUser {
		if u.ID == id {
			return c.JSON(http.StatusOK, u)
		}
	}
	return c.JSON(http.StatusNotFound, sdk.ErrNotFound)
}

func handlerAddUser(c sdk.Context) error {
	var obj User
	if err := c.BindValidate(&obj); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	total = total + 1
	nextID := total
	user := &User{
		ID:       nextID,
		FullName: obj.FullName,
		Coin:     obj.Coin,
		Age:      obj.Age,
		Address:  obj.Address,
	}
	lstUser = append(lstUser, user)
	return c.JSON(http.StatusCreated, user)
}

func handlerUpdateUser(c sdk.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, sdk.NewError("REQ_INVALID", "ID not found"))
	}

	var obj User
	if err := c.BindValidate(&obj); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if obj.ID <= 0 {
		return c.JSON(http.StatusBadRequest, sdk.NewError("ID_INVALID", "ID invalid"))
	}

	for i, u := range lstUser {
		if u.ID == id {
			lstUser[i].Address = obj.Address
			lstUser[i].Age = obj.Age
			lstUser[i].Coin = obj.Coin
			lstUser[i].FullName = obj.FullName
			return c.JSON(http.StatusAccepted, lstUser[i])
		}
	}
	return c.JSON(http.StatusNotFound, sdk.ErrNotFound)
}

func handlerDeleteUser(c sdk.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, sdk.NewError("REQ_INVALID", "ID not found"))
	}

	for i, u := range lstUser {
		if u.ID == id {
			lstUser = append(lstUser[:i], lstUser[i+1:]...)
			return c.JSON(http.StatusAccepted, sdk.OKStatus)
		}
	}
	return c.JSON(http.StatusNotAcceptable, sdk.NewError("ERROR", "Action not accepted"))
}

var (
	lstUser []*User
	total   int
)

func init() {
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
func main() {
	app := sdk.New("demo", false, sdk.HTTPConfig{Host: "localhost", Port: 9090})
	app.GET("/test", handlerTest) // GET localhost:9090/test

	v1 := app.GROUP("/api/v1")
	v1.GET("/users", handlerGetUsers)
	v1.GET("/users/:id", handlerGetUser)
	v1.POST("/users", handlerAddUser)
	v1.PUT("/users/:id", handlerUpdateUser)
	v1.DELETE("/users/:id", handlerDeleteUser)

	v2 := app.GROUP("/api/v2")
	v2.GET("/users/self", handlerTest)
	v2.GET("/users/self/profiler", handlerTest)

	app.Run()
}
