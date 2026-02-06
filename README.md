## Go Serialization + Compression benchmarks  (AI generated)

If you want to serialize and then compress data - what will be the most effective library to use?

This benchmarks tries to answer that question, which may be useful to for passing simple repetitive data over network-bounded environments (public internet, slow networks, etc)

Benchmark is run with GOMAXPROCS=1 and results may be vastly different depending on the platform, so it's just for a reference only. Plus it's just AI written in few hours and sample data is limited. You probably want to do more robust benchmark for best candidates on the real-world data for your specific use-case.

Compression ratio is calculated relative to the size of Go struct fields.
Compressed struct contains is not only the data itself, but also wrapper envelope.

So far github.com/hamba/avro/v2 + github.com/klauspost/compress/zstd provide the best balance between compression size and speed

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

Fast And Compact
```bash
ID                             | Ratio (%)  | Speed      | Ser+Comp   | De+Deser   | Status    
------------------------------------------------------------------------------------------
Avro/ZstdFastest/len=100  | 43.37      | 127.99     | 194.28     | 375.09     | OK        
ShamatonMsgPack/ZstdFastest/len=100 | 49.40      | 73.19      | 126.50     | 173.68     | OK        
Msgp/ZstdFastest/len=100  | 49.40      | 96.71      | 160.30     | 243.81     | OK        
UgorjiCBOR/ZstdFastest/len=100 | 49.40      | 57.59      | 106.46     | 125.43     | OK        
UgorjiBinc/ZstdFastest/len=100 | 50.60      | 57.26      | 99.07      | 135.69     | OK        
UgorjiMsgPack/ZstdFastest/len=100 | 50.60      | 59.99      | 107.11     | 136.36     | OK        
UgorjiSimple/ZstdFastest/len=100 | 50.60      | 57.65      | 99.19      | 137.66     | OK        
Avro/Snappy/len=100       | 53.01      | 146.35     | 281.46     | 304.88     | OK        
Avro/StdSnappy/len=100    | 54.22      | 276.09     | 494.12     | 625.69     | OK        
Avro/Lz4/len=100          | 55.42      | 149.28     | 264.67     | 342.42     | OK        
Avro/S2/len=100           | 55.42      | 296.76     | 558.42     | 633.33     | OK        
Sonic/Snappy/len=100      | 62.65      | 54.51      | 110.92     | 107.19     | OK        
JsonIterFastest/Snappy/len=100 | 62.65      | 51.27      | 154.76     | 76.68      | OK        
EasyJSON/Snappy/len=100   | 62.65      | 68.00      | 152.94     | 122.43     | OK        
Msgp/Snappy/len=100       | 63.86      | 152.53     | 268.32     | 353.46     | OK        
ShamatonMsgPack/Snappy/len=100 | 63.86      | 104.77     | 191.22     | 231.76     | OK        
MsgPack/Snappy/len=100    | 63.86      | 67.12      | 158.14     | 116.63     | OK        
UgorjiSimple/Snappy/len=100 | 65.06      | 75.29      | 130.67     | 177.64     | OK        
Sonic/S2/len=100          | 65.06      | 61.02      | 157.55     | 99.59      | OK        
UgorjiBinc/Snappy/len=100 | 65.06      | 73.19      | 131.28     | 165.41     | OK        
EasyJSON/S2/len=100       | 65.06      | 82.14      | 209.73     | 135.01     | OK        
JsonIterFastest/Lz4/len=100 | 65.06      | 50.10      | 165.48     | 71.86      | OK        
JsonIterFastest/S2/len=100 | 65.06      | 51.56      | 185.58     | 71.40      | OK        
JsonIter/S2/len=100       | 65.06      | 58.17      | 210.55     | 80.38      | OK        
GoccyJSON/S2/len=100      | 65.06      | 50.01      | 106.94     | 93.95      | OK        
UgorjiMsgPack/Snappy/len=100 | 65.06      | 83.39      | 156.53     | 178.45     | OK        
CBORCanonical/Snappy/len=100 | 65.06      | 63.62      | 179.59     | 98.51      | OK        
Sonic/Lz4/len=100         | 65.06      | 55.39      | 132.86     | 94.98      | OK        
CBOR/Snappy/len=100       | 65.06      | 63.01      | 188.02     | 94.76      | OK        
UgorjiCBOR/Snappy/len=100 | 65.06      | 76.49      | 151.63     | 154.36     | OK        
EasyJSON/Lz4/len=100      | 65.06      | 60.80      | 150.18     | 102.17     | OK        
UgorjiSimple/StdSnappy/len=100 | 66.27      | 86.34      | 166.97     | 178.77     | OK        
EasyJSON/StdSnappy/len=100 | 66.27      | 80.18      | 221.41     | 125.70     | OK        
Sonic/StdSnappy/len=100   | 66.27      | 67.69      | 171.23     | 111.95     | OK        
Msgp/Lz4/len=100          | 66.27      | 120.02     | 252.90     | 228.44     | OK        
Msgp/StdSnappy/len=100    | 66.27      | 188.68     | 394.37     | 361.77     | OK        
JsonIterFastest/StdSnappy/len=100 | 66.27      | 51.34      | 183.18     | 71.33      | OK        
ShamatonMsgPack/Lz4/len=100 | 66.27      | 86.48      | 168.10     | 178.09     | OK        
JsonIter/StdSnappy/len=100 | 66.27      | 54.77      | 184.19     | 77.96      | OK        
MsgPack/StdSnappy/len=100 | 66.27      | 71.03      | 177.31     | 118.49     | OK        
UgorjiMsgPack/StdSnappy/len=100 | 66.27      | 95.64      | 198.58     | 184.50     | OK        
MsgPack/Lz4/len=100       | 66.27      | 57.98      | 140.16     | 98.88      | OK        
ShamatonMsgPack/StdSnappy/len=100 | 66.27      | 115.64     | 239.84     | 223.31     | OK        
UgorjiMsgPack/Lz4/len=100 | 67.47      | 78.98      | 166.21     | 150.50     | OK        
CBOR/StdSnappy/len=100    | 67.47      | 65.74      | 203.15     | 97.19      | OK        
CBORCanonical/StdSnappy/len=100 | 67.47      | 65.74      | 234.32     | 91.38      | OK        
MsgPack/S2/len=100        | 67.47      | 73.90      | 179.30     | 125.71     | OK        
CBOR/S2/len=100           | 67.47      | 64.13      | 230.21     | 88.90      | OK        
CBOR/Lz4/len=100          | 67.47      | 55.84      | 172.91     | 82.47      | OK        
UgorjiMsgPack/S2/len=100  | 67.47      | 93.93      | 218.58     | 164.70     | OK        
CBORCanonical/Lz4/len=100 | 67.47      | 57.31      | 174.67     | 85.31      | OK        
CBORCanonical/S2/len=100  | 67.47      | 73.91      | 244.18     | 106.00     | OK        
UgorjiSimple/Lz4/len=100  | 67.47      | 69.34      | 139.35     | 138.02     | OK        
UgorjiSimple/S2/len=100   | 67.47      | 84.36      | 162.37     | 175.58     | OK        
ShamatonMsgPack/S2/len=100 | 67.47      | 108.81     | 247.97     | 193.89     | OK        
UgorjiCBOR/S2/len=100     | 67.47      | 88.47      | 188.41     | 166.79     | OK        
UgorjiCBOR/Lz4/len=100    | 67.47      | 74.58      | 162.75     | 137.65     | OK        
UgorjiCBOR/StdSnappy/len=100 | 67.47      | 90.97      | 196.56     | 169.35     | OK        
UgorjiBinc/StdSnappy/len=100 | 67.47      | 81.68      | 158.07     | 169.02     | OK        
UgorjiBinc/Lz4/len=100    | 67.47      | 67.36      | 132.91     | 136.58     | OK        
Msgp/S2/len=100           | 67.47      | 175.78     | 382.80     | 325.03     | OK        
UgorjiBinc/S2/len=100     | 67.47      | 92.04      | 188.54     | 179.84     | OK        
Avro/StdSnappy/len=5      | 90.48      | 229.93     | 400.74     | 539.44     | OK        
Avro/S2/len=5             | 90.48      | 152.94     | 323.03     | 290.48     | OK        
Avro/Snappy/len=5         | 92.86      | 154.98     | 283.49     | 341.90     | OK        
Avro/Lz4/len=5            | 94.05      | 93.12      | 167.80     | 209.25     | OK        
Avro/None/len=100         | 98.80      | 421.43     | 1033.41    | 711.63     | OK        
Avro/None/len=5           | 100.00     | 398.75     | 918.60     | 704.61     | OK  
...
```


