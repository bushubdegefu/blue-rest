package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var clientJsTemplate = `
import { toast } from "sonner";
import { authService } from "./authService";

const API_BASE_URL = "/api/v1";

// Default headers for all requests
const defaultHeaders = {
  "Content-Type": "application/json",
};

// Add auth token if available
const getAuthHeaders = () => {
  const token = authService.getToken();
  return token ? { "X-APP-TOKEN": token } : {};
};

// Generic API request function
const apiRequest = async (endpoint, options = {}) => {
  try {
    const url = {{.BackTick}}${{ "{" }}API_BASE_URL}${{ "{" }}endpoint}{{.BackTick}};
    const headers = {
      ...defaultHeaders,
      ...getAuthHeaders(),
      ...options.headers,
    };

    const config = {
      ...options,
      headers,
    };

    const response = await fetch(url, config);

    // Check if response is JSON
    const contentType = response.headers.get("content-type");
    const isJson = contentType && contentType.includes("application/json");

    const data = isJson ? await response.json() : await response.text();

    if (!response.ok) {
      // Handle unauthorized specifically to redirect to login
      if (response.status === 401) {
        authService.logout();
        window.location.href = "/login";
        throw new Error("Session expired. Please login again.");
      }

      throw new Error(isJson ? data.details || "An error occurred" : "An error occurred");
    }

    return data;
  } catch (error) {
    console.error("API Error:", error);
    toast.error(error.message || "Failed to connect to the server");
    throw error;
  }
};

export const api = {
  get: (endpoint, params = {}) => {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        queryParams.append(key, value);
      }
    });

    const queryString = queryParams.toString();
    const url = queryString ? {{.BackTick}}${{ "{" }}endpoint}?${{ "{" }}queryString}{{.BackTick}} : endpoint;

    return apiRequest(url, { method: "GET" });
  },

  post: (endpoint, data) => {
    return apiRequest(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  patch: (endpoint, data) => {
    return apiRequest(endpoint, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  },

  delete: (endpoint) => {
    return apiRequest(endpoint, { method: "DELETE" });
  },
};

`
var indexJsTemplate = `
export { api } from './client';
export { statsService } from './statsService';
{{- range .Models}}
export { {{.LowerName}}Service } from './{{.LowerName}}Service';
{{- end }}

`

var authJsTemplate = `
import { api } from "./client";
import { toast } from "sonner";

export const authService = {
  login: async (credentials) => {
    try {
      	const response = await api.post("/{{ .AppName | replaceString }}/login", {
	        email:      credentials?.username, // Using email field for username
	        password:   credentials?.password,
	        grant_type: credentials?.grant_type || "authorization_code",
	        token:     credentials?.token_type || "Bearer"
        });

        const user = await api.post("/{{ .AppName | replaceString }}/login", {
	        "email":"tokendecode@mail.com",
	        "password":"123456",
	        "grant_type":"token_decode",
	        "token": response.data.access_token
	        });




      if (response && response.data) {
        localStorage.setItem("isAuthenticated", "true");
        localStorage.setItem("app-token", response.data.access_token);
        localStorage.setItem("refresh-token", response.data.refresh_token);

        // Store basic user info
        localStorage.setItem("user", JSON.stringify(user?.data));

        return { success: true };
      } else {
        throw new Error("Invalid response from server");
      }
    } catch (error) {
      toast.error(error.message || "Login failed");
      throw error;
    }
  },

  logout: () => {
    localStorage.removeItem("isAuthenticated");
    localStorage.removeItem("app-token");
    localStorage.removeItem("refresh-token");
    localStorage.removeItem("user");
    return { success: true };
  },

  isAuthenticated: () => {
    return localStorage.getItem("isAuthenticated") === "true";
  },

  getCurrentUser: () => {
    const userStr = localStorage.getItem("user");
    return userStr ? JSON.parse(userStr) : null;
  },

  getToken: () => {
    return localStorage.getItem("app-token");
  },

  getRefreshToken: () => {
    return localStorage.getItem("refresh-token");
  }
};

`

