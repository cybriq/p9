package forkhash

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDivHash(t *testing.T) {

	expected := []string{
		"84a5f89d2ef27b86c83bdd76e6dfd4cf7d35f4179ec89719559413896c95459a",
		"0789b6358b631bcc3317d24520aaf83062a5fb70a83b6404e276fc015172dbb4",
		"cb3acf43707850372cf2c60ccf99cf64fc21a620b368c91361426181ab527043",
		"d83026157f2e25206fc5f09c5204aa67e3e589cc9aa7e9a691c60386c8e5ec3e",
		"7e77d5bb63651ad674627b2d1deddd99a77b9bd8e6119679173b41dd823fce25",
		"1082113ce80e7366c0acdc752527a974625a9db74d5758ecf12d5d91b57e366c",
		"61276256fbdf8a397fc64237843813883b3d0bb7e9ed8d43eb64e8ba602d83fd",
		"99c19dfe054c2845644628ef5d3a44aea911cf37f114ad3f4e6176739c19e35b",
		"0df6361b8d345ad1ea739b2f7b58b00e362e342ab9251ad66f69c30aab4de2a7",
		"42622f9e8814672f8ccee3cbbf43840b8ecafbfe378d017ce87b34fa70f61377",
		"37ad9ed372439be9ed3708a6543f6ff318a0ed6f52f1d138cfe7fed5119329ac",
		"a7dbca64c01da3cf87300f8ec248a8b209ebc749f3eb2c99f321a459d3f477c6",
		"343d23784e7318c70bf32e0bfbac999ec790340763e6fb504268355beee09d8d",
		"ffef48a3df0c032233f674c8c8b83d3e7e573ff99cf6f736143d2997bbdae4d9",
		"d98bd9e2340ff87528ecb47cefdf9c3e6fae2a660167fc6438e6d813d8b24fa3",
		"e5dd1af5bc41304d18d13297c10fded398d5cdeadcb0f96ba68ee3ae43f99cb2",
		"a9e24d2589bce260a3ddf84949ffb02e8a6dc26cdc693411693869bb81def946",
		"b6beec01cd58e8055164c23ff9fc468c583379487d2e18616f68a848a0021109",
		"063d94131d3e86519aeac759cc2ae6f9b7ffd468edeecc82660d34f2a7ace019",
		"ae98f1ac67a5213e1f7199eeedd7f0917797f63dda0d36933da69c10c3990193",
		"266c6c0d9cfb73a41d94cb0ac5a48fbef5d83ec816917a1c307d9115f0511ebe",
		"705f1187c42212e3683d420baa8a0c1c4678788f8024a10e0f3f9077e0a27f92",
		"45941452fec84e10814f59df23512b7fcb012979856ec34af6e93fe439f38298",
		"26e5093d569668ca9148537ea5a946e3edf43adc79557aa7e7d27c9ed3149095",
		"253cdb3556fefa9c21462c9ac60e279780ca69297a9b63d8018515f6f273a103",
		"e5623df4d1cbc89d0c92f647187eadf85d4d6b36bd87eba89a9ee0e754306b66",
		"0b08f68d2efec39760e2b7f180c3ea707e758c9f9ad5027b8db0f8deac76a87f",
		"a4ce35b109561ac93561b43ab35be401cf6dc471f2bdba355a30b6018be564a4",
		"f1e340799d0b296335f31db7288cbbd6f04112da09f6cc4bcb9ec2d94f5e4ccc",
		"91b41609f4492fa3e676be770a60cde43d91dc3d0a61e36954d942facda778c2",
		"379f0c7ade26503bc12fff887e9a643a9744fa478ad0af6bc06f4abf8ea9dab2",
		"ba9d0bc45add56a9818e27a5ad3cf4cc1afdfd9243656ee5808b67e04d907535",
	}
	o := "\n\texpected := []string{\n"
	// Standard block size is 106 bytes for a validator block. If another hash is
	// added to the block for proposals or other linked content such as a merkle
	// root for a transaction payload, the result is about 30% more processing time
	// per operation.
	empty := Blake3([]byte{})
	empty = append(empty, empty...)
	empty = append(empty, empty[:Len+10]...)
	for i := 0; i < 32; i++ {
		h := DivHash4(empty)
		expect, err := hex.DecodeString(expected[i])
		if err != nil {
			t.Fatalf("error decoding hex string '%s': %v", expected[i], err)
		}
		if string(h) != string(expect) {
			t.Fatalf("")
		}
		o += fmt.Sprintf("\t\t\"%s\",\n", hex.EncodeToString(h))
		empty = h
		empty = append(empty, empty...)
		empty = append(empty, empty[:Len+10]...)
	}
	o += "\t}\n"
	t.Log("generated:\n", o)
}

func benchmarkDivHash(reps int, b *testing.B) {

	// Maximum block length is 138 bytes for not yet implemented proposal blocks but
	// this will test the block length defined in pkg/blocks/block.go for a
	// validator token. When proposal blocks are implemented they have a longer time
	// per operation.
	empty := Blake3([]byte{})
	empty = append(empty, empty...)
	empty = append(empty, empty[:32]...)
	empty = append(empty, empty[:10]...)

	for i := 0; i < b.N; i++ {

		hash := DivHash(empty, reps)
		empty = hash
		empty = append(empty, empty...)
		empty = append(empty, empty[:32]...)
		empty = append(empty, empty[:10]...)
	}
}

// Operation time for DivHash is proportional to input size. Thus, for a given
// data size a repetition count is chosen based on the benchmark below:
//
// $ go test -bench=. ./pkg/proof/.
// goos: linux
// goarch: amd64
// pkg: github.com/cybriq/kismet/pkg/proof
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkDivHash1-16               10000            226545 ns/op
// BenchmarkDivHash2-16                 763           1815126 ns/op
// BenchmarkDivHash3-16                  82          15908234 ns/op
// BenchmarkDivHash4-16                   7         148824596 ns/op
// BenchmarkDivHash5-16                   1        1325919974 ns/op
// PASS
// ok      github.com/cybriq/kismet/pkg/proof      14.759s
//
// In order to render a reasonably integer number per second operations, 4
// repetitions will be selected.

func BenchmarkDivHash1(b *testing.B) { benchmarkDivHash(1, b) }
func BenchmarkDivHash2(b *testing.B) { benchmarkDivHash(2, b) }
func BenchmarkDivHash3(b *testing.B) { benchmarkDivHash(3, b) }
func BenchmarkDivHash4(b *testing.B) { benchmarkDivHash(4, b) }
func BenchmarkDivHash5(b *testing.B) { benchmarkDivHash(5, b) }
