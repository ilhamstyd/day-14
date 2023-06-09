package main

import (
	"context"
	"day-10/connection"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	ID           int
	Image        string
	ProjectName  string
	StartDate    string
	EndDate      string
	Duration     int
	Description  string
	Author       string
	Technologies []string
}

var dataProject = []Project{
	{
		ProjectName:  "gak punya duit rek",
		StartDate:    "07/06/2023",
		EndDate:      "08/06/2023",
		Duration:     1,
		Description:  "punya duit pusing gak punya duit lebih pusing",
		Author:       "wa doyok",
		Technologies: []string{"Node JS"},
	},
	{
		ProjectName:  "gak punya duit rek",
		StartDate:    "07/06/2023",
		EndDate:      "08/06/2023",
		Duration:     1,
		Description:  "punya duit pusing gak punya duit lebih pusing",
		Author:       "wa bewok",
		Technologies: []string{"Node JS"},
	},
	{
		ProjectName:  "gak punya duit rek",
		StartDate:    "07/06/2023",
		EndDate:      "09/06/2023",
		Duration:     3,
		Description:  "punya duit pusing gak punya duit lebih pusing",
		Author:       "wa kumis",
		Technologies: []string{"Node JS"},
	},
}

func main() {

	connection.DatabaseConnect()

	e := echo.New()

	e.Static("/public", "public")

	e.GET("/hello", helloword)
	e.GET("/", home)
	e.GET("/addProject", addProject)
	e.GET("/projeect-detail/:id", projectDetail)
	e.GET("/contactMe", contactMe)

	e.POST("/edit-project/:id", editProject)
	e.POST("/delete-project/:id", deleteProject)
	e.POST("/addFormProject", addFormProject)

	e.Logger.Fatal(e.Start("localhost:8000"))
}

func helloword(c echo.Context) error {
	return c.String(http.StatusOK, "helloworld")
}

func home(c echo.Context) error {

	item, _ := connection.Conn.Query(context.Background(), "SELECT id, description, image, name_project, technologies, start_date, end_date, duration FROM tb_projects")

	var result []Project
	for item.Next() {
		var each = Project{}

		err := item.Scan(&each.ID, &each.Description, &each.Image, &each.ProjectName, &each.Technologies, &each.StartDate, &each.EndDate, &each.Duration)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		each.Author = "ASTA"
		result = append(result, each)
	}

	projects := map[string]interface{}{
		"projects": result,
	}
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func addProject(c echo.Context) error {
	var template, err = template.ParseFiles("views/addProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return template.Execute(c.Response(), nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	for index, item := range dataProject {
		if id == index {
			ProjectDetail = Project{
				ProjectName:  item.ProjectName,
				StartDate:    item.StartDate,
				EndDate:      item.EndDate,
				Duration:     item.Duration,
				Description:  item.Description,
				Author:       item.Author,
				Technologies: item.Technologies,
			}
		}
	}

	item := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var template, err = template.ParseFiles("views/add-project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return template.Execute(c.Response(), item)
}

func contactMe(c echo.Context) error {
	var template, err = template.ParseFiles("views/contact-me.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return template.Execute(c.Response(), nil)
}

func addFormProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDateStr := c.FormValue("startDate")
	endDateStr := c.FormValue("endDate")
	description := c.FormValue("desc")
	technologies := c.Request().Form["technologies"]

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return err
	}

	duration := int(endDate.Sub(startDate).Hours() / 24)

	fmt.Println("Project name: ", projectName)
	fmt.Println("Start date: ", startDate)
	fmt.Println("End date: ", endDate)
	fmt.Println("Description: ", description)
	fmt.Println("Technologies: ", strings.Join(technologies, ", "))
	fmt.Println("Duration: ", duration, "hari")

	var newProject = Project{
		ProjectName:  projectName,
		StartDate:    startDateStr,
		EndDate:      endDateStr,
		Duration:     duration,
		Description:  description,
		Author:       "ASTA",
		Technologies: technologies,
	}

	dataProject = append(dataProject, newProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(delete echo.Context) error {
	i, _ := strconv.Atoi(delete.Param("id"))

	fmt.Println("index : ", i)

	dataProject = append(dataProject[:i], dataProject[i+1:]...)

	return delete.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(edit echo.Context) error {
	id, _ := strconv.Atoi(edit.Param("id"))
	fmt.Println("index : ", id)

	dataProject = append(dataProject[:id], dataProject[id+1:]...)
	return edit.Redirect(http.StatusMovedPermanently, "/addProject")
}
