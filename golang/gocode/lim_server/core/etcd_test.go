package core

import (
	"context"
	"fmt"
	"testing"
)

func TestEtcd(t *testing.T) {
	client := InitEtcd("127.0.0.1:2379")
	_, err := client.Put(context.Background(), "auth_api", "127.0.0.1:20021")

	getRes, err := client.Get(context.Background(), "auth_api")
	if err == nil && len(getRes.Kvs) > 0 {
		fmt.Println(string(getRes.Kvs[0].Value))
	}

}
