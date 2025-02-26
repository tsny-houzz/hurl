# hurl

Houzz Curl!

Be sure to set these environment variables in your shell:
```
STG_HOUZZ_USER
STG_HOUZZ_PASS
```

# Install

[Find the latest release here](https://github.com/tsny-houzz/hurl/releases) and put the binary in your `PATH`

# Usage

```
NAME:
   hurl - (houzz-curl) Curl substitute for stghouzz routing and testing

USAGE:
   EXAMPLE: hurl -b -c codespace=tsny http://prismic-cms-main.stghouzz.stg-main-eks.stghouzz.com/prismic-cms

DESCRIPTION:
   stghouzz requires basic http auth, this app handles those via env vars STG_HOUZZ_USER and STG_HOUZZ_PASS

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -d          Display only specific headers
   -c value    Set a cookie in the format 'name=value'. Defaults to 'jkdebug=value' if '=' is missing.
   -b          Whether to print the final response body to stdout
   -v          Verbose
   --no-auth   Don't use basic auth with env vars: STG_HOUZZ_USER and STG_HOUZZ_PASS
   --mc        Mimic a browser user-agent
   --help, -h  show help

---
❯ hurl -d -c jeff2610 https://www.stghouzz.com/learn/blog/free-fence-contract-template
Using default cookie: jkdebug=jeff2610

+ https://www.stghouzz.com/learn/blog/free-fence-contract-template
- https 404 Not Found ⚠️
< H: Hz-Serverid: jukwaa-main-master20250226064643033be368c4-5777dfbc57-kdw9f

---

❯ hurl -d https://www.stghouzz.com/products/furniture\?page\=new

+ https://www.stghouzz.com/products/furniture?page=new
- https 200 OK ✅
< H: X-Istio-Vs: prismic-cms-vs-delegate
< H: Hz-Serverid: prismic-cms-main-master20250224173337d93640525b-c58d8bff8-54gtb
```
