package main

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/contrib/static"
  "github.com/gin-gonic/gin"
)

func ShowEntriesHandler(c *gin.Context) {
  c.Header("Content-Type", "application/json")
  c.JSON(http.StatusOK, entries)
}

func LikeEntry(c *gin.Context) {
  // confirm entry ID sent is valid
  if entryid, err := strconv.Atoi(c.Param("entryID")); err == nil {
    // find entry, and increment likes
    for i := 0; i < len(entries); i++ {
      if entries[i].ID == entryid {
        entries[i].Likes += 1
      }
    }

    // return a pointer to the updated entries list
    c.JSON(http.StatusOK, &entries)
  } else {
    // entry ID is invalid
    c.AbortWithStatus(http.StatusNotFound)
  }
}

type Entry struct {
  ID      int     `json:"id" binding:"required"`
  Likes   int     `json:"likes"`
  Entry   string  `json:"Entry" binding:"required"`
}

var entries = []Entry{
  Entry{1, 0, "Lorem Ipsum is simply dummy text of the printing and typesetting industry."},
  Entry{2, 0, "It is a long established fact that a reader will be distracted by the readable content."},
  Entry{3, 0, "There is no one who loves pain, who seeks after it and wants to have it, simply because it is pain.."},
  Entry{4, 0, "Can you help translate this site into a foreign language?"},
  Entry{5, 0, "Please email us with details if you can help.."},
  Entry{6, 0, "here are now a set of mock banners available here in three colours and in a range of banner sizes."},
  Entry{7, 0, "Thank you for your support."},
}

func main() {
  // Set the router as the default one shipped with Gin
  router := gin.Default()

  // Serve frontend static files
  router.Use(static.Serve("/", static.LocalFile("./views", true)))

  // Setup route group for the API
  api := router.Group("/")
  {
    api.GET("/", func(c *gin.Context) {
      c.JSON(http.StatusOK, gin.H {
        "message": "pong",
      })
    })
  }

  api.GET("/entries", ShowEntriesHandler)
  api.POST("/entries/like/:EntryID", LikeEntry)

  // Start and run the server
  router.Run(":3000")
}
