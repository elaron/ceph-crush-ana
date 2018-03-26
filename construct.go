package main

import (
	"errors"
	"fmt"
	"sync"
)

/*
[rename]
default:EBS_SZAZ1
[add]
room:XL_IDC6 XL_IDC7
rack:vrack1 vrack2
[move]
root=EBS_SZAZ1: XL_IDC6 XL_IDC7
room=XL_IDC6: vrack1 vrack2
room:=XL_IDC7: vrack3
rack=vrack1: HOSTNAME1 HOSTNAME2
*/

type RenameRep struct{
	OrigName string `json:"orig_name"`
	TargetName string `json:"target_name"`
}

type AddRep struct{
	Type string `json:"type"`
	NewItems []string `json:"new_items"`
}

type MoveRep struct {
	TargetType string `json:"target_type"`
	TargetName string `json:"target_name"`
	SourceNames []string `json:"source_names"`
}

type ConstructRepresent struct {
	AddOps []AddRep `json:"add_ops"`
	MoveOps []MoveRep `json:"move_ops"`
}

type Cluster interface {
	AddNode(name string, id int64, nodeType string)(int64, string)
	MoveNode(nodeId int64, targetNodeId int64) error
	RenameNode(nodeName string, newName string)
	New(rep *ConstructRepresent)
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

func (forest *Forest)New(rep *ConstructRepresent){
	for _, addRep := range rep.AddOps {
		for _, itemName := range addRep.NewItems {
			id, name := forest.AddNode(itemName, addRep.Type)
			fmt.Println("[New]add node:",id,name)
		}

	}

	for _, moveRep := range rep.MoveOps {
		for _, itemName := range moveRep.SourceNames {
			forest.MoveNode(itemName,moveRep.TargetName)
			fmt.Println("[New]move node:",itemName,moveRep.TargetName)
		}
	}
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

func searchFatherNode(node, fatherNode *Node, nodeName string) (*Node, *Node, bool) {
	if node.Name == nodeName {
		return fatherNode, node, true
	}

	for _, childNode := range node.Children {
		father, me, findFather := searchFatherNode(childNode, node, nodeName)
		if true == findFather {
			return father, me, true
		}
	}
	return nil, node, false
}

func searchNode(node *Node, targetName string) (*Node, bool) {
	if node.Name == targetName {
		return node, true
	}

	for _, childNode := range node.Children {
		find, _ := searchNode(childNode, targetName)
		if nil != find {
			return find, true
		}
	}
	return nil,false
}

func (forest *Forest) MoveNode(nodeName string, targetName string) error  {
	var findFather, findTarget bool
 var fatherNode, sourceNode,  targetNode *Node
 for _, root := range forest.Roots {
 	if false != findFather && false != findTarget {
 		break
	}

 	if false == findFather {
		fatherNode, sourceNode, findFather = searchFatherNode(root, nil, nodeName)
	}

	if false == findTarget {
		targetNode, findTarget = searchNode(root,targetName)
	}
 }

 if false == findTarget {
 	msg := fmt.Sprintf("Can't find targetNode(name=%s).", targetName)
 	return errors.New(msg)
 }
	if false == findFather {
		msg := fmt.Sprintf("Can't find fatherNode of node(id=%d).", nodeName)
		return errors.New(msg)
	}

 if nil == sourceNode {
	 msg := fmt.Sprintf("Can't find Node(id=%d).", nodeName)
	 return errors.New(msg)
 }

 targetNode.Children[sourceNode.Id] = sourceNode
 fmt.Println("mv node success", *targetNode)

 if nil == fatherNode {
	 var index int
	 for i, root := range forest.Roots {
		 if root.Name == nodeName {
			 index = i
			 break
		 }
	 }
	 forest.Roots = append(forest.Roots[0:index], forest.Roots[index+1:]...)
	 return  nil
 }

	fatherNode.c_lock.Lock()
	delete(fatherNode.Children, sourceNode.Id)
	fatherNode.c_lock.Unlock()
	return nil
}

func (forest *Forest) RenameNode(origName string, newName string)  {
	var node *Node
	for _, root := range forest.Roots {
		if node == nil {
			node, _ = searchNode(root,origName)
		}else{
			break
		}
	}

	if nil == node {
		fmt.Println("Can't find node id=", origName)
		return
	}

	node.Name = newName
}
