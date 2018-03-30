package main

import (
	"sync"
	"errors"
	"fmt"
)

type SelectStep struct {
	Num int `json:"num"`
	Type string `json:"type"`
}

type RuleRepresent struct{
	Steps []SelectStep `json:"steps"`
}

type Rule interface {
	Init(rep RuleRepresent)
	GetSteps() []SelectStep
	DoRule(cluster Cluster) []*Node
}

type SimpleRule struct{
	rep RuleRepresent
}

func (rule *SimpleRule)Init(rep RuleRepresent)  {
	rule.rep = rep
}

func (rule SimpleRule)GetSteps() []SelectStep {
	return rule.rep.Steps
}

func (rule SimpleRule)DoRule(cluster Cluster) []*Node {
	var matchedItems, currentItems []*Node

	steps := rule.GetSteps()
	matchedItems = cluster.GetRoots()
	for _, step := range steps {
		currentItems = matchedItems
		matchedItems = matchedItems[:0]
		for _, item := range currentItems {
			availableItems := collectAvailableItems(item, step.Type)
			matchedItems = append(matchedItems, selectItem(availableItems, step.Num)...)
		}
	}

	return matchedItems
}

type RuleMng struct {
	ruleMap map[string]Rule
	lock sync.RWMutex
}

func (ruleMng *RuleMng) AddRule(ruleName string, rule Rule) error {
	ruleMng.lock.Lock()
	_, ok := ruleMng.ruleMap[ruleName]
	if true == ok {
		return errors.New("Rule already exist" + ruleName)
	}
	ruleMng.ruleMap[ruleName] = rule
	ruleMng.lock.Unlock()
	return nil
}

func (ruleMng *RuleMng) GetRule(ruleName string) (Rule, error)  {
	ruleMng.lock.RLock()
	defer ruleMng.lock.RUnlock()
	rule, ok := ruleMng.ruleMap[ruleName]
	if false == ok {
		msg := "Can't get rule" + ruleName
		return rule, errors.New(msg)
	}
	return rule,nil
}

func (ruleMng *RuleMng) DoRule(ruleName string, cluster Cluster) []*Node {
	var matchedItems []*Node
	rule, err := ruleMng.GetRule(ruleName)
	if nil != err {
		fmt.Println("[DoRule] fail, err =", err)
		return matchedItems
	}

	matchedItems = rule.DoRule(cluster)
	return matchedItems
}

func collectAvailableItems(node *Node, targetType string) []*Node {
	result := []*Node{}
	if targetType == node.Type {
		result = append(result, node)
	}else{
		for _,child := range node.Children {
			result = append(result, collectAvailableItems(child, targetType)...)
		}
	}

	return result
}

func selectItem(set []*Node, num int) []*Node {
	result := []*Node{}
	for i, node := range set {
		if i < num {
			result = append(result, node)
		}else{
			break
		}
	}
	return result
}