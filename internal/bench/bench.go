//go:generate go run github.com/mailru/easyjson/easyjson -all bench.go
//go:generate go run github.com/tinylib/msgp -file bench.go
package bench

import (
	"bytes"
	stdjson "encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	ddzstd "github.com/DataDog/zstd"
	"github.com/andybalholm/brotli"
	"github.com/bytedance/sonic"
	"github.com/clbanning/mxj/v2"
	"github.com/dsnet/compress/bzip2"
	"github.com/fxamacker/cbor/v2"
	goccyjson "github.com/goccy/go-json"
	goccyyaml "github.com/goccy/go-yaml"
	golsnappy "github.com/golang/snappy"
	"github.com/hamba/avro/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/compress/zstd"
	"github.com/klauspost/pgzip"
	"github.com/linkedin/goavro/v2"
	pelletiertoml "github.com/pelletier/go-toml/v2"
	"github.com/pierrec/lz4/v4"
	"github.com/segmentio/encoding/json"
	shamatonmsgpack "github.com/shamaton/msgpack/v2"
	"github.com/ugorji/go/codec"
	"github.com/ulikunitz/xz"
	"github.com/valyala/gozstd"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/zeebo/bencode"
	"go.mongodb.org/mongo-driver/bson"
	stdyaml "gopkg.in/yaml.v3"
)

// Data Structures

type Book struct {
	Rank             int    `json:"rank" yaml:"rank" toml:"rank" bson:"rank" avro:"rank" xml:"rank"`
	Book             string `json:"book" yaml:"book" toml:"book" bson:"book" avro:"book" xml:"book"`
	Author           string `json:"author" yaml:"author" toml:"author" bson:"author" avro:"author" xml:"author"`
	OriginalLanguage string `json:"original_language" yaml:"original_language" toml:"original_language" bson:"original_language" avro:"original_language" xml:"original_language"`
	FirstPublished   string `json:"first_published" yaml:"first_published" toml:"first_published" bson:"first_published" avro:"first_published" xml:"first_published"`
	ApproximateSales string `json:"approximate_sales" yaml:"approximate_sales" toml:"approximate_sales" bson:"approximate_sales" avro:"approximate_sales" xml:"approximate_sales"`
	Genre            string `json:"genre" yaml:"genre" toml:"genre" bson:"genre" avro:"genre" xml:"genre"`
}

func (b Book) Size() int {
	return 8 + len(b.Book) + len(b.Author) + len(b.OriginalLanguage) + len(b.FirstPublished) + len(b.ApproximateSales) + len(b.Genre)
}

type Payload struct {
	Total    int    `json:"Total" yaml:"Total" toml:"Total" bson:"Total" avro:"Total" xml:"Total"`
	Elements []Book `json:"Elements" yaml:"Elements" toml:"Elements" bson:"Elements" avro:"Elements" xml:"Elements"`
}

var avroSchema avro.Schema
var goavroCodec *goavro.Codec

func init() {
	var err error
	avroSchema, err = avro.Parse(`{
		"type": "record",
		"name": "Payload",
		"fields": [
			{"name": "Total", "type": "int"},
			{
				"name": "Elements",
				"type": {
					"type": "array",
					"items": {
						"type": "record",
						"name": "Book",
						"fields": [
							{"name": "rank", "type": "int"},
							{"name": "book", "type": "string"},
							{"name": "author", "type": "string"},
							{"name": "original_language", "type": "string"},
							{"name": "first_published", "type": "string"},
							{"name": "approximate_sales", "type": "string"},
							{"name": "genre", "type": "string"}
						]
					}
				}
			}
		]
	}`)
	if err != nil {
		panic(err)
	}

	goavroCodec, err = goavro.NewCodec(`{
		"type": "record",
		"name": "Payload",
		"fields": [
			{"name": "Total", "type": "int"},
			{
				"name": "Elements",
				"type": {
					"type": "array",
					"items": {
						"type": "record",
						"name": "Book",
						"fields": [
							{"name": "rank", "type": "int"},
							{"name": "book", "type": "string"},
							{"name": "author", "type": "string"},
							{"name": "original_language", "type": "string"},
							{"name": "first_published", "type": "string"},
							{"name": "approximate_sales", "type": "string"},
							{"name": "genre", "type": "string"}
						]
					}
				}
			}
		]
	}`)
	if err != nil {
		panic(err)
	}
}

