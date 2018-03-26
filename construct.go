package main

import (
	"errors"
	"fmt"
	"sync"
)

type Cluster interface {
	AddNode(name string, id int64, nodeType string)(int64, string)
	MoveNode(nodeId int64, targetNodeId int64) error
	RenameNode(nodeName string, newName string)
}

type Node struct {
	Name string `json:"name"`
	Id int64 `json:"id"`
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Weight float64 `json:"weight"`
	Children map[int64]*Node `json:"children,omitempty"`
	c_lock sync.RWMutex
}

type Forest struct{
	Roots []*Node `json:"treeList"`
}

var (
	osdId int64 = 0
	bucketId int64 = -1
	)

func getNewOsdId() int64 {
	newId := osdId
	osdId += 1
	return newId
}

func getNewBucketId() int64 {
	newId := bucketId
	bucketId -= 1
	return newId
}

func (forest *Forest) AddNode(name, nodeType string) (int64, string) {
	node := &Node{
		Name: name,
		Type: nodeType,
	}

	if nodeType == "osd" {
		node.Id = getNewOsdId()
		node.Name = fmt.Sprintf("osd.%d", node.Id)
	}else{
		node.Id = getNewBucketId()
		node.Children = make(map[int64]*Node)
	}

	forest.Roots = append(forest.Roots, node)
	return node.Id, node.Name
}

func searchFatherNode(node, fatherNode *Node, nodeId int64) (*Node, *Node, bool) {
	if node.Id == nodeId {
		return fatherNode, node, true
	}

	for _, childNode := range node.Children {
		father, me, findFather := searchFatherNode(childNode, node, nodeId)
		if true == findFather {
			return father, me, true
		}
	}
	return nil, node, false
}

func searchNode(node *Node, targetId int64) (*Node, bool) {
	if node.Id == targetId {
		return node, true
	}

	for _, childNode := range node.Children {
		find, _ := searchNode(childNode, targetId)
		if nil != find {
			return find, true
		}
	}
	return nil,false
}

func (forest *Forest) MoveNode(nodeId int64, targetNodeId int64) error  {
	var findFather, findTarget bool
 var fatherNode, sourceNode,  targetNode *Node
 for _, root := range forest.Roots {
 	if false != findFather && false != findTarget {
 		break
	}

 	if false == findFather {
		fatherNode, sourceNode, findFather = searchFatherNode(root, nil, nodeId)
	}

	if false == findTarget {
		targetNode, findTarget = searchNode(root,targetNodeId)
	}
 }

 if false == findTarget {
 	msg := fmt.Sprintf("Can't find targetNode(id=%d).", targetNodeId)
 	return errors.New(msg)
 }
	if false == findFather {
		msg := fmt.Sprintf("Can't find fatherNode of node(id=%d).", nodeId)
		return errors.New(msg)
	}

 if nil == sourceNode {
	 msg := fmt.Sprintf("Can't find Node(id=%d).", nodeId)
	 return errors.New(msg)
 }

 targetNode.Children[nodeId] = sourceNode
 fmt.Println("mv node success", *targetNode)

 if nil == fatherNode {
	 var index int
	 for i, root := range forest.Roots {
		 if root.Id == nodeId {
			 index = i
			 break
		 }
	 }
	 forest.Roots = append(forest.Roots[0:index], forest.Roots[index+1:]...)
	 return  nil
 }

	fatherNode.c_lock.Lock()
	delete(fatherNode.Children, nodeId)
	fatherNode.c_lock.Unlock()
	return nil
}

func (forest *Forest) RenameNode(nodeId int64, newName string)  {
	var node *Node
	for _, root := range forest.Roots {
		if node == nil {
			node, _ = searchNode(root,nodeId)
		}else{
			break
		}
	}

	if nil == node {
		fmt.Println("Can't find node id=", nodeId)
		return
	}

	node.Name = newName
}
