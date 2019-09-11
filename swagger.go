package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*******************api*******************/
func toSwaggerJson(app *XmlApplication)error{
	defer func(){
		if err := recover(); err != nil{
			fmt.Println(err)
		}
	}()
	swagger := xmlAppToSwagger(app)
	bytess, err := json.Marshal(swagger)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}else{
		err = os.MkdirAll("files", os.ModePerm)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		swaggerFileName := "files/" + app.ApplicationName + ".json"
		os.Remove(swaggerFileName)
		err = ioutil.WriteFile(swaggerFileName, bytess, os.ModePerm)
		if err != nil{
			fmt.Println(err)
			return err
		}
		return nil
	}
}

func xmlAppToSwagger(app *XmlApplication)*Swagger{
	swagger := &Swagger{}
	swagger.Swagger = "2.0"
	swagger.Host = "api.qianqiusoft.com"
	swagger.BasePath = "/api/v1"
	swagger.Schemes = []string{"http", "https"}
	swagger.Definitions = xmlTableBeansToDefinitions(&app.Tables, &app.Beans)
	swagger.Tags = xmlTableBeansToTags(&app.Tables, &app.Beans)
	swagger.Paths = xmlControllersToPaths(&app.Controllers)
	swagger.Info = xmlAppToSwaggerInfo(app)

	return swagger
}

func xmlAppToSwaggerInfo(app *XmlApplication)SwaggerInfo{
	wi := SwaggerInfo{}
	wi.Description = app.Desc
	wi.Title = app.Desc
	wi.Contact = Contact{"developer@qianqiusoft.com"}
	wi.License = License{"Apache 2.0", "http://www.apache.org/licenses/LICENSE-2.0.html"}
	wi.TermsOfService = "http://www.qianqiusoft.com"
	wi.Version = "1.0.0"
	return wi
}

func xmlTableBeansToDefinitions(table *XmlTables, beans *XmlBeans)map[string]StructInfo{
	Definitions := make(map[string]StructInfo)
	if table != nil{
		for i := range table.TableList{
			tableInfo := table.TableList[i]
			upperName := HandleUpperCasePrefix(tableInfo.Name)
			si := StructInfo{}
			si.Type = "object"
			si.Xml = XmlInfo{upperName}
			si.Properties = make(map[string]Property)
			for ci := range tableInfo.ColumnList{
				property := Property{Type: tableInfo.ColumnList[ci].Type, Description:tableInfo.ColumnList[ci].Caption}
				si.Properties[tableInfo.ColumnList[ci].Name] = property
			}
			si.Description = tableInfo.Desc
			Definitions[upperName] = si
		}
	}

	if beans != nil{
		for i := range beans.BeanList{
			bean := beans.BeanList[i]
			upperName := HandleUpperCasePrefix(bean.Name)
			si := StructInfo{}
			si.Type = "object"
			si.Xml = XmlInfo{upperName}
			si.Properties = make(map[string]Property)
			for ci := range bean.PropList{
				p := Property{Type: bean.PropList[ci].Type, Description:bean.PropList[ci].Caption}
				si.Properties[bean.PropList[ci].Name] = p
			}
			si.Description = bean.Desc
			Definitions[upperName]  = si
		}
	}
	si := StructInfo{}
	si.Type = "object"
	si.Properties = make(map[string]Property)
	si.Properties["code"] = Property{Type: "int", Description:"状态码"}
	si.Properties["msg"] = Property{Type:"string", Description:"相关信息"}
	si.Properties["data"] = Property{Type:"interface{}", Description:"返回值，不同接口返回不同对象"}
	Definitions["ApiResponse"] = si

	return Definitions
}

func xmlTableBeansToTags(table *XmlTables, beans *XmlBeans)[]Tag{
	tags := []Tag{}
	if table != nil{
		for i := range table.TableList{
			tableInfo := table.TableList[i]
			tag := Tag{tableInfo.Name, tableInfo.Desc}
			tags = append(tags, tag)
		}
	}

	if beans != nil{
		for i := range beans.BeanList{
			bean := beans.BeanList[i]
			tag := Tag{bean.Name, bean.Desc}
			tags = append(tags, tag)
		}
	}

	return tags
}

