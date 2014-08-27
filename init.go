package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

// Page defines what all pages should have
type Page struct {
	Title   string
	Content map[string][]string
}

// Index shows the index and landing splash/banner
// TODO: Handle banner here
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p := Page{"Go Vigilant!", nil}
	RenderPage("index", "index.tmpl", p, w)
}

// HostStat shows local machine sensor information (windows only atm)
// TODO: Add linux support (/proc & /sys)
func HostStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	type Sensor struct {
		Name       string
		Identifier string
		SensorType string
		Parent     string
		Value      *float32
		Index      int
	}

	// init COM, oh yeah
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer", nil, "Root\\OpenHardwareMonitor")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Sensor")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	dst := make(map[string][]string)

	for i := 0; i < count; i++ {
		// item is a SWbemObject, but really a type Sensor struct
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		defer item.Release()

		name, _ := oleutil.GetProperty(item, "Name")
		// For some reason OLE won't output floats :(
		// value, _ := oleutil.GetProperty(item, "Value")
		sensorType, _ := oleutil.GetProperty(item, "SensorType")

		// This _might_ combine _every_ sensor with the same name into one, Oopsies
		dst[name.Value().(string)] = append(dst[name.Value().(string)], name.Value().(string))
		dst[name.Value().(string)] = append(dst[name.Value().(string)], "12") // Fake it til you make it
		dst[name.Value().(string)] = append(dst[name.Value().(string)], sensorType.Value().(string))
	}

	p := Page{"Host status", dst}
	RenderPage("index", "host.tmpl", p, w)
}

// RenderPage takes a template and a Page struct and constructs & render a webpage
func RenderPage(page string, tmpl string, p Page, w http.ResponseWriter) {
	t := template.New(page)
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = t.Execute(w, p)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/host", HostStat)
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")

}
