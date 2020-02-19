package main

import (
	"encoding/base64"
	"encoding/json"
    "fmt"
    "crypto/tls"
	"io/ioutil"
	"strings"
    "net/http"
    "crypto/x509"
)

func main() {

//    tfvars := "ocp/terraform.tfvars.json"
//    url := "https://192.168.111.5:22623/config/master"

    // Load tfvars to retrieve the URL and CA cert
	tfVarsB, err := ioutil.ReadFile("ocp/terraform.tfvars.json")
	if err != nil {
		panic(err)
	}
    var tfVars map[string]interface{}
	if err := json.Unmarshal(tfVarsB, &tfVars); err != nil {
		panic(err)
	}
	//fmt.Println(tfVars)
	//json.Unmarshal([]byte(str), &res)
    var ignConf map[string]interface{}
	if err := json.Unmarshal([]byte(tfVars["ignition_master"].(string)), &ignConf); err != nil {
		panic(err)
    }
	fmt.Println(ignConf)
    url := ignConf["ignition"].(map[string]interface{})["config"].(map[string]interface{})["append"].([]interface{})[0].(map[string]interface{})["source"].(string)
	fmt.Println(url)
	caCertRaw := ignConf["ignition"].(map[string]interface{})["security"].(map[string]interface{})["tls"].(map[string]interface{})["certificateAuthorities"].([]interface{})[0].(map[string]interface{})["source"].(string)
    caCertB64 := strings.TrimPrefix(caCertRaw, "data:text/plain;charset=utf-8;base64,")
	fmt.Println(caCertB64)
	caCert, err := base64.StdEncoding.DecodeString(caCertB64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("CACERT %q\n", caCert)

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
            },
        },
    }

    // Get the data
    resp, err := client.Get(url)
    if err != nil {
		panic(err)
    }
    defer resp.Body.Close()
    fullIgn, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		panic(err)
    }
	fmt.Println(string(fullIgn))
}
