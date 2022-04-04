package merkletree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type merkleRootResponse struct {
	Root string `json:"root"`
}
type merkleProofResponse struct {
	MerkleProof []string `json:"merkle_proof"`
	Root        string   `json:"root"`
}

func GetMerkleRoot(leaves string) (string, error) {
	url := os.Getenv("AWS_SERVERLESS_HOST") + "/get-merkle-root"
	method := "POST"

	payload := strings.NewReader(`{"leaves":` + leaves + `}`)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	rawMerkleRoot, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var merkleRoot merkleRootResponse
	err = json.Unmarshal([]byte(rawMerkleRoot), &merkleRoot)
	if err != nil {
		return "", err
	}
	return string(merkleRoot.Root), nil
}

func GetMerkleProof(leaves string, target string) (string, string, error) {
	url := os.Getenv("AWS_SERVERLESS_HOST") + "/get-merkle-proof"
	method := "POST"

	payload := strings.NewReader(`{"leaves":` + leaves + `,"targetLeaf":` + target + `}`)

	fmt.Println(payload)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	fmt.Println(res)
	defer res.Body.Close()

	rawMerkleProof, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	fmt.Println(string(rawMerkleProof))
	var merkleProof merkleProofResponse
	//"json: cannot unmarshal array into Go struct field merkleProofResponse.merkle_proof of type string"
	err = json.Unmarshal([]byte(rawMerkleProof), &merkleProof)
	if err != nil {
		return "", "", err
	}
	proof := strings.Join(merkleProof.MerkleProof, ",")
	return string(merkleProof.Root), "[" + proof + "]", nil
}
