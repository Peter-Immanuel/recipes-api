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
	"os"
	"recipe-microservice/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Global variable
var recipeHandler *handlers.RecipeHandler

func init() {
	// Check for .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Check for database URI
	dbUri := os.Getenv("MONGODB_URI")
	if dbUri == "" {
		log.Fatal("Database connection string cannot be empty.")
	}

	// Create a context to use when querying mongodb instance
	ctx := context.Background()

	// Create the mongoDB Client to use in creating a handler
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB Instance")

	collection := client.Database("demo").Collection("recipes")

	recipeHandler = handlers.NewRecipeHandler(ctx, collection)

}

/*


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
*/
func main() {
	router := gin.Default()
	router.GET("/recipes", recipeHandler.ListRecipesHandler)
	router.POST("/recipes", recipeHandler.CreateRecipeHandler)

	router.GET("/recipes/:id", recipeHandler.GetRecipeByIDHandler)
	//router.PATCH("/recipes/:id", UpdateRecipeHandler)
	// router.DELETE("/recipes/:id", DeleteRecipeHandler)

	router.GET("/recipes/search", recipeHandler.SearchRecipeHandler)
	router.Run()
}
