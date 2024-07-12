package handlers

import (
	"log"
	"main/models"
	"net/http"
	"prototiger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var imageDisplayUpgrader = websocket.Upgrader{
	CheckOrigin: func(req *http.Request) bool {
		return true
	},
}

func ImageDisplayHandler(c *gin.Context) {
	conn, err := imageDisplayUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	models.NodeGenerator.NbNodes = 50
	tree, err := models.NodeGenerator.GenerateRandomTree()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		// mt, message, err := c.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }

		// updateRandomTree(tree)
		for _, node := range tree.Neighbours {
			node.Status = models.NodeGenerator.GenerateRandomStatus()
			message, err := proto.Marshal(tree.MarshalableTree())
			if err != nil {
				log.Println(err)
			}

			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}

			time.Sleep(time.Millisecond * 100)
		}
	}
}

type routeScadaServiceServer struct {
	prototiger.UnimplementedScadaServiceServer
	// savedFeatures []*prototiger.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	// routeNotes map[string][]*prototiger.RouteNote
}

func (s *routeGuideServer) UpdateTree(rect *prototiger.Rectangle, stream prototiger.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}