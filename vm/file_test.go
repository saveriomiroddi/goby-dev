package vm

import (
	"testing"
)

func TestFileDeletion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "file"

		File.open("/tmp/out1.txt", "w", 0755)
		File.open("/tmp/out2.txt", "w", 0755)
		File.open("/tmp/out3.txt", "w", 0755)

		File.delete("/tmp/out1.txt", "/tmp/out2.txt", "/tmp/out3.txt")
		`, 3},
		{`
		require "file"

		File.open("/tmp/out.txt", "w", 0755)
		File.delete("/tmp/out.txt")
		File.exist("/tmp/out.txt")
		`, false},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestFileWrite(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "file"

		l = 0
		File.open("/tmp/out.txt", "w", 0755) do |f|
		  l = f.write("12345")
		end

		l
		`, 5},
		{`
		require "file"

		File.open("/tmp/out.txt", "w", 0755) do |f|
		  f.write("Goby is awesome!!!")
		end

		File.new("/tmp/out.txt").read
		`, "Goby is awesome!!!"},
		{`
		require "file"

		File.open("/tmp/out.txt", "w", 0755)
		File.new("/tmp/out.txt").size
		`, 0},
		{`
		require "file"

		File.open("/tmp/out.txt", "w", 0755)
		File.exist("/tmp/out.txt")
		`, true},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

// TODO: Add failed tests
func TestFileObject(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "file"

		f = File.new("../test_fixtures/file_test/size.gb")
		f.name
		`, "../test_fixtures/file_test/size.gb"},
		{`
		require "file"

		f = File.new("../test_fixtures/file_test/size.gb")
		f.size
		`, 22},
		{`
		require "file"

		f = File.new("../test_fixtures/file_test/size.gb")
		f.close
		`, nil},
		{`
		require "file"

		f = File.new("../test_fixtures/file_test/size.gb")
		f.read
		`, "this file's size is\n22"},
		{`
		require "file"

		file = ""
		File.open("../test_fixtures/file_test/size.gb", "r", 0755) do |f|
	 	  file = f.read
	 	end
	 	file
		`, "this file's size is\n22"},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestExtnameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require "file"
		File.extname("loop.gb")
		`, ".gb"},
		{`
		require "file"
		File.extname("text.txt")
		`, ".txt"},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestBasenameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require "file"
		File.basename("/home/goby/plugin/test.gb")
		`, "test.gb"},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSplitMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []Object
	}{
		{`
		require "file"
		File.split("/home/goby/plugin/test.gb")
		`, []Object{initStringObject("/home/goby/plugin/"), initStringObject("test.gb")}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testArrayObject(t, i, evaluated, vm.initArrayObject(tt.expected))
	}
}

func TestJoinMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require "file"
		File.join("test1", "test2", "test3")
		`, "test1/test2/test3"},
		{`
		require "file"
		File.join("goby", "plugin")
		`, "goby/plugin"},
		{`
		require "file"
		File.join("plugin")
		`, "plugin"},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSizeMethod(t *testing.T) {
	input := `
	require "file"

	File.size("../test_fixtures/file_test/size.gb")
	`

	evaluated := testEval(t, input)
	checkExpected(t, 0, evaluated, 22)
}

//@TODO add test for chmod form a847c8b41f29657b380c1731ec36a660dbf49bc4
