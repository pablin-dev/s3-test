package entity

import "testing"

func TestYamlFile(t *testing.T) {
	yf := YamlFile{ID: "test", Version: 1, Expression: "1+1"}
	if yf.ID != "test" {
		t.Errorf("expected ID test, got %s", yf.ID)
	}
}
