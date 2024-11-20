package ui

import (
	"cupboard/internal/db"
	"cupboard/internal/model"
	"fmt"
	"net/http"

	"github.com/julvo/htmlgo"
	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type SignupPage struct {
	DB *db.FirestoreDB
}

func (p *SignupPage) GET(c echo.Context) error {
	return p.getPage(c, map[string]string{})
}

func toHousehold(c echo.Context) model.Household {
	h := model.Household{
		Id:      ulid.Make().String(),
		Head:    toPerson("hoh", c),
		Members: []model.Person{},
	}
	for i := 0; i < 5; i++ {
		prefix := fmt.Sprintf("person%d", i)
		val := c.FormValue(prefix + "FirstName")
		if val != "" {
			h.Members = append(h.Members, toPerson(prefix, c))
		}
	}
	return h
}

func toPerson(prefix string, c echo.Context) model.Person {
	return model.Person{
		FirstName:    c.FormValue(prefix + "FirstName"),
		LastName:     c.FormValue(prefix + "LastName"),
		Email:        c.FormValue(prefix + "Email"),
		Street:       c.FormValue(prefix + "Street"),
		City:         c.FormValue(prefix + "City"),
		PostalCode:   c.FormValue(prefix + "Zip"),
		Phone:        c.FormValue(prefix + "Phone"),
		Gender:       c.FormValue(prefix + "Gender"),
		DOB:          c.FormValue(prefix+"DobYear") + "-" + c.FormValue(prefix+"DobMonth") + "-" + c.FormValue(prefix+"DobDay"),
		Race:         c.FormValue(prefix + "Race"),
		Language:     c.FormValue(prefix + "Language"),
		Relationship: c.FormValue(prefix + "Relationship"),
	}
}

func (p *SignupPage) POST(c echo.Context) error {
	ctx := c.Request().Context()

	rb := GetResourceBundle(c)
	errs := p.validate(c, rb)
	if len(errs) == 0 {
		household := toHousehold(c)

		// save signup data
		if err := p.DB.AddHousehold(ctx, household); err != nil {
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf("Failed to save household: %v", err))
		}

		// success page
		page :=
			Html5_(
				Head_(
					Meta(Attr(a.Charset("UTF-8"))),
					Meta(Attr(a.Name("viewport"), a.Content("width=device-width, initial-scale=1.0"))),
					Title_(HTML(rb.Get("signup.title"))),
					Link(Attr(a.Rel("stylesheet"), a.Href("https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"))),
				),
				Body_(
					Div(Attr(a.Class("container my-5")),
						Img(Attr(a.Src("/static/img/mend-logo.png"), a.Alt("Logo"), a.Width("300"), a.Class("mb-2"))),

						H1(Attr(a.Class("text-center")), Text(rb.Get("misc.thankyou"))),
						P(Attr(a.Class("text-center")), Text(rb.Get("signup.success"))),
					)))
		return c.HTML(200, string(page))
	} else {
		return p.getPage(c, errs)
	}
}

func (p *SignupPage) validate(c echo.Context, rb *ResourceBundle) ValidationErrors {
	requiredFields := []string{"hohFirstName", "hohLastName", "hohDobMonth", "hohDobDay", "hohDobYear"}
	errs := ValidationErrors{}

	for _, f := range requiredFields {
		val := c.FormValue(f)
		if val == "" {
			errs[f] = rb.Get("misc.fieldrequired")
		}
	}
	return errs
}

func (p *SignupPage) getPage(c echo.Context, errs ValidationErrors) error {
	fb := &FormBuilder{Errs: errs, C: c}
	rb := GetResourceBundle(c)
	page :=
		Html5_(
			Head_(
				Meta(Attr(a.Charset("UTF-8"))),
				Meta(Attr(a.Name("viewport"), a.Content("width=device-width, initial-scale=1.0"))),
				Title_(HTML(rb.Get("signup.title"))),
				Link(Attr(a.Rel("stylesheet"), a.Href("https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"))),
			),
			Body_(
				Div(Attr(a.Class("container my-5")),
					Img(Attr(a.Src("/static/img/mend-logo.png"), a.Alt("Logo"), a.Width("300"), a.Class("mb-2"))),

					H1(Attr(a.Class("text-center")), Text(rb.Get("signup.title"))),
					P(Attr(a.Class("text-center")), Text(rb.Get("signup.intro"))),

					Div(Attr(a.Class("text-center")),
						A(Attr(a.Href("?lang=en")), Text("English")),
						Text(" | "),
						A(Attr(a.Href("?lang=es")), Text("Español")),
					),

					Form(Attr(a.Action("/signup"), a.Method("POST")), p.formBody(fb, rb, errs)...),
				),
			))
	return c.HTML(200, string(page))
}

