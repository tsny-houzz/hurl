# hurl

Houzz Curl!

```
❯ hurl -h
Usage of hurl:
  -b    Whether to print the final response body to stdout
  -c string
        Set a cookie in the format 'name=value'. Defaults to 'jkdebug=value' if '=' is missing.
  -d    Display only specific headers

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
