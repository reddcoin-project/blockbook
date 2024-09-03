//go:build unittest

package reddcoin

import (
	"encoding/hex"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "RkcB2uchtVVr7XktEXCXt9yjZcU62sHaF3"},
			want:    "76a91483836ef412f917f9ec855886f3e54fb136a6a1d188ac",
			wantErr: false,
		},
		{
			name:    "P2PKH2",
			args:    args{address: "RpczZwPpxf8k3ZKzJ9pxYy15ZDLFujP6GQ"},
			want:    "76a914af8b9de50896e5dc49dbb0f907e3b5205cdebd4f88ac",
			wantErr: false,
		},
		{
			name:    "P2PKH3",
			args:    args{address: "RkocjvgtdCXyeTxb3kd6wavtfJJ1TLy4ZK"},
			want:    "76a91485ad77e21eea957d4c6cc856049beefbf554f9b788ac",
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{address: "39Z1TR3G4mid5vjGv1cjxWNPXFHh9iirt4"},
			want:    "a914563d49c30fe5675335009c8bc9541bec26bc964387",
			wantErr: false,
		},
		{
			name:    "P2SH2",
			args:    args{address: "333Q3fpyhNdN29nNeMMXxRyLgiaepSncfB"},
			want:    "a9140ed2e98a0028efe83ed03b5bcbcfafb2a646392a87",
			wantErr: false,
		},
		{
			name:    "P2SH3",
			args:    args{address: "3PUYrbw1uFZoxbXUnjkEsd8uUqqVYshbHE"},
			want:    "a914eef72e369fdd1067b0c6fac831fb46bc8482d6e887",
			wantErr: false,
		},
	}
	parser := NewReddcoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddrDescFromAddress_Testnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "mwPVTmvuJuo5R8kPXV5T3BGHXPv2ME8QnW"},
			want:    "76a914ae18b3c7272cf18867787589850289e8f0bb5a4388ac",
			wantErr: false,
		},
		{
			name:    "P2PKH2",
			args:    args{address: "mnQJrRvmQhDhNfkMo8HmEkdapKuSX7MstF"},
			want:    "76a9144b8721406854be352f08505851b9496a9deb34eb88ac",
			wantErr: false,
		},
		{
			name:    "P2PKH3",
			args:    args{address: "mxgyZ2BhLJb9MaVoyRbScQEpruQj5tf3JD"},
			want:    "76a914bc5f591dfe6bf95832dbc3007a9ab3ff4074762488ac",
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{address: "2NEK99vk81VsjdhSKG2QjPFHeFwYzzGerhH"},
			want:    "a914e7184532fc62256e41825d1e11a3bd4333d5f12287",
			wantErr: false,
		},
		{
			name:    "P2SH2",
			args:    args{address: "2NGEnNnknu9gCVzK2cvGqyWVCaEmULVYtXf"},
			want:    "a914fc3583631c775b866cd3e1dea470330b4bf3fc7287",
			wantErr: false,
		},
		{
			name:    "P2SH3",
			args:    args{address: "2MsVhxVMoywgbCSFaEjDgQwMy9HpZDeab2M"},
			want:    "a91402bd2bbdd42b26d5d26f8f5e775a10c29a0a474987",
			wantErr: false,
		},
	}
	parser := NewReddcoinParser(GetChainParams("test"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc_Mainnet(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{script: "76a9140af8f2eaeacfc2aabb090c70473fcbee18763c7988ac"},
			want:    []string{"RZcp3pFoinWtdZhirH8kEfamuwRrXNpNW6"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH2",
			args:    args{script: "76a91471733496fc8f9e9ae6564ce561889bfb9c9ee30488ac"},
			want:    []string{"RixfRv2N9xFqpyFWJ9z7i8zHSxHS6Tx7G1"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH3",
			args:    args{script: "76a914ee1f51e41b1359b1f3d208f4f7df6627f4e1eb0788ac"},
			want:    []string{"RvKsQRZCAtU9nHCEfKaDfq6BVvZb2uRj1m"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{script: "a9146c9e50891a1e2d8be4391107db5768e5165779cc87"},
			want:    []string{"3BbLZ9n2eW7UcVWFa8CXShVMPbpZ8K9gHN"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH2",
			args:    args{script: "a91482efe2d1a5fc1c07447b1d517ab8dd5686ae54d187"},
			want:    []string{"3DdM945sQDmWXU6LuaRkdoRcX9P5KR32HN"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH3",
			args:    args{script: "a914a0934dfdd9b857b0236d6a021e010d2a64d1a60987"},
			want:    []string{"3GL4W9GrLHmi9J5Ey7TvNnAaK1EL1HK5kJ"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "OP_RETURN ascii",
			args:    args{script: "6a0461686f6a"},
			want:    []string{"OP_RETURN (ahoj)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN hex",
			args:    args{script: "6a072020f1686f6a20"},
			want:    []string{"OP_RETURN 2020f1686f6a20"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewReddcoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc_Testnet(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{script: "76a914c54ce2d90520a9ca79a7d443aa34b81a8216095588ac"},
			want:    []string{"myWBWz9hSXkAnHHqG1jwVktETXA57zDevH"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH2",
			args:    args{script: "76a9144b7fbea726a2450cdaae9170fa95ad3a94d5589588ac"},
			want:    []string{"mnQA1JRqmqAoQq1hyDiR3C3anFgtdTgj1r"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH3",
			args:    args{script: "76a914d70c34c0e5023de19377998a9af44b4154442c2d88ac"},
			want:    []string{"n182CPaQYGRuKEzUhGjhX5q4ysCnjk3ZUj"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH1",
			args:    args{script: "a914af43b48e9dc2e585caeb87837eb29fbdce307c8887"},
			want:    []string{"2N9DwM6pFmg62T8PicvbDELcXVXe7V3pNj6"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH2",
			args:    args{script: "a9146dca0704700de2fbc00730f214235b70afa840be87"},
			want:    []string{"2N3FjfJbFvzcMBHjiVvJH8gDXijC7iGnxHi"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2SH3",
			args:    args{script: "a9144cfbe3e9fbbb971ceeecbb613b1e806c4c76d9cf87"},
			want:    []string{"2MzGH6WNNzwuFPTgb9r8YzRnwymko3iepWs"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "OP_RETURN ascii",
			args:    args{script: "6a0c48656c6c6f20746865726521"},
			want:    []string{"OP_RETURN (Hello there!)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN hex",
			args:    args{script: "6a072020f1686f6a20"},
			want:    []string{"OP_RETURN 2020f1686f6a20"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewReddcoinParser(GetChainParams("test"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	testTx1       bchain.Tx
	testTxPacked1 = "0a2084eaf9490120f9aa35afee1916c0f0f76100f4e8fc1aaac12a37395a6cac3aee12b50102000000013ead574198141c42cf318a9e9322d3e5e8a4ba01e95e1865516bd8a0858fa1c00100000049483045022100ec316d44d2d01bce93ea5dc599348b079da1cda0b5fe0b862917461e6b09453e022034ec51ae397ce41a2975fd4256c6bd703edb4e7ce3cc011529d96af264218e9401ffffffff02000000000000000000fb948935b7250000232102adb55bb1659cb210938d34e9341fbee1812f4d77f6bebcb994296f34cb70ae4eac00000000435eea5b18c3bca9df0528a0cb980132751220c0a18f85a0d86b5165185ee901baa4e8e5d322939e8a31cf421c14984157ad3e18012249483045022100ec316d44d2d01bce93ea5dc599348b079da1cda0b5fe0b862917461e6b09453e022034ec51ae397ce41a2975fd4256c6bd703edb4e7ce3cc011529d96af264218e940128ffffffff0f3a0222003a530a0625b7358994fb10011a232102adb55bb1659cb210938d34e9341fbee1812f4d77f6bebcb994296f34cb70ae4eac222252704c6b596364514c5663674e456f4d346f417a4268787039386e453643517941594002"

	testTx2       bchain.Tx
	testTxPacked2 = "0a20be9f06f31023ecaaa68d69d8331ad980726b7dc908140d6bc575004845f1e86f128c02020000000101403e0ee807cdbad22d0b199b5ca6c06e4490b5e36817227302ecfc1ecee756010000004847304402201630f05fff280bf50b5f05341a2577886e13bb42ff361f268a1c8e12bb40a4ce02206f658f0193521882a91a8ecada33b42ca1af3d0b14b3a4a1b009db2c4913f72501ffffffff04000000000000000000009f9cef8a850000232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305acfe4ea9ef8a850000232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305aced32df2a02000000232103c8fc5c87f00bcc32b5ce5c036957f8befeff05bf4d88d2dcde720249f78d9313ac000000007bc5825f18fb8a8bfc0528e0cfd5013274122056e7ce1efcec0273221768e3b590446ec0a65c9b190b2dd2bacd07e80e3e40011801224847304402201630f05fff280bf50b5f05341a2577886e13bb42ff361f268a1c8e12bb40a4ce02206f658f0193521882a91a8ecada33b42ca1af3d0b14b3a4a1b009db2c4913f7250128ffffffff0f3a0222003a510a0408c0799c10011a232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305ac22225263334a4561457a74544e70344665565172326373376439656b416538647959746e3a530a06858aefa94efe10021a232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305ac22225263334a4561457a74544e70344665565172326373376439656b416538647959746e3a520a05022adf32ed10031a232103c8fc5c87f00bcc32b5ce5c036957f8befeff05bf4d88d2dcde720249f78d9313ac2222526d687a6a324770745a786b4b424d7162554c36566a466358386e70446e654158524002"

	testTx1_Testnet       bchain.Tx
	testTxPacked1_Testnet = "0a2044aa5aab3c39fe571761fbb610839eb6580ce498344c763345c096695192b6aa12e0010200000001ab177cbfb6535547e34e3d0bb188da55bf7e98e357a055411c4b86a8310ac816010000004847304402202dee7d44acf8977fa13a9e6c171ae961ea0517b029a6f9747aee16ad3da7c9e3022072c4effff83aff224e65c04547d17f059916d2092d0bb9b3c338041054599c7801ffffffff0300000000000000000015e50b72e41b00002321029d612e6edb2cd614635f0b8c2a610913bf27a7dd77da3ab588a55e0b32254116acc9b6a85400000000232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac00000000fd85746218fd8bd2930628f093093274122016c80a31a8864b1c4155a057e3987ebf55da88b10b3d4ee3475553b6bf7c17ab1801224847304402202dee7d44acf8977fa13a9e6c171ae961ea0517b029a6f9747aee16ad3da7c9e3022072c4effff83aff224e65c04547d17f059916d2092d0bb9b3c338041054599c780128ffffffff0f3a0222003a530a061be4720be51510011a2321029d612e6edb2cd614635f0b8c2a610913bf27a7dd77da3ab588a55e0b32254116ac22226d67714c795056327657395548665748525048563934343174486846595969465a373a510a0454a8b6c910021a232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac22226d714c7741393334776a636d39777043566a7a7664557562316f77375632436748384002"

	testTx2_Testnet       bchain.Tx
	testTxPacked2_Testnet = "0a20288553c5fe45d3bbfdfaed3d45ccf5f2a770806b4fc75a8cb8805c5db7ffbfa1128c0202000000012d946a8737941c28c59449c2b0848b38d504ac4d34d77a1d4bda16b23c155530010000004847304402205f78d460e134f9208eae78935b156cffa226416648902953f508c1e4f0c40abe022001bfca53f6c0b0b788b831ff1fcd1b66a0c489922e87684b9b71d46f216e7ce601ffffffff04000000000000000000400ff7a9c80000002321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac866f01aac80000002321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac7b0ab75500000000232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac000000009c3a1164189cf5c4a00628a0c21e327212203055153cb216da4b1d7ad7344dac04d5388b84b0c24994c5281c9437876a942d224847304402205f78d460e134f9208eae78935b156cffa226416648902953f508c1e4f0c40abe022001bfca53f6c0b0b788b831ff1fcd1b66a0c489922e87684b9b71d46f216e7ce60128ffffffff0f3a0222003a520a05c8a9f70f4010011a2321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac22226d7341476f456e336863746641755868664c425952676f41724a5a33446f725345443a520a05c8aa016f8610021a2321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac22226d7341476f456e336863746641755868664c425952676f41724a5a33446f725345443a510a0455b70a7b10021a232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac22226d714c7741393334776a636d39777043566a7a7664557562316f77375632436748384002"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "02000000013ead574198141c42cf318a9e9322d3e5e8a4ba01e95e1865516bd8a0858fa1c00100000049483045022100ec316d44d2d01bce93ea5dc599348b079da1cda0b5fe0b862917461e6b09453e022034ec51ae397ce41a2975fd4256c6bd703edb4e7ce3cc011529d96af264218e9401ffffffff02000000000000000000fb948935b7250000232102adb55bb1659cb210938d34e9341fbee1812f4d77f6bebcb994296f34cb70ae4eac00000000435eea5b",
		Blocktime: 1542086211,
		Time:      1542086211,
		Txid:      "84eaf9490120f9aa35afee1916c0f0f76100f4e8fc1aaac12a37395a6cac3aee",
		LockTime:  0,
		Version:   2,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100ec316d44d2d01bce93ea5dc599348b079da1cda0b5fe0b862917461e6b09453e022034ec51ae397ce41a2975fd4256c6bd703edb4e7ce3cc011529d96af264218e9401",
				},
				Txid:     "c0a18f85a0d86b5165185ee901baa4e8e5d322939e8a31cf421c14984157ad3e",
				Vout:     1,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
					Addresses: []string{
						"",
					},
				},
			},
			{
				ValueSat: *big.NewInt(41468807451899),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2102adb55bb1659cb210938d34e9341fbee1812f4d77f6bebcb994296f34cb70ae4eac",
					Addresses: []string{
						"RpLkYcdQLVcgNEoM4oAzBhxp98nE6CQyAY",
					},
				},
			},
		},
	}

	testTx2 = bchain.Tx{
		Hex:       "020000000101403e0ee807cdbad22d0b199b5ca6c06e4490b5e36817227302ecfc1ecee756010000004847304402201630f05fff280bf50b5f05341a2577886e13bb42ff361f268a1c8e12bb40a4ce02206f658f0193521882a91a8ecada33b42ca1af3d0b14b3a4a1b009db2c4913f72501ffffffff04000000000000000000009f9cef8a850000232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305acfe4ea9ef8a850000232103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305aced32df2a02000000232103c8fc5c87f00bcc32b5ce5c036957f8befeff05bf4d88d2dcde720249f78d9313ac000000007bc5825f",
		Blocktime: 1602405755,
		Time:      1602405755,
		Txid:      "be9f06f31023ecaaa68d69d8331ad980726b7dc908140d6bc575004845f1e86f",
		LockTime:  0,
		Version:   2,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402201630f05fff280bf50b5f05341a2577886e13bb42ff361f268a1c8e12bb40a4ce02206f658f0193521882a91a8ecada33b42ca1af3d0b14b3a4a1b009db2c4913f72501",
				},
				Txid:     "56e7ce1efcec0273221768e3b590446ec0a65c9b190b2dd2bacd07e80e3e4001",
				Vout:     1,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
					Addresses: []string{
						"",
					},
				},
			},
			{
				ValueSat: *big.NewInt(146831772),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305ac",
					Addresses: []string{
						"Rc3JEaEztTNp4FeVQr2cs7d9ekAe8dyYtn",
					},
				},
			},
			{
				ValueSat: *big.NewInt(146831772831486),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103395b5ca66998efe7821a97a742d57c7e508bea6124b1dbd69e9780242b4c8305ac",
					Addresses: []string{
						"Rc3JEaEztTNp4FeVQr2cs7d9ekAe8dyYtn",
					},
				},
			},
			{
				ValueSat: *big.NewInt(9309205229),
				N:        3,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103c8fc5c87f00bcc32b5ce5c036957f8befeff05bf4d88d2dcde720249f78d9313ac",
					Addresses: []string{
						"Rmhzj2GptZxkKBMqbUL6VjFcX8npDneAXR",
					},
				},
			},
		},
	}

	testTx1_Testnet = bchain.Tx{
		Hex:       "0200000001ab177cbfb6535547e34e3d0bb188da55bf7e98e357a055411c4b86a8310ac816010000004847304402202dee7d44acf8977fa13a9e6c171ae961ea0517b029a6f9747aee16ad3da7c9e3022072c4effff83aff224e65c04547d17f059916d2092d0bb9b3c338041054599c7801ffffffff0300000000000000000015e50b72e41b00002321029d612e6edb2cd614635f0b8c2a610913bf27a7dd77da3ab588a55e0b32254116acc9b6a85400000000232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac00000000fd857462",
		Blocktime: 1651803645,
		Time:      1651803645,
		Txid:      "44aa5aab3c39fe571761fbb610839eb6580ce498344c763345c096695192b6aa",
		LockTime:  0,
		Version:   2,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402202dee7d44acf8977fa13a9e6c171ae961ea0517b029a6f9747aee16ad3da7c9e3022072c4effff83aff224e65c04547d17f059916d2092d0bb9b3c338041054599c7801",
				},
				Txid:     "16c80a31a8864b1c4155a057e3987ebf55da88b10b3d4ee3475553b6bf7c17ab",
				Vout:     1,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
					Addresses: []string{
						"",
					},
				},
			},
			{
				ValueSat: *big.NewInt(30667979875605),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21029d612e6edb2cd614635f0b8c2a610913bf27a7dd77da3ab588a55e0b32254116ac",
					Addresses: []string{
						"mgqLyPV2vW9UHfWHRPHV9441tHhFYYiFZ7",
					},
				},
			},
			{
				ValueSat: *big.NewInt(1420342985),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac",
					Addresses: []string{
						"mqLwA934wjcm9wpCVjzvdUub1ow7V2CgH8",
					},
				},
			},
		},
	}

	testTx2_Testnet = bchain.Tx{
		Hex:       "02000000012d946a8737941c28c59449c2b0848b38d504ac4d34d77a1d4bda16b23c155530010000004847304402205f78d460e134f9208eae78935b156cffa226416648902953f508c1e4f0c40abe022001bfca53f6c0b0b788b831ff1fcd1b66a0c489922e87684b9b71d46f216e7ce601ffffffff04000000000000000000400ff7a9c80000002321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac866f01aac80000002321036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac7b0ab75500000000232103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac000000009c3a1164",
		Blocktime: 1678850716,
		Time:      1678850716,
		Txid:      "288553c5fe45d3bbfdfaed3d45ccf5f2a770806b4fc75a8cb8805c5db7ffbfa1",
		LockTime:  0,
		Version:   2,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402205f78d460e134f9208eae78935b156cffa226416648902953f508c1e4f0c40abe022001bfca53f6c0b0b788b831ff1fcd1b66a0c489922e87684b9b71d46f216e7ce601",
				},
				Txid:     "3055153cb216da4b1d7ad7344dac04d5388b84b0c24994c5281c9437876a942d",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "",
					Addresses: []string{
						"",
					},
				},
			},
			{
				ValueSat: *big.NewInt(861845000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac",
					Addresses: []string{
						"msAGoEn3hctfAuXhfLBYRgoArJZ3DorSED",
					},
				},
			},
			{
				ValueSat: *big.NewInt(861845680006),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21036c27da387708f00ac6a1372b5e3082332d2bd4157334920e2b678c8a2225a96eac",
					Addresses: []string{
						"msAGoEn3hctfAuXhfLBYRgoArJZ3DorSED",
					},
				},
			},
			{
				ValueSat: *big.NewInt(1438059131),
				N:        2,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103d4b22ae69b0ff7554f4c343cd213d00fd5131466cc21d8ebfab97c52ec9a00c9ac",
					Addresses: []string{
						"mqLwA934wjcm9wpCVjzvdUub1ow7V2CgH8",
					},
				},
			},
		},
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *ReddcoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "reddcoin-1",
			args: args{
				tx:        testTx1,
				height:    2500000,
				blockTime: 1542086211,
				parser:    NewReddcoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
		{
			name: "reddcoin-2",
			args: args{
				tx:        testTx2,
				height:    3500000,
				blockTime: 1602405755,
				parser:    NewReddcoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_PackTx_Testnet(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *ReddcoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "reddcoin-testnet-1",
			args: args{
				tx:        testTx1_Testnet,
				height:    150000,
				blockTime: 1651803645,
				parser:    NewReddcoinParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    testTxPacked1_Testnet,
			wantErr: false,
		},
		{
			name: "reddcoin-testnet-2",
			args: args{
				tx:        testTx2_Testnet,
				height:    500000,
				blockTime: 1678850716,
				parser:    NewReddcoinParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    testTxPacked2_Testnet,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *ReddcoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "reddcoin-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewReddcoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   2500000,
			wantErr: false,
		},
		{
			name: "reddcoin-2",
			args: args{
				packedTx: testTxPacked2,
				parser:   NewReddcoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx2,
			want1:   3500000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want = %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want = %v", got1, tt.want1)
			}
		})
	}
}

func Test_UnpackTx_Testnet(t *testing.T) {
	type args struct {
		packedTx string
		parser   *ReddcoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "reddcoin-testnet-1",
			args: args{
				packedTx: testTxPacked1_Testnet,
				parser:   NewReddcoinParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    &testTx1_Testnet,
			want1:   150000,
			wantErr: false,
		},
		{
			name: "reddcoin-testnet-2",
			args: args{
				packedTx: testTxPacked2_Testnet,
				parser:   NewReddcoinParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    &testTx2_Testnet,
			want1:   500000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want = %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want = %v", got1, tt.want1)
			}
		})
	}
}
