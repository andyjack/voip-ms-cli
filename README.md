## Easily block a number

I use [voip.ms](https://voip.ms) to operate my home phone.  I'd like an easy
way to block numbers using the [voip.ms
API](https://www.voip.ms/m/apidocs.php).

I was interested in doing a My First Golang project, so many thanks to
https://github.com/stancarney/govoipms, which showed me the way and was a
valuable resource (i.e. I copied a lot from there.)

## config setup

Create a `voip-ms-cli` in your config dir.  E.g.:

```
mkdir -p ~/.config/voip-ms-cli
```

On linux you might be using `XDG_CONFIG_HOME`, create `voip-ms-cli` dir there
instead.

In this new dir, you'll need to create a config.toml with your email and API
password from voip.ms.  See `config-example.toml`.

## Sample Usage

Block a number

```
go install
voip-ms-cli block-number 18001234567 -note "No more telemarketer"
```

Print balance
```
voip-ms-cli show-balance
```

Other commands:
```
voip-ms-cli show-recent
voip-ms-cli block-recent
```

<!--
 vim: tw=78
-->
