package fac

import (
	"regexp/syntax"
	"testing"

	"github.com/awslabs/ssosync/internal/aws"
	"github.com/awslabs/ssosync/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMatchAWSGroups(t *testing.T) {
	tests := []struct {
		name           string
		awsGroupMatch  string
		inputGroups    []*aws.Group
		expectedGroups []*aws.Group
		expectedErr    error
	}{
		{
			name:          "correctly matches all groups with default group match",
			awsGroupMatch: config.DefaultAwsGroupMatch,
			inputGroups: []*aws.Group{
				{DisplayName: "aws-group-A"},
				{DisplayName: "aws-group-B"},
				{DisplayName: "aws-group-C"},
				{DisplayName: "aws-meow-meow"},
			},
			expectedGroups: []*aws.Group{
				{DisplayName: "aws-group-A"},
				{DisplayName: "aws-group-B"},
				{DisplayName: "aws-group-C"},
				{DisplayName: "aws-meow-meow"},
			},
			expectedErr: nil,
		},
		{
			name:          "correctly matches selected groups",
			awsGroupMatch: "aws-group-.*",
			inputGroups: []*aws.Group{
				{DisplayName: "aws-group-A"},
				{DisplayName: "aws-group-B"},
				{DisplayName: "aws-group-C"},
				{DisplayName: "aws-meow-meow"},
			},
			expectedGroups: []*aws.Group{
				{DisplayName: "aws-group-A"},
				{DisplayName: "aws-group-B"},
				{DisplayName: "aws-group-C"},
			},
			expectedErr: nil,
		},
		{
			name:          "returns an error when input groups empty",
			awsGroupMatch: "aws-group-.*",
			inputGroups:   []*aws.Group{},
			expectedErr:   NoAWSGroupsErr,
		},
		{
			name:          "returns an error when input groups nil",
			awsGroupMatch: "aws-group-*",
			inputGroups:   []*aws.Group{},
			expectedErr:   NoAWSGroupsErr,
		},
		{
			name:          "returns an error when regex invalid",
			awsGroupMatch: "[^0-1",
			inputGroups:   []*aws.Group{{DisplayName: "aws-group-A"}},
			expectedErr: BadRegexError{
				Message: "can't compile regex [^0-1",
				Err:     &syntax.Error{Code: syntax.ErrMissingBracket, Expr: "[^0-1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups, err := MatchAWSGroups(tt.inputGroups, tt.awsGroupMatch)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedGroups, groups)
		})
	}
}
