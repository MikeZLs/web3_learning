package pwd

import (
	"fmt"
	"testing"
)

func TestHashPwd(t *testing.T) {
	hash := HashPwd("1234")
	fmt.Println(hash)
}

func TestCheckPwd(t *testing.T) {
	ok := CheckPwd("$2a$04$xacK/wa9fnPgKus2cTQIxOUVxLOA55zL/JfQjcgUZl5qSBKu7XiTW", "1234")
	fmt.Println(ok)

}
