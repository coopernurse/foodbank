package ui

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func GetResourceBundle(c echo.Context) *ResourceBundle {
	lang := c.QueryParam("lang")
	if lang == "" {
		lang = c.FormValue("lang")
	}

	res, ok := resources[lang]
	if !ok {
		lang = "en"
		res = resources[lang]
	}

	return &ResourceBundle{Lang: lang, Resources: res}
}

var resources = map[string]map[string]string{
	"en": {
		"signup.title":        "Community Cupboard Sign-Up Form",
		"signup.intro":        `This information is helpful in providing our services. None of your information will be shared.`,
		"signup.success":      `We have saved your information. Please ask for a shopping sheet from a staff member.`,
		"signup.hoh":          "Head of Household",
		"signup.othermembers": "Others Living in the Household",
		"misc.firstname":      "First Name",
		"misc.lastname":       "Last Name",
		"misc.address":        "Address",
		"misc.city":           "City",
		"misc.zipcode":        "ZIP Code",
		"misc.email":          "Email",
		"misc.phone":          "Phone",
		"misc.gender":         "Gender",
		"misc.male":           "Male",
		"misc.female":         "Female",
		"misc.prefernottosay": "Prefer not to say",
		"misc.dob":            "Date of Birth",
		"misc.month":          "Month",
		"misc.day":            "Day",
		"misc.year":           "Year",
		"misc.primarylang":    "Primary Language",
		"misc.english":        "English",
		"misc.spanish":        "Spanish",
		"misc.other":          "Other",
		"misc.relationship":   "Relationship",
		"misc.child":          "Child",
		"misc.grandchild":     "Grandchild",
		"misc.spouse":         "Spouse",
		"misc.parent":         "Parent",
		"misc.grandparent":    "Grandparent",
		"misc.sibling":        "Sibling",
		"misc.friend":         "Friend",
		"misc.race":           "Race",
		"misc.race.white":     "White/Anglo",
		"misc.race.latino":    "Latina/Latino",
		"misc.race.black":     "Black/Afr. American",
		"misc.race.asian":     "Asian/Pacific Islander",
		"misc.submit":         "Submit",
		"misc.thankyou":       "Thank You",
		"misc.fieldrequired":  "This field is required",
	},
	"es": {
		"signup.title":        "Formulario de Registro de Community Cupboard",
		"signup.intro":        `Esta información es útil para proporcionar nuestros servicios. Ninguna de su información será compartida.`,
		"signup.success":      `Hemos guardado su información. Por favor, solicite una hoja de compras a un miembro del personal.`,
		"signup.hoh":          "Cabeza de Familia",
		"signup.othermembers": "Otras Personas en el Hogar",
		"misc.firstname":      "Nombre",
		"misc.lastname":       "Apellido",
		"misc.address":        "Dirección",
		"misc.city":           "Ciudad",
		"misc.zipcode":        "Código Postal",
		"misc.email":          "Correo Electrónico",
		"misc.phone":          "Teléfono",
		"misc.gender":         "Género",
		"misc.male":           "Hombre",
		"misc.female":         "Mujer",
		"misc.prefernottosay": "Prefiere no decir",
		"misc.dob":            "Fecha de Nacimiento",
		"misc.month":          "Mes",
		"misc.day":            "Día",
		"misc.year":           "Año",
		"misc.primarylang":    "Idioma Principal",
		"misc.english":        "Inglés",
		"misc.spanish":        "Español",
		"misc.other":          "Otro",
		"misc.relationship":   "Relación",
		"misc.child":          "Hijo/a",
		"misc.grandchild":     "Nieto/a",
		"misc.spouse":         "Esposo/a",
		"misc.parent":         "Padre/Madre",
		"misc.grandparent":    "Abuelo/a",
		"misc.sibling":        "Hermano/a",
		"misc.friend":         "Amigo/a",
		"misc.race":           "Raza",
		"misc.race.white":     "Blanco/Anglo",
		"misc.race.latino":    "Latino/Latina",
		"misc.race.black":     "Negro/Afroamericano",
		"misc.race.asian":     "Asiático/Isleño del Pacífico",
		"misc.submit":         "Enviar",
		"misc.thankyou":       "Gracias",
		"misc.fieldrequired":  "Este campo es obligatorio",
	},
}

type ResourceBundle struct {
	Lang      string
	Resources map[string]string
}

func (r *ResourceBundle) Get(key string) string {
	v, ok := r.Resources[key]
	if !ok {
		return fmt.Sprintf("resource not found. lang=%s key=%s", r.Lang, key)
	}
	return v
}