func xmlControllersToPaths(controller *XmlControllers)map[string]map[string]Method{
	paths := make(map[string]map[string]Method)
	for i := range controller.ControllerList{
		ctrl := controller.ControllerList[i]
		for ai := range ctrl.Apis{
			pstr := "/" + ctrl.Name + "/" + ctrl.Apis[ai].Name
			paths[pstr] = make(map[string]Method)
			sm := Method{}
			sm.Tags = []string{ctrl.Name}
			sm.Description = ctrl.Apis[ai].Desc
			sm.Consumes = []string{"application/json"}
			sm.Produces = []string{"application/json"}
			sm.OperationId = ctrl.Apis[ai].Name + HandleUpperCasePrefix(ctrl.Name)
			sm.Summary = ctrl.Apis[ai].Desc
			sm.Responses = make(map[string]HttpStatus)
			sm.Responses["200"] = HttpStatus{"操作成功", Schema{"#/definitions/ApiResponse"}}
			sm.Responses["500"] = HttpStatus{"内部操作失败", Schema{"#/definitions/ApiResponse"}}
			sm.Responses["404"] = HttpStatus{"接口不存在", Schema{"#/definitions/ApiResponse"}}
			sm.Responses["400"] = HttpStatus{"没有权限", Schema{"#/definitions/ApiResponse"}}
			sm.Parameters = []Parameter{}
			for pi := range ctrl.Apis[ai].ParamList{
				xp := ctrl.Apis[ai].ParamList[pi]
				sp := Parameter{}
				sp.Name = xp.Name
				sp.Description = xp.Desc
				if xp.Ref != ""{
					sp.In = "body"
					tn := HandleUpperCasePrefix(strings.TrimPrefix(xp.Ref, "$"))
					fmt.Println("----------------------------------->", tn)
					sp.Schema =  make(map[string]interface{})
					if strings.HasSuffix(tn, " array"){
						fmt.Println("==================================>", strings.TrimSuffix(tn, " array"))
						sp.Schema["type"] = "array"
						items := make(map[string]interface{})
						items["$ref"] = "#/definitions/" + strings.TrimSuffix(tn, " array")
						sp.Schema["items"] = items
					}else{
						sp.Schema["$ref"] = "#/definitions/" + tn
					}
				}else {
					sp.In = "query"
					sp.Type = xp.Type
				}
				sm.Parameters = append(sm.Parameters, sp)
			}
			paths[pstr][strings.ToLower(ctrl.Apis[ai].Method)] = sm
		}
	}
	return paths
}

/**
 * 字符串首字母转化为大写 ios_bbbbbbbb -> IosBbbbbbbbb
 */
func HandleUpperCasePrefix(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

/*******************entity*************************/
type Swagger struct {
	Swagger string `json:"swagger"`
	Info SwaggerInfo `json:"info"`
	Host string `json:"host"`
	BasePath string `json:"basePath"`
	Tags []Tag `json:"tags"`
	Schemes []string `json:"schemes"`
	Paths map[string]map[string]Method `json:"paths"`
	Definitions map[string]StructInfo `json:"definitions"`
}

type SwaggerInfo struct {
	Description string `json:"description"`
	Version string `json:"version"`
	Title string `json:"title"`
	TermsOfService string `json:"termsOfService"`
	Contact Contact `json:"contact"`
	License License `json:"license"`
}

type Contact struct{
	Email string `json:"email"`
}

type License struct{
	Name string `json:"name"`
	Url string `json:"url"`
}

type Tag struct{
	Name string `json:"name"`
	Description string `json:"description"`
}

type Method struct {
	Tags []string `json:"tags"`
	Summary string `json:"summary"`
	Description string `json:"description"`
	OperationId string `json:"operationId"`
	Consumes []string `json:"consumes"`
	Produces []string `json:"produces"`
	Parameters []Parameter `json:"parameters"`
	Responses map[string]HttpStatus `json:"responses"`
}

type Parameter struct{
	In string `json:"in"`
	Name string `json:"name"`
	Description string `json:"description"`
	Required bool `json:"required"`
	Schema map[string]interface{} `json:"schema"`
	Type string `json:"type"`
}

type Schema struct {
	Ref string `json:"$ref"`
}

type HttpStatus struct {
	Description string `json:"description"`
	Schema Schema `json:"schema"`
}

type StructInfo struct {
	Type string `json:"type"`
	Properties map[string]Property `json:"properties"`
	Xml XmlInfo `json:"xml"`
	Description string `json:"description"`
}

type Property struct {
	Type string `json:"type"`
	Description string `json:"description"`
	//Format string `format`
	//Enum []string `json:"enum"`
	//Default interface{} `json:"default"`
}

type XmlInfo struct {
	Name string `json:"name"`
}