func (p *SignupPage) formBody(fb *FormBuilder, rb *ResourceBundle, errs ValidationErrors) []HTML {
	h := []htmlgo.HTML{H2(Attr(a.Class("my-4")), Text(rb.Get("signup.hoh")))}
	h = append(h, Input(Attr(a.Type("hidden"), a.Name("lang"), a.Value(rb.Lang))))
	h = append(h, p.personForm("hoh", true, fb, rb, errs)...)
	h = append(h, H2(Attr(a.Class("my-4")), Text(rb.Get("signup.othermembers"))))
	for i := 0; i < 5; i++ {
		h = append(h, p.personForm(fmt.Sprintf("person%d", i), false, fb, rb, errs)...)
		h = append(h, Hr_())
	}
	h = append(h, Div(Attr(a.Class("text-center mt-4")),
		Button(Attr(a.Class("btn btn-primary"), a.Type("submit")), Text(rb.Get("misc.submit"))),
	))
	return h
}

func (p *SignupPage) personForm(prefix string, headOfHousehold bool, fb *FormBuilder, rb *ResourceBundle, errs ValidationErrors) []HTML {
	h := []htmlgo.HTML{
		Div(Attr(a.Class("form-row")),
			fb.InputDiv("col-md-6", prefix+"FirstName", rb.Get("misc.firstname")),
			fb.InputDiv("col-md-6", prefix+"LastName", rb.Get("misc.lastname")),
		),
	}
	if headOfHousehold {
		h = append(h, Div(Attr(a.Class("form-row")),
			fb.InputDiv("col-md-6", prefix+"Street", rb.Get("misc.address")),
			fb.InputDiv("col-md-4", prefix+"City", rb.Get("misc.city")),
			fb.InputDiv("col-md-2", prefix+"Zip", rb.Get("misc.zipcode")),
		))
		h = append(h, Div(Attr(a.Class("form-row")),
			fb.InputDiv("col-md-6", prefix+"Email", rb.Get("misc.email")),
			fb.InputDiv("col-md-6", prefix+"Phone", rb.Get("misc.phone")),
		))
	}

	monthClass, monthErrEl := fb.GetFormClassAndValidationElem(prefix + "DobMonth")
	dayClass, dayErrEl := fb.GetFormClassAndValidationElem(prefix + "DobDay")
	yearClass, yearErrEl := fb.GetFormClassAndValidationElem(prefix + "DobYear")

	h = append(h, Div(Attr(a.Class("form-row")),
		fb.SelectDiv("col-md-6", prefix+"Gender", rb.Get("misc.gender"), []ValueLabel{
			// {Value: "male", Label: rb.Get("misc.male")},
			// {Value: "female", Label: rb.Get("misc.female")},
			// {Value: "optout", Label: rb.Get("misc.prefernottosay")},
		}),
		Div(Attr(a.Class("form-group col-md-6")),
			Label_(Text(rb.Get("misc.dob"))),
			Div(Attr(a.Class("form-row")),
				Div(Attr(a.Class("col")),
					Select(Attr(a.Class(monthClass), a.Name(prefix+"DobMonth")),
						fb.selectOptions(prefix+"DobMonth", monthValueLabels(rb.Get("misc.month")))...,
					),
					monthErrEl,
				),
				Div(Attr(a.Class("col")),
					Select(Attr(a.Class(dayClass), a.Name(prefix+"DobDay")),
						fb.selectOptions(prefix+"DobDay", dayValueLabels(rb.Get("misc.day")))...,
					),
					dayErrEl,
				),
				Div(Attr(a.Class("col")),
					Select(Attr(a.Class(yearClass), a.Name(prefix+"DobYear")),
						fb.selectOptions(prefix+"DobYear", yearValueLabels(rb.Get("misc.year")))...,
					),
					yearErrEl,
				),
			),
		),
	))

	var col2 htmlgo.HTML
	if headOfHousehold {
		col2 = fb.SelectDiv("col-md-6", prefix+"Language", rb.Get("misc.primarylang"), []ValueLabel{
			// {Value: "english", Label: rb.Get("misc.english")},
			// {Value: "spanish", Label: rb.Get("misc.spanish")},
			// {Value: "other", Label: rb.Get("misc.other")},
		})
	} else {
		col2 = fb.SelectDiv("col-md-6", prefix+"Relationship", rb.Get("misc.relationship"), []ValueLabel{
			// {Value: "child", Label: rb.Get("misc.child")},
			// {Value: "grandchild", Label: rb.Get("misc.grandchild")},
			// {Value: "spouse", Label: rb.Get("misc.spouse")},
			// {Value: "parent", Label: rb.Get("misc.parent")},
			// {Value: "grandparent", Label: rb.Get("misc.grandparent")},
			// {Value: "sibling", Label: rb.Get("misc.sibling")},
			// {Value: "friend", Label: rb.Get("misc.friend")},
			// {Value: "other", Label: rb.Get("misc.other")},
		})
	}

	h = append(h, Div(Attr(a.Class("form-row")),
		fb.SelectDiv("col-md-6", prefix+"Race", rb.Get("misc.race"), []ValueLabel{
			// {Value: "white", Label: rb.Get("misc.race.white")},
			// {Value: "latino", Label: rb.Get("misc.race.latino")},
			// {Value: "black", Label: rb.Get("misc.race.black")},
			// {Value: "asian", Label: rb.Get("misc.race.asian")},
			// {Value: "other", Label: rb.Get("misc.other")},
		}),
		col2,
	))

	return h
}
