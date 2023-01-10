package jetton

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"strings"
	"testing"
	"time"
)

var api = func() *ton.APIClient {
	client := liteclient.NewConnectionPool()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.AddConnectionsFromConfigUrl(ctx, "https://ton-blockchain.github.io/testnet-global.config.json")
	if err != nil {
		panic(err)
	}

	return ton.NewAPIClient(client)
}()

func TestJettonMasterClient_GetJettonData(t *testing.T) {
	cli := NewJettonMasterClient(api, address.MustParseAddr("EQAbMQzuuGiCne0R7QEj9nrXsjM7gNjeVmrlBZouyC-SCLlO"))

	data, err := cli.GetJettonData(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if data.AdminAddr.String() != "EQDIhshKpDtDt-uTawnlGIq-cNijBgB7jyYprczoAOXTWLwf" {
		t.Fatal("admin addr diff")
	}

	if data.TotalSupply.Uint64() < 20000000000000000 || data.TotalSupply.Uint64() > 21000000000000000 {
		t.Fatal("supply diff")
	}

	if data.Content.(*nft.ContentOffchain).URI != "ipfs://QmeebZm4sCmXGdM9iiD8TtYDyQwyFa3Dm27hxonVQyFT4P" {
		t.Fatal("content diff")
	}
}

func TestJettonMasterClient_GetWalletAddress(t *testing.T) {
	cli := NewJettonMasterClient(api, address.MustParseAddr("EQAbMQzuuGiCne0R7QEj9nrXsjM7gNjeVmrlBZouyC-SCLlO"))
	ctx := api.Client().StickyContext(context.Background())

	data, err := cli.GetJettonWallet(ctx, address.MustParseAddr("EQC9bWZd29foipyPOGWlVNVCQzpGAjvi1rGWF7EbNcSVClpA"))
	if err != nil {
		t.Fatal(err)
	}

	if data.Address().String() != "EQCRosA17HVslmZKvqcM1aCNZgbzc4plnHpP29wHU_SWlblY" {
		t.Fatal("owner diff")
	}

	b, err := data.GetBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if b.String() != "22686.666348532" {
		t.Fatal("balance diff:", b.String())
	}
}

func TestJettonMasterClient_Transfer(t *testing.T) {
	cli := NewJettonMasterClient(api, address.MustParseAddr("EQAbMQzuuGiCne0R7QEj9nrXsjM7gNjeVmrlBZouyC-SCLlO"))

	ctx := api.Client().StickyContext(context.Background())

	w := getWallet(api)
	log.Println("test wallet:", w.Address().String())

	tokenWallet, err := cli.GetJettonWallet(ctx, w.Address())
	if err != nil {
		t.Fatal(err)
	}

	b, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("token balance:", b.String())

	fmt.Println("Transferring tokens...")

	amt := tlb.MustFromTON("1.15")
	to := address.MustParseAddr("EQD4vUD2PYRLQd0mSwjmnnWSpeulTjZoFypJVUJAyJoUbrRu")
	transferPayload, err := tokenWallet.BuildTransferPayload(to, amt, tlb.MustFromTON("0"), nil)
	if err != nil {
		panic(err)
	}

	msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)

	err = w.Send(ctx, msg, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Burning tokens...")
	burnPayload, err := tokenWallet.BuildBurnPayload(amt, to)
	if err != nil {
		panic(err)
	}

	msg = wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), burnPayload)

	err = w.Send(ctx, msg, true)
	if err != nil {
		panic(err)
	}

	b2, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if b.NanoTON().Uint64() == b2.NanoTON().Uint64() {
		t.Fatal("balance was not changed after burn")
	}

	want := b.NanoTON().Uint64() - amt.NanoTON().Uint64()*2
	got := b2.NanoTON().Uint64()
	if want != got {
		t.Fatal("balance not expected, want ", want, "got", got)
	}
}

func getWallet(api *ton.APIClient) *wallet.Wallet {
	words := strings.Split("cement secret mad fatal tip credit thank year toddler arrange good version melt truth embark debris execute answer please narrow fiber school achieve client", " ")
	w, err := wallet.FromSeed(api, words, wallet.V3)
	if err != nil {
		panic(err)
	}
	return w
}
