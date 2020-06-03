package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"image/jpeg"
	"lukechampine.com/jsteg"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"math/rand"
	"time"
	"strings"
)

type FUTUREResponse struct {
	Images []string
}

type FUTUREResult struct {
	Result FUTUREResponse
}

func get_random_image(client http.Client) (string) {

	rand.Seed(time.Now().UnixNano())

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

	imageUrlComponents := []string{"http://wearebuildingthefuture.com/_midnightcypher?query=", left[rand.Intn(93)]}
	imageUrl := strings.Join(imageUrlComponents, "")

	response, error1 := client.Get(imageUrl)
	if error1 != nil {
		log.Fatal(error1)
	}
	defer response.Body.Close()

	responseData, error2 := ioutil.ReadAll(response.Body)
	if error2 != nil {
		log.Fatal(error2)
	}

	futureResult := new(FUTUREResult)
	json.Unmarshal([]byte(responseData), futureResult)
	listOfImages := futureResult.Result.Images
	numberOfImages := len(listOfImages)
	return listOfImages[rand.Intn(numberOfImages)]
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.MaxMultipartMemory = 8 << 20

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		hidden, _ := jsteg.Reveal(file)
		c.HTML(http.StatusOK, "decyphered.tmpl.html", gin.H{
			"recoveredMessage": string(hidden),
		})
	})

	router.POST("/encode", func(c *gin.Context) {
		futureImageUrl := get_random_image(client)

		response, err := client.Get(futureImageUrl)
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		for ok := true; ok; ok = ( !(strings.HasSuffix(response.Header.Get("Content-Type"), "jpeg")) ) {
			futureImageUrl = get_random_image(client)
			response, err = client.Get(futureImageUrl)
			if err != nil || response.StatusCode != http.StatusOK {
				c.Status(http.StatusServiceUnavailable)
				return
			}
			reader = response.Body
		}

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="image.jpg"`,
		}

		out := new(bytes.Buffer)
		img, _ := jpeg.Decode(reader)
		data := []byte(c.PostForm("secretMessage"))
		jsteg.Hide(out, img, data, nil)

		c.DataFromReader(http.StatusOK, int64(out.Len()), "image/jpeg", out, extraHeaders)
	})

	router.Static("/.well-known/pki-validation", "2FF0A0CC6029BB26BB49BEDD95CE23F8.txt")

	router.Run(":" + port)
}
