## Easily block a number

I use [voip.ms](https://voip.ms) to operate my home phone.  I'd like an easy
way to block numbers using the [voip.ms
API](https://www.voip.ms/m/apidocs.php).

I was interested in doing a My First Golang project, so many thanks to
https://github.com/stancarney/govoipms, which showed me the way and was a
valuable resource (i.e. I copied a lot from there.)

## Sample Usage

Block a number

```
go build
./voip-ms-cli -block-number 18001234567 -note "No more telemarkter"
```

Print balance
```
./voip-ms-cli -print-balance
```

<!--
 vim: tw=78
-->
