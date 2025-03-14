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

func CommonCMD() {

	// running go mod tidy finally
	if err := exec.Command("go", "get", "-u").Run(); err != nil {
		fmt.Printf("error go get: %v \n", err)
	}

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
}

`