func ClientAuthAndIndexJSFrame() {
	// Create the models directory if it does not exist
	// #################################################
	err := os.MkdirAll("api", os.ModePerm)
	if err != nil {
		panic(err)
	}
	// ####################################################
	//  rabbit template
	clientjs_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(clientJsTemplate)
	if err != nil {
		panic(err)
	}

	clientjs_file, err := os.Create("api/client.js")
	if err != nil {
		panic(err)
	}
	defer clientjs_file.Close()

	err = clientjs_tmpl.Execute(clientjs_file, RenderData)
	if err != nil {
		panic(err)
	}
	// ####################################################
	indexjs_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(indexJsTemplate)
	if err != nil {
		panic(err)
	}

	indexjs_file, err := os.Create("api/index.js")
	if err != nil {
		panic(err)
	}
	defer indexjs_file.Close()

	err = indexjs_tmpl.Execute(indexjs_file, RenderData)
	if err != nil {
		panic(err)
	}
	// ####################################################
	authjs_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(authJsTemplate)
	if err != nil {
		panic(err)
	}

	authjs_file, err := os.Create("api/authService.js")
	if err != nil {
		panic(err)
	}
	defer authjs_file.Close()

	err = authjs_tmpl.Execute(authjs_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var serviceJSTemplate = `
import { api } from "./client";

export const {{.LowerName}}Service = {
  // Get a paginated list of {{.LowerName}}s
  get{{.Name}}s: (params = {}) => {
    const { page = 1, size = 10 } = params;
    return api.get("/{{ .AppName | replaceString }}/{{.LowerName}}", params);
  },

  // Get a specific {{.LowerName}} by ID
  get{{.Name}}ById: ({{.LowerName}}Id) => {
    return api.get({{.BackTick}}/{{ .AppName | replaceString }}/{{.LowerName}}/${{ "{" }}{{.LowerName}}Id{{ "}" }}{{.BackTick}});
  },

  // Create a new {{.LowerName}}
  create{{.Name}}: ({{.LowerName}}Data) => {
    return api.post("/{{ .AppName | replaceString }}/{{.LowerName}}", {{.LowerName}}Data);
  },

  // Update a {{.LowerName}}
  update{{.Name}}: (data) => {
    return api.patch({{.BackTick}}/{{ .AppName | replaceString }}/{{.LowerName}}/${{ "{" }}data?.{{.LowerName}}Id{{ "}" }}{{.BackTick}}, data?.{{.LowerName}}Data);
  },

  // Delete a {{.LowerName}}
  delete{{.Name}}: ({{.LowerName}}Id) => {
    return api.delete({{.BackTick}}/{{ .AppName | replaceString }}/{{.LowerName}}/${{ "{" }}{{.LowerName}}Id{{ "}" }}{{.BackTick}});
  },

{{- range .Relations }}
{{- if .MtM}}
//###############################################
// Now realationshipQeury Endpoints(Many to Many)
//###############################################
	get{{.ParentName}}{{.FieldName}}: (data)=>{
		return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}},{ page: data?.page, size: data?.size });
	},

	// Get {{.LowerFieldName}}s that can be assigned to a {{.LowerFieldName}}
	getAvailable{{.FieldName}}sFor{{.ParentName}}: ({{.LowerParentName}}Id) => {
	    return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}complement{{.LowerParentName}}/${{ "{" }}{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},
	// Get permissions that can be assigned to a {{.LowerFieldName}}
	getAttached{{.FieldName}}sFor{{.ParentName}}: ({{.LowerParentName}}Id) => {
	    return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}noncomplement{{.LowerParentName}}/${{ "{" }}{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},

	add{{.FieldName}}{{.ParentName}}: (data) => {
		return api.post({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerFieldName}}Id{{ "}" }}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},

	delete{{.FieldName}}{{.ParentName}}: (data) => {
		return api.delete({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerFieldName}}Id{{ "}" }}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},
{{- end}}
{{- end }}
{{- range .Relations }}
{{- if .OtM}}
//###############################################
// Now realationshipQeury Endpoints(one to Many)
//###############################################
	get{{.ParentName}}{{.FieldName}}: (data)=>{
		return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}},{ page: data?.page, size: data?.size });
	},

	// Get {{.LowerFieldName}}s that can be assigned to a {{.LowerFieldName}}
	getAvailable{{.FieldName}}sFor{{.ParentName}}: ({{.LowerParentName}}Id) => {
	    return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}complement{{.LowerParentName}}/${{ "{" }}{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},
	// Get permissions that can be assigned to a {{.LowerFieldName}}
	getAttached{{.FieldName}}sFor{{.ParentName}}: ({{.LowerParentName}}Id) => {
	    return api.get({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}noncomplement{{.LowerParentName}}/${{ "{" }}{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},

	add{{.FieldName}}{{.ParentName}}: (data) => {
		return api.post({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerFieldName}}Id{{ "}" }}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},

	delete{{.FieldName}}{{.ParentName}}: (data) => {
		return api.delete({{.BackTick}}/{{ $.AppName | replaceString }}/{{.LowerFieldName}}{{.LowerParentName}}/${{ "{" }}data?.{{.LowerFieldName}}Id{{ "}" }}/${{ "{" }}data?.{{.LowerParentName}}Id{{ "}" }}{{.BackTick}});
	},

{{- end }}
{{- end }}

}

`

func ClientJSFrame() {

	// ############################################################
	curd_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(serviceJSTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("api", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range RenderData.Models {

		model_name := strings.ToLower(model.Name)
		folder_path := fmt.Sprintf("api/%vService.js", model_name)
		curd_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}
		err = curd_tmpl.Execute(curd_file, model)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		curd_file.Close()

	}

}

// var tsTypesTemplate = `
// `

// var formsTemplate = `
// `

// func ClientTypesAndForms() {

// }
