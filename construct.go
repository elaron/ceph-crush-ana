package main

import (
	"errors"
	"fmt"
	"sync"
	"encoding/json"
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

type AddBucketRep struct{
	Type string `json:"type"`
	NewItems []string `json:"new_items"`
}

type AddOsdRep struct{
	TargetHost string `json:"target_host"`
	OsdNum int32 `json:"osd_num"`
}

type MoveRep struct {
	TargetType string `json:"target_type"`
	TargetName string `json:"target_name"`
	SourceNames []string `json:"source_names"`
}

type ConstructRepresent struct {
	AddOps []AddBucketRep `json:"add_ops"`
	MoveOps []MoveRep     `json:"move_ops"`
}

type Cluster interface {
	AddOsds(rep *AddOsdRep)[]string
	AddNode(name, nodeType string) (int, string)
	MoveNode(nodeName string, targetName string) error
	RenameNode(origName string, newName string)
	New(rep *ConstructRepresent)
	GetRoots() []*Node
}

type Node struct {
	Name string `json:"name"`
	Id int `json:"id"`
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Weight float64 `json:"weight"`
	Children map[int]*Node `json:"children,omitempty"`
	c_lock sync.RWMutex
}

type Forest struct{
	osdId int
	bucketId int
	Roots []*Node `json:"treeList"`
}

func (forest *Forest) getNewOsdId() int {
	newId := forest.osdId
	forest.osdId += 1
	return newId
}

func (forest *Forest) getNewBucketId() int {
	if forest.bucketId == 0{
		forest.bucketId = -1
	}
	newId := forest.bucketId
	forest.bucketId -= 1
	return newId
}

func (forest *Forest) New(rep *ConstructRepresent){
	for _, addRep := range rep.AddOps {
		for _, itemName := range addRep.NewItems {
			id, name := forest.AddNode(itemName, addRep.Type)
			fmt.Println("[New]add node:",id,name)
		}

	}
	for _, moveRep := range rep.MoveOps {
		forest.MoveOp(&moveRep)
	}

}

func (forest *Forest) MoveOp(moveRep * MoveRep){
	for _, itemName := range moveRep.SourceNames {
		forest.MoveNode(itemName,moveRep.TargetName)
		fmt.Println("[New]move node:",itemName,moveRep.TargetName)
	}
}

func (forest *Forest) AddOsds(rep *AddOsdRep)[]string{
	var i int32
	result := make([]string, rep.OsdNum)
	for i = 0; i < rep.OsdNum; i++ {
		_, name := forest.AddNode("", "osd")
		result[i] = name
	}

	forest.MoveOp(&MoveRep{"host",rep.TargetHost, result})
	return result
}

func (forest *Forest) AddNode(name, nodeType string) (int, string) {
	node := &Node{
		Name: name,
		Type: nodeType,
	}

	if nodeType == "osd" {
		node.Id = forest.getNewOsdId()
		node.Name = fmt.Sprintf("osd.%d", node.Id)
	}else{
		node.Id = forest.getNewBucketId()
		node.Children = make(map[int]*Node)
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
 fmt.Println("mv node success", targetNode)

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

func (forest *Forest)ToJson() string {
	b, err := json.Marshal(&forest)
	if nil != err {
		fmt.Printf("[ToJson] fail, err=", err)
		return ""
	}
	return string(b)
}

func (forest *Forest) GetRoots() []*Node {
	return forest.Roots
}