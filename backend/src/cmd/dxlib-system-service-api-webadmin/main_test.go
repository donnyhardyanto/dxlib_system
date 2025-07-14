package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/donnyhardyanto/dxlib/utils/crypto/datablock"
	"github.com/donnyhardyanto/dxlib/utils/crypto/x25519"
	"github.com/donnyhardyanto/dxlib/utils/http/client"
	"github.com/donnyhardyanto/dxlib/utils/lv"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain")
	go func() {
		main()
	}()
	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestUnitTest(t *testing.T) {
	t.Run("DataBlock Serialization/Deserialization test", func(t *testing.T) {
		var err error
		// Test the DataBlock serialization and deserialization
		lvData := "This is a test data"
		lvDataAsBytes := []byte(lvData)
		lvDataLength := uint32(len(lvDataAsBytes))

		lvDataBlock := lv.LV{
			Length: lvDataLength,
			Value:  lvDataAsBytes,
		}

		dataBlock := datablock.DataBlock{
			Time:   lv.LV{},
			Nonce:  lv.LV{},
			PreKey: lv.LV{},
			Data:   lvDataBlock,
			DataHash: lv.LV{
				Length: 0,
				Value:  nil,
			},
		}

		aLV, err := dataBlock.AsLV()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		newDataBlock, err := datablock.NewDataBlockFromLV(aLV)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		if string(newDataBlock.Nonce.Value) != string(dataBlock.Nonce.Value) {
			t.Log("Nonce is not the same")
			t.FailNow()
			return
		}

		if newDataBlock.Data.Length != dataBlock.Data.Length {
			t.Log("Data length is not the same")
			t.FailNow()
			return
		}

		if string(newDataBlock.Data.Value) != string(dataBlock.Data.Value) {
			t.Log("Data value is not the same")
			t.FailNow()
			return
		}
	})
}

