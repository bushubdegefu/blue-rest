package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func TestFrameFiber() {

	test_tmpl, err := template.New("RenderData").Parse(testTemplateFiber)
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
	test_app_tmpl, err := template.New("RenderData").Parse(testAppTemplate)
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

	// ###################################################################
}

var testTemplateFiber = `
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
		post_data: models.{{.Name}}Post{
			Name:        "New one",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 2",
		description: "post {{.Name}} 2",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "New two",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 3",
		description: "post {{.Name}} 3",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "New three",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 4",
		description: "post {{.Name}} 4",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "New four",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 5",
		description: "post {{.Name}} 5",
		route:       groupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "Name four",
			Description: "Description of Name one",
		},
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
		patch_data: models.{{.Name}}Patch{
			Name:        "Name one updated",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 1",
		description: "patch {{.Name}}s- 1",
		route:       groupPath+"/{{.LowerName}}/2",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name two updated",
			Description: "Description of Name one for test one updated",
		},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 1",
		description: "patch {{.Name}}s- 1",
		route:       groupPath+"/{{.LowerName}}/1000",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name two updated",
			Description: "Description of Name one for test one updated",
		},
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

	// First test case
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

func Test{{.Name}}operations(t *testing.T) {
	// creating database for test
	models.InitDatabase()
	defer models.CleanDatabase()
	setupUserTestApp()

	// test {{.Name}} Post Operations
	for _, test := range tests{{.Name}}sPost {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)
			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// test {{.Name}} Patch Operations
	for _, test := range tests{{.Name}}sPatchID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)
			req := httptest.NewRequest(http.MethodPatch,test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")
			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// test {{.Name}} Get batch
	for _, test := range tests{{.Name}}sGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// test Get Single {{.Name}} test cases
	for _, test := range tests{{.Name}}sGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}


	// test {{.Name}} Delete Operations
	t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {
		test_route := fmt.Sprintf("%v/%v/:%v",groupPath,{{.LowerName}},3)
		req_delete := httptest.NewRequest(http.MethodDelete, test_route, bytes.NewReader(post_data))

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 200, resp.StatusCode, test.description+"deleteing")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {
			test_route_1 := fmt.Sprintf("%v/%v/:%v",groupPath,{{.LowerName}},1000000)
			req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req_delete.Header.Set("Content-Type", "application/json")

			resp, _ := TestApp.Test(req_delete)
			assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
			})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {
		test_route_2 := fmt.Sprintf("%v/%v/:%v",groupPath,{{.LowerName}}, "$$$")
		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, bytes.NewReader(post_data))

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := TestApp.Test(req_delete)

		assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
	})

}

`

var testAppTemplate = `
package tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"{{.ProjectName}}/controllers"
)

var (
	TestApp   *fiber.App
	groupPath = "/api/v1"
)

func setupUserTestApp() {
	godotenv.Load(".test.env")
	TestApp = fiber.New()
	manager.SetupRoutes(TestApp)
}

func nextFunc(contx *fiber.Ctx) error {
	contx.Next()
	return nil
}

`
