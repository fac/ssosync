// Package fac contains custom functions for additional operations on data.
package fac

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/awslabs/ssosync/internal/aws"
	log "github.com/sirupsen/logrus"
)

// NoAWSGroupsErr indicates no AWS groups were received.
var NoAWSGroupsErr = errors.New("received no AWS groups")

// BadRegexError represents a regex compilation error.
type BadRegexError struct {
	Message string
	Err     error
}

func (e BadRegexError) Error() string {
	return e.Message
}

// MatchAWSGroups will filter out the AWS groups that don't match the regex.
// Returns an error on failure, a list of AWS groups that match on success.
func MatchAWSGroups(awsGroups []*aws.Group, matchRegex string) ([]*aws.Group, error) {
	if len(awsGroups) == 0 {
		return nil, NoAWSGroupsErr
	}

	awsGroupRegex, err := regexp.Compile(matchRegex)
	if err != nil {
		return nil, BadRegexError{Message: fmt.Sprintf("can't compile regex %s", matchRegex), Err: err}
	}

	matchedGroups := make([]*aws.Group, 0)
	for _, group := range awsGroups {
		if awsGroupRegex.FindStringIndex(group.DisplayName) == nil {
			log.Infof("AWS group %s will not be included in sync", group.DisplayName)
			continue
		}

		matchedGroups = append(matchedGroups, group)
	}

	return matchedGroups, nil
}
