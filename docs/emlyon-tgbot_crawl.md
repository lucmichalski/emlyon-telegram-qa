## emlyon-tgbot crawl

Crawl pages from the website

### Synopsis

Crawl pages from the website in all languages or specific one...

```
emlyon-tgbot crawl [flags]
```

### Options

```
      --cache-dir string    Cache output path. (default "./shared/data/emlyon")
      --cache-reset         Reset http cache
  -h, --help                help for crawl
  -j, --jobs int            Parallel jobs (default 4)
  -l, --languages strings   Language to crawl (valid: en,fr,cn) (default [en])
      --no-cache            Disable http cache
```

### Options inherited from parent commands

```
  -n, --db-name string   database name. (default "emlyon-tgbot")
  -d, --debug            debug mode.
  -g, --generate-doc     generate documentation.
  -t, --truncate         truncate database.
  -v, --verbose          verbose output.
```

### SEE ALSO

* [emlyon-tgbot](emlyon-tgbot.md)	 - A smart telegram bot about EM Lyon news and knowledge base.

###### Auto generated by spf13/cobra on 21-Sep-2020