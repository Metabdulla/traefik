package cryptSign

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T){
	 r:= RequestBodyRaw{
			AppKey:"10002",
			BizContent:"SRr6MzYOzEdlr+LZ7uHqZM+rzfTsfuScPaUqL/d6wIuHGrZZ9Q4uQQ==",
			ServiceName:"individual-underwriting",
			//SignValue:"67a78402de8a44a778d60039afe9ff54",
			Timestamp:"1496819504250",
			Version:"1.0.0",

	}
	 data,_ := json.Marshal(&r)
	 data2 :=[]byte(`{"appKey":"10002","bizContent":"SRr6MzYOzEdlr+LZ7uHqZM+rzfTsfuScPaUqL/d6wIuHGrZZ9Q4uQQ==","serviceName":"individual-underwriting","version":"1.0.0","timestamp":"1496819504250"}`)
	fmt.Println(string(data))
	fmt.Println(md5Hash(data))
	str , err := getSortedKeyValue(data2)
	if err!=nil {
		t.Fatal(err)
	}
	fmt.Println(str)
	fmt.Println(md5Hash([]byte(str)))

}
