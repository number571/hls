package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"
	tu "github.com/number571/go-peer/settings/testutils"
)

var (
	settings = tu.NewSettings()
	mapping  = map[string]string{
		serviceInHLS: fmt.Sprintf("http://%s", serviceOutHLS),
	}
)

const (
	akeySize   = 1024
	addressHLS = "localhost:9571"

	serviceInHLS  = "hidden-echo-service"
	serviceOutHLS = "localhost:8080"

	sPrivKeyHLS = `Priv(go-peer\rsa){3082025E02010002818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C702030100010281810083866C4CA38EDAACE6B62A69A8C5682FD24F136A2E081C34F5AFB89372737AE3D052000317C7A2C9164180DD8E09E53C94F88341DFA8BD275E594CBAB9D4B008E1FE2D613D35202E841858BC665C0338221F34D9F143D60A5C2C4625459DAD3C0E3592F6B32D3E4105AB713CE42E73C44F10687954402E7A8D2952CC1C4589B9024100EB57FACD3A75AD2BA0BBF2152BF3760CFE45F78731C3B98D770DD790082753E4697CE8927112632F7BE86880121F4A08880DE7C45D16EB8D76E72214768D517B024100D8821E8FCF0DC72C5DD63A4A39CCC42601E1553022B75C9D01EF3DA2F706081E694C98E684BA5482B6E5975F6B371DE6E81AD42CB74A7CFA52A6D5522E0D4625024100BFDD211DF17C006AE206778CE520FDEC07DC98B9424BF3D92DE73E07316E86895FAAB29CB8CC29CA8B74E4C50C812FC516CE675602226E750D2BCFEFE8DABB4302407843AF1E4AF16855A8BA3B1EC8048A606262FCA30465BE3828BEF009FA158BA4F8F0E76E05044BB5604B204E8C8BCD3C5A69ACBA3A06526DEA4369F380493751024100EAE2ADA08E39EB52C8314DBB0F16A087DD9AE4BA3DADCD3BE515EA4193F83E62066ECDB3BE47CB377AE7F5480141FF60C20AAE818B3CDFAAA6244D97FB09FF0D}`
	sPubKeyHLS  = `Pub(go-peer\rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`
)

// client -> HLS -> server -\
// client <- HLS <- server -/
func TestService(t *testing.T) {
	// server
	srv := startServerHTTP(t)
	defer srv.Close()

	// service
	node := startNodeHLS(t)
	defer node.Close()

	// client
	time.Sleep(100 * time.Millisecond)
	startClientHLS(t)
}

// SERVER

func startServerHTTP(t *testing.T) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:    serviceOutHLS,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testEchoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Message string `json:"message"`
	}

	var resp struct {
		Echo  string `json:"echo"`
		Error int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Error = 1
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Echo = req.Message
	json.NewEncoder(w).Encode(resp)
}

// HLS

func startNodeHLS(t *testing.T) nt.Node {
	privKey := cr.LoadPrivKeyByString(sPrivKeyHLS)
	client := lc.NewClient(privKey, settings)

	node := nt.NewNode(client).
		Handle([]byte(PatternHLS), testRouteHLS)

	go func() {
		err := node.Listen(addressHLS)
		if err != nil {
			t.Error(err)
		}
	}()

	return node
}

func testRouteHLS(client lc.Client, msg lc.Message) []byte {
	request := new(Request)

	deserialize(msg.Body.Data, request)
	if request == nil {
		return nil
	}

	addr, ok := mapping[request.Host]
	if !ok {
		return nil
	}

	req, err := http.NewRequest(
		request.Method,
		addr+request.Path,
		bytes.NewReader(request.Body),
	)
	if err != nil {
		return nil
	}

	for key, val := range request.Head {
		req.Header.Add(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}

// CLIENT

func startClientHLS(t *testing.T) {
	priv := cr.NewPrivKey(akeySize)
	client := lc.NewClient(priv, settings)

	node := nt.NewNode(client).
		Handle([]byte(PatternHLS), nil)

	err := node.Connect(addressHLS)
	if err != nil {
		t.Error(err)
	}

	msg := lc.NewMessage(
		[]byte(PatternHLS),
		serialize(&Request{
			Host:   serviceInHLS,
			Path:   "/echo",
			Method: "GET",
			Head: map[string]string{
				"Content-Type": "application/json",
			},
			Body: []byte(`{"message": "hello, world!"}`),
		}),
	)

	pubKey := cr.LoadPubKeyByString(sPubKeyHLS)
	route := lc.NewRoute(pubKey, nil, nil)

	res, err := node.Broadcast(route, msg)
	if err != nil {
		t.Error(err)
	}

	if string(res) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		t.Error("result does not match")
	}
}
