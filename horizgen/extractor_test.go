package horizgen_test

import (
	"testing"

	"github.com/Brandon2255p/ehext/horizgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	i := `
	// TestEvent comment
	// multiline
	type TestEvent struct {
		TestAy  string
		TestBee int
	}
	
	// TestEvent comment
	type TestEvent struct {
		TestAy  string
		TestBee int
		MetaData ` + "`eh:\"optional\"`" + `
	}
	
	
	type TestEvent struct {
		TestAy  string
		TestBee int
		TestSlice []int
	}
`

	classes := horizgen.ExtractClasses(i)
	require.Len(t, classes, 3)
	require.Equal(t,
		`// TestEvent comment
	// multiline
	type TestEvent struct {
		TestAy  string
		TestBee int
	}`, classes[0])
}

func TestEventName(t *testing.T) {
	i := `// TestEvent comment
	type TestEvent struct {
		TestAy  string
		TestBee int
	}`
	name := horizgen.ExtractName(i)
	assert.Equal(t, "TestEvent", name)
}

func TestCommandName(t *testing.T) {
	i := `// TestCommand comment
	type TestCommand struct {
		TestAy  string
		TestBee int
	}`
	name := horizgen.ExtractName(i)
	assert.Equal(t, "TestCommand", name)
}

func TestAggregateName(t *testing.T) {
	i := `// TestAggregate comment
	type TestAggregate struct {
		TestAy  string
		TestBee int
	}`
	name := horizgen.ExtractName(i)
	assert.Equal(t, "TestAggregate", name)
}

func TestComment(t *testing.T) {
	i := `// TestEvent comment
	// multiline 1
	// multiline 2
	type TestEvent struct {
		TestAy  string
		TestBee int
	}`
	comment := horizgen.ExtractComment(i)
	assert.Equal(t, `// TestEvent comment
	// multiline 1
	// multiline 2`, comment)
}

func TestNoComment(t *testing.T) {
	i := `type TestEvent struct {
		TestAy  string
		TestBee int
	}`
	comment := horizgen.ExtractComment(i)
	assert.Equal(t, "", comment)
}

func TestMembers(t *testing.T) {
	i := `type TestEvent struct {
		TestAy  string
		privateMember  string
		TestBee int
		SomeArr	[]int
	}`
	comment := horizgen.ExtractPublicMembers(i)
	assert.Len(t, comment, 4)
	assert.Equal(t, comment[0], horizgen.PublicMember{Name: "TestAy", Type: "string"})
}
