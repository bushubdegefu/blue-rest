package temps

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"
)

func CommonFrame() {

	// ############################################################
	common_tmpl, err := template.New("RenderData").Parse(commonTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("common", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("common/common.go")
	if err != nil {
		panic(err)
	}
	defer common_file.Close()

	RenderData.BackTick = "`"
	err = common_tmpl.Execute(common_file, RenderData)
	if err != nil {
		panic(err)
	}

}

func CommonTargetJSON(project_name string) {

	// ############################################################
	targetJSON_tmpl, err := template.New("RenderData").Parse(targetJSONTemplate)
	if err != nil {
		panic(err)
	}

	targetJSON_file, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}
	defer targetJSON_file.Close()

	RenderData.ProjectName = project_name
	RenderData.AppName = "Change App Name"
	err = targetJSON_tmpl.Execute(targetJSON_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func CommonTargetAuthJSON(project_name string, app_name string) {

	// ############################################################
	targetJSON_tmpl, err := template.New("RenderData").Parse(blueAuthJSONTemplate)
	if err != nil {
		panic(err)
	}

	targetJSON_file, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}
	defer targetJSON_file.Close()

	RenderData.ProjectName = project_name
	RenderData.AppName = app_name
	err = targetJSON_tmpl.Execute(targetJSON_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func CommonCMDInit() {
	time.Sleep(2 * time.Second)
	// running go mod tidy finally
	if err := exec.Command("go", "get", "-u").Run(); err != nil {
		fmt.Printf("error go get: %v \n", err)
	}
}

func CommonCMD() {

	time.Sleep(2 * time.Second)
	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error tidy: %v \n", err)
	}
}

func CommonModInit(project_module string) {
	// running go mod tidy finally
	if err := exec.Command("go", "mod", "init", project_module).Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}
}

var commonTemplate = `
package common

import (
	"math"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResponseHTTP struct {
	Success bool        {{.BackTick}}json:"success"{{.BackTick}}
	Data    interface{} {{.BackTick}}json:"data"{{.BackTick}}
	Message string      {{.BackTick}}json:"details"{{.BackTick}}
}

type ResponsePagination struct {
	Success bool        {{.BackTick}}json:"success"{{.BackTick}}
	Items   interface{} {{.BackTick}}json:"data"{{.BackTick}}
	Message string      {{.BackTick}}json:"details"{{.BackTick}}
	Total   uint        {{.BackTick}}json:"total"{{.BackTick}}
	Page    uint        {{.BackTick}}json:"page"{{.BackTick}}
	Size    uint        {{.BackTick}}json:"size"{{.BackTick}}
	Pages   uint        {{.BackTick}}json:"pages"{{.BackTick}}
}

func Pagination(db *gorm.DB, queryModel interface{}, responseObjectModel interface{}, page uint, size uint, tracer context.Context) (ResponsePagination, error) {
	//  protection against requesting large amount of data
	//  set to 50
	var update_size uint
	if size > 50 {
		size = 50
	}

	count_channel := make(chan int64)
	str_chann := make(chan string)
	defer func() {
		close(str_chann)
		close(count_channel)
	}()
	var offset int64 = int64(page-1) * int64(update_size)
	//finding count value
	go func(comm <-chan int64) {
		var local_counter int64
		if tracer != nil {
			db.WithContext(tracer).Select("*").Model(&queryModel).Count(&local_counter)
		} else {
			db.Select("*").Model(&queryModel).Count(&local_counter)
		}
		count_channel <- local_counter

	}(count_channel)
	//  set offset value for page One
	var response_page int64
	go func(comm <-chan string) {
		if page == 1 {
			if tracer != nil {
				db.WithContext(tracer).Order("id asc").Limit(int(size)).Offset(0).Preload(clause.Associations).Find(&responseObjectModel)

			} else {
				db.Order("id asc").Limit(int(size)).Offset(0).Preload(clause.Associations).Find(&responseObjectModel)

			}

			response_page = 1
		} else {
			if tracer != nil {
				db.WithContext(tracer).Order("id asc").Limit(int(size)).Offset(int(offset)).Preload(clause.Associations).Find(&responseObjectModel)

			} else {
				db.Order("id asc").Limit(int(size)).Offset(int(offset)).Preload(clause.Associations).Find(&responseObjectModel)
			}
			// response_channel <- loc_resp
			response_page = int64(page)
		}
		str_chann <- "completed"
	}(str_chann)
	count := <-count_channel
	response_obj := <-str_chann
	pages := math.Ceil(float64(count) / float64(size))

	result := ResponsePagination{
		Success: true,
		Items:   responseObjectModel,
		Message: response_obj,
		Total:   uint(count),
		Page:    uint(response_page),
		Size:    uint(size),
		Pages:   uint(pages),
	}
	return result, nil
}

func PaginationPureModel(db *gorm.DB, queryModel interface{}, responseObjectModel interface{}, page uint, size uint, tracer context.Context) (ResponsePagination, error) {
	if size > 100 {
		size = 100
	}
	count_channel := make(chan int64)
	str_chann := make(chan string)
	defer func() {
		close(str_chann)
		close(count_channel)
	}()

	var offset int64 = int64(page-1) * int64(size)
	//finding count value
	go func(comm <-chan int64) {
		var local_counter int64
		if tracer != nil {
			db.WithContext(tracer).Select("*").Model(&queryModel).Count(&local_counter)
		} else {
			db.Select("*").Model(&queryModel).Count(&local_counter)
		}
		count_channel <- local_counter

	}(count_channel)
	//  set offset value for page One
	var response_page int64
	go func(comm <-chan string) {
		if page == 1 {
			if tracer != nil {
				db.WithContext(tracer).Model(&queryModel).Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			} else {
				db.Model(&queryModel).Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			}
			response_page = 1
		} else {
			if tracer != nil {
				db.WithContext(tracer).Model(&queryModel).Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			} else {
				db.Model(&queryModel).Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			}
			// response_channel <- loc_resp
			response_page = int64(page)
		}
		str_chann <- "completed"
	}(str_chann)

	count := <-count_channel
	response_obj := <-str_chann
	pages := math.Ceil(float64(count) / float64(size))
	result := ResponsePagination{
		Success: true,
		Items:   responseObjectModel,
		Message: response_obj,
		Total:   uint(count),
		Page:    uint(response_page),
		Size:    uint(size),
		Pages:   uint(pages),
	}
	return result, nil
}

func PaginationPureModelSearch(db *gorm.DB, queryModel interface{}, responseObjectModel interface{}, page uint, size uint, tracer context.Context, searchTerm map[string]interface{}) (ResponsePagination, error) {
	if size > 100 {
		size = 100
	}

	count_channel := make(chan int64)
	str_chann := make(chan string)
	defer func() {
		close(str_chann)
		close(count_channel)
	}()

	var offset int64 = int64(page-1) * int64(size)

	// Create a base query
	query := db.Model(&queryModel)

	// Apply search filters dynamically based on the provided searchTerm map
	if len(searchTerm) > 0 {
		// Loop through the searchTerm map and dynamically apply filters
		for key, value := range searchTerm {
			// Apply "LIKE" condition for string fields (or exact match for other types)
			if strValue, ok := value.(string); ok  && strValue != ""{
				query = query.Or(key+" LIKE ?", "%"+strValue+"%")
			}
		}
	}

	// Finding count value
	go func(comm <-chan int64) {
		var local_counter int64
		if tracer != nil {
			query.WithContext(tracer).Select("*").Count(&local_counter)
		} else {
			query.Select("*").Count(&local_counter)
		}
		count_channel <- local_counter
	}(count_channel)

	// Set offset value for pagination
	var response_page int64
	go func(comm <-chan string) {
		if page == 1 {
			if tracer != nil {
				query.WithContext(tracer).Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			} else {
				query.Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			}
			response_page = 1
		} else {
			if tracer != nil {
				query.WithContext(tracer).Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			} else {
				query.Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			}
			response_page = int64(page)
		}
		str_chann <- "completed"
	}(str_chann)

	count := <-count_channel
	response_obj := <-str_chann
	pages := math.Ceil(float64(count) / float64(size))

	result := ResponsePagination{
		Success: true,
		Items:   responseObjectModel,
		Message: response_obj,
		Total:   uint(count),
		Page:    uint(response_page),
		Size:    uint(size),
		Pages:   uint(pages),
	}

	return result, nil
}

// Generic function to filter the map based on a list of allowed keys.
func FilterMapByKeys(input map[string]any, allowedKeys []string) map[string]any {
	filtered := make(map[string]any)

	for _, key := range allowedKeys {
		if value, ok := input[key]; ok {
			filtered[key] = value
		}
	}

	return filtered
}

func PaginationPureModelFilterOneToMany(db *gorm.DB, queryModel interface{}, responseObjectModel interface{}, otm_string string, parentID uint,
	page uint, size uint, tracer context.Context) (ResponsePagination, interface{}, error) {
	if size > 100 {
		size = 100
	}
	count_channel := make(chan int64)
	str_chann := make(chan string)
	defer func() {
		close(str_chann)
		close(count_channel)
	}()

	var offset int64 = int64(page-1) * int64(size)
	//finding count value
	go func(comm <-chan int64) {
		var local_counter int64
		if tracer != nil {
			db.WithContext(tracer).Select("*").Model(&queryModel).Where(otm_string, parentID).Count(&local_counter)
		} else {
			db.WithContext(tracer).Select("*").Model(&queryModel).Where(otm_string, parentID).Count(&local_counter)
		}
		count_channel <- local_counter

	}(count_channel)
	//  set offset value for page One
	var response_page int64
	go func(comm <-chan string) {
		if page == 1 {
			if tracer != nil {
				db.WithContext(tracer).Model(&queryModel).Where(otm_string, parentID).Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			} else {
				db.WithContext(tracer).Model(&queryModel).Where(otm_string, parentID).Order("id asc").Limit(int(size)).Offset(0).Find(&responseObjectModel)
			}
			response_page = 1
		} else {
			if tracer != nil {
				db.WithContext(tracer).Model(&queryModel).Where(otm_string, parentID).Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			} else {
				db.WithContext(tracer).Model(&queryModel).Where(otm_string, parentID).Order("id asc").Limit(int(size)).Offset(int(offset)).Find(&responseObjectModel)
			}
			// response_channel <- loc_resp
			response_page = int64(page)
		}
		str_chann <- "completed"
	}(str_chann)

	count := <-count_channel
	response_obj := <-str_chann
	pages := math.Ceil(float64(count) / float64(size))
	result := ResponsePagination{
		Success: true,
		Items:   responseObjectModel,
		Message: response_obj,
		Total:   uint(count),
		Page:    uint(response_page),
		Size:    uint(size),
		Pages:   uint(pages),
	}
	return result, responseObjectModel, nil
}
`

var targetJSONTemplate = `
{
  "project_name": "{{.ProjectName}}",
  "app_name": "changeappname",
  "models": [
    {
      "name": "Role",
      "rln_model": ["User$mtm$user_roles", "Feature$otm"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$true$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Description",
          "type": "string",
          "annotation": "gorm:\"not null; \" json:\"description,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Active",
          "type": "bool",
          "annotation": "gorm:\"default:true; constraint:not null;\" json:\"active\"",
          "curd_flag": "true$false$true$true$false$false"
        },
        {
          "name": "Users",
          "type": "[]User",
          "annotation": "gorm:\"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"users,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        },
        {
          "name": "Features",
          "type": "[]Feature",
          "annotation": "gorm:\"foreignkey:RoleID; constraint:OnUpdate:CASCADE; OnDelete:SET NULL;\" json:\"features,omitempty\"",
          "curd_flag": "true$false$false$false$true$false"
        },
        {
          "name": "Pages",
          "type": "[]Page",
          "annotation": "gorm:\"many2many:page_roles; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\" json:\"pages,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        },

        {
          "name": "AppID",
          "type": "sql.NullInt64",
          "annotation": "gorm:\"foreignkey:AppID OnDelete:SET NULL\" json:\"app,omitempty\" swaggertype:\"number\"",
          "curd_flag": "false$false$false$false$true$false"
        }
      ]
    },
    {
      "name": "App",
      "rln_model": ["Role$otm"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$true$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "UUID",
          "type": "string",
          "annotation": "gorm:\"constraint:not null; unique; type:string;\" json:\"uuid\"",
          "curd_flag": "true$false$false$false$true$false"
        },
        {
          "name": "Active",
          "type": "bool",
          "annotation": "gorm:\"constraint:not null;\" json:\"active\"",
          "curd_flag": "true$false$true$true$false$false"
        },
        {
          "name": "Description",
          "type": "string",
          "annotation": "gorm:\"not null; \" json:\"description,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Roles",
          "type": "[]Role",
          "annotation": "gorm:\"association_foreignkey:AppID constraint:OnUpdate:SET NULL OnDelete:SET NULL\" json:\"roles,omitempty\"",
          "curd_flag": "true$false$*gorm.Modelfalse$false$true$false"
        }
      ]
    },
    {
      "name": "User",
      "rln_model": ["Role$mtm$user_roles"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Email",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"email,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Password",
          "type": "string",
          "annotation": "gorm:\"not null;\" json:\"password,omitempty\"",
          "curd_flag": "false$false$false$true$false$false"
        },
        {
          "name": "DateRegistred",
          "type": "time.Time",
          "annotation": "gorm:\"constraint:not null; default:current_timestamp;\" json:\"date_registered,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Disabled",
          "type": "bool",
          "annotation": "gorm:\"default:true; constraint:not null;\" json:\"disabled\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "UUID",
          "type": "string",
          "annotation": "gorm:\"constraint:not null; unique; type:string;\" json:\"uuid\"",
          "curd_flag": "true$false$false$false$true$false"
        },
        {
          "name": "Roles",
          "type": "[]Role",
          "annotation": "gorm:\"many2many:user_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"roles,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "Feature",
      "rln_model": ["Endpoint$otm"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Description",
          "type": "string",
          "annotation": "gorm:\"not null; \" json:\"description,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Active",
          "type": "bool",
          "annotation": "gorm:\"default:true; constraint:not null;\" json:\"active\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "RoleID",
          "type": "sql.NullInt64",
          "annotation": "gorm:\"foreignkey:RoleID OnDelete:SET NULL\" json:\"role,omitempty\" swaggertype:\"number\"",
          "curd_flag": "false$false$false$false$true$false"
        },
        {
          "name": "Endpoints",
          "type": "[]Endpoint",
          "annotation": "gorm:\"association_foreignkey:FeatureID constraint:OnUpdate:CASCADE,OnDelete:SET NULL\" json:\"endpoints,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "Endpoint",
      "rln_model": [],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "RoutePath",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"route_path,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Method",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"method,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Description",
          "type": "string",
          "annotation": "gorm:\"not null; \" json:\"description,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "FeatureID",
          "type": "sql.NullInt64",
          "annotation": "gorm:\"foreignkey:FeatureID default:NULL;,OnDelete:SET NULL;\" json:\"feature_id,omitempty\" swaggertype:\"number\"",
          "curd_flag": "true$false$false$true$true$false"
        }
      ]
    },
    {
      "name": "Page",
      "rln_model": [],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Active",
          "type": "bool",
          "annotation": "gorm:\"default:true; constraint:not null;\" json:\"active\"",
          "curd_flag": "true$false$true$true$false$false"
        },
        {
          "name": "Description",
          "type": "string",
          "annotation": "gorm:\"not null; \" json:\"description,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Roles",
          "type": "[]Role",
          "annotation": "gorm:\"many2many:page_roles; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"roles,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "JWTSalt",
      "rln_model": [],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "SaltA",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"salt_a,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "SaltB",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"salt_b,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        }
      ]
    }
  ]
}`

var blueAuthJSONTemplate = `
{
  "project_name": "{{.ProjectName}}",
  "app_name": "{{.AppName}}",
  "models": [
    {
      "name": "Permission",
      "rln_model": ["User$mtm$user_permissions"],
      "search_fields": ["codename"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "ContentTypeID",
          "type": "sql.NullInt64",
          "annotation": "gorm:\"foreignkey:ContentTypeID default:NULL;,OnDelete:SET NULL;\" json:\"content_type_id,omitempty\" swaggertype:\"number\"",
          "curd_flag": "true$false$false$true$true$false"
        },
        {
          "name": "Codename",
          "type": "string",
          "annotation": "gorm:\"not null; unique; \" json:\"codename,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Users",
          "type": "[]User",
          "annotation": "gorm:\"many2many:user_permissions; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"users,omitempty\"",
          "curd_flag": "false$false$false$false$false$true"
        },
        {
          "name": "Groups",
          "type": "[]Group",
          "annotation": "gorm:\"many2many:group_permissions; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"groups,omitempty\"",
          "curd_flag": "false$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "User",
       "search_fields": ["username","email","first_name","last_name"],
      "rln_model": [
        "Permission$mtm$user_permissions",
        "Group$mtm$user_groups"
      ],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Password",
          "type": "string",
          "annotation": "gorm:\"not null;\" json:\"password,omitempty\"",
          "curd_flag": "false$true$true$true$false$false"
        },
        {
          "name": "LastLogin",
          "type": "time.Time",
          "annotation": "gorm:\"constraint:not null; default:current_timestamp;\" json:\"last_login,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "IsSuperuser",
          "type": "bool",
          "annotation": "gorm:\"default:false; constraint:not null;\" json:\"is_superuser\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Username",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"username,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "FirstName",
          "type": "string",
          "annotation": "gorm:\"constraint:not null; type:string;\" json:\"first_name\"",
          "curd_flag": "true$true$true$false$true$false"
        },
        {
          "name": "LastName",
          "type": "string",
          "annotation": "gorm:\"constraint:not null; type:string;\" json:\"last_name\"",
          "curd_flag": "true$true$true$false$true$false"
        },
        {
          "name": "Email",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"email,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "IsStaff",
          "type": "bool",
          "annotation": "gorm:\"default:false; constraint:not null;\" json:\"is_staff\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "IsActive",
          "type": "bool",
          "annotation": "gorm:\"default:true; constraint:not null;\" json:\"is_active\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Groups",
          "type": "[]Group",
          "annotation": "gorm:\"many2many:user_groups; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"groups,omitempty\"",
          "curd_flag": "fasle$false$false$false$false$true"
        },
        {
          "name": "Permissions",
          "type": "[]Permission",
          "annotation": "gorm:\"many2many:user_permissions; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"permissions,omitempty\"",
          "curd_flag": "fasle$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "ContentType",
      "rln_model": [],
       "search_fields": [],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "AppLabel",
          "type": "string",
          "annotation": "gorm:\"not null;\" json:\"app_label,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Model",
          "type": "string",
          "annotation": "gorm:\"not null;\" json:\"model,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        }
      ]
    },
    {
      "name": "Group",
      "rln_model": ["Permission$mtm$group_permissions"],
       "search_fields": ["name"],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "Name",
          "type": "string",
          "annotation": "gorm:\"not null;\" json:\"name,omitempty\"",
          "curd_flag": "true$true$true$true$false$false"
        },
        {
          "name": "Users",
          "type": "[]User",
          "annotation": "gorm:\"many2many:user_groups; constraint:OnUpdate:CASCADE; OnDelete:CASCADE;\" json:\"users,omitempty\"",
          "curd_flag": "true$false$false$false$false$true"
        }
      ]
    },
    {
      "name": "JWTSalt",
      "rln_model": [],
       "search_fields": [],
      "fields": [
        {
          "name": "ID",
          "type": "uint",
          "annotation": "gorm:\"primaryKey;autoIncrement:true\" json:\"id,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "SaltA",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"salt_a,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        },
        {
          "name": "SaltB",
          "type": "string",
          "annotation": "gorm:\"not null; unique;\" json:\"salt_b,omitempty\"",
          "curd_flag": "true$false$false$false$false$false"
        }
      ]
    }
  ]
}

`
