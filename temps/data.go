package temps

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
	TableName       string  `json:"table_name"`
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
	NormalModelName  string `json:"normal_model_name"`
	ModelName        string `json:"model_name"`
	Name             string `json:"name"`
	LowerName        string `json:"lower_name"`
	Type             string `json:"type"`
	UpperType        string `json:"upper_type"`
	Annotation       string `json:"annotation"`
	MongoAnnotation  string `json:"mongo_annotation"`
	CurdFlag         string `json:"curd_flag"`
	Get              bool   `json:"get"`
	Post             bool   `json:"post"`
	Patch            bool   `json:"patch"`
	Put              bool   `json:"put"`
	OtM              bool   `json:"otm"`
	MtM              bool   `json:"mtm"`
	ProjectName      string `json:"project_name"`
	AppName          string `json:"app_name"`
	BackTick         string `json:"back_tick"`
	RandomFeildValue string `json:"random_field_value"`
}

var RenderData Data

// Generate random data for different field types
func generateRandomValue(fieldType string) string {
	switch fieldType {
	case "string":
		return fmt.Sprintf("\"%s\"", randomString(10)) // Random string of length 10
	case "int", "int32", "int64":
		return fmt.Sprintf("%d", rand.Intn(1000)) // Random int between 0 and 1000
	case "float64":
		return fmt.Sprintf("%f", rand.Float64()*100.0) // Random float between 0.0 and 100.0
	case "bool":
		return fmt.Sprintf("%t", rand.Intn(2) == 0) // Random bool
	case "time.Time":
		// Generate a random date/time within a certain range (e.g., last year to now)
		start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Now()
		duration := end.Sub(start)
		randomDuration := time.Duration(rand.Int63n(int64(duration)))
		randomTime := start.Add(randomDuration)
		return fmt.Sprintf("\"%s\"", randomTime.Format(time.RFC3339))
	case "ID":

		return fmt.Sprintf("%v", rand.Intn(1000))
	case "sql.NullInt64":
		// Randomly decide if it should be valid or NULL
		return fmt.Sprintf("%f", rand.Float64()*100.0)
	default:
		return "\"\""
	}
}

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
			RenderData.Models[i].Fields[j].NormalModelName = RenderData.Models[i].Name
			RenderData.Models[i].Fields[j].LowerName = strings.ToLower(RenderData.Models[i].Fields[j].Name)
			RenderData.Models[i].Fields[j].UpperType = CapitalizeFirstLetter(RenderData.Models[i].Fields[j].Type)
			RenderData.Models[i].Fields[j].Get, _ = strconv.ParseBool(cf[0])
			RenderData.Models[i].Fields[j].Post, _ = strconv.ParseBool(cf[1])
			RenderData.Models[i].Fields[j].Patch, _ = strconv.ParseBool(cf[2])
			RenderData.Models[i].Fields[j].Put, _ = strconv.ParseBool(cf[3])
			RenderData.Models[i].Fields[j].AppName = RenderData.AppName
			RenderData.Models[i].Fields[j].ProjectName = RenderData.ProjectName
			RenderData.Models[i].Fields[j].BackTick = "`"
			RenderData.Models[i].Fields[j].RandomFeildValue = generateRandomValue(RenderData.Models[i].Fields[j].Type)

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
			if len(rmf) > 2 {
				cur_relation.TableName = rmf[2]
			}
			cur_relation.ParentFields = RenderData.Models[i].Fields
			rl_list = append(rl_list, cur_relation)
			RenderData.Models[i].Relations = rl_list
		}

	}
	return nil
}
