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
ID                             | Ratio (%)  | Speed      | Ser+Comp   | De+Deser   | Status    
------------------------------------------------------------------------------------------
Avro/Gozstd/len=100       | 40.96      | 124.34     | 189.89     | 360.22     | OK        
Avro/DataDogZstd/len=100  | 40.96      | 77.51      | 98.97      | 357.36     | OK        
Avro/DataDogZstdFastest/len=100 | 42.17      | 159.08     | 261.58     | 406.00     | OK        
Avro/GozstdFastest/len=100 | 42.17      | 166.93     | 275.39     | 423.86     | OK        
Avro/ZstdFastest/len=100  | 43.37      | 131.46     | 202.42     | 375.00     | OK        
EasyJSON/Gozstd/len=100   | 46.99      | 62.75      | 133.42     | 118.48     | OK        
EasyJSON/DataDogZstdFastest/len=100 | 46.99      | 66.03      | 150.82     | 117.44     | OK        
EasyJSON/GozstdFastest/len=100 | 46.99      | 74.56      | 172.63     | 131.26     | OK        
Msgp/DataDogZstd/len=100  | 46.99      | 53.30      | 67.74      | 250.06     | OK        
Sonic/GozstdFastest/len=100 | 46.99      | 56.81      | 128.97     | 101.55     | OK        
Sonic/DataDogZstdFastest/len=100 | 46.99      | 50.53      | 117.09     | 88.89      | OK        
MsgPack/Gozstd/len=100    | 48.19      | 56.99      | 110.28     | 117.93     | OK        
UgorjiCBOR/Gozstd/len=100 | 48.19      | 64.88      | 121.65     | 139.05     | OK        
EasyJSON/ZstdFastest/len=100 | 48.19      | 58.97      | 118.05     | 117.83     | OK        
ShamatonMsgPack/Gozstd/len=100 | 48.19      | 74.94      | 128.25     | 180.27     | OK        
CBOR/Gozstd/len=100       | 48.19      | 55.94      | 138.68     | 93.76      | OK        
CBORCanonical/Gozstd/len=100 | 48.19      | 59.00      | 139.27     | 102.37     | OK        
Msgp/Gozstd/len=100       | 48.19      | 106.17     | 179.49     | 259.92     | OK        
UgorjiMsgPack/GozstdFastest/len=100 | 49.40      | 69.73      | 133.10     | 146.47     | OK        
UgorjiCBOR/DataDogZstdFastest/len=100 | 49.40      | 70.84      | 141.59     | 141.78     | OK        
Msgp/DataDogZstdFastest/len=100 | 49.40      | 119.93     | 219.67     | 264.15     | OK        
UgorjiSimple/Gozstd/len=100 | 49.40      | 65.57      | 112.51     | 157.16     | OK        
Msgp/GozstdFastest/len=100 | 49.40      | 121.15     | 221.40     | 267.57     | OK        
CBORCanonical/DataDogZstdFastest/len=100 | 49.40      | 62.85      | 159.81     | 103.58     | OK        
MsgPack/GozstdFastest/len=100 | 49.40      | 60.47      | 133.44     | 110.58     | OK        
ShamatonMsgPack/GozstdFastest/len=100 | 49.40      | 82.80      | 156.16     | 176.25     | OK        
MsgPack/DataDogZstdFastest/len=100 | 49.40      | 63.05      | 130.31     | 122.15     | OK        
CBOR/ZstdFastest/len=100  | 49.40      | 50.55      | 116.06     | 89.54      | OK        
ShamatonMsgPack/DataDogZstdFastest/len=100 | 49.40      | 82.86      | 149.60     | 185.72     | OK        
UgorjiBinc/Gozstd/len=100 | 49.40      | 67.85      | 116.83     | 161.86     | OK        
CBORCanonical/GozstdFastest/len=100 | 49.40      | 57.14      | 148.95     | 92.70      | OK        
CBORCanonical/ZstdFastest/len=100 | 49.40      | 50.76      | 115.43     | 90.60      | OK        
CBOR/GozstdFastest/len=100 | 49.40      | 60.00      | 162.83     | 95.00      | OK        
CBOR/DataDogZstdFastest/len=100 | 49.40      | 61.61      | 160.36     | 100.05     | OK        
UgorjiMsgPack/Gozstd/len=100 | 49.40      | 59.47      | 106.35     | 134.93     | OK        
UgorjiMsgPack/DataDogZstdFastest/len=100 | 49.40      | 70.77      | 139.53     | 143.62     | OK        
UgorjiCBOR/GozstdFastest/len=100 | 49.40      | 69.10      | 142.11     | 134.49     | OK        
UgorjiCBOR/ZstdFastest/len=100 | 49.40      | 58.63      | 108.30     | 127.84     | OK        
ShamatonMsgPack/ZstdFastest/len=100 | 50.60      | 65.99      | 112.66     | 159.27     | OK        
MsgPack/ZstdFastest/len=100 | 50.60      | 53.00      | 105.21     | 106.79     | OK        
UgorjiSimple/ZstdFastest/len=100 | 50.60      | 61.34      | 105.47     | 146.60     | OK        
UgorjiSimple/GozstdFastest/len=100 | 50.60      | 68.09      | 121.70     | 154.59     | OK        
UgorjiMsgPack/ZstdFastest/len=100 | 50.60      | 58.74      | 105.97     | 131.80     | OK        
UgorjiBinc/GozstdFastest/len=100 | 50.60      | 71.92      | 134.59     | 154.44     | OK        
UgorjiBinc/ZstdFastest/len=100 | 50.60      | 59.76      | 102.54     | 143.21     | OK        
Msgp/ZstdFastest/len=100  | 50.60      | 87.92      | 144.71     | 224.04     | OK        
UgorjiBinc/DataDogZstdFastest/len=100 | 50.60      | 73.43      | 131.74     | 165.91     | OK        
UgorjiSimple/DataDogZstdFastest/len=100 | 50.60      | 73.24      | 133.34     | 162.52     | OK        
Avro/Snappy/len=100       | 53.01      | 216.53     | 337.59     | 603.86     | OK        
Avro/StdSnappy/len=100    | 55.42      | 268.71     | 488.07     | 597.89     | OK        
Avro/S2/len=100           | 55.42      | 213.90     | 448.38     | 409.01     | OK        
Avro/Lz4/len=100          | 56.63      | 182.89     | 334.56     | 403.40     | OK        
EasyJSON/Snappy/len=100   | 62.65      | 61.86      | 123.59     | 123.85     | OK        
Sonic/Snappy/len=100      | 62.65      | 56.07      | 125.34     | 101.46     | OK        
JsonIter/Snappy/len=100   | 62.65      | 51.73      | 152.32     | 78.33      | OK        
UgorjiSimple/Snappy/len=100 | 65.06      | 76.18      | 136.71     | 172.03     | OK        
ShamatonMsgPack/Snappy/len=100 | 65.06      | 101.39     | 188.76     | 219.07     | OK        
Sonic/S2/len=100          | 65.06      | 64.23      | 160.46     | 107.11     | OK        
JsonIterFastest/Lz4/len=100 | 65.06      | 50.49      | 174.43     | 71.07      | OK        
UgorjiMsgPack/Snappy/len=100 | 65.06      | 79.37      | 146.27     | 173.52     | OK        
JsonIter/Lz4/len=100      | 65.06      | 52.50      | 168.46     | 76.26      | OK        
JsonIter/S2/len=100       | 65.06      | 54.17      | 190.58     | 75.69      | OK        
Msgp/Snappy/len=100       | 65.06      | 128.26     | 244.69     | 269.54     | OK        
CBOR/Snappy/len=100       | 65.06      | 67.61      | 193.49     | 103.92     | OK        
EasyJSON/Lz4/len=100      | 65.06      | 68.34      | 181.53     | 109.61     | OK        
UgorjiBinc/Snappy/len=100 | 65.06      | 77.21      | 132.28     | 185.48     | OK        
MsgPack/Snappy/len=100    | 65.06      | 67.75      | 138.33     | 132.78     | OK        
Sonic/Lz4/len=100         | 65.06      | 55.84      | 128.46     | 98.77      | OK        
GoccyJSON/S2/len=100      | 65.06      | 53.59      | 110.80     | 103.77     | OK  
...
```