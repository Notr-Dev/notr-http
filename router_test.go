package notrhttp

import (
	"testing"
)

func TestMatchPath(t *testing.T) {
	tests := []struct {
		path     string
		pattern  string
		expected bool
		params   map[string]string
	}{
		{"/", "/", true, map[string]string{}},
		{"/users", "/users", true, map[string]string{}},
		{"/users/123", "/users/{id}", true, map[string]string{"id": "123"}},
		{"/users/123/posts", "/users/{id}/posts", true, map[string]string{"id": "123"}},
		{"/users/123/posts/456", "/users/{id}/posts/{postId}", true, map[string]string{"id": "123", "postId": "456"}},
		{"/users/123", "/users", false, map[string]string{}},
		{"/users", "/users/{id}", false, map[string]string{}},
		{"/users/123/posts", "/users/{id}/comments", false, map[string]string{}},
		{"/users/123/posts/456", "/users/{id}/posts", false, map[string]string{}},
	}

	for _, test := range tests {
		t.Run(test.path+"_"+test.pattern, func(t *testing.T) {
			is, params := matchPath(test.path, test.pattern)
			if is != test.expected {
				t.Errorf("expected %v, got %v", test.expected, is)
			}
			if len(params) != len(test.params) {
				t.Errorf("expected params %v, got %v", test.params, params)
			}
			for key, value := range test.params {
				if params[key] != value {
					t.Errorf("expected param %s to be %s, got %s", key, value, params[key])
				}
			}
		})
	}
}
