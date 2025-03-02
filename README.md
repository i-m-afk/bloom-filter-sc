# bloom-filter-sc

A spell checker using bloom filter. This is my attempt for one of the many coding challenges by [John Crickett.](https://codingchallenges.fyi/challenges/challenge-bloom)
One can use this site https://hur.st/bloomfilter/?n=104334&p=1.0E-3&m=&k=15 for calculating their required size of bloom filter.

## Instructions to run

- In terminal run `go build && ./bloom-filter` this will run spell checker program. Make sure you have dictionary file `dict.txt` and save the bloom filter in directory.
- use `--dict <filename>` to use different dictionary.
- use `--bf <filename>` to use existing bloom filter.

[How it works](https://imrishav.vercel.app/blog/bloom-filter-spell-checker/)