Most Compact
```bash
ID                             | Ratio (%)  | Speed      | Ser+Comp   | De+Deser   | Status    
------------------------------------------------------------------------------------------
Avro/Brotli/len=100       | 38.55      | 12.13      | 13.73      | 104.44     | OK        
GoAvro/Brotli/len=100     | 38.55      | 5.83       | 8.64       | 17.92      | OK        
GoAvro/XZ/len=100         | 39.76      | 0.87       | 1.11       | 4.16       | OK        
Avro/XZ/len=100           | 39.76      | 0.94       | 1.14       | 5.44       | OK        
JsonIter/Brotli/len=100   | 40.96      | 7.64       | 9.69       | 36.15      | OK        
Sonic/Brotli/len=100      | 40.96      | 8.66       | 10.42      | 51.36      | OK        
GoAvro/ZstdBest/len=100   | 40.96      | 9.02       | 15.36      | 21.84      | OK        
GoAvro/Flate/len=100      | 40.96      | 8.22       | 14.84      | 18.45      | OK        
GoAvro/Zlib/len=100       | 40.96      | 8.43       | 15.48      | 18.52      | OK        
JSON/Brotli/len=100       | 40.96      | 6.96       | 9.16       | 28.93      | OK        
GoccyJSON/Brotli/len=100  | 40.96      | 8.26       | 10.01      | 47.39      | OK        
GoAvro/Pzip/len=100       | 40.96      | 4.40       | 8.41       | 9.21       | OK        
GoAvro/Gzip/len=100       | 40.96      | 8.21       | 14.47      | 18.97      | OK        
GoccyYAML/Brotli/len=100  | 40.96      | 1.53       | 2.78       | 3.40       | OK        
SegmentioJSON/Brotli/len=100 | 40.96      | 6.76       | 8.81       | 29.15      | OK        
Avro/Gzip/len=100         | 40.96      | 26.59      | 32.47      | 146.77     | OK        
Avro/Pzip/len=100         | 40.96      | 6.97       | 13.06      | 14.97      | OK        
Avro/Zlib/len=100         | 40.96      | 29.15      | 37.38      | 132.31     | OK        
Avro/Flate/len=100        | 40.96      | 29.54      | 36.77      | 150.36     | OK        
Avro/ZstdBest/len=100     | 40.96      | 40.38      | 46.45      | 308.59     | OK        
JsonIterFastest/Brotli/len=100 | 40.96      | 8.56       | 10.79      | 41.43      | OK        
EasyJSON/Brotli/len=100   | 40.96      | 8.79       | 10.44      | 55.82      | OK        
BurntSushiTOML/Brotli/len=100 | 42.17      | 2.84       | 4.75       | 7.08       | OK        
XML/Brotli/len=100        | 42.17      | 3.95       | 7.85       | 7.94       | OK        
GoAvro/Zstd/len=100       | 42.17      | 7.42       | 11.96      | 19.51      | OK        
YAML/Brotli/len=100       | 42.17      | 2.52       | 4.44       | 5.86       | OK        
Avro/Zstd/len=100         | 42.17      | 21.73      | 24.46      | 194.80     | OK        
TOML/Brotli/len=100       | 42.17      | 4.88       | 6.14       | 23.77      | OK        
UgorjiBinc/Brotli/len=100 | 43.37      | 9.21       | 10.72      | 65.21      | OK        
BurntSushiTOML/XZ/len=100 | 43.37      | 0.83       | 1.08       | 3.56       | OK        
JSON/XZ/len=100           | 43.37      | 0.88       | 1.09       | 4.45       | OK        
Msgp/Brotli/len=100       | 43.37      | 9.67       | 10.95      | 82.33      | OK        
MsgPack/Brotli/len=100    | 43.37      | 8.99       | 10.81      | 53.43      | OK        
TOML/XZ/len=100           | 43.37      | 0.79       | 0.98       | 3.99       | OK        
UgorjiCBOR/Brotli/len=100 | 43.37      | 8.86       | 10.32      | 62.43      | OK        
UgorjiSimple/Brotli/len=100 | 43.37      | 9.15       | 10.62      | 65.92      | OK        
EasyJSON/XZ/len=100       | 43.37      | 0.86       | 1.07       | 4.56       | OK        
ShamatonMsgPack/Brotli/len=100 | 43.37      | 9.77       | 11.30      | 71.76      | OK        
SegmentioJSON/XZ/len=100  | 43.37      | 0.88       | 1.10       | 4.49       | OK        
UgorjiMsgPack/Brotli/len=100 | 43.37      | 8.79       | 10.19      | 63.86      | OK        
CBOR/Brotli/len=100       | 43.37      | 9.07       | 11.04      | 50.81      | OK        
GoccyYAML/ZstdBest/len=100 | 43.37      | 1.67       | 3.25       | 3.45       | OK        
CBORCanonical/Brotli/len=100 | 43.37      | 9.37       | 11.49      | 50.70      | OK        
Avro/ZstdFastest/len=100  | 43.37      | 115.28     | 192.42     | 287.59     | OK        
GoccyYAML/XZ/len=100      | 43.37      | 0.59       | 0.84       | 1.96       | OK        
JsonIter/XZ/len=100       | 43.37      | 0.84       | 1.07       | 3.90       | OK        
GoAvro/ZstdFastest/len=100 | 43.37      | 11.42      | 23.69      | 22.03      | OK        
JsonIterFastest/XZ/len=100 | 43.37      | 0.92       | 1.15       | 4.56       | OK        
GoccyJSON/XZ/len=100      | 43.37      | 0.90       | 1.12       | 4.49       | OK        
Sonic/XZ/len=100          | 43.37      | 0.84       | 1.04       | 4.56       | OK        
JsonIter/Gzip/len=100     | 44.58      | 16.89      | 26.38      | 46.95      | OK        
YAML/ZstdBest/len=100     | 44.58      | 3.07       | 5.54       | 6.88       | OK        
JSON/Gzip/len=100         | 44.58      | 13.31      | 21.93      | 33.89      | OK
...
```