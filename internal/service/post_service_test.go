package service

import "testing"

func TestCreatePostService(t *testing.T) {
	s := NewPostService(nil)

	testCases := []string{
		"",
		" ",
		"\n",
	}

	for _, tc := range testCases {
		err := s.CreatePostService(1, tc, "テスト")
		if err == nil {
			t.Fatalf("タイトル [%s]、エラーのはずですが、結果は nil でした", tc)
		}

		expectedMsg := "タイトルが入力されていません"
		if err != nil && err.Error() != expectedMsg {
			t.Errorf("エラーメッセージが相違します: [%s]のはずですが、 [%s]でした", expectedMsg, err.Error())
		}
	}
}
