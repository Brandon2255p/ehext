package horizgen

import (
	"regexp"
	"strings"
)

// // {{.Event}} {{.Description}}
// type {{.Event}} struct {
// 	Name string
// 	Age  int
// }

// ExtractClasses will pull comments, struct and members from a go file
func ExtractClasses(input string) []string {
	pattern := regexp.MustCompile(`(?P<comment>\/\/\s*[\w\s\d]+)*type[\s\w]+struct\s+{[\[\]\n\t\w\s:"\\.` + "`" + `]+}`)
	return pattern.FindAllString(input, -1)
}

// ExtractName will take a full function input such as:
// 		`// TestEvent comment
// 		// multiline 1
// 		// multiline 2
//		type TestEvent struct {
// 			TestAy  string
// 			TestBee int
// 		}`
//
// and produce: "TestEvent"
// NOTE: event names must end in Event
func ExtractName(intput string) string {
	pattern := regexp.MustCompile(`type[\s]+([\w]+[Event|Command|Aggregate]+)[\s]+struct`)
	match := pattern.FindStringSubmatch(intput)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

// ExtractCommandName will take a full function input such as:
// 		`// TestCommand comment
// 		// multiline 1
// 		// multiline 2
//		type TestCommand struct {
// 			TestAy  string
// 			TestBee int
// 		}`
//
// and produce: "TestEvent"
// NOTE: event names must end in Event
func ExtractCommandName(intput string) string {
	pattern := regexp.MustCompile(`[\w]+Command`)
	return pattern.FindString(intput)
}

// ExtractComment will take a full function input such as:
// 		`// TestEvent comment
// 		// multiline 1
// 		// multiline 2
//		type TestEvent struct {
// 			TestAy  string
// 			TestBee int
// 		}`
//
// and produce:
// 		`// TestEvent comment
// 		// multiline 1
// 		// multiline 2`
//
// if there is no comment then an empty string is returned
func ExtractComment(input string) string {
	pattern := regexp.MustCompile(`(?P<all_comments>(\/\/\s*[\w\s\d]+)*)type`)
	matches := pattern.FindStringSubmatch(input)
	names := pattern.SubexpNames()
	result := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return strings.Trim(result["all_comments"], "\n\t")
}

type PublicMember struct {
	Type string
	Name string
}

func ExtractPublicMembers(input string) []PublicMember {
	pattern := regexp.MustCompile(`{[\n\r\t]+(?P<member>(?P<name>[\w]+)[\s\t]*(?P<type>[\[\]\w]+)[\s\n\r\t]+)+}`)
	patternMember := regexp.MustCompile(`(?P<name>[\w]+)[\s\t]*(?P<type>[\[\]\w]+)`)
	matches := pattern.FindStringSubmatch(input)
	members := patternMember.FindAllStringSubmatch(matches[0], -1)
	results := make([]PublicMember, len(members), cap(members))
	for i, match := range members {
		member := PublicMember{
			Name: match[1],
			Type: match[2],
		}
		results[i] = member
	}
	return results
}
