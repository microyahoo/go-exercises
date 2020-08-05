package main

import (
	"fmt"
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
	insertionSort(rules)
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
}

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
		for j := i; j > a && AlertLevelMap[sl[j].Level] < AlertLevelMap[sl[j-1].Level]; j-- {
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
