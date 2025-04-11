package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Incidente struct {
	ID              int    `json:"id" bson:"id"`
	Empleado        string `json:"empleado" bson:"empleado"`
	TipoEquipo      string `json:"tipo_equipo" bson:"tipo_equipo"`
	DetalleProblema string `json:"detalle_problema" bson:"detalle_problema"`
	DiaProblema     string `json:"dia_problema" bson:"dia_problema"`
	Estado          string `json:"estado" bson:"estado"`
}

var coll *mongo.Collection

func connectToMongoDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://Lipito:Rotom3-18@apiincidentes.5zad9yj.mongodb.net/?retryWrites=true&w=majority&appName=ApiIncidentes").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println("Error al conectar con MongoDB:", err)
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("Error al hacer ping a MongoDB:", err)
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	coll = client.Database("db_incidentes").Collection("incidentes")
}

func Estados(estado string) bool {
	validEstados := []string{"pendiente", "en proceso", "resuelto"}
	for _, valid := range validEstados {
		if estado == valid {
			return true
		}
	}
	return false
}

func getIncidentes(c *gin.Context) {
	var incidentes []Incidente
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error al obtener los datos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener los datos"})
		return
	}
	defer func() {
		if err := cursor.Close(context.TODO()); err != nil {
			fmt.Println("Error al cerrar el cursor:", err)
		}
	}()

	if err = cursor.All(context.TODO(), &incidentes); err != nil {
		fmt.Println("Error al decodificar los datos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo decodificar los datos"})
		return
	}
	c.JSON(http.StatusOK, incidentes)
}

func postIncidente(c *gin.Context) {
	var newIncidente Incidente
	if err := c.BindJSON(&newIncidente); err == nil {
		if !Estados(newIncidente.Estado) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido. Valores permitidos: pendiente, en proceso, resuelto"})
			return
		}
		fmt.Printf("Datos recibidos: %+v\n", newIncidente)
		newIncidente.ID = int(time.Now().Unix())
		_, err := coll.InsertOne(context.TODO(), newIncidente)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar el incidente"})
			return
		}
		c.JSON(http.StatusCreated, newIncidente)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos invalidos"})
	}
}

func getIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var incidente Incidente
	err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&incidente)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el incidente"})
		return
	}
	c.JSON(http.StatusOK, incidente)
}

func putIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedData struct {
		Estado string `json:"estado"`
	}
	if err := c.BindJSON(&updatedData); err == nil {
		if !Estados(updatedData.Estado) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido. Valores permitidos: pendiente, en proceso, resuelto"})
			return
		}

		var existingIncidente Incidente
		err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&existingIncidente)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el incidente"})
			return
		}

		existingIncidente.Estado = updatedData.Estado
		_, err = coll.UpdateOne(context.TODO(), bson.D{{"id", id}}, bson.D{{"$set", existingIncidente}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallido al actualizar el incidente"})
			return
		}
		c.JSON(http.StatusOK, existingIncidente)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
	}
}

func deleteIncidenteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := coll.DeleteOne(context.TODO(), bson.D{{"id", id}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el incidente"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Incidente eliminado"})
}

func main() {
	connectToMongoDB()

	router := gin.Default()

	// Habilitar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // Cambia esto si tu frontend está en otro origen
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/incidentes", getIncidentes)
	router.POST("/incidentes", postIncidente)
	router.GET("/incidentes/:id", getIncidenteById)
	router.PUT("/incidentes/:id", putIncidenteById)
	router.DELETE("/incidentes/:id", deleteIncidenteById)

	router.Run("localhost:8080")
}
