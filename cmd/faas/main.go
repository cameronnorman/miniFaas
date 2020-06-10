package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Server    *echo.Echo
	functions map[string]string
}

func (a *App) Run() {
	a.Server.Logger.Fatal(a.Server.Start(":1323"))
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	validator := func(token string, c echo.Context) (bool, error) {
		return token == os.Getenv("TOKEN"), nil
	}
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:TOKEN",
		Validator: validator,
	}))

	app := &App{
		Server:    e,
		functions: map[string]string{},
	}
	e.POST("/create", app.saveUpload)
	e.GET("/all", app.listFunctions)
	e.POST("/:function_name/run", app.runFunction)
	app.Run()
}

func (a *App) runFunction(c echo.Context) error {
	functionName := c.Param("function_name")
	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		fmt.Println(err)
		return err
	}
	command := "./bin/" + functionName
	requestInBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestBodyString := string(requestInBytes)
	cmd := exec.Command(command, requestBodyString)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, string(out))
}

func (a *App) listFunctions(c echo.Context) error {
	files, err := ioutil.ReadDir("./bin")

	if err != nil {
	}

	functions := []string{}
	for _, f := range files {
		functions = append(functions, f.Name())
	}

	return c.JSON(http.StatusOK, functions)
}

func (a *App) saveUpload(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")
	//email := c.FormValue("email")

	//------------
	// Read files
	//------------

	// Multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Source
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("bin/" + name)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	cmd := exec.Command("chmod", "+x", "bin/"+name)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(out)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Upload successful"})
}

func updateMovieName(c echo.Context) error {
	authToken := c.QueryParam("token")
	if authToken != os.Getenv("TOKEN") {
		return c.String(http.StatusUnauthorized, "Incorrect token")
	}
	renameURL := "https://api.put.io/v2/files/rename"
	token := "token 7XWGGF4PPUBHAB2L3Z2R"
	u, _ := url.ParseRequestURI(renameURL)

	requestBody := map[string]interface{}{}
	if err := c.Bind(&requestBody); err != nil {
		fmt.Println(err)
		return err
	}
	fileID := requestBody["file_id"]
	re := regexp.MustCompile(`(.*)\s\((\d*)\).*`)
	result := re.FindStringSubmatch(requestBody["name"].(string))
	name := fmt.Sprintf("%s - %s", result[1], result[2])
	data := url.Values{}
	data.Set("name", name)
	data.Set("file_id", fileID.(string))
	urlStr := u.String()

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(resp)
		return err
	}

	return nil
}
