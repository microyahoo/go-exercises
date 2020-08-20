package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("vim-go")
	group := "osd"
	rule1 := &AlertRule{
		Level: AlertLevelInfo,
		Group: group,
	}
	rule2 := &AlertRule{
		Level: AlertLevelCritical,
		Group: group,
	}
	rule3 := &AlertRule{
		Level: AlertLevelWarning,
		Group: group,
	}
	rule4 := &AlertRule{
		Level: AlertLevelError,
		Group: group,
	}
	var rules []*AlertRule
	rules = append(rules, rule1, rule2, rule3, rule4)
	fmt.Println(rules)
	fmt.Printf("********insertionSort************\n")
	insertionSort(rules)
	fmt.Println(rules)
	fmt.Printf("********Sort************\n")
	sort.Slice(rules, func(i, j int) bool {
		return AlertLevelMap[rules[i].Level] > AlertLevelMap[rules[j].Level]
	})
	fmt.Println(rules)

	rulesMap := make(map[string]map[string]*AlertRule)
	for _, rule := range rules {
		if _, ok := rulesMap[rule.Group]; !ok {
			rulesMap[rule.Group] = make(map[string]*AlertRule)
		}
		rulesMap[rule.Group][rule.Level] = rule
	}
	fmt.Println(rulesMap)
	var infoRule *AlertRule
	for group, ruleMap := range rulesMap {
		for level, rule := range ruleMap {
			if rule.TriggerMode == "" {
				rule.TriggerValue = level
			}
			if rule.Level == AlertLevelInfo {
				infoRule = rule
			}
			ruleMap[level] = rule
		}
		rulesMap[group] = ruleMap
	}
	fmt.Println(rulesMap)
	infoRule.TriggerValue = "xxx"
	fmt.Println(rulesMap)

	var testMap = make(map[string]bool)
	fmt.Println(testMap["x"])

	fmt.Println("===============================================================\n")
	rule1 = &AlertRule{
		ResourceType: "osd",
		Type:         "status",
		Level:        AlertLevelInfo,
		Group:        group,
	}
	rule2 = &AlertRule{
		ResourceType: "osd",
		Type:         "status",
		Level:        AlertLevelCritical,
		Group:        group,
	}
	rule3 = &AlertRule{
		ResourceType: "osd",
		Type:         "status",
		Level:        AlertLevelWarning,
		Group:        group,
	}
	rule4 = &AlertRule{
		ResourceType: "osd",
		Type:         "capacity",
		Level:        AlertLevelError,
		Group:        group,
	}
	rule5 := &AlertRule{
		ResourceType: "osd",
		Type:         "capacity",
		Level:        AlertLevelWarning,
		Group:        group,
	}
	fmt.Println("len(rules)= ", len(rules))
	fmt.Println("cap(rules)= ", cap(rules))
	rules = rules[:0]
	fmt.Println("len(rules)= ", len(rules))
	fmt.Println("cap(rules)= ", cap(rules))
	rules = append(rules, rule1, rule2, rule3, rule4, rule5)
	fmt.Println("len(rules)= ", len(rules))
	fmt.Println("cap(rules)= ", cap(rules))

	orderedAlertRuleMap := make(map[string][]*AlertRule)
	allRules := rules[:]
	for _, rule := range allRules {
		key := alertRuleResourceTypeKey(rule)
		_, ok := orderedAlertRuleMap[key]
		if !ok {
			orderedAlertRuleMap[key] = make([]*AlertRule, 0)
		}
		orderedAlertRuleMap[key] = append(orderedAlertRuleMap[key], rule)
	}
	for key, rules := range orderedAlertRuleMap {
		insertionSort(rules)
		orderedAlertRuleMap[key] = rules
	}
	for key, rules := range orderedAlertRuleMap {
		fmt.Printf("key = %s, alerts = %#v\n", key, rules)
		for _, rule := range rules {
			fmt.Println("\t", rule)
		}
	}
}

func alertRuleResourceTypeKey(rule *AlertRule) string { return rule.ResourceType + "/" + rule.Type }

// AlertRule defines the alert rule
type AlertRule struct {
	ID            int64
	Type          string
	Group         string
	ResourceType  string
	TriggerValue  string
	TriggerMode   string
	TriggerPeriod uint64
	Level         string
	Enabled       bool
	AlarmID       string
}

func (rule *AlertRule) String() string {
	return fmt.Sprintf("{Level=%s, TriggerValue=%s}", rule.Level, rule.TriggerValue)
}

func insertionSort(sl []*AlertRule) {
	a, b := 0, len(sl)
	for i := a + 1; i < b; i++ {
		for j := i; j > a && AlertLevelMap[sl[j].Level] > AlertLevelMap[sl[j-1].Level]; j-- {
			sl[j], sl[j-1] = sl[j-1], sl[j]
		}
	}
}

// Define alert level
const (
	AlertLevelInfo     = "info"
	AlertLevelWarning  = "warning"
	AlertLevelError    = "error"
	AlertLevelCritical = "critical"
)

// AlertLevelMap is used to compare the level of alerts
var AlertLevelMap = map[string]int{
	AlertLevelInfo:     0,
	AlertLevelWarning:  1,
	AlertLevelError:    2,
	AlertLevelCritical: 3,
}
