package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func TestFrameEcho() {

	// ############################################################

	test_tmpl, err := template.New("RenderData").Parse(testTemplateEcho)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range RenderData.Models {

		folder_path := fmt.Sprintf("tests/%v_controller_test.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		test_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = test_tmpl.Execute(test_file, model)
		if err != nil {
			panic(err)
		}
		test_file.Close()

	}

	// ###################################################################
	test_app_tmpl, err := template.New("RenderData").Parse(tempEchoCoverTemplate)
	if err != nil {
		panic(err)
	}
	//
	test_app_file, err := os.Create("tests/test_app.go")
	if err != nil {
		panic(err)
	}

	err = test_app_tmpl.Execute(test_app_file, RenderData)
	if err != nil {
		panic(err)
	}
	defer test_app_file.Close()
}

var testTemplateEcho = `
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"{{.ProjectName}}/models"
	"{{.ProjectName}}/controllers"
)

// go test -coverprofile=coverage.out ./...
// go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

// ##########################################################################
var tests{{.Name}}sPost = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	{{.LowerName}}_id      string          //path param
	post_data    models.{{.Name}}Post // patch_data
	expectedCode int             // expected HTTP status code
}{
	{
		name:        "post {{.Name}} - 1",
		description: "post {{.Name}} 1",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 2",
		description: "post {{.Name}} 2",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 3",
		description: "post {{.Name}} 3",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 4",
		description: "post {{.Name}} 4",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 5",
		description: "post {{.Name}} 5",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{},
		expectedCode: 500,
	},
}

// ##########################################################################
var tests{{.Name}}sPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	patch_data   models.{{.Name}}Patch // patch_data
	expectedCode int              // expected HTTP status code
}{
	{
		name:        "patch {{.Name}}s- 1",
		description: "patch {{.Name}}s- 1",
		route:       groupPath+"/{{.LowerName}}/1",
		patch_data: models.{{.Name}}Patch{},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 1",
		description: "patch {{.Name}}s- 1",
		route:       groupPath+"/{{.LowerName}}/2",
		patch_data: models.{{.Name}}Patch{},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 1",
		description: "patch {{.Name}}s- 1",
		route:       groupPath+"/{{.LowerName}}/1000",
		patch_data: models.{{.Name}}Patch{},
		expectedCode: 500,
	},

}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var tests{{.Name}}sGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	{
		name:         "get {{.Name}}s- 1",
		description:  "get {{.Name}}s- 1",
		route:        groupPath+"/{{.LowerName}}?page=1&size=10",
		expectedCode: 200,
	},
	{
		name:         "get {{.Name}}s - 2",
		description:  "get {{.Name}}s- 2",
		route:        groupPath+"/{{.LowerName}}?page=0&size=-5",
		expectedCode: 400,
	},
	{
		name:         "get {{.Name}}s- 3",
		description:  "get {{.Name}}s- 3",
		route:        groupPath+"/{{.LowerName}}?page=1&size=0",
		expectedCode: 400,
	},
}

// ##############################################################
var tests{{.Name}}sGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	{
		name:         "get {{.Name}}s By ID  1",
		description:  "get {{.Name}}s By ID  1",
		route:        groupPath+"/{{.LowerName}}/1",
		expectedCode: 200,
	},
	{
		name:         "get {{.Name}}s By ID  2",
		description:  "get {{.Name}}s By ID  2",
		route:        groupPath+"/{{.LowerName}}/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get {{.Name}}s By ID  3",
		description:  "get {{.Name}}s By ID  3",
		route:        groupPath+"/{{.LowerName}}/1000",
		expectedCode: 404,
	},
}

func TestPost{{.Name}}Operations(t *testing.T) {


	//  test  test Post  {{.Name}} operations
	for _, test := range tests{{.Name}}sPost {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)
			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// Now testing the Get{{.Name}}s funciton
			controllers.Post{{.Name}}(echo_contx)


			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)
		})
	}


	// test Get {{.Name}} By ID cases
	for index, test := range tests{{.Name}}sGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("{{.LowerName}}_id")
			echo_contx.SetParamValues(fmt.Sprintf("%v",index+1))

			// Now testing the Get{{.Name}}s funciton
			controllers.Get{{.Name}}ByID(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

	// test {{.Name}} Patch Operations
	for index, test := range tests{{.Name}}sPatchID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)
			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("{{.LowerName}}_id")
			echo_contx.SetParamValues(fmt.Sprintf("%v",index+1))

			// Now testing the Get{{.Name}}s funciton
			controllers.Patch{{.Name}}(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

	// test {{.Name}} Get batch test cases
	for _, test := range tests{{.Name}}sGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)
			// Now testing the Get{{.Name}}s funciton
			controllers.Get{{.Name}}s(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}


	// test {{.Name}} Delete Operations
	t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {
		test_route := fmt.Sprintf("%v/:%v", groupPath,"{{.LowerName}}","{{.LowerName}}_id")
		req_delete := httptest.NewRequest(http.MethodDelete, test_route,nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		//  this is the response recorder
		resp := httptest.NewRecorder()

		//  create echo context to test the app function
		echo_contx_del := TestApp.NewContext(req_delete, resp)
		echo_contx_del.SetPath(test_route)

		// seting path paramenters
		// path_value := fmt.Sprintf("%v/%v/:%v", groupPath,"{{.LowerName}}","{{.LowerName}}_id")
		echo_contx_del.SetParamNames("{{.LowerName}}_id")
		echo_contx_del.SetParamValues(fmt.Sprintf("%v",3))

		// Now testing the Get{{.Name}}s funciton
		controllers.Delete{{.Name}}(echo_contx_del)
		assert.Equalf(t, 200, resp.Result().StatusCode, "deleteing {{.LowerName}}")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {
		test_route := fmt.Sprintf("%v/%v/:%v", groupPath, "{{.LowerName}}", "{{.LowerName}}_id")
		req_delete := httptest.NewRequest(http.MethodDelete, test_route, nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		//  this is the response recorder
		resp := httptest.NewRecorder()

		//  create echo context to test the app function
		echo_contx_del := TestApp.NewContext(req_delete, resp)
		echo_contx_del.SetPath(test_route)

		// seting path paramenters
		//path_value := fmt.Sprintf("%v", 2000)
		echo_contx_del.SetParamNames("{{.LowerName}}_id")
		echo_contx_del.SetParamValues(fmt.Sprintf("%v",2000))

		// Now testing the Get{{.Name}}s funciton
		controllers.Delete{{.Name}}(echo_contx_del)
		assert.Equalf(t, 500, resp.Result().StatusCode, "deleteing {{.LowerName}}")
	})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {
		test_route := fmt.Sprintf("%v/%v/:%v", groupPath,"{{.LowerName}}","{{.LowerName}}_id")
		req_delete := httptest.NewRequest(http.MethodDelete, test_route,nil)

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")

		//  this is the response recorder
		resp := httptest.NewRecorder()

		//  create echo context to test the app function
		echo_contx_del := TestApp.NewContext(req_delete, resp)
		echo_contx_del.SetPath(test_route)

		// seting path paramenters
		path_value := fmt.Sprintf("%v", "@@")
		echo_contx_del.SetParamNames("{{.LowerName}}_id")
		echo_contx_del.SetParamValues(path_value)

		// Now testing the Get{{.Name}}s funciton
		controllers.Delete{{.Name}}(echo_contx_del)
		assert.Equalf(t, 500, resp.Result().StatusCode, "deleteing {{.LowerName}}")
	})

}

`
var tempEchoCoverTemplate = `
package tests

import (
	"github.com/labstack/echo/v4"
	"github.com/joho/godotenv"
	"{{.ProjectName}}/controllers"
)

var (
	TestApp  *echo.Echo
	groupPath = "/api/v1"
)

func setupUserTestApp() {
	godotenv.Load(".test.env")
	TestApp = echo.New()
	manager.SetupRoutes(TestApp)
}
`
