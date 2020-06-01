package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"io"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"math/rand"
	"time"
	"strings"
)

type Urls struct {
	Body string
	Domain string
	Header string
	Language string
	Url string
}

type FUTUREResponse struct {
	Answer string
	Chatbot int
	Corrected string
	Images []string
	Map string
	NRes int
	Reply string
	SmallSummary string
	Time float64
	Urls []Urls
}

type FUTUREResult struct {
	Result FUTUREResponse
}

func main() {

	left := [...]string{
		"admiring",
		"adoring",
		"affectionate",
		"agitated",
		"amazing",
		"angry",
		"awesome",
		"blissful",
		"boring",
		"brave",
		"clever",
		"cocky",
		"compassionate",
		"competent",
		"condescending",
		"confident",
		"cranky",
		"dazzling",
		"determined",
		"distracted",
		"dreamy",
		"eager",
		"ecstatic",
		"elastic",
		"elated",
		"elegant",
		"eloquent",
		"epic",
		"fervent",
		"festive",
		"flamboyant",
		"focused",
		"friendly",
		"frosty",
		"gallant",
		"gifted",
		"goofy",
		"gracious",
		"happy",
		"hardcore",
		"heuristic",
		"hopeful",
		"hungry",
		"infallible",
		"inspiring",
		"jolly",
		"jovial",
		"keen",
		"kind",
		"laughing",
		"loving",
		"lucid",
		"mystifying",
		"modest",
		"musing",
		"naughty",
		"nervous",
		"nifty",
		"nostalgic",
		"objective",
		"optimistic",
		"peaceful",
		"pedantic",
		"pensive",
		"practical",
		"priceless",
		"quirky",
		"quizzical",
		"relaxed",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"stoic",
		"stupefied",
		"suspicious",
		"tender",
		"thirsty",
		"trusting",
		"unruffled",
		"upbeat",
		"vibrant",
		"vigilant",
		"vigorous",
		"wizardly",
		"wonderful",
		"xenodochial",
		"youthful",
		"zealous",
		"zen",
	}

	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/", func(c *gin.Context) {
		imageUrlComponents := []string{"http://wearebuildingthefuture.com/_answer?query=", left[rand.Intn(93)]}
		imageUrl := strings.Join(imageUrlComponents, "")

    		response, err := http.Get(imageUrl)
    		if err != nil {
        		log.Fatal(err)
    		}
    		defer response.Body.Close()
 
    		responseData, err := ioutil.ReadAll(response.Body)
    		if err != nil {
        		log.Fatal(err)
    		}
		
		responseString := string(responseData)
		fmt.Printf(responseString)
 
		futureResult := new(FUTUREResult)
		json.Unmarshal([]byte(responseData), futureResult)
		fmt.Println(futureResult)
		listOfImages := futureResult.Result.Images
		numberOfImages := len(listOfImages) 
		futureImageUrl := listOfImages[rand.Intn(numberOfImages)]

		img, _ := os.Create("image.jpg")
    		defer img.Close()

    		resp, _ := http.Get(futureImageUrl)
    		defer resp.Body.Close()

    		b, _ := io.Copy(img, resp.Body)
    		fmt.Println("File size: ", b)	

		//extraHeaders := map[string]string{
			//"Content-Disposition": `attachment; filename="gopher.png"`,
		//}

		//c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)

		c.JSON(http.StatusOK, gin.H{"response": c.PostForm("secretMessage"), "rand": futureImageUrl, "fileSize": b})
	})

	router.Static("/.well-known/pki-validation", "2FF0A0CC6029BB26BB49BEDD95CE23F8.txt")

	router.Run(":" + port)
}
