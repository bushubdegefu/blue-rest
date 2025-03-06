package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func CurdFrameEcho() {

	// ############################################################

	curd_tmpl, err := template.New("RenderData").Parse(curdTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("controllers", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range RenderData.Models {

		folder_path := fmt.Sprintf("controllers/%v_controller.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		curd_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = curd_tmpl.Execute(curd_file, model)
		if err != nil {
			panic(err)
		}
		curd_file.Close()

	}

}

var curdTemplate = `
package controllers

import (
	{{- range .Relations }}
	{{- if .MtM}}
	"math"
	{{- end}}
	{{- end}}

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"{{.ProjectName}}/common"
	"{{.ProjectName}}/models"
	"{{.ProjectName}}/observe"
)

// Get{{.Name}}is a function to get a {{.Name}}s by ID
// @Summary Get {{.Name}}s
// @Description Get {{.Name}}s
// @Tags {{.Name}}s
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.{{.Name}}Get}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}} [get]
func Get{{.Name}}s(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	//  Geting dbsession
     db := contx.Get("db").(*gorm.DB)

	//  querying result with pagination using gorm function
	result, err := common.PaginationPureModel(db, models.{{.Name}}{}, []models.{{.Name}}{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all {{.Name}}.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.JSON(http.StatusOK, result)
}

// Get{{.Name}}ByID is a function to get a {{.Name}}s by ID
// @Summary Get {{.Name}} by ID
// @Description Get {{.LowerName}} by ID
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Get}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [get]
func Get{{.Name}}ByID(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)


	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  Geting dbsession
     db := contx.Get("db").(*gorm.DB)

	// Preparing and querying database using Gorm
	var {{.LowerName}}s_get models.{{.Name}}Get
	var {{.LowerName}}s models.{{.Name}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.Name}}{}).Where("id = ?", id).First(&{{.LowerName}}s); res.Error != nil {
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode({{.LowerName}}s, &{{.LowerName}}s_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success got one {{.LowerName}}.",
		Data:    &{{.LowerName}}s_get,
	})
}

// Add {{.Name}} to data
// @Summary Add a new {{.Name}}
// @Description Add {{.Name}}
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Add {{.Name}}"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}} [post]
func Post{{.Name}}(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_{{.LowerName}} := new(models.{{.Name}}Post)

	//first parse request data
	if err := contx.Bind(&posted_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  initiate -> {{.LowerName}}
	{{.LowerName}} := new(models.{{.Name}})
	{{- range .Fields}} {{- if .Post}}
	{{.ModelName}}.{{.Name}} = posted_{{.ModelName}}.{{.Name}}
	{{- end}}
	{{- end}}

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "{{.Name}} Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} created successfully.",
		Data:    {{.LowerName}},
	})
}

// Patch {{.Name}} to data
// @Summary Patch {{.Name}}
// @Description Patch {{.Name}}
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Patch {{.Name}}"
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [patch]
func Patch{{.Name}}(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_{{.LowerName}} := new(models.{{.Name}}Patch)
	if err := contx.Bind(&patch_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var {{.LowerName}} models.{{.Name}}
	{{.LowerName}}.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).First(&{{.LowerName}}, {{.LowerName}}.ID).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerName}}).UpdateColumns(*patch_{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Return  success response
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} updated successfully.",
		Data:    {{.LowerName}},
	})
}

// Delete{{.Name}}s function removes a {{.LowerName}} by ID
// @Summary Remove {{.Name}} by ID
// @Description Remove {{.LowerName}} by ID
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [delete]
func Delete{{.Name}}(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)

	// get deleted {{.LowerName}} attributes to return
	var {{.LowerName}} models.{{.Name}}

	// validate path params
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()

	// first getting {{.LowerName}} and checking if it exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", id).First(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Delete the {{.LowerName}}
	if err := db.Delete(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error deleting {{.LowerName}}",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} deleted successfully.",
		Data:    {{.LowerName}},
	})
}
// ################################################################
// Relationship Based Endpoints
// ################################################################
{{- range .Relations }}
{{- if .MtM}}

// Add {{.FieldName}} to {{.ParentName}}
// @Summary Add {{.ParentName}} to {{.FieldName}}
// @Description Add {{.FieldName}} {{.ParentName}}
// @Tags {{.FieldName}}{{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [post]
func Add{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate path params
	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer.Tracer).Where(" id = ? ",uint({{.LowerParentName}}_id)).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  {{.LowerParentName}}ending assocation
	var {{.LowerFieldName}} models.{{.FieldName}}
	if err := db.WithContext(tracer.Tracer).Where(" id = ? ",uint({{.LowerFieldName}}_id)).First(&{{.LowerFieldName}}); err.Error != nil {
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Append(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "{{.ParentName}}ending {{.ParentName}} Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Creating a {{.LowerParentName}} {{.FieldName}}.",
		Data:    {{.LowerParentName}},
	})
}

// Get {{.FieldName}}s of {{.ParentName}}
// @Summary Get {{.ParentName}} to {{.FieldName}}
// @Description Get {{.FieldName}} {{.ParentName}}
// @Tags {{.FieldName}}{{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.{{.FieldName}}Get}
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [post]
func Get{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
	db := contx.Get("db").(*gorm.DB)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	// place holder for total number of items
	var total int64

	var	{{.LowerFieldName}}s []models.{{.FieldName}}
	join_string := "INNER JOIN {{.TableName}} ur ON {{.LowerFieldName}}s.id = ur.{{.LowerFieldName}}_id"
	filter_string := "{{.LowerParentName}}_id = ?"


	//  to make sure no more that 50 items will be queried per request
	if Limit > 100 {
		Limit = 100
	}


	// validate path params
	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// {{.LowerParentName}} to be added
	//getting total number of items
	if err := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Joins(join_string).Where(filter_string, {{.LowerParentName}}_id).Count(&total); err != nil {
			return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
				Success: false,
				Message: err.Error.Error(),
				Data:    nil,
			})
	}

	//  actual result query
	if err := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Joins(join_string).Where(filter_string, {{.LowerParentName}}_id).Order("id asc").Limit(int(Limit)).Offset(int(Page - 1)).Find(&{{.LowerFieldName}}s); err != nil {
			return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
				Success: false,
				Message: err.Error.Error(),
				Data:    nil,
			})
	}

	pages := math.Ceil(float64(total) / float64(Limit))
	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponsePagination{
		Success: true,
		Items:   {{.LowerFieldName}}s,
		Message: "working",
		Total:   uint(total),
		Page:    uint(Page),
		Size:    uint(Limit),
		Pages:   uint(pages),
	})
}

// Delete {{.ParentName}} to {{.FieldName}}
// @Summary Add {{.ParentName}}
// @Description Delete {{.FieldName}} {{.ParentName}}
// @Tags {{.FieldName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.ParentName}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [delete]
func Delete{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
	db := contx.Get("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil || {{.LowerFieldName}}_id == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil || {{.LowerParentName}}_id == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// fetching {{.LowerParentName}} to be deleted
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}.ID = uint({{.LowerParentName}}_id)
	if res := db.Find(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fettchng {{.LowerFieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	{{.LowerFieldName}}.ID = uint({{.LowerFieldName}}_id)
	if err := db.Find(&{{.LowerFieldName}}); err.Error != nil {
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	// removing {{.LowerParentName}}
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Delete(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNonAuthoritativeInfo, common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected H{{.LowerParentName}}ened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Removing a {{.LowerParentName}} from {{.LowerFieldName}}.",
		Data:    {{.LowerParentName}},
	})
}

{{- end}}
{{- end}}

{{- range .Relations }}
{{- if .OtM}}
// Get {{.FieldName}}s of {{.ParentName}}
// @Summary Get  {{.FieldName}}s of {{.ParentName}}
// @Description Get {{.FieldName}}s of {{.ParentName}}
// @Tags {{.ParentName}}{{.FieldName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.{{.FieldName}}Get}
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [post]
func Get{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
	db := contx.Get("db").(*gorm.DB)

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}


	// validate path params
	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// {{.LowerParentName}} to be added

	_, result, err := common.PaginationPureModelFilterOneToMany(db, models.{{.FieldName}}{}, []models.{{.FieldName}}{}, "{{.LowerParentName}}_id = ?", uint({{.LowerParentName}}_id), uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
			return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
	}


	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, result)
}


// Add {{.FieldName}} to {{.ParentName}}
// @Summary Add {{.FieldName}} to {{.ParentName}}
// @Description Add {{.FieldName}} to {{.ParentName}}
// @Tags {{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.ParentName}} ID"
// @Param {{.LowerParentName}}_id query int true "{{.FieldName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerParentName}}{{.LowerFieldName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [patch]
func Add{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching relation element
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerFieldName}} to be added
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.QueryParam("{{.LowerParentName}}_id"))
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer.Tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Append(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Error Adding Record",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Adding a {{.FieldName}} to {{.ParentName}}.",
		Data:    {{.LowerParentName}},
	})
}

// Delete {{.ParentName}} {{.FieldName}}
// @Summary Delete {{.ParentName}} {{.FieldName}}
// @Description Delete {{.ParentName}} {{.FieldName}}
// @Tags {{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.ParentName}} ID"
// @Param {{.LowerParentName}}_id query int true "{{.FieldName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerParentName}}{{.LowerFieldName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [delete]
func Delete{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	//  Geting tracer
		tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
		db := contx.Get("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting {{.FieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.QueryParam("{{.LowerParentName}}_id"))
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing {{.FieldName}} From {{.ParentName}}
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Delete(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Deleteing a {{.FieldName}} From {{.ParentName}}.",
		Data:    {{.LowerParentName}},
	})
}


{{- end}}
{{- end}}


`
