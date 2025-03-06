package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func TestFrameEcho() {

	// ############################################################

	test_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(testTemplateEcho)
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
	test_app_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(tempEchoCoverTemplate)
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
}

var testTemplateEcho = `
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
		"strings"
		"net/http"
		"net/http/httptest"
		"testing"
		"io"


		"github.com/stretchr/testify/assert"
		"{{.ProjectName}}/controllers"
		"{{.ProjectName}}/observe"
		"{{.ProjectName}}/database"
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
			route:       fmt.Sprintf("%v/{{.LowerName}}",testsetting.GroupPath),
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
			name:        "post {{.Name}} 2",
			description: "post {{.Name}} 2",
			route:       fmt.Sprintf("%v/{{.LowerName}}",testsetting.GroupPath),
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
		{
			name:        "post {{.Name}} 3",
			description: "post {{.Name}} 3",
			route:       fmt.Sprintf("%v/{{.LowerName}}",testsetting.GroupPath),
			post_data: models.{{.Name}}Post{
				{{- range .Fields}} {{- if .Post}}
					{{- if eq .Type "uint" }}
					  {{.Name}}: {{randomUInt}},
					{{- else if eq .Type "time.Time" }}
					  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
					{{- else if eq .Type "string" }}
					  {{.Name}}: "{{randomString}} 3",  // Ensure quotes for string values
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
			route:       fmt.Sprintf("%v/{{.LowerName}}",testsetting.GroupPath),
			post_data: models.{{.Name}}Post{
				{{- range .Fields}} {{- if .Post}}
				{{- if eq .Type "uint" }}
				  {{.Name}}: {{randomUInt}},
				{{- else if eq .Type "time.Time" }}
				  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
				{{- else if eq .Type "string" }}
				  {{.Name}}: "{{randomString}} 4",  // Ensure quotes for string values
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
			route:       fmt.Sprintf("%v/{{.LowerName}}",testsetting.GroupPath),
			post_data: models.{{.Name}}Post{
				{{- range .Fields}} {{- if .Post}}
				{{- if eq .Type "uint" }}
				  {{.Name}}: {{randomUInt}},
				{{- else if eq .Type "time.Time" }}
				  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
				{{- else if eq .Type "string" }}
				  {{.Name}}: "{{randomString}} 5",  // Ensure quotes for string values
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
			route:       fmt.Sprintf("%v/{{.LowerName}}/2",testsetting.GroupPath),
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
			route:       fmt.Sprintf("%v/{{.LowerName}}/2",testsetting.GroupPath),
			patch_data: models.{{.Name}}Patch{
				{{- range .Fields}} {{- if .Patch}}
				{{- if eq .Type "uint" }}
				  {{.Name}}: {{randomUInt}},
				{{- else if eq .Type "time.Time" }}
				  {{.Name}}: time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second),
				{{- else if eq .Type "string" }}
				  {{.Name}}: "{{randomString}}",  // Ensure quotes for string values
				{{- else if or (eq .Type "int") (eq .Name "float64") (eq .Name "uint") }}
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
			route:       fmt.Sprintf("%v/{{.LowerName}}/1000",testsetting.GroupPath),
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
			route:        fmt.Sprintf("%v/{{.LowerName}}?page=1&size=10",testsetting.GroupPath),
			expectedCode: 200,
		},
		{
			name:         "get {{.Name}}s - 2",
			description:  "get {{.Name}}s- 2",
			route:        fmt.Sprintf("%v/{{.LowerName}}?page=0&size=-5",testsetting.GroupPath),
			expectedCode: 400,
		},
		{
			name:         "get {{.Name}}s- 3",
			description:  "get {{.Name}}s- 3",
			route:        fmt.Sprintf("%v/{{.LowerName}}?page=1&size=0",testsetting.GroupPath),
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
			route:        fmt.Sprintf("%v/{{.LowerName}}/1",testsetting.GroupPath),
			expectedCode: 200,
		},
		{
			name:         "get {{.Name}}s By ID  2",
			description:  "get {{.Name}}s By ID  2",
			route:        fmt.Sprintf("%v/{{.LowerName}}/-1",testsetting.GroupPath),
			expectedCode: 404,
		},
		// Second test case
		{
			name:         "get {{.Name}}s By ID  3",
			description:  "get {{.Name}}s By ID  3",
			route:        fmt.Sprintf("%v/{{.LowerName}}/1000",testsetting.GroupPath),
			expectedCode: 404,
		},
	}


	func TestPost{{.Name}}Operations(t *testing.T) {
		testsetting.SetupTestApp()
		defer models.CleanDatabase(true)

		//  test  test Post  {{.Name}} operations
		for _, test := range tests{{.Name}}sPost {
			t.Run(test.name, func(t *testing.T) {
				//  changing post data to json
				post_data, _ := json.Marshal(test.post_data)
				req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

				// Add specfic headers if needed as below
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-APP-TOKEN", "test")

				test_description := fmt.Sprintf("%v to %v", test.description,req.RequestURI)

				//  this is the response recorder
				resp := httptest.NewRecorder()

				//  create echo context to test the app function
				echo_contx := testsetting.TestApp.NewContext(req, resp)
				// setting db and tracer mock context
				// Inject mock dependencies into the context
				tracer, span := observe.EchoAppSpanner(echo_contx, fmt.Sprintf("%v-root", test.route))

				// echo_contx.Set("tracer", mockTracer) // Mock tracer
				echo_contx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
				mockDB, _ := database.ReturnSession()
				echo_contx.Set("db", mockDB)

				echo_contx.SetPath(test.route)

				// Now testing the Get{{.Name}}s funciton
				controllers.Post{{.Name}}(echo_contx)

				// Read the response body
				body, err := io.ReadAll(resp.Result().Body)
				if err != nil {
					t.Fatalf("Error reading response body: %v", err)
				}

				// Printing response body and status code for debugging (better structured output)
				t.Logf("Test Name: %s", test.name)
				t.Logf("Request URI: %s", req.RequestURI)
				t.Logf("Response Body: %s", string(body))
				t.Logf("Expected Status Code: %d", test.expectedCode)
				t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

				//  Finally asserting test cases
				assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test_description)
			})
		}


		// test Get {{.Name}} By ID cases
		for _, test := range tests{{.Name}}sGetByID {
			t.Run(test.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, test.route, nil)
				// Add specfic headers if needed as below
				// req.Header.Set("X-APP-TOKEN", "hi")
				req.Header.Set("X-APP-TOKEN", "test")
				test_description := fmt.Sprintf("%v to %v", test.description,req.RequestURI)
				//  this is the response recorder
				resp := httptest.NewRecorder()

				//  create echo context to test the app function
				echo_contx := testsetting.TestApp.NewContext(req, resp)
				// setting db and tracer mock context
				// Inject mock dependencies into the context
				tracer, span := observe.EchoAppSpanner(echo_contx, fmt.Sprintf("%v-root", test.route))

				// echo_contx.Set("tracer", mockTracer) // Mock tracer
				echo_contx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
				mockDB, _ := database.ReturnSession()
				echo_contx.Set("db", mockDB)

				echo_contx.SetPath(test.route)

				// extrating path parm value
				path_value := strings.Split(test.route, "/")
				lastItem := path_value[len(path_value)-1]

				// seting path paramenters
				echo_contx.SetParamNames("{{.LowerName}}_id")
				echo_contx.SetParamValues(fmt.Sprintf("%v",lastItem))

				// Now testing the Get{{.Name}}s funciton
				controllers.Get{{.Name}}ByID(echo_contx)

				// Read the response body
				body, err := io.ReadAll(resp.Result().Body)
				if err != nil {
					t.Fatalf("Error reading response body: %v", err)
				}

				// Printing response body and status code for debugging (better structured output)
				t.Logf("Test Name: %s", test.name)
				t.Logf("Request URI: %s", req.RequestURI)
				t.Logf("Response Body: %s", string(body))
				t.Logf("Expected Status Code: %d", test.expectedCode)
				t.Logf("Actual Status Code: %d", resp.Result().StatusCode)


				//  Finally asserting test cases
				assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test_description)

			})
		}

		// test {{.Name}} Patch Operations
		for _, test := range tests{{.Name}}sPatchID {
			t.Run(test.name, func(t *testing.T) {
				//  changing post data to json
				patch_data, _ := json.Marshal(test.patch_data)
				req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

				// Add specfic headers if needed as below
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-APP-TOKEN", "test")
				test_description := fmt.Sprintf("%v to %v", test.description,req.RequestURI)
				//  this is the response recorder
				resp := httptest.NewRecorder()

				//  create echo context to test the app function
				echo_contx := testsetting.TestApp.NewContext(req, resp)
				// setting db and tracer mock context
				// Inject mock dependencies into the context
				tracer, span := observe.EchoAppSpanner(echo_contx, fmt.Sprintf("%v-root", test.route))

				// echo_contx.Set("tracer", mockTracer) // Mock tracer
				echo_contx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
				mockDB, _ := database.ReturnSession()
				echo_contx.Set("db", mockDB)

				echo_contx.SetPath(test.route)

				// extrating path parm value
				path_value := strings.Split(test.route, "/")
				lastItem := path_value[len(path_value)-1]

				// seting path paramenters
				echo_contx.SetParamNames("{{.LowerName}}_id")
				echo_contx.SetParamValues(fmt.Sprintf("%v",lastItem))

				// Now testing the Get{{.Name}}s funciton
				controllers.Patch{{.Name}}(echo_contx)

				// Read the response body
				body, err := io.ReadAll(resp.Result().Body)
				if err != nil {
					t.Fatalf("Error reading response body: %v", err)
				}

				// Printing response body and status code for debugging (better structured output)
				t.Logf("Test Name: %s", test.name)
				t.Logf("Request URI: %s", req.RequestURI)
				t.Logf("Response Body: %s", string(body))
				t.Logf("Expected Status Code: %d", test.expectedCode)
				t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

				//  Finally asserting test cases
				assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test_description)

			})
		}

		// test {{.Name}} Get batch test cases
		for _, test := range tests{{.Name}}sGet {
			t.Run(test.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, test.route, nil)
				// Add specfic headers if needed as below
				req.Header.Set("X-APP-TOKEN", "test")
				test_description := fmt.Sprintf("%v to %v", test.description,req.RequestURI)

				//  this is the response recorder
				resp := httptest.NewRecorder()

				//  create echo context to test the app function
				echo_contx := testsetting.TestApp.NewContext(req, resp)
				// setting db and tracer mock context
				// Inject mock dependencies into the context
				tracer, span := observe.EchoAppSpanner(echo_contx, fmt.Sprintf("%v-root", test.route))

				// echo_contx.Set("tracer", mockTracer) // Mock tracer
				echo_contx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
				mockDB, _ := database.ReturnSession()
				echo_contx.Set("db", mockDB)

				echo_contx.SetPath(test.route)
				// Now testing the Get{{.Name}}s funciton
				controllers.Get{{.Name}}s(echo_contx)

				// Read the response body
				body, err := io.ReadAll(resp.Result().Body)
				if err != nil {
					t.Fatalf("Error reading response body: %v", err)
				}

				// Printing response body and status code for debugging (better structured output)
				t.Logf("Test Name: %s", test.name)
				t.Logf("Request URI: %s", req.RequestURI)
				t.Logf("Response Body: %s", string(body))
				t.Logf("Expected Status Code: %d", test.expectedCode)
				t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

				//  Finally asserting test cases
				assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test_description)

			})
		}


		// test {{.Name}} Delete Operations
		t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {
			test_route := fmt.Sprintf("%v/%v/:%v", testsetting.GroupPath,"{{.LowerName}}","{{.LowerName}}_id")
			req_delete := httptest.NewRequest(http.MethodDelete, test_route,nil)

			// Add specfic headers if needed as below
			req_delete.Header.Set("Content-Type", "application/json")
			req_delete.Header.Set("X-APP-TOKEN", "test")
			test_description := fmt.Sprintf("deleteing {{.LowerName}} via path: %v",req_delete.RequestURI)

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx_del := testsetting.TestApp.NewContext(req_delete, resp)
			// setting db and tracer mock context
			// Inject mock dependencies into the context
			tracer, span := observe.EchoAppSpanner(echo_contx_del, fmt.Sprintf("%v-root", test_route))

			// echo_contx_del.Set("tracer", mockTracer) // Mock tracer
			echo_contx_del.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
			mockDB, _ := database.ReturnSession()
			echo_contx_del.Set("db", mockDB)
			echo_contx_del.SetPath(test_route)

			// seting path paramenters
			// path_value := fmt.Sprintf("%v/%v/:%v", testsetting.GroupPath,"{{.LowerName}}","{{.LowerName}}_id")
			echo_contx_del.SetParamNames("{{.LowerName}}_id")
			echo_contx_del.SetParamValues(fmt.Sprintf("%v",3))

			// Now testing the Get{{.Name}}s funciton
			controllers.Delete{{.Name}}(echo_contx_del)

			// Read the response body
			body, err := io.ReadAll(resp.Result().Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", "Delete {{.Name}} for sucess")
			t.Logf("Request URI: %s", req_delete.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", 200)
			t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

			assert.Equalf(t, 200, resp.Result().StatusCode, test_description)
		})

		t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {
			test_route := fmt.Sprintf("%v/%v/:%v", testsetting.GroupPath, "{{.LowerName}}", "{{.LowerName}}_id")
			req_delete := httptest.NewRequest(http.MethodDelete, test_route, nil)

			// Add specfic headers if needed as below
			req_delete.Header.Set("Content-Type", "application/json")
			req_delete.Header.Set("X-APP-TOKEN", "test")
			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx_del := testsetting.TestApp.NewContext(req_delete, resp)
			// setting db and tracer mock context
			// Inject mock dependencies into the context
			tracer, span := observe.EchoAppSpanner(echo_contx_del, fmt.Sprintf("%v-root", test_route))

			// echo_contx_del.Set("tracer", mockTracer) // Mock tracer
			echo_contx_del.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
			mockDB, _ := database.ReturnSession()
			echo_contx_del.Set("db", mockDB)
			echo_contx_del.SetPath(test_route)

			// seting path paramenters
			//path_value := fmt.Sprintf("%v", 2000)
			echo_contx_del.SetParamNames("{{.LowerName}}_id")
			echo_contx_del.SetParamValues(fmt.Sprintf("%v",2000))
			test_description := fmt.Sprintf("deleteing {{.LowerName}} via path: %v",req_delete.RequestURI)

			// Now testing the Get{{.Name}}s funciton
			controllers.Delete{{.Name}}(echo_contx_del)

			// Read the response body
			body, err := io.ReadAll(resp.Result().Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", "Delete {{.Name}} that does not exist")
			t.Logf("Request URI: %s", req_delete.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", 404)
			t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

			assert.Equalf(t, 404, resp.Result().StatusCode, test_description)
		})

		t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {
			test_route := fmt.Sprintf("%v/%v/:%v", testsetting.GroupPath,"{{.LowerName}}","{{.LowerName}}_id")
			req_delete := httptest.NewRequest(http.MethodDelete, test_route,nil)

			// Add specfic headers if needed as below
			req_delete.Header.Set("Content-Type", "application/json")
			req_delete.Header.Set("X-APP-TOKEN", "test")
			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx_del := testsetting.TestApp.NewContext(req_delete, resp)
			// setting db and tracer mock context
			// Inject mock dependencies into the context
			tracer, span := observe.EchoAppSpanner(echo_contx_del, fmt.Sprintf("%v-root", test_route))

			// echo_contx.Set("tracer", mockTracer) // Mock tracer
			echo_contx_del.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})
			mockDB, _ := database.ReturnSession()
			echo_contx_del.Set("db", mockDB)
			echo_contx_del.SetPath(test_route)

			// seting path paramenters
			path_value := fmt.Sprintf("%v", "@@")
			echo_contx_del.SetParamNames("{{.LowerName}}_id")
			echo_contx_del.SetParamValues(path_value)
			test_description := fmt.Sprintf("deleteing {{.LowerName}} via path: %v",req_delete.RequestURI)

			// Now testing the Get{{.Name}}s funciton
			controllers.Delete{{.Name}}(echo_contx_del)

			// Read the response body
			body, err := io.ReadAll(resp.Result().Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Printing response body and status code for debugging (better structured output)
			t.Logf("Test Name: %s", "Delete Request Path that is not valid")
			t.Logf("Request URI: %s", req_delete.RequestURI)
			t.Logf("Response Body: %s", string(body))
			t.Logf("Expected Status Code: %d", 400)
			t.Logf("Actual Status Code: %d", resp.Result().StatusCode)

			assert.Equalf(t, 400, resp.Result().StatusCode,test_description )
		})

	}

	`
var tempEchoCoverTemplate = `
	package testsetting

	import (
		"github.com/labstack/echo/v4"
		"github.com/joho/godotenv"
		"{{.ProjectName}}/manager"
		"{{.ProjectName}}/models"
	)

	var (
		TestApp  *echo.Echo
		GroupPath = "/api/v1"
	)

	func SetupTestApp() {
		godotenv.Load(".test.env")
		models.InitDatabase(true)
		TestApp = echo.New()
		manager.SetupRoutes(TestApp,true)
	}
	`
