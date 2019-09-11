package main

import (
	"encoding/xml"
	"fmt"
)

func toXmlApp(bytess []byte)(*XmlApplication, error){
	xmlApp := &XmlApplication{}
	err := xml.Unmarshal(bytess, xmlApp)
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}else{
		return xmlApp, nil
	}
}

type XmlApplication struct {
	ApplicationName string `xml:"name,attr"`
	PackageName string `xml:"packagename,attr"`
	Desc         string `xml:"desc,attr"`
	Controllers XmlControllers `xml:"controllers"`
	Tables      XmlTables      `xml:"tables"`
	Beans 		XmlBeans       `xml:"beans"`
}

type XmlControllers struct {
	ControllerList []XmlController `xml:"controller"`
}

type XmlController struct {
	Name string `xml:"name,attr"`
	Desc         string `xml:"desc,attr"`
	SkipLogin    bool `xml:"skip_login,attr"`
	Apis []XmlApi `xml:"api"`
	ApplicationName string `xml:"-"`
	PackageName string `xml:"-"`
}

type XmlApi struct {
	Name string `xml:"name,attr"`
	Desc string `xml:"desc,attr"`
	Method string `xml:"method,attr"`
	Function string `xml:"function,attr"`//page,tree
	Table string `xml:"table,attr"`
	ParamList []XmlApiParam `xml:"param"`
	Return XmlReturn `xml:"return"`
}

type XmlApiParam struct {
	Name string `xml:"name,attr"`
	TransType string `xml:"trans-type,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr"`
	Ref string `xml:"ref,attr"`
	Must bool `xml:"must,attr"`
	DefaultValue string `xml:"default-value,attr"`
}

type XmlReturn struct {
	Success XmlSuccess `xml:"success"`
	Failure XmlFailure `xml:"failure"`
}

type XmlSuccess struct {
	Ref string `xml:"ref,attr"`
	Desc string `xml:"desc,attr"`
}

type XmlFailure struct {
	Ref string `xml:"ref,attr"`
	Desc string `xml:"desc,attr"`
}

//
type XmlTables struct {
	TableList []XmlTable `xml:"table"`
}

type XmlTable struct {
	Name           string `xml:"name,attr"`
	Desc           string `xml:"desc,attr"`
	ImportDateTime bool
	ColumnList        []XmlColumn `xml:"column"`
}

type XmlColumn struct {
	Name         string `xml:"name,attr"`
	Caption      string `xml:"caption,attr"`
	IsNull       bool   `xml:"isNull,attr"`
	IsPK         bool   `xml:"isPK,attr"`
	IsUnique     bool   `xml:"isUnique,attr"`
	Size         int    `xml:"size,attr"`
	Type         string `xml:"type,attr"`
	DbType       string `xml:"dbtype,attr"`
	DefaultValue string `xml:"default-value,attr"`
}

type XmlBeans struct {
	BeanList []XmlBean `xml:"bean"`
}

type XmlBean struct {
	Name           string `xml:"name,attr"`
	Desc           string `xml:"desc,attr"`
	Inher		 string `xml:"inher,attr"`
	ImportDateTime bool
	PropList        []XmlProp `xml:"prop"`
}

type XmlProp struct {
	Name         string `xml:"name,attr"`
	Caption      string `xml:"caption,attr"`
	Type         string `xml:"type,attr"`
	DefaultValue string `xml:"default-value,attr"`
}
