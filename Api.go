package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Incidente struct {
	ID              int    `json:"id"`
	Empleado        string `json:"empleado"`
	TipoEquipo      string `json:"tipo_equipo"`
	DetalleProblema string `json:"detalle_problema"`
	DiaProblema     string `json:"dia_problema"`
}

var posts = []Incidente{
	{
		ID:              1,
		Empleado:        "Juan Pérez",
		TipoEquipo:      "Computadora",
		DetalleProblema: "No enciende la pantalla",
		DiaProblema:     "2023-10-15",
	},
	{
		ID:              2,
		Empleado:        "María Gómez",
		TipoEquipo:      "Impresora",
		DetalleProblema: "Atascamiento de papel frecuente",
		DiaProblema:     "2023-10-16",
	},
	{
		ID:              3,
		Empleado:        "Carlos Ruiz",
		TipoEquipo:      "Red",
		DetalleProblema: "Conexión WiFi intermitente",
		DiaProblema:     "2023-10-17",
	},
	{
		ID:              4,
		Empleado:        "Ana López",
		TipoEquipo:      "Computadora",
		DetalleProblema: "Teclado no responde",
		DiaProblema:     "2023-10-18",
	},
	{
		ID:              5,
		Empleado:        "Pedro Sánchez",
		TipoEquipo:      "Teléfono",
		DetalleProblema: "No tiene tono de llamada",
		DiaProblema:     "2023-10-19",
	},
}

func getIncidentes(c *gin.Context) {
	c.JSON(http.StatusOK, posts)
}

func postIncidente(c *gin.Context) {
	var newIncidente Incidente
	if err := c.BindJSON(&newIncidente); err == nil {
		newIncidente.ID = len(posts) + 1
		posts = append(posts, newIncidente)
		c.JSON(http.StatusCreated, newIncidente)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
	}
}

func getIncidenteById(c *gin.Context) {
	id := c.Param("id")
	for _, post := range posts {
		if strconv.Itoa(post.ID) == id {
			c.JSON(http.StatusOK, post)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Incidente not found"})
}

func main() {
	router := gin.Default()
	router.GET("/incidentes", getIncidentes)
	router.POST("/incidentes", postIncidente)
	router.GET("/incidentes/:id", getIncidenteById)
	router.Run("localhost:8080")
}
