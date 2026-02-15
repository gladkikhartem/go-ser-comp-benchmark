## Go Serialization + Compression benchmarks  (AI generated)

If you want to serialize and then compress data - what will be the most effective library to use?

This benchmarks tries to answer that question, which may be useful to for passing simple repetitive data over network-bounded environments (public internet, slow networks, etc)

Benchmark is run with GOMAXPROCS=1 and results may be vastly different depending on the platform, so it's just for a reference only. Plus it's just AI written in few hours and sample data is limited. You probably want to do more robust benchmark for best candidates on the real-world data for your specific use-case.

Compression ratio is calculated relative to the size of Go struct fields.
Compressed struct contains is not only the data itself, but also wrapper envelope.

So far Avro + Zstd provide the best balance between compression size and speed

Full benchmarks results without exclusion of 'slow' libraries are located in full_results.txt
```go

type Book struct {
	Rank             int    
	Book             string 
	Author           string 
	OriginalLanguage string 
	FirstPublished   string 
	ApproximateSales string 
	Genre            string 
}
func (b Book) Size() int {
	return 8 + len(b.Book) + len(b.Author) + len(b.OriginalLanguage) + len(b.FirstPublished) + len(b.ApproximateSales) + len(b.Genre)
}

type Payload struct {
	Total    int    
	Elements []Book 
}
```

Fast & Compact Sample
```bash
Excluded 617 results with speed < 50MB/s
ID                                       | Ratio (%)  | Speed      | Ser+Comp   | De+Deser  
------------------------------------------------------------------------------------------
Avro/DataDogZstd/len=100                 | 41         | 77         | 97         | 374       
Avro/DataDogZstdFastest/len=100          | 42         | 157        | 260        | 394       
Avro/GozstdFastest/len=100               | 42         | 165        | 271        | 425       
Avro/Gozstd/len=100                      | 42         | 130        | 207        | 347       
Avro/ZstdFastest/len=100                 | 43         | 111        | 174        | 310       
GOB/DataDogZstd/len=100                  | 45         | 53         | 79         | 163       
GOB/Gozstd/len=100                       | 46         | 74         | 138        | 159       
Sonic/Gozstd/len=100                     | 47         | 52         | 112        | 98        
Msgp/DataDogZstd/len=100                 | 47         | 51         | 63         | 265       
EasyJSON/Gozstd/len=100                  | 47         | 64         | 132        | 124       
GOB/GozstdFastest/len=100                | 47         | 81         | 167        | 157       
GOB/DataDogZstdFastest/len=100           | 47         | 84         | 161        | 176       
CBORCanonical/Gozstd/len=100             | 48         | 60         | 143        | 104       
JsonIter/DataDogZstdFastest/len=100      | 48         | 51         | 162        | 75        
GOB/ZstdFastest/len=100                  | 48         | 74         | 143        | 154       
Sonic/GozstdFastest/len=100              | 48         | 54         | 119        | 101       
Sonic/DataDogZstdFastest/len=100         | 48         | 55         | 124        | 100       
UgorjiCBOR/Gozstd/len=100                | 48         | 64         | 114        | 146       
EasyJSON/ZstdFastest/len=100             | 48         | 59         | 117        | 120       
EasyJSON/DataDogZstdFastest/len=100      | 48         | 65         | 149        | 114       
EasyJSON/GozstdFastest/len=100           | 48         | 73         | 175        | 125       
CBOR/Gozstd/len=100                      | 48         | 59         | 146        | 100       
Msgp/GozstdFastest/len=100               | 49         | 120        | 233        | 248       
CBOR/ZstdFastest/len=100                 | 49         | 56         | 129        | 100       
CBORCanonical/DataDogZstdFastest/len=100 | 49         | 61         | 152        | 103       
CBORCanonical/GozstdFastest/len=100      | 49         | 61         | 155        | 101       
CBORCanonical/ZstdFastest/len=100        | 49         | 56         | 128        | 101       
MsgPack/Gozstd/len=100                   | 49         | 59         | 116        | 118       
Msgp/DataDogZstdFastest/len=100          | 49         | 114        | 222        | 233       
MsgPack/GozstdFastest/len=100            | 49         | 64         | 140        | 118       
MsgPack/DataDogZstdFastest/len=100       | 49         | 62         | 133        | 116       
Msgp/Gozstd/len=100                      | 49         | 112        | 186        | 282       
CBOR/GozstdFastest/len=100               | 49         | 67         | 174        | 108       
CBOR/DataDogZstdFastest/len=100          | 49         | 61         | 154        | 101       
UgorjiSimple/Gozstd/len=100              | 49         | 65         | 110        | 160       
UgorjiCBOR/DataDogZstdFastest/len=100    | 49         | 70         | 138        | 140       
UgorjiCBOR/GozstdFastest/len=100         | 49         | 73         | 145        | 145       
UgorjiBinc/Gozstd/len=100                | 49         | 67         | 121        | 150       
UgorjiCBOR/ZstdFastest/len=100           | 49         | 59         | 107        | 132       
ShamatonMsgPack/Gozstd/len=100           | 49         | 81         | 139        | 195       
UgorjiMsgPack/DataDogZstdFastest/len=100 | 49         | 73         | 144        | 147       
UgorjiMsgPack/GozstdFastest/len=100      | 49         | 73         | 146        | 145       
UgorjiMsgPack/Gozstd/len=100             | 49         | 66         | 120        | 147       
ShamatonMsgPack/GozstdFastest/len=100    | 49         | 85         | 164        | 178       
ShamatonMsgPack/DataDogZstdFastest/len=100 | 49         | 83         | 150        | 183       
MsgPack/ZstdFastest/len=100              | 51         | 54         | 109        | 108       
UgorjiSimple/GozstdFastest/len=100       | 51         | 68         | 128        | 146       
UgorjiBinc/DataDogZstdFastest/len=100    | 51         | 71         | 130        | 158       
UgorjiSimple/ZstdFastest/len=100         | 51         | 61         | 105        | 144       
ShamatonMsgPack/ZstdFastest/len=100      | 51         | 69         | 116        | 170       
UgorjiBinc/GozstdFastest/len=100         | 51         | 72         | 143        | 147       
Msgp/ZstdFastest/len=100                 | 51         | 88         | 150        | 215       
UgorjiSimple/DataDogZstdFastest/len=100  | 51         | 74         | 140        | 156       
UgorjiMsgPack/ZstdFastest/len=100        | 51         | 62         | 111        | 139       
UgorjiBinc/ZstdFastest/len=100           | 51         | 61         | 105        | 147       
Avro/Snappy/len=100                      | 53         | 224        | 352        | 617       
Avro/S2/len=100                          | 55         | 220        | 391        | 504       
Avro/StdSnappy/len=100                   | 55         | 235        | 443        | 500       
Avro/Lz4/len=100                         | 57         | 200        | 353        | 462       
GOB/Snappy/len=100                       | 59         | 95         | 202        | 181       
GOB/StdSnappy/len=100                    | 60         | 111        | 256        | 195       
GOB/Lz4/len=100                          | 61         | 94         | 205        | 174       
GOB/S2/len=100                           | 61         | 113        | 249        | 206
...
```