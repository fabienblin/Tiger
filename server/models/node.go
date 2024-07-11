package models

import (
	"fmt"
	"math/rand"
	"prototiger"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var NodeGenerator NodeFactory

/**
 * The prototigre Node definition doesn't allow to serialize a data tree that might have cyclic references
 * Therefor the neighbours must be a list of id in protobuf but can be a list of pointers in golang
 */
type Node struct {
	prototiger.Node
	Neighbours map[string]*Node
}

func (n *Node) MarshalableTree() []*prototiger.Node {
	if n == nil {
		return nil
	}

	marshalable := []*prototiger.Node{}

	for _, node := range n.Neighbours {
		copyNode := &prototiger.Node{
			Id:           node.Id,
			Type:         node.Type,
			Status:       node.Status,
			Measurements: node.Measurements,
			Commands:     node.Commands,
			X:            node.X,
			Y:            node.Y,
		}

		// add neighbours to copyNode
		if node.Type == prototiger.Node_NODE_GROUP {
			copyNode.Neighbours = node.MarshalableTree()
		}

		for _, copyNeighbour := range node.Neighbours {
			copyNode.Neighbours = append(copyNode.Neighbours, &prototiger.Node{Id: copyNeighbour.Id})
		}
		marshalable = append(marshalable, copyNode)
	}

	return marshalable
}

func (n *Node) UpdateRandomTree() {
	if n.Type == prototiger.Node_NODE_GROUP {
		for i := range n.Neighbours {
			n.Neighbours[i].UpdateRandomTree()
		}
	}
	n.Status = NodeGenerator.GenerateRandomStatus()
}

type NodeFactory struct {
	NbNodes uint
}

func (n *NodeFactory) GenerateRandomNode() *Node {
	node := new(Node)
	node.Id = uuid.New().String()
	enumTypeLen := len(prototiger.Node_Type_name)
	node.Type = prototiger.Node_Type(rand.Intn(enumTypeLen-1) + 1) // random type must not be NODE_GROUP(0)
	node.Neighbours = map[string]*Node{}
	node.Status = NodeGenerator.GenerateRandomStatus()
	node.X = rand.Float32()
	node.Y = rand.Float32()
	return node
}

func (n *NodeFactory) GenerateRandomTree() (*Node, error) {
	if (n.NbNodes == 0) {
		return nil, fmt.Errorf("unable to generate a tree of %s nodes", n.NbNodes)
	}

	root := NodeGenerator.GenerateRandomNode()
	root.Type = prototiger.Node_NODE_GROUP

	var nodeList []*Node

	// generate tree nodes
	for i := 0; i < int(n.NbNodes); i++ {
		newNode := NodeGenerator.GenerateRandomNode()
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
				linkedNeighbour = nodeList[rand.Intn(int(n.NbNodes))]
			}
			nodeList[i].Neighbours[linkedNeighbour.Id] = linkedNeighbour

			// reverse link
			linkedNeighbour.Neighbours[nodeList[i].Id] = nodeList[i]
		}
	}

	root.Neighbours[nodeList[0].Id] = nodeList[0]
	return root, nil
}

func (n *NodeFactory) GenerateRandomStatus() *prototiger.Status {
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
