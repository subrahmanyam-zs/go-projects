package Department

import "testing"

func TestValidatePost(t *testing.T) {
	testCases := []struct {
		desc string
		input
		expectedOutput []byte
	}{
		{"Valid input", []byte(`{"id":1, "name":"HR","floorNo": 1}`), []byte(`{"Id":1,"Name":"HR","FloorNo":1}`)},
		{"Invalid input", []byte(`{"id":0,"name":"Tech","floorNo":2}`), []byte("Invalid id")},
		{"for Unmarshal error", []byte(`{"id":"2","name":"hr","floorNo":2}`), []byte("Unmarshal Error")},
	}

}
