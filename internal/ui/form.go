package ui

import (
	"github.com/julvo/htmlgo"
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/labstack/echo/v4"
)

type FormBuilder struct {
	Errs ValidationErrors
	C    echo.Context
}

func (f *FormBuilder) InputDiv(class string, name string, label string) HTML {
	inputClass, errorEl := f.GetFormClassAndValidationElem(name)
	val := f.C.FormValue(name)
	return Div(Attr(a.Class("form-group "+class)),
		Label(Attr(a.For(name)), Text(label)),
		Input(Attr(a.Type("text"), a.Class(inputClass), a.Name(name), a.Id(name), a.Value(val))),
		errorEl,
	)
}

func (f *FormBuilder) SelectDiv(class string, name string, label string, vals []ValueLabel) HTML {
	inputClass, errorEl := f.GetFormClassAndValidationElem(name)
	return Div(Attr(a.Class("form-group "+class)),
		Label(Attr(a.For(name)), Text(label)),
		Select(Attr(a.Class(inputClass), a.Name(name), a.Id(name)), f.selectOptions(name, vals)...),
		errorEl,
	)
}

func (f *FormBuilder) GetFormClassAndValidationElem(name string) (string, HTML) {
	valErr, hasErr := f.Errs[name]
	inputClass := "form-control"
	var errorEl HTML
	if hasErr {
		inputClass += " is-invalid"
		errorEl = Div(Attr(a.Class("invalid-feedback")), Text(valErr))
	}
	return inputClass, errorEl
}

func (f *FormBuilder) selectOptions(name string, vals []ValueLabel) []HTML {
	out := make([]htmlgo.HTML, len(vals))
	val := f.C.FormValue(name)
	for i, v := range vals {
		attrs := []a.Attribute{a.Value(v.Value)}
		if v.Value == val {
			attrs = append(attrs, a.Selected("selected"))
		}
		out[i] = Option(Attr(attrs...), Text(v.Label))
	}
	return out
}

type ValidationErrors map[string]string

type ValueLabel struct {
	Value string
	Label string
}

func yearValueLabels(label string) []ValueLabel {
	//v := make([]ValueLabel, 101)
	v := make([]ValueLabel, 1)
	v[0] = ValueLabel{Value: "", Label: label}
	// year := time.Now().Year() + 1
	// for i := 1; i < len(v); i++ {
	// 	y := year - i
	// 	v[i] = ValueLabel{
	// 		Value: strconv.Itoa(y),
	// 		Label: fmt.Sprintf("%02d", y),
	// 	}
	// }
	return v
}

func monthValueLabels(label string) []ValueLabel {
	return numberValueLabels(12, label)
}

func dayValueLabels(label string) []ValueLabel {
	return numberValueLabels(31, label)
}

func numberValueLabels(max int, label string) []ValueLabel {
	//v := make([]ValueLabel, max+1)
	v := make([]ValueLabel, 1)
	v[0] = ValueLabel{Value: "", Label: label}
	// for i := 1; i < len(v); i++ {
	// 	v[i] = ValueLabel{
	// 		Value: strconv.Itoa(i),
	// 		Label: fmt.Sprintf("%02d", i),
	// 	}
	// }
	return v
}
