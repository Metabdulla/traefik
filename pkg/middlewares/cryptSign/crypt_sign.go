package cryptSign

import (
	"context"
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/util"
	"io/ioutil"
	"net/http"
)

const (
	typeName = "EncryptSign"
	ENCRYPT_TYPE_RSA= "rsa"
	ENCRYPT_TYPE_3DES="3des"
)


type ResponseBody struct {
	ResponseBodyRaw
	SignValue  string `json:"signValue"`
}

//RequestBodyRaw  make sure that the filed is sorted
type ResponseBodyRaw struct {
	Code  string `json:"code"`
	Message  string `json:"message"`
	BizContent  string `json:"bizContent"`
	Timestamp   string `json:"timestamp"`
}

func md5Hash( data []byte) string {
	//todo
	h:= md5.New()
	h.Write(data)
	sum:=  h.Sum(nil)
	return hex.EncodeToString(sum)
}


func (c * EncryptSign)getSortedKeyValue( data []byte  ) ( string , error ){
	m:= make(map[string]interface{})
	err := json.Unmarshal(data,&m)
	if err!=nil {
		return "", err
	}
	if len(m) ==0 {
		return  "",fmt.Errorf("data is nil")
	}
	delete(m, c.SignFiledName)
	sortedData  ,err := json.Marshal(&m)
	return string(sortedData),err
}

const  DefaultSignFiledName = "signValue"
const DefaultEncryptFieldName  ="bizContent"

type EncryptSign struct {
	name string
	SignFiledName    string
	EncryptFieldName string
	next http.Handler
	Type  string
}

func New(ctx context.Context, name string , next http.Handler) (http.Handler ,error){
	middlewares.GetLogger(ctx, name, typeName).Debug("Creating middleware")
	return &EncryptSign{
		name: name,
		next:next,
	},nil
}

func (c *EncryptSign) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := middlewares.GetLogger(req.Context(), c.name, typeName)
	var err error
	var bodyData []byte
	if req.Body==nil {
		err = fmt.Errorf("body is nil")
	}else {
		bodyData ,err = ioutil.ReadAll(req.Body)
	}

	if err!=nil {
		logger.Debug("read body faild")
		c.Error(rw,"102","验签失败",400 )
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(bodyData,&m)
	if err!=nil {
		logger.Debug("read body faild")
		c.Error(rw,"102","验签失败",400 )
		return
	}

	c.next.ServeHTTP(rw,req)
}

//encryptAndSign for client
func (c *EncryptSign)encryptAndSign(bodyData []byte , encType ,publicKey string,  privKey  *rsa.PrivateKey  ) {


}
//checkSignAndDecrypt for
func (c *EncryptSign)checkSignAndDecrypt(bodyData []byte , encType ,publicKey string,  privKey  *rsa.PrivateKey  ) {


}

//encrypt for client
func (c *EncryptSign)encrypt(bodyData []byte, encType, publicKey string ){

}

//decrypt
func (c *EncryptSign)sign(data []byte  ,encType string,privKey  *rsa.PrivateKey){

}

//decrypt
func (c *EncryptSign)decrypt(data []byte  ,encType ,publicKey string  ,privKey  *rsa.PrivateKey){

}

func (c *EncryptSign)checkSign(data []byte  ,encType,publicKey string ) {

}

func (c *EncryptSign)Error(rw http.ResponseWriter, bizCode string ,message string , httpCode int  ) {
   data :=  ResponseBody{
   ResponseBodyRaw:ResponseBodyRaw{
   	 Code:bizCode,
   	 Message: message,
   },
   }
   util.JSON(rw,httpCode,&data)
}
