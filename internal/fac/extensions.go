// Package fac contains custom functions for additional operations on data.
package fac

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/awslabs/ssosync/internal/aws"
)

// ErrNoAWSGroups indicates no AWS groups were received.
var ErrNoAWSGroups = errors.New("received no AWS groups")

// ErrorBadRegex represents a regex compilation error.
type ErrorBadRegex struct {
	Message string
	Err     error
}

func (e ErrorBadRegex) Error() string {
	return e.Message
}

// MatchAWSGroups will filter out the AWS groups that don't match the regex.
// Returns an error on failure, a list of AWS groups that match on success.
func MatchAWSGroups(awsGroups []*aws.Group, matchRegex string) ([]*aws.Group, error) {
	if awsGroups == nil || len(awsGroups) == 0 {
		return nil, ErrNoAWSGroups
	}

	awsGroupRegex, err := regexp.Compile(matchRegex)
	if err != nil {
		return nil, ErrorBadRegex{Message: fmt.Sprintf("can't compile regex %s", matchRegex), Err: err}
	}

	matchedGroups := make([]*aws.Group, 0)
	for _, group := range awsGroups {
		if awsGroupRegex.FindStringIndex(group.DisplayName) != nil {
			matchedGroups = append(matchedGroups, group)
		}
	}

	return matchedGroups, nil
}
