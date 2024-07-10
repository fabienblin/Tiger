package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"prototiger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	var tree *prototiger.Node = generateRandomTree(50)

	for {
		// mt, message, err := c.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }

		// updateRandomTree(tree)
		for _, node := range tree.Neighbours {
			node.Status = generateRandomStatus()
			message, err := json.Marshal(marshalableTree(tree))
			if err != nil {
				log.Println(err)
			}
		
			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}
	
			time.Sleep(time.Millisecond*100)
		}
	}
}

func generateRandomNode() (*prototiger.Node) {
	node := new(prototiger.Node)
	node.Id = uuid.New().String()
	enumTypeLen := len(prototiger.Node_Type_name)
	node.Type = prototiger.Node_Type(rand.Intn(enumTypeLen-1) + 1) // random type must not be NODE_GROUP(0)
	node.Neighbours = map[string]*prototiger.Node{}
	node.Status = generateRandomStatus()
	node.X = rand.Float32()
	node.Y = rand.Float32()
	return node
}

func generateRandomTree(n int) (*prototiger.Node) {
	root := generateRandomNode()
	root.Type = prototiger.Node_NODE_GROUP

	var nodeList []*prototiger.Node

	// generate tree nodes
	for i := 0; i < n; i++ {
		newNode := generateRandomNode()
		nodeList = append(nodeList, newNode)
		root.Neighbours[newNode.Id] = newNode
	}

	// link tree nodes
	for i := range nodeList {
		nbLinks := rand.Intn(3) //rand.Intn(int(math.Sqrt(float64(n)))) + 1
		// link to neighbours
		for j := 0; j < nbLinks; j++ {
			// link to any other random node
			linkedNeighbour := nodeList[i]
			for linkedNeighbour == nodeList[i] {
				linkedNeighbour = nodeList[rand.Intn(n)]
			}
			nodeList[i].Neighbours[linkedNeighbour.Id] = linkedNeighbour

			// reverse link
			linkedNeighbour.Neighbours[nodeList[i].Id] = nodeList[i]
		}
	}

	root.Neighbours[nodeList[0].Id] = nodeList[0]
	return root
}

func generateRandomStatus() (*prototiger.Status) {
	status := new(prototiger.Status)
	enumTypeLen := len(prototiger.Status_Type_name)
	status.Type = prototiger.Status_Type(rand.Intn(enumTypeLen))
	if rand.Int()%2 == 0 {
		status.IsActive = true
	}
	if rand.Int()%2 == 0 {
		status.IsAlarmed = true
	}
	status.Timestamp = timestamppb.Now()

	return status
}

func marshalableTree(tree *prototiger.Node) (map[string]*prototiger.Node) {
	if tree == nil {
		return nil
	}

	marshalable := map[string]*prototiger.Node{}

	for _, node := range tree.Neighbours {
		copyNode := *node
		if copyNode.Type == prototiger.Node_NODE_GROUP {
			subTree := marshalableTree(&copyNode)
			marshalable[copyNode.Id].Neighbours = subTree
		} else {
			for i := range copyNode.Neighbours {
				copyNode.Neighbours[i] = &prototiger.Node{Id: i}
			}
			marshalable[copyNode.Id] = &copyNode
		}
	}

	return marshalable
}

func updateRandomTree(node *prototiger.Node) {
	if node.Type == prototiger.Node_NODE_GROUP {
		for i := range node.Neighbours {
			updateRandomTree(node.Neighbours[i])
		}
	}
	node.Status = generateRandomStatus()
}
