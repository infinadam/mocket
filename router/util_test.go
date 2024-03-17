package router

import (
	"regexp"
	"testing"
)

// should match a string with no groupings
func TestRouterUtilMatch(t *testing.T) {
	re := regexp.MustCompile(`test`)
	matched, groups := match(re, "this is a test")

	if !matched {
		t.Error("expected a match, got none")
	}

	if len(groups) != 0 {
		t.Error("expected no groups")
	}
}

// should match with a grouping
func TestRouterUtilMatchGroup(t *testing.T) {
	re := regexp.MustCompile(`(?P<test>\d+)`)
	matched, groups := match(re, "test 123 string")

	if !matched {
		t.Error("expected a match, got none")
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}

	if groups["test"] != "123" {
		t.Errorf("expected \"test\" to be \"123\", was %q", groups["test"])
	}
}

// should not match
func TestRouterUtilMatchFails(t *testing.T) {
	re := regexp.MustCompile(`nothing`)
	matched, groups := match(re, "test string")

	if matched {
		t.Error("expected not to match")
	}

	if len(groups) != 0 {
		t.Error("expected no groups")
	}
}
