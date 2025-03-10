package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func TestFrameFiber() {

	test_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(testTemplateFiber)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// #################################################
	err = os.MkdirAll("testsetting", os.ModePerm)
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
	test_app_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(testAppTemplate)
	if err != nil {
		panic(err)
	}
	//
	test_app_file, err := os.Create("testsetting/settings.go")
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
	{{- $break := false }}
	{{- range .Fields}}
		{{- if eq .Type "time.Time" }}
	"math/rand"
	"time"
		{{- $break = true }}
		{{- end}}
	{{- end}}
	"io"


	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"{{.ProjectName}}/models"
	"{{.ProjectName}}/testsetting"

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
		route:       testsetting.GroupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			{{- range .Fields}} {{- if .Post}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 2",
		description: "post {{.Name}} 2",
		route:       testsetting.GroupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			{{- range .Fields}} {{- if .Post}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 3",
		description: "post {{.Name}} 3",
		route:       testsetting.GroupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			{{- range .Fields}} {{- if .Post}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 4",
		description: "post {{.Name}} 4",
		route:       testsetting.GroupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			{{- range .Fields}} {{- if .Post}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: {{.RandomFeildValue}},  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "post {{.Name}} 5",
		description: "post {{.Name}} 5",
		route:       testsetting.GroupPath+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			{{- range .Fields}} {{- if .Post}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: {{.RandomFeildValue}},  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
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
		route:       testsetting.GroupPath+"/{{.LowerName}}/1",
		patch_data: models.{{.Name}}Patch{
			{{- range .Fields}} {{- if .Patch}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 2",
		description: "patch {{.Name}}s- 2",
		route:       testsetting.GroupPath+"/{{.LowerName}}/2",
		patch_data: models.{{.Name}}Patch{
			{{- range .Fields}} {{- if .Patch}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 200,
	},
	{
		name:        "patch {{.Name}}s- 3",
		description: "patch {{.Name}}s- 3",
		route:       testsetting.GroupPath+"/{{.LowerName}}/1000",
		patch_data: models.{{.Name}}Patch{
			{{- range .Fields}} {{- if .Patch}}
			{{- if eq .Type "uint" }}
			  {{.Name}}: {{randomUInt}},
			{{- else if eq .Type "time.Time" }}
			  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
			{{- else if eq .Type "string" }}
			  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
			{{- else if or (eq .Type "int") (eq .Type "float64") (eq .Type "uint") }}
			  {{.Name}}: {{randomUInt}},  // Numeric types
			{{- else if eq .Type "bool" }}
			  {{.Name}}: {{randomBool}},  // Assume RandomFeildValue is a string "true"/"false" and convert
			{{- else}}
			  {{.Name}}: {{.RandomFeildValue}},  // Default fallback
			{{- end}}
			{{- end}}
			{{- end}}
		},
		expectedCode: 404,
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
		route:        testsetting.GroupPath+"/{{.LowerName}}?page=1&size=10",
		expectedCode: 200,
	},
	{
		name:         "get {{.Name}}s - 2",
		description:  "get {{.Name}}s- 2",
		route:        testsetting.GroupPath+"/{{.LowerName}}?page=0&size=-5",
		expectedCode: 400,
	},
	{
		name:         "get {{.Name}}s- 3",
		description:  "get {{.Name}}s- 3",
		route:        testsetting.GroupPath+"/{{.LowerName}}?page=1&size=0",
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
		route:        testsetting.GroupPath+"/{{.LowerName}}/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get {{.Name}}s By ID  2",
		description:  "get {{.Name}}s By ID  2",
		route:        testsetting.GroupPath+"/{{.LowerName}}/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get {{.Name}}s By ID  3",
		description:  "get {{.Name}}s By ID  3",
		route:        testsetting.GroupPath+"/{{.LowerName}}/1000",
		expectedCode: 404,
	},
}

func Test{{.Name}}operations(t *testing.T) {
	// creating database for test
	testsetting.SetupTestApp()
	defer models.CleanDatabase(true)


	// test {{.Name}} Post Operations
	for _, test := range tests{{.Name}}sPost {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)
			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-APP-TOKEN", "hi")

			// Sending the request using the test app
			resp, err := testsetting.TestApp.Test(req, -1)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}
			defer resp.Body.Close() // Ensure we close the response body after reading it

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", test.name)
			t.Logf("Request URI: %s", req.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", test.expectedCode)
			t.Logf("Actual Status Code: %d", resp.StatusCode)

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
			req.Header.Set("X-APP-TOKEN", "hi")

			// Sending the request using the test app
			resp, err := testsetting.TestApp.Test(req, -1)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}
			defer resp.Body.Close() // Ensure we close the response body after reading it

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", test.name)
			t.Logf("Request URI: %s", req.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", test.expectedCode)
			t.Logf("Actual Status Code: %d", resp.StatusCode)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

	// test {{.Name}} Get batch
	for _, test := range tests{{.Name}}sGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			req.Header.Set("X-APP-TOKEN", "hi")

			// Sending the request using the test app
			resp, err := testsetting.TestApp.Test(req, -1)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}
			defer resp.Body.Close() // Ensure we close the response body after reading it

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", test.name)
			t.Logf("Request URI: %s", req.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", test.expectedCode)
			t.Logf("Actual Status Code: %d", resp.StatusCode)


			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		})
	}

	// test Get Single {{.Name}} test cases
	for _, test := range tests{{.Name}}sGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
			req.Header.Set("X-APP-TOKEN", "hi")
			// Sending the request using the test app
			resp, err := testsetting.TestApp.Test(req, -1)
			if err != nil {
				t.Fatalf("Error making request: %v", err)
			}
			defer resp.Body.Close() // Ensure we close the response body after reading it

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", test.name)
			t.Logf("Request URI: %s", req.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", test.expectedCode)
			t.Logf("Actual Status Code: %d", resp.StatusCode)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}


	// test {{.Name}} Delete Operations
	t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {
		test_route := fmt.Sprintf("%v/%v/%v",testsetting.GroupPath,"{{.LowerName}}",3)
		req_delete := httptest.NewRequest(http.MethodDelete, test_route,nil)
		req_delete.Header.Set("X-APP-TOKEN", "hi")

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := testsetting.TestApp.Test(req_delete)
		defer resp.Body.Close() // Ensure we close the response body after reading it

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		// Printing response body and status code for debugging (better structured output)
		t.Logf("Test Name: %s", "Delete {{.Name}} for success")
		t.Logf("Request URI: %s", req_delete.RequestURI)
		t.Logf("Response Body: %s", string(body))
		t.Logf("Expected Status Code: %d", 200)
		t.Logf("Actual Status Code: %d", resp.StatusCode)

		assert.Equalf(t, 200, resp.StatusCode, "deleteing {{.LowerName}}")
	})

	t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {
			test_route_1 := fmt.Sprintf("%v/%v/%v",testsetting.GroupPath,"{{.LowerName}}",1000000)
			req_delete := httptest.NewRequest(http.MethodDelete, test_route_1,nil)
			req_delete.Header.Set("X-APP-TOKEN", "hi")

			// Add specfic headers if needed as below
			req_delete.Header.Set("Content-Type", "application/json")
			resp, _ := testsetting.TestApp.Test(req_delete)
			defer resp.Body.Close() // Ensure we close the response body after reading it

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", "Deleting {{.Name}} that does not exist")
			t.Logf("Request URI: %s", req_delete.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", 404)
			t.Logf("Actual Status Code: %d", resp.StatusCode)

			assert.Equalf(t, 404, resp.StatusCode, "deleteing {{.LowerName}}")
			})

	t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {
		test_route_2 := fmt.Sprintf("%v/%v/%v",testsetting.GroupPath,"{{.LowerName}}", "$$$")
		req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, nil)
		req_delete.Header.Set("X-APP-TOKEN", "hi")

		// Add specfic headers if needed as below
		req_delete.Header.Set("Content-Type", "application/json")
		resp, _ := testsetting.TestApp.Test(req_delete)
		defer resp.Body.Close() // Ensure we close the response body after reading it

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		// Printing response body and status code for debugging (better structured output)
		t.Logf("Test Name: %s", "Deleting With Invalid Path")
		t.Logf("Request URI: %s", req_delete.RequestURI)
		t.Logf("Response Body: %s", string(body))
		t.Logf("Expected Status Code: %d", 400)
		t.Logf("Actual Status Code: %d", resp.StatusCode)

		assert.Equalf(t, 400, resp.StatusCode, "deleteing {{.LowerName}}")
	})

}

`

var testAppTemplate = `
package testsetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"{{.ProjectName}}/manager"
	"{{.ProjectName}}/models"
)

var (
	TestApp   *fiber.App
	GroupPath = "/api/v1"
)

func SetupTestApp() {
	godotenv.Load(".test.env")
	models.InitDatabase(true)
	TestApp = fiber.New()
	manager.SetupRoutes(TestApp,true)
}


`
