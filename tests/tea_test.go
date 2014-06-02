package test

/**
	This tests basic Tea functionality as well as Amber's compilation.
**/

import (
	"bufio"
	"bytes"
	"github.com/vqtran/tea"
	"github.com/vqtran/tea/engines"
	"html/template"
	"testing"
)

func Test_SetGetEngine(t *testing.T) {
	// If try to set unsupported engine
	err := tea.SetEngine("doesnotexist")
	if err == nil {
		t.Fatal(err.Error())
	}
	// Correct engine
	err = tea.SetEngine("amber")
	if err != nil {
		t.Fatal(err.Error())
	}
	// Make sure its set correctly
	if *tea.GetEngine() != engines.Amber {
		t.Fatal("Engine not set correctly.")
	}
}

func Test_CompileAndGet(t *testing.T) {
	// Nonexistent directory
	err := tea.Compile("blah/", tea.Options{".amber", true})
	if err == nil {
		t.Fatal("Did not return error when directory not found.")
	}

	// Correct directory
	err = tea.Compile("amber_templates", tea.Options{".amber", true})
	if err != nil {
		t.Fatal("Returned error when should not have.")
	}

	// Nonexistent template
	val, ok := tea.Get("blah")
	if ok || val != nil {
		t.Fatal("Returned error when should not have.")
	}

	// All compilations/keys set correctly
	val1, ok1 := tea.Get("test1")
	if !ok1 || val1 == nil {
		t.Fatal("File not compiled or stored in map correctly.")
	}

	val2, ok2 := tea.Get("test2")
	if !ok2 || val2 == nil {
		t.Fatal("File not compiled or stored in map correctly.")
	}

	val3, ok3 := tea.Get("more/test3")
	if !ok3 || val3 == nil {
		t.Fatal("File not compiled or stored in map correctly.")
	}

	val4, ok4 := tea.Get("more/more/test4")
	if !ok4 || val4 == nil {
		t.Fatal("File not compiled or stored in map correctly.")
	}

	// Make sure file parsing is the same
	var doc1, doc4 bytes.Buffer
	val1.(*template.Template).Execute(&doc1, nil)
	val4.(*template.Template).Execute(&doc4, nil)
	if doc1.String() != doc4.String() {
		t.Fatal("Compiled files do not match when they should.")
	}

	var doc2, doc3 bytes.Buffer
	val2.(*template.Template).Execute(&doc2, nil)
	val3.(*template.Template).Execute(&doc3, nil)
	if doc2.String() != doc3.String() {
		t.Fatal("Compiled files do not match when they should.")
	}

	// Check against engine's CompileFile
	compilefile, err := engines.Amber.CompileFile("amber_templates/test1.amber")
	if err != nil {
		t.Fatal(err.Error())
	}
	var doc5 bytes.Buffer
	compilefile.(*template.Template).Execute(&doc5, nil)
	if doc1.String() != doc5.String() || doc4.String() != doc5.String() {
		t.Fatal("Compiled templates do not much engine's CompileFile result.")
	}
}

func Test_Delete(t *testing.T) {
	tea.Delete("test1")
	_, ok := tea.Get("test1")
	if ok {
		t.Fatal("Delete function does not properly remove key from map.")
	}
}

func Test_Clear(t *testing.T) {
	tea.Clear()
	if len(tea.GetCache()) != 0 {
		t.Fatal("Clear does not remove all elements.")
	}
}

// Test functionality of nonrecursive load
func Test_Nonrecursive(t *testing.T) {
	err := tea.Compile("amber_templates", tea.Options{".amber", false})
	if err != nil {
		t.Fatal("Returned error when should not have.")
	}

	// Shouldn't be anything other than test1 and test2
	for k := range tea.GetCache() {
		if k != "test1" && k != "test2" {
			t.Fatal("Non-recursive search not correct.")
		}
	}

	// Test1 and test2 should also be in there.
	val1, ok1 := tea.Get("test1")
	if val1 == nil || !ok1 {
		t.Fatal("Non-recursive search did not load in all files.")
	}

	val2, ok2 := tea.Get("test2")
	if val2 == nil || !ok2 {
		t.Fatal("Non-recursive search did not load in all files.")
	}
}

// Test Functionality of Render
func Test_Render(t *testing.T) {
	err := tea.Compile("amber_templates", tea.Options{".amber", true})
	if err != nil {
		t.Fatal("Returned error when should not have.")
	}
	var b bytes.Buffer
   writer := bufio.NewWriter(&b)

	// Test for error condition
	err = tea.Render(writer, "blah", nil)
	if err == nil {
		t.Fatal("Did not return error when it should have.")
	}

	// Usual condition
	data := map[string]string{"Name":"World"}
	err = tea.Render(writer, "test1", data)
	if err != nil {
		t.Fatal("Error rendering valid template.")
	}
}
