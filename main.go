package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	connection "personal-web/connnection"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AddProject struct {
	ID          int
	StartDate   string
	EndDate     string
	Description string
	Image       string
	Author      string
	Name        string
}

func main() {
	connection.DatabaseConnect()

	e := echo.New()

	//serve static files from public directory
	e.Static("/public", "public")

	// ROuting
	e.GET("/Hello", helloworld)
	e.GET("/about", about)
	e.GET("/", home)
	e.GET("/contact", contactMe)
	e.GET("/addProject", addProject)
	e.POST("/add-add-project-detail", addaddprojectdetail)
	e.GET("/add-project-detail/:id", addProjectDetail)
	e.GET("/delete-addProject/:id", deleteaddProject)

	e.Logger.Fatal(e.Start("localhost:3000"))

}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func about(c echo.Context) error {
	return c.String(http.StatusOK, "ini adalah about anjay")
}

func home(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		// fmt.Println("tidak ada")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, start_date, end_date, description, image, name FROM tb_projects")
	fmt.Println(data)

	var result []AddProject

	for data.Next() {
		var each = AddProject{}

		err := data.Scan(&each.ID, &each.StartDate, &each.EndDate, &each.Description, &each.Image, &each.Name)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		each.Author = "ilham setyadji"

		result = append(result, each)
	}

	addProjects := map[string]interface{}{
		"AddProject": result,
	}

	return tmpl.Execute(c.Response(), addProjects)
}

func contactMe(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func addProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func addProjectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/add-project-detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	var AddProjectDetail = AddProject{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, start_date, end_date, description, image, name FROM tb_projects WHERE id = $1", id).Scan(&AddProjectDetail.ID, &AddProjectDetail.StartDate, &AddProjectDetail.EndDate, &AddProjectDetail.Description, &AddProjectDetail.Image, &AddProjectDetail.Name)

	AddProjectDetail.Author = "Ilham"

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"AddProject": AddProjectDetail,
	}

	return tmpl.Execute(c.Response(), data)
}

func addaddprojectdetail(c echo.Context) error {
	start_date := c.FormValue("start_date")
	end_date := c.FormValue("end_date")
	name := c.FormValue("name")
	description := c.FormValue("description")
	image := "image.png"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(start_date, end_date, name, description, image) VALUES ($1, $2, $3, $4, $5)", start_date, end_date, name, description, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteaddProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // id = 0 string => 0 int

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/addProject")
}
