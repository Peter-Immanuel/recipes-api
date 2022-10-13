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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Global variable

var ctx context.Context
var err error
var client *mongo.Client

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"  bson:"PublishedAt"`
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

func RecipesCollection() *mongo.Collection {
	return client.Database("demo").Collection("recipes")
}

func NewRecipeHandler(c *gin.Context) {

	collection := client.Database("demo").Collection("recipes")
	var recipe Recipe //serialize request payload with struct
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while inserting a new recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func ListRecipeHandler(c *gin.Context) {
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}

	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	collection := client.Database("demo").Collection("recipes")
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var newRecipe Recipe
	err := c.ShouldBindJSON(&newRecipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	filter := bson.D{{"_id", id}}
	newValue := bson.D{{"$set", bson.D{
		{"name", newRecipe.Name},
		{"nags", newRecipe.Tags},
		{"ingredients", newRecipe.Ingredients},
		{"instructions", newRecipe.Instructions},
	}}}

	result, err := collection.UpdateOne(ctx, filter, newValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

/*
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

*/
func GetSpecificRecipeHandler(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	filter := bson.D{{"_id", id}}
	collection := RecipesCollection()

	recipe := collection.FindOne(ctx, filter)

	fmt.Println()
	fmt.Println(recipe.Err())
	fmt.Println()

	if recipe == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Sorry Resource Not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, recipe)
	return
}

func main() {
	router := gin.Default()
	router.GET("/recipes", ListRecipeHandler)
	router.POST("/recipes", NewRecipeHandler)

	router.GET("/recipes/:id", GetSpecificRecipeHandler)
	router.PATCH("/recipes/:id", UpdateRecipeHandler)
	// router.DELETE("/recipes/:id", DeleteRecipeHandler)

	// router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
}