func TestACombineV1(t *testing.T) {
	t.Run("LV Combine test", func(t *testing.T) {
		lvA, err := lv.NewLV([]byte("abc"))
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvB, err := lv.NewLV([]byte("def"))
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvC, err := lv.CombineLV(lvA, lvB)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvCAsBytes, err := lvC.MarshalBinary()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvCx1 := lv.LV{}
		err = lvCx1.UnmarshalBinary(lvCAsBytes)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		lvDs, err := lvC.Expand()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvD1 := lvDs[0]
		lvD2 := lvDs[1]

		x1 := bytes.Equal(lvD1.Value, lvA.Value)
		if x1 != true {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		x2 := bytes.Equal(lvD2.Value, lvB.Value)
		if x2 != true {
			t.Log(err.Error())
			t.FailNow()
			return
		}

	})
}
func TestAPI1(t *testing.T) {

	APISystemProtocol := "http://"

	APISystemAddress := "127.0.0.1:15000"

	var index string

	t.Run("User login test", func(t *testing.T) {
		var err error
		// Test the SelfLogin function

		t.Log("1:")

		edA0PublicKeyAsBytes, edA0PrivateKeyAsBytes, err := ed25519.GenerateKey(nil)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		ecdhA1PublicKeyAsBytes, ecdhA1PrivateKeyAsBytes, err := x25519.GenerateKeyPair()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		ecdhA2PublicKeyAsBytes, ecdhA2PrivateKeyAsBytes, err := x25519.GenerateKeyPair()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		ecdhA1PublicKeyAsHexString := hex.EncodeToString(ecdhA1PublicKeyAsBytes[:])
		ecdhA2PublicKeyAsHexString := hex.EncodeToString(ecdhA2PublicKeyAsBytes[:])
		edA1PublicKeyAsHexString := hex.EncodeToString(edA0PublicKeyAsBytes[:])

		_, r, err := client.HTTPClientReadAll("POST", APISystemProtocol+APISystemAddress+"/self/prekey", nil, utils.JSON{
			"a0": edA1PublicKeyAsHexString,
			"a1": ecdhA1PublicKeyAsHexString,
			"a2": ecdhA2PublicKeyAsHexString,
		})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		if r.StatusCode != 200 {
			t.Log("Status code is not 200")
			t.FailNow()
			return
		}

		responseDataAsJSON, err := r.BodyAsJSON()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		index = responseDataAsJSON["i"].(string)
		edB0PublicKeyAsHexString := responseDataAsJSON["b0"].(string)
		ecdhB1PublicKeyAsHexString := responseDataAsJSON["b1"].(string)
		ecdhB2PublicKeyAsHexString := responseDataAsJSON["b2"].(string)

		edB0PublicKeyAsBytes, err := hex.DecodeString(edB0PublicKeyAsHexString)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		ecdhB1PublicKeyAsBytes, err := hex.DecodeString(ecdhB1PublicKeyAsHexString)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		ecdhB2PublicKeyAsBytes, err := hex.DecodeString(ecdhB2PublicKeyAsHexString)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		sharedKey1AsBytes, err := x25519.ComputeSharedSecret(ecdhA1PrivateKeyAsBytes[:], ecdhB1PublicKeyAsBytes)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		sharedKey2AsBytes, err := x25519.ComputeSharedSecret(ecdhA2PrivateKeyAsBytes[:], ecdhB2PublicKeyAsBytes)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		t.Log("2:")

		userLogin := os.Getenv("TEST_USER_LOGIN")
		password := os.Getenv("TEST_USER_PASSWORD")

		lvUserLogin, err := lv.NewLV([]byte(userLogin))
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		lvPasswordHash, err := lv.NewLV([]byte(password))
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		dataBlockEnvelopeAsHexString, err := datablock.PackLVPayload(index, edA0PrivateKeyAsBytes, sharedKey1AsBytes, lvUserLogin, lvPasswordHash)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		_, r, err = client.HTTPClientReadAll("POST", APISystemProtocol+APISystemAddress+"/self/login", nil, utils.JSON{
			"i": index,
			"d": dataBlockEnvelopeAsHexString,
		})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		if r.StatusCode != 200 {
			t.Logf("Status code is not 200 but %d", r.StatusCode)
			t.FailNow()
			return
		}

		responseDataAsJSON, err = r.BodyAsJSON()
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		dataBlockEnvelopeAsHexString = responseDataAsJSON["d"].(string)

		lvPayloadElements, err := datablock.UnpackLVPayload(index, edB0PublicKeyAsBytes, sharedKey2AsBytes, dataBlockEnvelopeAsHexString)
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}

		lvSessionKey := lvPayloadElements[0]

		sessionKey := string(lvSessionKey.Value)
		if sessionKey == "" {
			t.Log("session_key not found in parameter")
			t.FailNow()
			return
		}

		_, r, err = client.HTTPClientReadAll("POST", APISystemProtocol+APISystemAddress+"/self/logout", client.HTTPHeader{
			"Authorization": "Bearer " + sessionKey,
		}, utils.JSON{})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		if r.StatusCode != 200 {
			t.Logf("Status code is not 200, but %d", r.StatusCode)
			t.FailNow()
			return
		}

		_, r, err = client.HTTPClientReadAll("POST", APISystemProtocol+APISystemAddress+"/self/logout", client.HTTPHeader{
			"Authorization": "Bearer " + sessionKey,
		}, utils.JSON{})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		if r.StatusCode == 200 {
			t.Logf("Status code is should not 200, but %d", r.StatusCode)
			t.FailNow()
			return
		}

		time.Sleep(2 * time.Minute)
		_, r, err = client.HTTPClientReadAll("POST", APISystemProtocol+APISystemAddress+"/self/logout", client.HTTPHeader{
			"Authorization": "Bearer " + sessionKey,
		}, utils.JSON{})
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
			return
		}
		if r.StatusCode == 200 {
			t.Logf("Status code is should not 200, but %d", r.StatusCode)
			t.FailNow()
			return
		}
	})
}
