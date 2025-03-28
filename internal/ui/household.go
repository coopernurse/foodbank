package ui

import (
	"fmt"
	"foodbank/internal/db"
	"net/http"

	. "github.com/julvo/htmlgo"
	a "github.com/julvo/htmlgo/attributes"
	"github.com/labstack/echo/v4"
)

type HouseholdListPage struct {
	DB *db.FirestoreDB
}

func (p *HouseholdListPage) GET(c echo.Context) error {
	ctx := c.Request().Context()

	deleteID := c.QueryParam("delete")
	if deleteID != "" {
		if err := p.DB.DeleteHousehold(ctx, deleteID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	households, _, err := p.DB.GetHouseholds(ctx, 50, "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	rows := make([]HTML, len(households))
	for i, h := range households {
		if h.Id != "" {
			rows[i] = Tr_(
				Td_(HTML(h.Created())),
				Td_(HTML(h.Head.LastName)),
				Td_(HTML(h.Head.FirstName)),
				Td_(HTML(FormatDOB(h.Head.DOB))),
				Td_(A(Attr(a.Href(fmt.Sprintf("/household/%s", h.Id))), HTML("view"))),
				Td_(A(Attr(a.Href(fmt.Sprintf("/households?delete=%s", h.Id)),
					a.Onclick("{.}", "return confirm('Are you sure you want to delete this household?')"),
				), HTML("delete"))),
			)
		}
	}

	page := Html5_(
		Head_(
			Meta(Attr(a.Charset("UTF-8"))),
			Meta(Attr(a.Name("viewport"), a.Content("width=device-width, initial-scale=1.0"))),
			Title_(HTML("Households")),
			Link(Attr(a.Rel("stylesheet"), a.Href("https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"))),
		),
		Body_(
			FontScalingStyle("1.1rem"),
			Div(Attr(a.Class("container my-5")),
				Img(Attr(a.Src("/static/img/mend-logo.png"), a.Alt("Logo"), a.Width("300"), a.Class("mb-2"))),

				H1_(HTML("Household Signups")),
				Table(Attr(a.Class("table table-striped")),
					Thead_(
						Th_(HTML("Created")),
						Th_(HTML("First Name")),
						Th_(HTML("Last Name")),
						Th_(HTML("Date of Birth")),
						Th_(HTML("View")),
						Th_(HTML("Delete")),
					),
					Tbody_(rows...)))))

	return c.HTML(http.StatusOK, string(page))
}

type HouseholdDetailPage struct {
	DB *db.FirestoreDB
}

func (p *HouseholdDetailPage) GET(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	household, err := p.DB.GetHouseholdByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to retrieve household with id %s: %v", id, err),
		})
	}

	page := Html5_(
		Head_(
			Meta(Attr(a.Charset("UTF-8"))),
			Meta(Attr(a.Name("viewport"), a.Content("width=device-width, initial-scale=1.0"))),
			Title_(HTML("Households")),
			Link(Attr(a.Rel("stylesheet"), a.Href("https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"))),
		),
		Body_(
			FontScalingStyle("1.1rem"),
			Style_(Text(``)),
			Div(Attr(a.Class("container my-5")),
				Img(Attr(a.Src("/static/img/mend-logo.png"), a.Alt("Logo"), a.Width("300"), a.Class("mb-2"))),

				// Household head details
				H2_(HTML("Head of Household")),
				Table(Attr(a.Class("table table-bordered")),
					Tbody_(
						Tr_(Td_(HTML("Date Created")), Td_(HTML(household.Created()))),
						Tr_(Td_(HTML("First Name")), Td_(HTML(household.Head.FirstName))),
						Tr_(Td_(HTML("Last Name")), Td_(HTML(household.Head.LastName))),
						Tr_(Td_(HTML("Date of Birth")), Td_(HTML(FormatDOB(household.Head.DOB)))),
						Tr_(Td_(HTML("Gender")), Td_(HTML(household.Head.Gender))),
						Tr_(Td_(HTML("Race")), Td_(HTML(household.Head.Race))),
						Tr_(Td_(HTML("Language")), Td_(HTML(household.Head.Language))),
						Tr_(Td_(HTML("Email")), Td_(HTML(household.Head.Email))),
						Tr_(Td_(HTML("Phone")), Td_(HTML(household.Head.Phone))),
						Tr_(Td_(HTML("Address")), Td_(HTML(fmt.Sprintf("%s, %s, %s %s",
							household.Head.Street, household.Head.City, household.Head.State, household.Head.PostalCode)))),
					),
				),
				// Household members details
				H2_(HTML("Other Household Members")),
				Table(Attr(a.Class("table table-striped")),
					Thead_(
						Tr_(
							Th_(HTML("First Name")),
							Th_(HTML("Last Name")),
							Th_(HTML("Date of Birth")),
							Th_(HTML("Relationship")),
							Th_(HTML("Gender")),
							Th_(HTML("Race")),
						),
					),
					Tbody_(func() []HTML {
						rows := make([]HTML, len(household.Members))
						for i, member := range household.Members {
							rows[i] = Tr_(
								Td_(HTML(member.FirstName)),
								Td_(HTML(member.LastName)),
								Td_(HTML(FormatDOB(member.DOB))),
								Td_(HTML(member.Relationship)),
								Td_(HTML(member.Gender)),
								Td_(HTML(member.Race)),
							)
						}
						return rows
					}()...),
				),
			)))

	return c.HTML(http.StatusOK, string(page))
}
