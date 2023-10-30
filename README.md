# Write Your Own Compression Tool

This respository is a challenge from John Cricket Coding Challenges, [writing a compression tool](https://codingchallenges.fyi/challenges/challenge-huffman).


## How to use

```
go build
```

### Flags
```
-f input file
-o output file
-d decompress 
-v verbose
```

### Compress a text files

```
coding-challenge-compression -f sample.txt -o compressed_sample
```

### Decompress a file
```
coding-challenge-compression -d -f compressed_sample -o uncompressed_sample.txt
```

## Future improvements
* Create more test cases to test edge cases