# hurl

Houzz Curl!

Be sure to set these environment variables in your shell:
```
STG_HOUZZ_USER
STG_HOUZZ_PASS
```

# Usage

```
NAME:
   hurl - Curl substitute for stghouzz routing and testing

USAGE:
   EXAMPLE: hurl -b -c codespace=tsny http://prismic-cms-main.stghouzz.stg-main-eks.stghouzz.com/prismic-cms

DESCRIPTION:
   Basic auth is handled by env vars STG_HOUZZ_USER and STG_HOUZZ_PASS

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -d          Display only specific headers
   -c value    Set a cookie in the format 'name=value'. Defaults to 'jkdebug=value' if '=' is missing.
   -b          Whether to print the final response body to stdout
   -v          Verbose
   --no-auth   Don't use basic auth with env vars: STG_HOUZZ_USER and STG_HOUZZ_PASS
   --help, -h  show help

---

❯ hurl -d -c jeff2610 https://www.stghouzz.com/learn/blog/free-fence-contract-template
> jkdebug: jeff2610

> https://www.stghouzz.com/learn/blog/free-fence-contract-template
- https 302 Found
Hz-Serverid: prismic-cms-debug-jeff2610-6d696bbc56-jd67t
Location: https://www.stghouzz.com/pro-learn/blog/free-fence-contract-template

> https://www.stghouzz.com/pro-learn/blog/free-fence-contract-template
- https 404 Not Found
Hz-Serverid: jukwaa-main-master20241105073610b3a0a11e9e-7ffd5fb894-wrtss
```
