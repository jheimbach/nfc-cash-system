package auth

import "testing"

func Test_createRandomKey(t *testing.T) {
	t.Run("create some random keys", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			got := CreateRandomKey()
			if len(got) != 64 {
				t.Errorf("createRandomKey should've returned 64 char long string, got %d", len(got))
			}
			t.Log(got)
		}
	})
}
