// Recipes API
//This is a sample recipes API. You can find out more about
//
// Schemas: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Peter Bemshima
//
// Consumes:
// 	- application/json
//
// Produces:
// 	- application/json
//
// Swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Global variable

var ctx context.Context
var err error
var client *mongo.Client

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name`
	Tags         []string  `json:"tags`
	Ingredients  []string  `json:"ingredients`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt`
}

func (r *Recipe) setUp() {
	r.ID = xid.New().String()
	r.PublishedAt = time.Now()
}

var recipes []Recipe

func init() {

	// Check for .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// collect database uri including checks
	dbUri := os.Getenv("MONGODB_URI")
	if dbUri == "" {
		log.Fatal("Database connection string cannot be empty.")
	}

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB Instance")
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.setUp()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var newRecipe Recipe
	err := c.ShouldBindJSON(&newRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Sorry Recipe not found"})
		return
	}
	recipes[index] = newRecipe
	c.JSON(http.StatusOK, newRecipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not Found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1])
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe successfully deleted"})
}

func SearchRecipeHandler(c *gin.Context) {
	searchTag := c.Query("tag")
	resultList := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, tag := range recipes[i].Tags {
			if strings.EqualFold(tag, searchTag) {
				found = true
			}
		}
		if found {
			resultList = append(resultList, recipes[i])
		}
	}
	c.JSON(http.StatusOK, resultList)
}

func GetSpecificRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	found := false
	var result Recipe
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			found = true
			result = recipes[i]
			break
		}
	}
	if found {
		c.JSON(http.StatusOK, result)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found"})
}

func main() {
	router := gin.Default()
	router.GET("/recipes", ListRecipeHandler)
	router.POST("/recipes", NewRecipeHandler)

	router.GET("/recipes/:id", GetSpecificRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)

	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
}
