package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// CapitalizeFirstLetter capitalizes the first letter of the input string
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s // Return empty string if input is empty
	}

	// Convert the first letter to uppercase and concatenate with the rest of the string
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

// Define a struct to match the structure of your JSON data
type Data struct {
	ProjectName string  `json:"project_name"`
	AppName     string  `json:"app_name"`
	BackTick    string  `json:"back_tick"`
	Models      []Model `json:"models"`
}

type Model struct {
	Name        string         `json:"name"`
	LowerName   string         `json:"lower_name"`
	RlnModel    []string       `json:"rln_model"` // value to one of the models defined in the config json file
	BackTick    string         `json:"back_tick"`
	Fields      []Field        `json:"fields"`
	ProjectName string         `json:"project_name"`
	AppName     string         `json:"app_name"`
	Relations   []Relationship `json:"relations"`
}

type Relationship struct {
	ParentName      string  `json:"parent_name"`
	LowerParentName string  `json:"lower_parent_name"`
	FieldName       string  `json:"field_name"`
	LowerFieldName  string  `json:"lower_field_name"`
	ParentFields    []Field `json:"parent_fields"`
	ChildFields     []Field `json:"child_fields"`
	MtM             bool    `json:"mtm"`
	OtM             bool    `json:"otm"`
	MtO             bool    `json:"mto"`
	BackTick        string  `json:"back_tick"`
}

type Field struct {
	ModelName       string `json:"model_name"`
	Name            string `json:"name"`
	LowerName       string `json:"lower_name"`
	Type            string `json:"type"`
	UpperType       string `json:"upper_type"`
	Annotation      string `json:"annotation"`
	MongoAnnotation string `json:"mongo_annotation"`
	CurdFlag        string `json:"curd_flag"`
	Get             bool   `json:"get"`
	Post            bool   `json:"post"`
	Patch           bool   `json:"patch"`
	Put             bool   `json:"put"`
	OtM             bool   `json:"otm"`
	MtM             bool   `json:"mtm"`
	ProjectName     string `json:"project_name"`
	AppName         string `json:"app_name"`
	BackTick        string `json:"back_tick"`
}

var RenderData Data

func LoadData(file_name string) error {
	if file_name == "" {
		file_name = "config.json"
	}

	// Open the JSON file
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return err
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&RenderData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}
	// setting default value for config data file
	//  GetPostPatchPut
	// "Get$Post$Patch$Put$OtM$MtM"

	RenderData.BackTick = "`"
	for i := 0; i < len(RenderData.Models); i++ {
		RenderData.Models[i].LowerName = strings.ToLower(RenderData.Models[i].Name)
		RenderData.Models[i].AppName = RenderData.AppName
		RenderData.Models[i].ProjectName = RenderData.ProjectName
		RenderData.Models[i].BackTick = "`"

		//#####################################
		for j := 0; j < len(RenderData.Models[i].Fields); j++ {
			RenderData.Models[i].Fields[j].BackTick = "`"
			cf := strings.Split(RenderData.Models[i].Fields[j].CurdFlag, "$")
			RenderData.Models[i].Fields[j].ModelName = strings.ToLower(RenderData.Models[i].Name)
			RenderData.Models[i].Fields[j].LowerName = strings.ToLower(RenderData.Models[i].Fields[j].Name)
			RenderData.Models[i].Fields[j].UpperType = CapitalizeFirstLetter(RenderData.Models[i].Fields[j].Type)
			RenderData.Models[i].Fields[j].Get, _ = strconv.ParseBool(cf[0])
			RenderData.Models[i].Fields[j].Post, _ = strconv.ParseBool(cf[1])
			RenderData.Models[i].Fields[j].Patch, _ = strconv.ParseBool(cf[2])
			RenderData.Models[i].Fields[j].Put, _ = strconv.ParseBool(cf[3])
			RenderData.Models[i].Fields[j].AppName = RenderData.AppName
			RenderData.Models[i].Fields[j].ProjectName = RenderData.ProjectName
			RenderData.Models[i].Fields[j].BackTick = "`"

		}
		//#####################################
		rl_list := make([]Relationship, 0)
		for k := 0; k < len(RenderData.Models[i].RlnModel); k++ {
			rmf := strings.Split(RenderData.Models[i].RlnModel[k], "$")
			cur_relation := Relationship{
				ParentName:      RenderData.Models[i].Name,
				LowerParentName: RenderData.Models[i].LowerName,
				FieldName:       rmf[0],
				LowerFieldName:  strings.ToLower(rmf[0]),
				MtM:             rmf[1] == "mtm",
				OtM:             rmf[1] == "otm",
				MtO:             rmf[1] == "mto",
				BackTick:        "`",
			}
			cur_relation.ParentFields = RenderData.Models[i].Fields
			rl_list = append(rl_list, cur_relation)
			RenderData.Models[i].Relations = rl_list
		}

	}
	return nil
}
