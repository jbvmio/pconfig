package pconfig

import (
	"fmt"
	"regexp"

	"github.com/jbvmio/kafka"
	"go.uber.org/zap"
)

// DeleteCG deletes a consumer group.
func DeleteCG(client *kafka.KClient, group string, logger *zap.Logger) {
	logger.Info("delete consumer group called", zap.String(`group`, group))
	var found bool
	groups, errs := client.ListGroups()
	if len(errs) > 1 {
		for _, e := range errs {
			fmt.Println(e)
		}
		logger.Warn("error fetching existing group metadata")
		return
	}
	for _, g := range groups {
		if g == group {
			found = true
			break
		}
	}
	if found {
		logger.Warn("removing existing group", zap.String(`group`, group))
		err := client.RemoveGroup(group)
		if err != nil {
			logger.Warn("Error deleting existing group", zap.Error(err))
		}
	} else {
		logger.Info("consumer group not found", zap.String(`group`, group))
	}
}

// TopicsExist returns true if the given topic exists, otherwise false.
func TopicsExist(client *kafka.KClient, logger *zap.Logger, topics ...string) bool {
	logger.Info("validating topics", zap.Int(`number of topics`, len(topics)))
	var matched int
	regex := makeRegex(topics...)
	tMeta, err := client.GetTopicMeta()
	if err != nil {
		logger.Warn("error fetching topic metadata", zap.Error(err))
		return false
	}
	dupe := make(map[string]bool)
checkLoop:
	for _, t := range tMeta {
		if !dupe[t.Topic] {
			dupe[t.Topic] = true
			if regex.MatchString(t.Topic) {
				logger.Debug("topic found", zap.String(`topic`, t.Topic))
				matched++
			}
			if matched == len(topics) {
				break checkLoop
			}
		}
	}
	if matched == len(topics) {
		logger.Debug("all topics validated", zap.Strings(`topics`, topics))
		return true
	}
	logger.Warn("could not validate all topics", zap.Strings(`topics`, topics))
	return false
}

func makeRegex(terms ...string) *regexp.Regexp {
	var regStr string
	switch len(terms) {
	case 0:
		regStr = ""
	case 1:
		regStr = `^(` + terms[0] + `)$`
	default:
		regStr = `^(` + terms[0]
		for _, t := range terms[1:] {
			regStr += `|` + t
		}
		regStr += `)$`
	}
	return regexp.MustCompile(regStr)
}
