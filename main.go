package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"

  jwtmiddleware "github.com/auth0/go-jwt-middleware"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/contrib/static"
  "github.com/gin-gonic/gin"
)

type Response struct {
  Message string `json:"message"`
}

type Entry struct {
  ID      int     `json:"id" binding:"required"`
  Likes   int     `json:"likes"`
  Entry   string  `json:"Entry" binding:"required"`
}

type Jwks struct {
  Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
  Kty string   `json:"kty"`
  Kid string   `json:"kid"`
  Use string   `json:"use"`
  N   string   `json:"n"`
  E   string   `json:"e"`
  X5c []string `json:"x5c"`
}

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

var entries = []Entry{
  Entry{1, 0, "Lorem Ipsum is simply dummy text of the printing and typesetting industry."},
  Entry{2, 0, "It is a long established fact that a reader will be distracted by the readable content."},
  Entry{3, 0, "There is no one who loves pain, who seeks after it and wants to have it, simply because it is pain.."},
  Entry{4, 0, "Can you help translate this site into a foreign language?"},
  Entry{5, 0, "Please email us with details if you can help.."},
  Entry{6, 0, "here are now a set of mock banners available here in three colours and in a range of banner sizes."},
  Entry{7, 0, "Thank you for your support."},
}

var jwtMiddleWare *jwtmiddleware.JWTMiddleware

// authMiddleware intercepts the requests, and check for a valid jwt token
func authMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    // Get the client secret key
    err := jwtMiddleWare.CheckJWT(c.Writer, c.Request)
    if err != nil {
      // Token not found
      fmt.Println(err)
      c.Abort()
      c.Writer.WriteHeader(http.StatusUnauthorized)
      c.Writer.Write([]byte("Unauthorized"))
      return
    }
  }
}

func getPemCert(token *jwt.Token) (string, error) {
  cert := ""
  resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")
  if err != nil {
    return cert, err
  }
  defer resp.Body.Close()

  var jwks = Jwks{}
  err = json.NewDecoder(resp.Body).Decode(&jwks)

  if err != nil {
    return cert, err
  }

  x5c := jwks.Keys[0].X5c
  for k, v := range x5c {
    if token.Header["kid"] == jwks.Keys[k].Kid {
      cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
    }
  }

  if cert == "" {
    return cert, errors.New("unable to find appropriate key.")
  }

  return cert, nil
}

func main() {
  jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    aud := os.Getenv("AUTH0_API_AUDIENCE")
    checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
    if !checkAudience {
      return token, errors.New("Invalid audience.")
    }
    // verify iss claim
    iss := os.Getenv("AUTH0_DOMAIN")
    checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
    if !checkIss {
      return token, errors.New("Invalid issuer.")
    }

    cert, err := getPemCert(token)
    if err != nil {
      log.Fatalf("could not get cert: %+v", err)
    }

    result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
    return result, nil
  },
  SigningMethod: jwt.SigningMethodRS256,
})

// register our actual jwtMiddleware
jwtMiddleWare = jwtMiddleware
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

  api.GET("/entries", authMiddleware(), ShowEntriesHandler)
  api.POST("/entries/like/:EntryID", authMiddleware(), LikeEntry)

  // Start and run the server
  router.Run(":3000")
}
