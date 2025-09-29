package main

import (
	_ "crypto/sha256" // register hash algorithm
	"fmt"

	"github.com/opencontainers/go-digest"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DiffID digest.Digest
type ChainID digest.Digest

func main() {
	// "sha256:9a659e9da99d18005f614a3fb4f5edf4e42d541af26cd02947d624150536ca62"
	dgst := digest.NewDigestFromEncoded(digest.Algorithm(digest.SHA256), "9a659e9da99d18005f614a3fb4f5edf4e42d541af26cd02947d624150536ca62")
	fmt.Println(dgst)
	fmt.Println(len(dgst))
	fmt.Println(dgst.Validate())
	dgst = digest.NewDigestFromEncoded(digest.Algorithm(digest.SHA256), "4364d123be7ce44234ada8aaf447c44d4e067332bfe3409330c6b13b9c639235")
	fmt.Println(dgst)
	fmt.Println(len(dgst))
	fmt.Println(dgst.Validate())
	dgst = digest.NewDigestFromEncoded(digest.Algorithm(digest.SHA256), "8261c5794e4af988004cd99bc5a9d88806f9c9f1fa59f6c43b002e90b88e7a9c")
	fmt.Println(dgst)
	fmt.Println(len(dgst))
	fmt.Println(dgst.Validate())

	dgst = digest.NewDigestFromEncoded(digest.Algorithm(digest.SHA256), "68e656b251e67e8358bef8483ab0d51c6619f3e7a1a9f0e75838d41ff368f728")
	fmt.Println(dgst)
	fmt.Println(len(dgst))
	fmt.Println(dgst.Validate())

	fmt.Println("===============================")
	diffIDs := []DiffID{
		"sha256:ad00c0c1d948d1c75112e914e3610acd57755d57bd9864b6b3880e6de892ef1e",
		"sha256:b79ac58f353c0513558867675b8ca4be2631498182b80dcdbeefe713fabbc3fd",
		"sha256:18c4c015e3390518beb84419d5bfe94b55cc54e3062782c0c551b78220f268e7",
		"sha256:c6eef6a3993593efa1304fe7e3ce9e52af53925869801dbb769c5e36fcf74eca",
		"sha256:03db561303525431d105799a06fda15e9f4cfde92ff263d532b5b0be65a8a62f",
		"sha256:78c9748f2d2918ee4cf5120b8b5519a16fe3a2a6a625d23e54282f7a00d8a076",
		"sha256:e1732f90232ae81e058100b0a97135083bb6596accbf7cb16facf94a8cd17061",
		"sha256:6b733e4833c866355119eda8e9c13d03f87cd9e9fd3f51ba051dd7be46eb2138",
		"sha256:fe6828cc147ad881ac18e9e91b935cf2439196ec7d3dd9c920a4a9c965a421ef",
		"sha256:8e1349c4476850d9312afb11ad815eb66ae86d6d86f70cb8e8647153e5aa48ab",
	}
	fmt.Println(CreateChainID(diffIDs)) // sha256:31a12eef15201b21bdf74385eff4e3100c3e6fcc6c3e5eaeabd359684be7891c

	q, err := resource.ParseQuantity("-1")
	fmt.Println("Parse -1:", err) // 输出错误
	fmt.Println("quantiry: ", q)
}

// moby layer/layer.go
// CreateChainID returns ID for a layerDigest slice
func CreateChainID(dgsts []DiffID) ChainID {
	return createChainIDFromParent("", dgsts...)
}

func createChainIDFromParent(parent ChainID, dgsts ...DiffID) ChainID {
	if len(dgsts) == 0 {
		return parent
	}
	if parent == "" {
		return createChainIDFromParent(ChainID(dgsts[0]), dgsts[1:]...)
	}
	m := string(parent) + " " + string(dgsts[0])
	fmt.Println(m)
	// H = "H(n-1) SHA256(n)"
	dgst := digest.FromBytes([]byte(m))
	fmt.Println(dgst)
	return createChainIDFromParent(ChainID(dgst), dgsts[1:]...)
}
