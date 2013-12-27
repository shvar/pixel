package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
    "code.google.com/p/go-uuid/uuid"
)

type Pixel map[string]string

func create_pixel() Pixel {
	var pixel = Pixel{"_id": uuid.New(), "blob": "", "timestamp": time.Now().Format(time.UnixDate)}
	c := session.DB("pixels").C("potential")
	err := c.Insert(pixel)
    if err != nil {
        panic(err)
    }
	return pixel
}

func serve_pixel(r *http.Request, pixel_id string) {
	c := session.DB("pixels").C("potential")
	result := Pixel{}
	err = c.Find(bson.M{"_id": pixel_id}).One(&result)
	if err != nil {
        log.Printf("Can't serve pixel: %s", pixel_id)
    } else {
    	log.Printf("Serving pixel: %s", result["_id"])
    	c := session.DB("pixels").C("delivered")
    	err = c.Find(bson.M{"_id": pixel_id}).One(&result)
    	if err != nil {
          log.Printf("Mark pixel as delivered: %s", pixel_id)
          result["timestamp"] = time.Now().Format(time.UnixDate)
          c.Insert(result)
        } else {
          log.Printf("Sorry, pixel already served: %s", pixel_id)
        }
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[5:7] == "s/" {
		serve_pixel(r, r.URL.Path[7:])
	} else if r.URL.Path[5:9] == "reg/" {
		var new_pixel = create_pixel()
		fmt.Fprintf(w, "{'id':%s}", new_pixel["_id"])
		log.Printf("Issued pixel with UUID %s", new_pixel["_id"])
	} else {
    	log.Printf("Unknown path: %s", r.URL.Path[5:9])
	}
}

var session, err = mgo.Dial("localhost")

func main() {
    defer session.Close()
    http.HandleFunc("/pix/", handler)
    http.ListenAndServe("localhost:5700", nil)
}
