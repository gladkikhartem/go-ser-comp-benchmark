## Go Serialization + Compression benchmarks  (AI generated)

This a benchmark of a combination of serialization and compression speed in Go.
It may be useful to identify best combination of parameters for passing simple repetitive data for network-bounded environments (public internet, slow networks, etc)

Benchmark is run with GOMAXPROCS=1 and results may be vastly different depending on the platform, so it's just for a reference only. Plus it's just AI written - so double check results you are interested in.

Compression ratio is calculated relative to the size of Go struct fields.
Compressed struct contains is not only the data itself, but also wrapper envelope.

Full benchmarks are located in full_results.txt
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

```bash
Performing initial speed check (None compression)...

Excluded Marshalers (Speed < 20MB/s with None compression):
- YAML (2.80 MB/s)
- TOML (9.37 MB/s)
- Bencode (8.12 MB/s)
- XML (6.09 MB/s)
- GoccyYAML (1.31 MB/s)
- BurntSushiTOML (3.71 MB/s)
- GoAvro (10.34 MB/s)

Performing Avro speed check with all compressors...
Excluded Compressors (Speed < 10MB/s with Avro):
- Pzip (4.20 MB/s)
- Brotli (8.13 MB/s)
- XZ (1.04 MB/s)
- Bzip2 (4.89 MB/s)

Running benchmarks...

Excluded 267 results with speed < 50MB/s
ID                        | OrigSize   | MarshSize  | CompSize   | Speed      | Ser+Comp   | De+Deser   | Status    
-------------------------------------------------------------------------------------------------------------------
Avro/ZstdFastest/100      | 83         | 82         | 36         | 133.64     | 203.96     | 387.59     | OK        
EasyJSON/ZstdFastest/100  | 83         | 187        | 39         | 57.18      | 123.79     | 106.27     | OK        
Msgp/ZstdFastest/100      | 83         | 156        | 41         | 95.98      | 158.13     | 244.20     | OK        
CBORCanonical/ZstdFastest/100 | 83         | 160        | 41         | 51.63      | 126.23     | 87.36      | OK        
MsgPack/ZstdFastest/100   | 83         | 156        | 41         | 51.65      | 107.31     | 99.57      | OK        
CBOR/ZstdFastest/100      | 83         | 160        | 41         | 52.52      | 119.49     | 93.72      | OK        
ShamatonMsgPack/ZstdFastest/100 | 83         | 156        | 41         | 69.23      | 121.60     | 160.75     | OK        
UgorjiCBOR/ZstdFastest/100 | 83         | 160        | 41         | 58.11      | 101.66     | 135.66     | OK        
UgorjiMsgPack/ZstdFastest/100 | 83         | 159        | 42         | 60.31      | 112.53     | 129.95     | OK        
UgorjiSimple/ZstdFastest/100 | 83         | 173        | 42         | 59.84      | 104.33     | 140.34     | OK        
UgorjiBinc/ZstdFastest/100 | 83         | 164        | 42         | 60.56      | 104.79     | 143.51     | OK        
Avro/Snappy/100           | 83         | 82         | 44         | 168.16     | 297.01     | 387.64     | OK        
Avro/StdSnappy/100        | 83         | 82         | 45         | 223.55     | 377.96     | 547.21     | OK        
Avro/S2/100               | 83         | 82         | 46         | 301.88     | 568.93     | 643.12     | OK        
Avro/Lz4/100              | 83         | 82         | 46         | 167.60     | 291.55     | 394.23     | OK        
EasyJSON/Snappy/100       | 83         | 187        | 52         | 65.55      | 141.83     | 121.87     | OK        
Sonic/Snappy/100          | 83         | 187        | 52         | 54.26      | 127.10     | 94.69      | OK        
ShamatonMsgPack/Snappy/100 | 83         | 156        | 53         | 98.42      | 192.01     | 201.93     | OK        
Msgp/Snappy/100           | 83         | 156        | 53         | 144.23     | 258.89     | 325.66     | OK        
MsgPack/Snappy/100        | 83         | 156        | 53         | 74.61      | 164.17     | 136.77     | OK        
UgorjiMsgPack/Snappy/100  | 83         | 159        | 53         | 86.41      | 167.89     | 178.04     | OK        
GoccyJSON/S2/100          | 83         | 187        | 54         | 50.68      | 109.75     | 94.15      | OK        
JsonIter/S2/100           | 83         | 187        | 54         | 58.22      | 230.12     | 77.94      | OK        
JsonIterFastest/S2/100    | 83         | 187        | 54         | 58.41      | 226.56     | 78.71      | OK        
Sonic/StdSnappy/100       | 83         | 187        | 54         | 66.32      | 157.17     | 114.73     | OK        
Sonic/Lz4/100             | 83         | 187        | 54         | 54.04      | 129.61     | 92.67      | OK        
Sonic/S2/100              | 83         | 187        | 54         | 64.24      | 172.29     | 102.43     | OK        
UgorjiCBOR/Snappy/100     | 83         | 160        | 54         | 85.17      | 167.01     | 173.80     | OK        
JsonIter/StdSnappy/100    | 83         | 187        | 54         | 60.12      | 220.31     | 82.69      | OK        
JsonIter/Lz4/100          | 83         | 187        | 54         | 51.24      | 177.86     | 71.97      | OK        
GoccyJSON/StdSnappy/100   | 83         | 187        | 54         | 51.70      | 107.94     | 99.24      | OK        
UgorjiBinc/Snappy/100     | 83         | 164        | 54         | 73.69      | 134.48     | 162.99     | OK        
UgorjiSimple/Snappy/100   | 83         | 173        | 54         | 69.33      | 118.83     | 166.43     | OK        
EasyJSON/S2/100           | 83         | 187        | 54         | 80.24      | 245.03     | 119.30     | OK        
EasyJSON/Lz4/100          | 83         | 187        | 54         | 66.04      | 180.86     | 104.01     | OK        
JsonIterFastest/StdSnappy/100 | 83         | 187        | 54         | 53.33      | 196.30     | 73.23      | OK        
EasyJSON/StdSnappy/100    | 83         | 187        | 54         | 80.63      | 212.22     | 130.03     | OK        
CBOR/Snappy/100           | 83         | 160        | 54         | 67.89      | 192.77     | 104.80     | OK        
CBORCanonical/Snappy/100  | 83         | 160        | 54         | 61.03      | 163.03     | 97.55      | OK        
UgorjiCBOR/StdSnappy/100  | 83         | 160        | 55         | 93.64      | 200.46     | 175.71     | OK        
Msgp/Lz4/100              | 83         | 156        | 55         | 129.10     | 273.24     | 244.71     | OK        
ShamatonMsgPack/S2/100    | 83         | 156        | 55         | 105.74     | 240.77     | 188.53     | OK        
ShamatonMsgPack/Lz4/100   | 83         | 156        | 55         | 91.18      | 185.37     | 179.44     | OK        
ShamatonMsgPack/StdSnappy/100 | 83         | 156        | 55         | 116.74     | 243.45     | 224.29     | OK        
UgorjiSimple/StdSnappy/100 | 83         | 173        | 55         | 86.76      | 172.20     | 174.85     | OK        
UgorjiMsgPack/StdSnappy/100 | 83         | 159        | 55         | 94.04      | 204.90     | 173.80     | OK        
Msgp/S2/100               | 83         | 156        | 55         | 165.02     | 407.61     | 277.27     | OK        
UgorjiSimple/S2/100       | 83         | 173        | 55         | 94.53      | 196.26     | 182.37     | OK        
MsgPack/S2/100            | 83         | 156        | 55         | 74.18      | 198.80     | 118.35     | OK        
MsgPack/Lz4/100           | 83         | 156        | 55         | 63.45      | 154.06     | 107.89     | OK        
UgorjiBinc/StdSnappy/100  | 83         | 164        | 55         | 80.05      | 167.55     | 153.28     | OK        
MsgPack/StdSnappy/100     | 83         | 156        | 55         | 72.34      | 179.49     | 121.18     | OK        
Msgp/StdSnappy/100        | 83         | 156        | 55         | 201.11     | 441.02     | 369.68     | OK        
CBOR/StdSnappy/100        | 83         | 160        | 55         | 68.90      | 240.23     | 96.60      | OK        
UgorjiCBOR/Lz4/100        | 83         | 160        | 56         | 75.60      | 164.04     | 140.21     | OK        
CBOR/Lz4/100              | 83         | 160        | 56         | 63.27      | 183.45     | 96.57      | OK        
UgorjiBinc/S2/100         | 83         | 164        | 56         | 93.23      | 203.71     | 171.89     | OK        
UgorjiBinc/Lz4/100        | 83         | 164        | 56         | 75.41      | 148.93     | 152.77     | OK        
CBORCanonical/StdSnappy/100 | 83         | 160        | 56         | 76.45      | 227.97     | 115.03     | OK        
CBOR/S2/100               | 83         | 160        | 56         | 72.95      | 256.51     | 101.94     | OK        
CBORCanonical/Lz4/100     | 83         | 160        | 56         | 59.17      | 169.32     | 90.95      | OK        
CBORCanonical/S2/100      | 83         | 160        | 56         | 70.19      | 232.30     | 100.58     | OK        
UgorjiMsgPack/S2/100      | 83         | 159        | 56         | 88.24      | 202.15     | 156.60     | OK        
UgorjiMsgPack/Lz4/100     | 83         | 159        | 56         | 74.89      | 150.56     | 148.99     | OK        
UgorjiSimple/Lz4/100      | 83         | 173        | 56         | 72.58      | 145.84     | 144.48     | OK        
UgorjiCBOR/S2/100         | 83         | 160        | 56         | 90.79      | 190.56     | 173.40     | OK        
BSON/StdSnappy/100        | 83         | 200        | 65         | 50.10      | 130.39     | 81.35      | OK        
Avro/None/100             | 83         | 82         | 82         | 456.32     | 1110.53    | 774.61     | OK        
Avro/None/1               | 90         | 94         | 94         | 128.70     | 274.20     | 242.54     | OK        
Avro/StdSnappy/1          | 90         | 94         | 96         | 142.97     | 257.13     | 321.99     | OK        
Avro/S2/1                 | 90         | 94         | 97         | 64.42      | 112.25     | 151.21     | OK        
Avro/ZstdFastest/1        | 90         | 94         | 107        | 77.93      | 135.98     | 182.54     | OK        
Msgp/None/100             | 83         | 156        | 156        | 354.52     | 1331.21    | 483.20     | OK        
MsgPack/None/100          | 83         | 156        | 156        | 96.39      | 278.31     | 147.46     | OK        
ShamatonMsgPack/None/100  | 83         | 156        | 156        | 146.11     | 355.86     | 247.89     | OK        
UgorjiMsgPack/None/100    | 83         | 159        | 159        | 115.04     | 284.87     | 192.96     | OK        
CBORCanonical/None/100    | 83         | 160        | 160        | 80.69      | 407.02     | 100.64     | OK        
CBOR/None/100             | 83         | 160        | 160        | 95.64      | 441.46     | 122.09     | OK        
UgorjiCBOR/None/100       | 83         | 160        | 160        | 136.23     | 380.16     | 212.31     | OK        
UgorjiBinc/None/100       | 83         | 164        | 164        | 122.97     | 298.08     | 209.33     | OK        
UgorjiSimple/None/100     | 83         | 173        | 173        | 82.86      | 167.34     | 164.15     | OK        
Msgp/None/1               | 90         | 180        | 180        | 274.92     | 845.94     | 407.27     | OK        
ShamatonMsgPack/None/1    | 90         | 180        | 180        | 85.75      | 172.76     | 170.27     | OK        
MsgPack/None/1            | 90         | 180        | 180        | 56.03      | 158.10     | 86.78      | OK        
ShamatonMsgPack/Snappy/1  | 90         | 180        | 184        | 65.96      | 114.34     | 155.89     | OK        
Msgp/StdSnappy/1          | 90         | 180        | 184        | 166.21     | 327.46     | 337.54     | OK        
CBOR/None/1               | 90         | 184        | 184        | 59.48      | 220.98     | 81.38      | OK        
ShamatonMsgPack/StdSnappy/1 | 90         | 180        | 184        | 71.86      | 134.13     | 154.77     | OK        
CBORCanonical/None/1      | 90         | 184        | 184        | 71.30      | 282.00     | 95.42      | OK        
UgorjiCBOR/None/1         | 90         | 184        | 184        | 75.17      | 173.59     | 132.59     | OK        
UgorjiMsgPack/None/1      | 90         | 184        | 184        | 71.01      | 160.06     | 127.64     | OK        
ShamatonMsgPack/S2/1      | 90         | 180        | 184        | 70.94      | 127.22     | 160.36     | OK        
Msgp/Snappy/1             | 90         | 180        | 184        | 140.15     | 248.74     | 321.05     | OK        
Msgp/S2/1                 | 90         | 180        | 184        | 159.86     | 306.04     | 334.70     | OK        
JsonIterFastest/None/100  | 83         | 187        | 187        | 71.21      | 407.14     | 86.30      | OK        
JsonIter/None/100         | 83         | 187        | 187        | 67.13      | 350.94     | 83.01      | OK        
EasyJSON/None/100         | 83         | 187        | 187        | 114.24     | 427.53     | 155.89     | OK        
GoccyJSON/None/100        | 83         | 187        | 187        | 58.59      | 123.04     | 111.86     | OK        
Sonic/None/100            | 83         | 187        | 187        | 77.00      | 234.14     | 114.74     | OK        
CBOR/StdSnappy/1          | 90         | 184        | 188        | 58.22      | 168.48     | 88.96      | OK        
UgorjiCBOR/S2/1           | 90         | 184        | 188        | 61.32      | 123.43     | 121.84     | OK        
UgorjiMsgPack/StdSnappy/1 | 90         | 184        | 188        | 63.67      | 132.09     | 122.93     | OK        
UgorjiCBOR/Snappy/1       | 90         | 184        | 188        | 58.94      | 114.69     | 121.24     | OK        
CBORCanonical/Snappy/1    | 90         | 184        | 188        | 55.87      | 145.84     | 90.56      | OK        
CBOR/Snappy/1             | 90         | 184        | 188        | 55.99      | 148.17     | 90.00      | OK        
UgorjiMsgPack/Snappy/1    | 90         | 184        | 188        | 57.62      | 111.54     | 119.21     | OK        
UgorjiCBOR/StdSnappy/1    | 90         | 184        | 188        | 61.39      | 128.03     | 117.93     | OK        
CBORCanonical/S2/1        | 90         | 184        | 188        | 59.54      | 168.43     | 92.09      | OK        
UgorjiMsgPack/S2/1        | 90         | 184        | 188        | 62.87      | 125.85     | 125.64     | OK        
CBOR/S2/1                 | 90         | 184        | 188        | 59.69      | 171.71     | 91.50      | OK        
CBORCanonical/StdSnappy/1 | 90         | 184        | 188        | 55.22      | 168.79     | 82.08      | OK        
UgorjiBinc/None/1         | 90         | 189        | 189        | 73.74      | 162.76     | 134.83     | OK        
UgorjiBinc/StdSnappy/1    | 90         | 189        | 193        | 61.25      | 122.96     | 122.06     | OK        
UgorjiBinc/Snappy/1       | 90         | 189        | 193        | 57.06      | 108.35     | 120.55     | OK        
ShamatonMsgPack/ZstdFastest/1 | 90         | 180        | 193        | 52.24      | 90.49      | 123.55     | OK        
UgorjiBinc/S2/1           | 90         | 189        | 193        | 60.16      | 115.60     | 125.47     | OK        
Msgp/Lz4/1                | 90         | 180        | 199        | 54.80      | 94.41      | 130.60     | OK        
BSON/None/100             | 83         | 200        | 200        | 61.37      | 200.08     | 88.53      | OK        
UgorjiSimple/None/1       | 90         | 203        | 203        | 73.38      | 166.12     | 131.45     | OK        
UgorjiSimple/StdSnappy/1  | 90         | 203        | 206        | 60.75      | 117.50     | 125.76     | OK        
UgorjiSimple/Snappy/1     | 90         | 203        | 207        | 58.51      | 109.58     | 125.57     | OK        
UgorjiSimple/S2/1         | 90         | 203        | 207        | 61.09      | 118.31     | 126.30     | OK        
JsonIterFastest/None/1    | 90         | 218        | 218        | 50.54      | 184.60     | 69.59      | OK        
Sonic/None/1              | 90         | 218        | 218        | 54.61      | 126.83     | 95.91      | OK        
EasyJSON/None/1           | 90         | 218        | 218        | 90.42      | 237.03     | 146.18     | OK        
EasyJSON/StdSnappy/1      | 90         | 218        | 220        | 64.71      | 134.63     | 124.60     | OK        
EasyJSON/S2/1             | 90         | 218        | 222        | 72.54      | 152.68     | 138.20     | OK        
EasyJSON/Snappy/1         | 90         | 218        | 222        | 66.27      | 133.93     | 131.18     | OK        
EasyJSON/ZstdFastest/1    | 90         | 218        | 231        | 53.94      | 105.00     | 110.92     | OK
```