// Interfaces

type Marshaler interface {
	Name() string
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type Compressor interface {
	Name() string
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

// Implementations: Marshaler

type JSONMarshaler struct{}

func (m JSONMarshaler) Name() string { return "JSON" }
func (m JSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return stdjson.Marshal(v)
}
func (m JSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return stdjson.Unmarshal(data, v)
}

type MsgPackMarshaler struct{}

func (m MsgPackMarshaler) Name() string { return "MsgPack" }
func (m MsgPackMarshaler) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}
func (m MsgPackMarshaler) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

type CBORMarshaler struct{}

func (m CBORMarshaler) Name() string { return "CBOR" }
func (m CBORMarshaler) Marshal(v interface{}) ([]byte, error) {
	return cbor.Marshal(v)
}
func (m CBORMarshaler) Unmarshal(data []byte, v interface{}) error {
	return cbor.Unmarshal(data, v)
}

type CBORPreferredMarshaler struct {
	em cbor.EncMode
	dm cbor.DecMode
}

func NewCBORPreferredMarshaler() CBORPreferredMarshaler {
	em, _ := cbor.CanonicalEncOptions().EncMode()
	dm, _ := cbor.DecOptions{}.DecMode()
	return CBORPreferredMarshaler{em: em, dm: dm}
}
func (m CBORPreferredMarshaler) Name() string { return "CBORCanonical" }
func (m CBORPreferredMarshaler) Marshal(v interface{}) ([]byte, error) {
	return m.em.Marshal(v)
}
func (m CBORPreferredMarshaler) Unmarshal(data []byte, v interface{}) error {
	return m.dm.Unmarshal(data, v)
}

type YAMLMarshaler struct{}

func (m YAMLMarshaler) Name() string { return "YAML" }
func (m YAMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return stdyaml.Marshal(v)
}
func (m YAMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	return stdyaml.Unmarshal(data, v)
}

type TOMLMarshaler struct{}

func (m TOMLMarshaler) Name() string { return "TOML" }
func (m TOMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return pelletiertoml.Marshal(v)
}
func (m TOMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	return pelletiertoml.Unmarshal(data, v)
}

type BurntSushiTOMLMarshaler struct{}

func (m BurntSushiTOMLMarshaler) Name() string { return "BurntSushiTOML" }
func (m BurntSushiTOMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (m BurntSushiTOMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

type XMLMarshaler struct{}

func (m XMLMarshaler) Name() string { return "XML" }
func (m XMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
func (m XMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

type JsonIterMarshaler struct{}

func (m JsonIterMarshaler) Name() string { return "JsonIter" }
func (m JsonIterMarshaler) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}
func (m JsonIterMarshaler) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

type JsonIterFastestMarshaler struct {
	api jsoniter.API
}

func NewJsonIterFastestMarshaler() JsonIterFastestMarshaler {
	return JsonIterFastestMarshaler{api: jsoniter.ConfigFastest}
}
func (m JsonIterFastestMarshaler) Name() string { return "JsonIterFastest" }
func (m JsonIterFastestMarshaler) Marshal(v interface{}) ([]byte, error) {
	return m.api.Marshal(v)
}
func (m JsonIterFastestMarshaler) Unmarshal(data []byte, v interface{}) error {
	return m.api.Unmarshal(data, v)
}

type SonicMarshaler struct{}

func (m SonicMarshaler) Name() string { return "Sonic" }
func (m SonicMarshaler) Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}
func (m SonicMarshaler) Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}

type GoccyJSONMarshaler struct{}

func (m GoccyJSONMarshaler) Name() string { return "GoccyJSON" }
func (m GoccyJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return goccyjson.Marshal(v)
}
func (m GoccyJSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return goccyjson.Unmarshal(data, v)
}

type SegmentioJSONMarshaler struct{}

func (m SegmentioJSONMarshaler) Name() string { return "SegmentioJSON" }
func (m SegmentioJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (m SegmentioJSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type GoccyYAMLMarshaler struct{}

func (m GoccyYAMLMarshaler) Name() string { return "GoccyYAML" }
func (m GoccyYAMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return goccyyaml.Marshal(v)
}
func (m GoccyYAMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	return goccyyaml.Unmarshal(data, v)
}

type BSONMarshaler struct{}

func (m BSONMarshaler) Name() string { return "BSON" }
func (m BSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return bson.Marshal(v)
}
func (m BSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return bson.Unmarshal(data, v)
}

type MXJXMLMarshaler struct{}

func (m MXJXMLMarshaler) Name() string { return "MXJXML" }
func (m MXJXMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	// MXJ works with maps
	data, err := stdjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	mv, err := mxj.NewMapJson(data)
	if err != nil {
		return nil, err
	}
	return mv.Xml()
}
func (m MXJXMLMarshaler) Unmarshal(data []byte, v interface{}) error {
	mv, err := mxj.NewMapXml(data)
	if err != nil {
		return err
	}
	jb, err := mv.Json()
	if err != nil {
		return err
	}
	return stdjson.Unmarshal(jb, v)
}

type AvroMarshaler struct{}

func (m AvroMarshaler) Name() string { return "Avro" }
func (m AvroMarshaler) Marshal(v interface{}) ([]byte, error) {
	return avro.Marshal(avroSchema, v)
}
func (m AvroMarshaler) Unmarshal(data []byte, v interface{}) error {
	return avro.Unmarshal(avroSchema, data, v)
}

type GoAvroMarshaler struct{}

func (m GoAvroMarshaler) Name() string { return "GoAvro" }
func (m GoAvroMarshaler) Marshal(v interface{}) ([]byte, error) {
	// goavro requires map[string]interface{} for generic marshaling if not using specific types
	// This might be slow but let's see.
	// Actually, we can convert Payload to map using JSON as a shortcut for benchmark purposes
	// or manually.
	data, err := stdjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m2 map[string]interface{}
	if err := stdjson.Unmarshal(data, &m2); err != nil {
		return nil, err
	}
	return goavroCodec.BinaryFromNative(nil, m2)
}
func (m GoAvroMarshaler) Unmarshal(data []byte, v interface{}) error {
	native, _, err := goavroCodec.NativeFromBinary(data)
	if err != nil {
		return err
	}
	// Convert back to struct
	jb, err := stdjson.Marshal(native)
	if err != nil {
		return err
	}
	return stdjson.Unmarshal(jb, v)
}

type BencodeMarshaler struct{}

func (m BencodeMarshaler) Name() string { return "Bencode" }
func (m BencodeMarshaler) Marshal(v interface{}) ([]byte, error) {
	return bencode.EncodeBytes(v)
}
func (m BencodeMarshaler) Unmarshal(data []byte, v interface{}) error {
	return bencode.DecodeBytes(data, v)
}

type EasyJSONMarshaler struct{}

func (m EasyJSONMarshaler) Name() string { return "EasyJSON" }
func (m EasyJSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return v.(Payload).MarshalJSON()
}
func (m EasyJSONMarshaler) Unmarshal(data []byte, v interface{}) error {
	return v.(*Payload).UnmarshalJSON(data)
}

type MsgpMarshaler struct{}

func (m MsgpMarshaler) Name() string { return "Msgp" }
func (m MsgpMarshaler) Marshal(v interface{}) ([]byte, error) {
	p := v.(Payload)
	return p.MarshalMsg(nil)
}
func (m MsgpMarshaler) Unmarshal(data []byte, v interface{}) error {
	p := v.(*Payload)
	_, err := p.UnmarshalMsg(data)
	return err
}

type ShamatonMsgPackMarshaler struct{}

func (m ShamatonMsgPackMarshaler) Name() string { return "ShamatonMsgPack" }
func (m ShamatonMsgPackMarshaler) Marshal(v interface{}) ([]byte, error) {
	return shamatonmsgpack.Marshal(v)
}
func (m ShamatonMsgPackMarshaler) Unmarshal(data []byte, v interface{}) error {
	return shamatonmsgpack.Unmarshal(data, v)
}

type UgorjiMsgPackMarshaler struct {
	h *codec.MsgpackHandle
}

func NewUgorjiMsgPackMarshaler() UgorjiMsgPackMarshaler {
	return UgorjiMsgPackMarshaler{h: &codec.MsgpackHandle{}}
}
func (m UgorjiMsgPackMarshaler) Name() string { return "UgorjiMsgPack" }
func (m UgorjiMsgPackMarshaler) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, m.h).Encode(v)
	return b, err
}
func (m UgorjiMsgPackMarshaler) Unmarshal(data []byte, v interface{}) error {
	return codec.NewDecoderBytes(data, m.h).Decode(v)
}

type UgorjiCBORMarshaler struct {
	h *codec.CborHandle
}

func NewUgorjiCBORMarshaler() UgorjiCBORMarshaler {
	return UgorjiCBORMarshaler{h: &codec.CborHandle{}}
}
func (m UgorjiCBORMarshaler) Name() string { return "UgorjiCBOR" }
func (m UgorjiCBORMarshaler) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, m.h).Encode(v)
	return b, err
}
func (m UgorjiCBORMarshaler) Unmarshal(data []byte, v interface{}) error {
	return codec.NewDecoderBytes(data, m.h).Decode(v)
}

