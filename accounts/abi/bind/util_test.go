// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/accounts/abi/bind/backends"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/crypto"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
var (
	chainID      = big.NewInt(1)
	activeKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	activeAddr   = crypto.PubkeyToAddress(activeKey.PublicKey)
)

var waitDeployedTests = map[string]struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}{
	"successful deploy": {
		code:        `0161736db90aef0100789cbc585d90145596feeecdbad9d59554757517d0dd0a9a14b8fe6c530d0d34c8f257a040a30648cb2a616075565556754a5566939955d0ecd2d5dbb288babae8cbee8b6884eec68611acfb6484f8b20f138ccef8e2c484ce44cc8f118e46cc4410e36838e388334c9c9b593f2d88ce0433fd50f7dc73bef39d73cffde9aaf3fbfebb5e5fb80286576500d8666d823530c1676630c16630c11b0dd69810338d99460313c004586382912ed268481b2783a04903fc3d9e524cbbaeee37eba6eb0334493ee85abef9a0e54fee732cdb375d3052f7ec378de2155a355b2cee323c709a440f8cac1bcdeda84e41694fc76bf9601acb168bf798d36376c9418414a931dbf22da3621d37ffd1702d235f313d08b274ef32fd71d32e9a2e549af7ed32fd1d8eedbb46c1cf168baee9798141dde1d805c30f022ca0780f38e3be6bd96574912a32be3bbb2610bbef2eb866c1a99b2ea212bdcfb5ec165b37a9c4dd476a4625205bb4cbf4b71b15c32e983b5da7dac405547beba65b748d52b362942ba19a39a25b8976092162a24b44049b628c098581ab8b9406cbbe7f3aa235d01c7e12d5ba1eed56ab66d571a739b45c6ed234a67279c3331516cbe58a866fe44cbba8f064d12c540cd72ceeacd905df726c2cd4561426cd233533ef3887b1285230bc492c8e1cb62a15f4770545f431d0bd62a751a9e48dc2610cc62e8a2e0eed266cc77624ff6576763690e60229f9188dda238ccf64df3e3d3b8b73696445f6cd4064d9ee2c92a708a2232b7496ed6aceb33f6aa2bb74b461c9c7a511c9d3348e7289d27ec73bf8459b5fede01759a1ab57f28b6c972edab02c748c72e9ceb7e92c8bdbb8a44d9e9456b6338eec7352312b17967c22483ef9a41c59f22999a59c25ff4deab8aee891e4d3242fe3d9ff3e2da3dc16649e7c46d29e3b2d81c18c6f6b99759efc7752c5917dad23e8bcccd23cccadb12b8eecec57603a4ffeeb6c80929926cfc899d25c0ff9b4927936588bae249f23a947eba846163a1b0d22ea9497446863d76b5fb5431d4caccdc43b985856e8fc4a2696edd2d95577b055a3eb98a756e4c8ce46b51df4388c1b2573bc96d78fe99bf4691855a766fbbae5e9bee3e85ed5a854e01cb54d7723c29762232c4fb71d5f976a20fb66545b88abfdbddba5658a66be56ce5976c9b9a00011a98fd02b8d05581ee25e0cc7ff07c02f02584fc21aa1ec8e00ca26f248f235e21300a6b4801e0efcb39423a297018f4b390a750ac025fa27b0a6ab0620fa5b002a93ea6e06b011f539001b491c8a2e674057946719f00a11ac15ca4314d20cd1e8a5550088fd8c3e7e49496a02c04d004e11c7c70b76d19063f1bd72fc0c403c27c52f091d3dc8001ec42f907a5582787701f80f9a7e4e8b88cfb5c5932db1e70312bf20314939ac231dadec529008245309c04b04bcdc66badc669262af64e28ce43ea2ea276fa262ac6365d2f440334a848561642128f8210ec41ea58f3a0752ef86bbc6624de442426e03f02ac58bcb788be63ae4936d39c8a947ca926a80a8924daac54435463bc581fe47a98877b28115001e01f006cdb70c926d15eddb36dc5027baa8a46341cc6847cc40c68d27011c65c077896067fc4339d06ec77f234559b8082791b62ffe47d6123f6f893d3102dc23c51e12ef25869e7e12f7c99d5b42e2fd24ca1c1fa21c77cb4dbc18569eefe928bcb4e45968d9db615942359f25edbe76e5e53e3d11c2dbd5c7c0972ca8d9db94eac4520a2d38b095fcf348ece4c02100ef90b924972a0b50329b89fe27212731f0320f887e4ce6aa243ad7247228e14478afab4681ae36f01a8f872ad7b0cba647172705e066007fdff1240c01c8001801b0a543bf12c0660421ee03b01740313cde950edcd3009e07f002807374d2007c07c0055a7607ee08036a0c9861c0cbacaddfcb811c070c0e381c38c281531c38cd8167781b472e7d00e8b22f0fe7bbc3b11c8ea3e1980dc73d00e884ce01a083768901b4257773806aff1207a8b2173be2e015a55938239f77cd3aeb65b72452424924fb6f4cf4b203fdf42c41498ca5b60e6ed4fe41e30380b20e63292092629a463a40c42584261b34405d8cb19492086c5d3d121d5d0125b145d334a03bc37a599fba3a05c40468d032204de0b27510589061c1e4d6c1ad83ba06c405020571250282268272037ad67620924bd9ea14611ed20e6a406f681364ebebb0092045d303fdc1646107b368322f121dde8b835429ebfe0c9a19538a039d5991dfe05284cc1a7043b8a290951c6e9cbf0c7259d24a460396927f7b11c0ebd1e63fb48a659bb747e81fdae5e01dbd9488d3c082b3c0860f78a6eb0d1f352dcfb08e4f1a7679f82ea750ab9ab6ef0d979d61cf2d0c972d7fb296cf149cea70ddf60b9386650f979d9575db1f2e84df99bde1f657d996723883bf2abd53f3a76a3e50b7fd8a95cf4cd2823abe52e7a6e88b4075caaa989902e459169cd6adfc1c4c74c7964daacaf38d39c154e543e89c8131084e1749f998312162075f784cb072d3203d3f674cf4c6368bf81be256f513b16693b8e3ace856958f66e6b610db792154e50f4ce74a87d365c6443cb659a4d5a84ee0f767f4550138ae2a27b92ed2ea1ed1ab0e89f48b62449d13836745bc44c09fce94daaca7f87cd6d39c58978925eabb847d6f46170b9457584608f5acf2d90919e90390f27f18119c9db790a7796b214bd4f322ad8b5b3262ed79727a67461742f93f9651955f9fd8d4a411caffb2f33281086f926c6082c59e8dd63e8d065592167a59942d64794a685f08a61ccb9485f68460ca3409bf124c399e392cb4b288ab19a19504536b6d6f7a8f940b2077c53b23124acf278229473265915012738229ae14e3bf104cf132249544fc9c488c89e5c322aefae28eb908a39cfbc31a47b8f83b55f93e7471b35a127be6c4c09c489544bf5a16e23d91c88b252342a8be5a168977446a1f7104047d73ed4d7a92eb112e562a1766f64498882be7d919b1443d1bd63dc42d5095a75ab8f3a25f9d13dd6a592c239c2e52eaa8e856cfa865c15e5495971abae853951f42170bdf12ddcf0af65fed1ad05bacbc1fd4e0e09968ed53e5d668a96da6275af928308f538986a944f7cbba64a844fba5b8924a344e255a5912bd6fd122db1cf4bc2b17038e7b6588b56108bc1a89856f87e7bb858a6197f5bae97a9663eb1b32ab32abf4db7cb7661fd6d7ac5d7de7e8aadbff16577cf86b6ff5750b6eda1e3cfabd7dd89c4661d27061844d83ba51a999a859b63fb26e1455636acab2cbabd7ad59bf7ac3c8faf5a3a8387659971f35dbb3cab659d42ddb970ea36b839f24d86f1ea959ae898263172df9e33fb7dd712aa87a657841c7e3b0395d772bab4b85913a3a7eede018a67145eba009b8af5669c977596dc76cb1d8f91242b6166463a1d9566837158e236fda66c92a5886dbfc5905cf2ae73cf9e9cacf3a2689a2d983295aa512f2418305f8018fd846d564dfe337206c46b12bfa50fcab2d2825ec3e455a8d27d1ea39a91deda6aeab759aa2ed2653f755fa4bb1b0b5a4cdeb2a2d900da578bb979498d746ea093a48c9ab378f7adb7da3be2b5a46a9de5ceea8e1557305a352c9157cc7f5165ed1f059d4d9ef592cdb3dfdb2db33d06cf60cb67b3daf2b830fff539aca9ade98eef04b0fa50b8eedf986eda737968c8a670ea52d7baae67be98d0f1f1a4a0717269cf8d353e42e1ddc1aa5953e31d462a504ae49d742761c90748b35bc1f6972289ac7cc6248d01121384c6d97f0125dcb451eb6b6477039bec9c1fd731dbc6f76f89a5a96c2cdec2c24ede15fb22f57e30a4fc23c3adfad5d876d9997448b23507d8b2dba76daadbb911e4a1bb6634f579d9a778d1355348d62de34fc6f97770831eb5499138766ff040000ffff`,
		gas:         3000000,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
	"empty code": {
		code:        ``,
		gas:         300000,
		wantErr:     bind.ErrNoCodeAfterDeploy,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
}

func TestWaitDeployed(t *testing.T) {
	for name, test := range waitDeployedTests {
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(10000000000)},
			activeAddr: {Balance: big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18))},
		})

		// Create the transaction.
		var (
			tx  *types.Transaction
			err error
		)
		tx = types.NewContractCreation(0, big.NewInt(0), test.gas, big.NewInt(1), common.FromHex(test.code))
		tx, err = types.SignTx(tx, types.NewHubbleSigner(big.NewInt(1)), testKey)
		if err != nil {
			t.Fatalf("sign transaction failed: %v", err)
		}

		// Wait for it to get produced in the background.
		var (
			address common.Address
			packed  = make(chan struct{})
			ctx     = context.Background()
		)
		go func() {
			address, err = bind.WaitDeployed(ctx, backend, tx)
			close(packed)
		}()

		// Send and pack the transaction.
		backend.SendTransaction(ctx, tx)
		backend.Commit()

		select {
		case <-packed:
			if err != test.wantErr {
				t.Errorf("test %q: error mismatch: got %q, want %q", name, err, test.wantErr)
			}
			if address != test.wantAddress {
				t.Errorf("test %q: unexpected contract address %s", name, address.Hex())
			}
		case <-time.After(2 * time.Second):
			t.Errorf("test %q: timeout", name)
		}
	}
}