type UgorjiBincMarshaler struct {
	h *codec.BincHandle
}

func NewUgorjiBincMarshaler() UgorjiBincMarshaler {
	return UgorjiBincMarshaler{h: &codec.BincHandle{}}
}
func (m UgorjiBincMarshaler) Name() string { return "UgorjiBinc" }
func (m UgorjiBincMarshaler) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, m.h).Encode(v)
	return b, err
}
func (m UgorjiBincMarshaler) Unmarshal(data []byte, v interface{}) error {
	return codec.NewDecoderBytes(data, m.h).Decode(v)
}

type UgorjiSimpleMarshaler struct {
	h *codec.SimpleHandle
}

func NewUgorjiSimpleMarshaler() UgorjiSimpleMarshaler {
	return UgorjiSimpleMarshaler{h: &codec.SimpleHandle{}}
}
func (m UgorjiSimpleMarshaler) Name() string { return "UgorjiSimple" }
func (m UgorjiSimpleMarshaler) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, m.h).Encode(v)
	return b, err
}
func (m UgorjiSimpleMarshaler) Unmarshal(data []byte, v interface{}) error {
	return codec.NewDecoderBytes(data, m.h).Decode(v)
}

// Implementations: Compressor

type NoCompressor struct{}

func (c NoCompressor) Name() string { return "None" }
func (c NoCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}
func (c NoCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

type GzipCompressor struct{}

func (c GzipCompressor) Name() string { return "Gzip" }
func (c GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c GzipCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

type PzipCompressor struct{}

func (c PzipCompressor) Name() string { return "Pzip" }
func (c PzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := pgzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c PzipCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := pgzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

type ZlibCompressor struct{}

func (c ZlibCompressor) Name() string { return "Zlib" }
func (c ZlibCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c ZlibCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

type FlateCompressor struct{}

func (c FlateCompressor) Name() string { return "Flate" }
func (c FlateCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c FlateCompressor) Decompress(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	return io.ReadAll(r)
}

type ZstdCompressor struct{}

func (c ZstdCompressor) Name() string { return "Zstd" }
func (c ZstdCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c ZstdCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := zstd.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

type ZstdFastestCompressor struct {
	enc *zstd.Encoder
	dec *zstd.Decoder
}

func NewZstdFastestCompressor() ZstdFastestCompressor {
	enc, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	dec, _ := zstd.NewReader(nil)
	return ZstdFastestCompressor{enc: enc, dec: dec}
}
func (c ZstdFastestCompressor) Name() string { return "ZstdFastest" }
func (c ZstdFastestCompressor) Compress(data []byte) ([]byte, error) {
	return c.enc.EncodeAll(data, nil), nil
}
func (c ZstdFastestCompressor) Decompress(data []byte) ([]byte, error) {
	return c.dec.DecodeAll(data, nil)
}

type ZstdBestCompressor struct {
	enc *zstd.Encoder
	dec *zstd.Decoder
}

func NewZstdBestCompressor() ZstdBestCompressor {
	enc, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	dec, _ := zstd.NewReader(nil)
	return ZstdBestCompressor{enc: enc, dec: dec}
}
func (c ZstdBestCompressor) Name() string { return "ZstdBest" }
func (c ZstdBestCompressor) Compress(data []byte) ([]byte, error) {
	return c.enc.EncodeAll(data, nil), nil
}
func (c ZstdBestCompressor) Decompress(data []byte) ([]byte, error) {
	return c.dec.DecodeAll(data, nil)
}

type SnappyCompressor struct{}

func (c SnappyCompressor) Name() string { return "Snappy" }
func (c SnappyCompressor) Compress(data []byte) ([]byte, error) {
	// Snappy uses a different API style (Encode/Decode)
	return snappy.Encode(nil, data), nil
}
func (c SnappyCompressor) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

type StdSnappyCompressor struct{}

func (c StdSnappyCompressor) Name() string { return "StdSnappy" }
func (c StdSnappyCompressor) Compress(data []byte) ([]byte, error) {
	return golsnappy.Encode(nil, data), nil
}
func (c StdSnappyCompressor) Decompress(data []byte) ([]byte, error) {
	return golsnappy.Decode(nil, data)
}

type S2Compressor struct{}

func (c S2Compressor) Name() string { return "S2" }
func (c S2Compressor) Compress(data []byte) ([]byte, error) {
	return s2.Encode(nil, data), nil
}
func (c S2Compressor) Decompress(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}

type BrotliCompressor struct{}

func (c BrotliCompressor) Name() string { return "Brotli" }
func (c BrotliCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c BrotliCompressor) Decompress(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}

type Lz4Compressor struct{}

func (c Lz4Compressor) Name() string { return "Lz4" }
func (c Lz4Compressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := lz4.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c Lz4Compressor) Decompress(data []byte) ([]byte, error) {
	r := lz4.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}

type XZCompressor struct{}

func (c XZCompressor) Name() string { return "XZ" }
func (c XZCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c XZCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

type Bzip2Compressor struct{}

func (c Bzip2Compressor) Name() string { return "Bzip2" }
func (c Bzip2Compressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := bzip2.NewWriter(&buf, nil)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (c Bzip2Compressor) Decompress(data []byte) ([]byte, error) {
	r, err := bzip2.NewReader(bytes.NewReader(data), nil)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

type GozstdCompressorBest struct{}

func (c GozstdCompressorBest) Name() string { return "GozstdBest" }
func (c GozstdCompressorBest) Compress(data []byte) ([]byte, error) {
	return gozstd.CompressLevel(nil, data, 20), nil
}
func (c GozstdCompressorBest) Decompress(data []byte) ([]byte, error) {
	return gozstd.Decompress(nil, data)
}

type GozstdCompressorFastest struct{}

func (c GozstdCompressorFastest) Name() string { return "GozstdFastest" }
func (c GozstdCompressorFastest) Compress(data []byte) ([]byte, error) {
	return gozstd.CompressLevel(nil, data, 1), nil
}
func (c GozstdCompressorFastest) Decompress(data []byte) ([]byte, error) {
	return gozstd.Decompress(nil, data)
}

type GozstdCompressor struct{}

func (c GozstdCompressor) Name() string { return "Gozstd" }
func (c GozstdCompressor) Compress(data []byte) ([]byte, error) {
	return gozstd.Compress(nil, data), nil
}
func (c GozstdCompressor) Decompress(data []byte) ([]byte, error) {
	return gozstd.Decompress(nil, data)
}

type DataDogZstdCompressorBest struct{}

func (c DataDogZstdCompressorBest) Name() string { return "DataDogZstdBest" }
func (c DataDogZstdCompressorBest) Compress(data []byte) ([]byte, error) {
	return ddzstd.CompressLevel(nil, data, ddzstd.BestCompression)
}
func (c DataDogZstdCompressorBest) Decompress(data []byte) ([]byte, error) {
	return ddzstd.Decompress(nil, data)
}

type DataDogZstdCompressorFastest struct{}

func (c DataDogZstdCompressorFastest) Name() string { return "DataDogZstdFastest" }
func (c DataDogZstdCompressorFastest) Compress(data []byte) ([]byte, error) {
	return ddzstd.CompressLevel(nil, data, ddzstd.BestSpeed)
}
func (c DataDogZstdCompressorFastest) Decompress(data []byte) ([]byte, error) {
	return ddzstd.Decompress(nil, data)
}

type DataDogZstdCompressor struct{}

func (c DataDogZstdCompressor) Name() string { return "DataDogZstd" }
func (c DataDogZstdCompressor) Compress(data []byte) ([]byte, error) {
	return ddzstd.CompressLevel(nil, data, ddzstd.DefaultCompression)
}
func (c DataDogZstdCompressor) Decompress(data []byte) ([]byte, error) {
	return ddzstd.Decompress(nil, data)
}

// Benchmarking

type BenchmarkResult struct {
	ID             string
	OriginalSize   int
	MarshaledSize  int
	CompressedSize int
	SpeedMBs       float64
	SerCompSpeed   float64
	DecompDesSpeed float64
	Corrupted      bool
}

func runBenchmark(m Marshaler, c Compressor, k int, samples [][]Book) BenchmarkResult {
	// Warmup / Sanity check with first sample
	warmupPayload := Payload{
		Total:    k,
		Elements: samples[0],
	}
	b, err := m.Marshal(warmupPayload)
	if err != nil {
		panic(fmt.Sprintf("Marshal failed for %s: %v", m.Name(), err))
	}
	compressed, err := c.Compress(b)
	if err != nil {
		panic(fmt.Sprintf("Compress failed for %s: %v", c.Name(), err))
	}
	decompressed, err := c.Decompress(compressed)
	if err != nil {
		panic(fmt.Sprintf("Decompress failed for %s: %v", c.Name(), err))
	}
	var newPayload Payload
	if err := m.Unmarshal(decompressed, &newPayload); err != nil {
		panic(fmt.Sprintf("Unmarshal failed for %s: %v", m.Name(), err))
	}

	var totalOriginalSize int
	var totalMarshaledSize int
	var totalCompressedSize int
	var totalDuration time.Duration
	var totalSerCompDuration time.Duration
	var totalDecompDesDuration time.Duration
	corrupted := false

	loopN := 10 // run multiple times, to avoid flaky benchmarks due to time calculation being slow
	for _, elements := range samples {
		payload := Payload{
			Total:    k,
			Elements: elements,
		}

		// Calculate original size
		sampleOriginalSize := 0
		for _, item := range payload.Elements {
			sampleOriginalSize += item.Size()
		}
		totalOriginalSize += sampleOriginalSize

		// Run benchmark
		var marshaledSize, compressedSize int
		var p Payload
		start := time.Now()
		var comp []byte
		for i := 0; i < loopN; i++ {
			data, _ := m.Marshal(payload)
			marshaledSize = len(data)
			comp, _ = c.Compress(data)
		}
		totalSerCompDuration += time.Since(start)
		startDecomp := time.Now()

		for i := 0; i < loopN; i++ {
			compressedSize = len(comp)

			decomp, _ := c.Decompress(comp)
			p = Payload{}
			m.Unmarshal(decomp, &p)
		}
		totalDecompDesDuration += time.Since(startDecomp)

		totalMarshaledSize += marshaledSize
		totalCompressedSize += compressedSize

		if !corrupted && !reflect.DeepEqual(payload, p) {
			corrupted = true
		}
	}

	totalDuration = totalSerCompDuration + totalDecompDesDuration
	numSamples := len(samples)
	// Speed calculation: Total Original Size / Total Duration
	totalBytesProcessed := float64(totalOriginalSize) * float64(loopN)

	speed := (totalBytesProcessed / (1024 * 1024)) / totalDuration.Seconds()
	serCompSpeed := (totalBytesProcessed / (1024 * 1024)) / totalSerCompDuration.Seconds()
	decompDesSpeed := (totalBytesProcessed / (1024 * 1024)) / totalDecompDesDuration.Seconds()

	return BenchmarkResult{
		ID:             fmt.Sprintf("%s/%s/len=%d", m.Name(), c.Name(), k),
		OriginalSize:   totalOriginalSize / (numSamples * k),
		MarshaledSize:  totalMarshaledSize / (numSamples * k),
		CompressedSize: totalCompressedSize / (numSamples * k),
		SpeedMBs:       speed,
		SerCompSpeed:   serCompSpeed,
		DecompDesSpeed: decompDesSpeed,
		Corrupted:      corrupted,
	}
}

func Run(all bool) {
	// Load data
	file, err := os.ReadFile("internal/bench/sample.json")
	if err != nil {
		panic(err)
	}
	var dataset []Book
	if err := stdjson.Unmarshal(file, &dataset); err != nil {
		panic(err)
	}

	if len(dataset) == 0 {
		panic("No data in sample.json")
	}

	marshalers := []Marshaler{
		JSONMarshaler{},
		MsgPackMarshaler{},
		CBORMarshaler{},
		YAMLMarshaler{},
		TOMLMarshaler{},
		BSONMarshaler{},
		AvroMarshaler{},
		BencodeMarshaler{},
		XMLMarshaler{},
		JsonIterMarshaler{},
		SonicMarshaler{},
		GoccyJSONMarshaler{},
		GoccyYAMLMarshaler{},
		BurntSushiTOMLMarshaler{},
		ShamatonMsgPackMarshaler{},
		NewUgorjiMsgPackMarshaler(),
		NewUgorjiCBORMarshaler(),
		EasyJSONMarshaler{},
		MsgpMarshaler{},
		SegmentioJSONMarshaler{},
		GoAvroMarshaler{},
		NewJsonIterFastestMarshaler(),
		NewUgorjiBincMarshaler(),
		NewUgorjiSimpleMarshaler(),
		NewCBORPreferredMarshaler(),
	}

	compressors := []Compressor{
		NoCompressor{},
		GzipCompressor{},
		PzipCompressor{},
		ZlibCompressor{},
		FlateCompressor{},
		ZstdCompressor{},
		SnappyCompressor{},
		S2Compressor{},
		BrotliCompressor{},
		Lz4Compressor{},
		XZCompressor{},
		Bzip2Compressor{},
		NewZstdFastestCompressor(),
		NewZstdBestCompressor(),
		StdSnappyCompressor{},
		GozstdCompressor{},
		GozstdCompressorBest{},
		GozstdCompressorFastest{},
		DataDogZstdCompressor{},
		DataDogZstdCompressorBest{},
		DataDogZstdCompressorFastest{},
	}

	kValues := []int{1, 5, 100}

	// Generate 50 random samples of max(kValues, 100) items
	maxK := 100
	for _, k := range kValues {
		if k > maxK {
			maxK = k
		}
	}

	allSamples := make([][]Book, 50)
	for i := 0; i < 50; i++ {
		allSamples[i] = make([]Book, maxK)
		for j := 0; j < maxK; j++ {
			allSamples[i][j] = dataset[rand.Intn(len(dataset))]
		}
	}

	getSamples := func(k int) [][]Book {
		res := make([][]Book, 50)
		for i := 0; i < 50; i++ {
			res[i] = allSamples[i][:k]
		}
		return res
	}

	var results []BenchmarkResult
	excludedResultsCount := 0

	fmt.Println("Performing initial speed check (None compression)...")
	excludedMarshalers := make(map[string]bool)
	var excludedList []string
	initialSamples := getSamples(100)
	for _, m := range marshalers {
		// Use k=100 for a stable speed measurement
		res := runBenchmark(m, NoCompressor{}, 100, initialSamples)
		if !all && res.SpeedMBs < 20.0 {
			excludedMarshalers[m.Name()] = true
			excludedList = append(excludedList, fmt.Sprintf("%s (%.2f MB/s)", m.Name(), res.SpeedMBs))
		}
	}

	if len(excludedList) > 0 {
		fmt.Println("\nExcluded Marshalers (Speed < 20MB/s with None compression):")
		for _, item := range excludedList {
			fmt.Printf("- %s\n", item)
		}
	}

	fmt.Println("\nPerforming Avro speed check with all compressors...")
	var filteredCompressors []Compressor
	var excludedCompressorsList []string
	avroM := AvroMarshaler{}
	for _, c := range compressors {
		// Use k=100 for a stable speed measurement
		res := runBenchmark(avroM, c, 100, initialSamples)
		if !all && res.SpeedMBs < 10.0 {
			excludedCompressorsList = append(excludedCompressorsList, fmt.Sprintf("%s (%.2f MB/s)", c.Name(), res.SpeedMBs))
			continue
		}
		filteredCompressors = append(filteredCompressors, c)
	}

	if len(excludedCompressorsList) > 0 {
		fmt.Println("Excluded Compressors (Speed < 10MB/s with Avro):")
		for _, item := range excludedCompressorsList {
			fmt.Printf("- %s\n", item)
		}
	}
	compressors = filteredCompressors

	fmt.Println("\nRunning benchmarks...")

	for _, k := range kValues {
		samples := getSamples(k)
		for _, m := range marshalers {
			if excludedMarshalers[m.Name()] {
				continue
			}
			for _, c := range compressors {
				res := runBenchmark(m, c, k, samples)
				if !all && res.SpeedMBs < 50.0 {
					excludedResultsCount++
					continue
				}
				results = append(results, res)
				// fmt.Printf("Completed: %s\n", res.ID)
			}
		}
	}

	// Sort by Speed descending
	sort.Slice(results, func(i, j int) bool {
		return float64(results[i].CompressedSize) < float64(results[j].CompressedSize)
	})

	// Output
	if excludedResultsCount > 0 {
		fmt.Printf("\nExcluded %d results with speed < 50MB/s\n", excludedResultsCount)
	}

	fmt.Printf("%-30s | %-10s | %-10s | %-10s | %-10s | %-10s\n", "ID", "Ratio (%)", "Speed", "Ser+Comp", "De+Deser", "Status")
	fmt.Println(strings.Repeat("-", 90))
	for _, r := range results {
		status := "OK"
		if r.Corrupted {
			status = "CORRUPTED"
		}
		ratio := float64(r.CompressedSize) / float64(r.OriginalSize) * 100
		fmt.Printf("%-25s | %-10.2f | %-10.2f | %-10.2f | %-10.2f | %-10s\n",
			r.ID, ratio, r.SpeedMBs, r.SerCompSpeed, r.DecompDesSpeed, status)
	}

}
