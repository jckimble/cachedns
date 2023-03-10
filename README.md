# cachedns
![build status](https://github.com/jckimble/cachedns/actions/workflows/build.yml/badge.svg?branch=master)

A Caching DNS Server with domain filters and preloading blacklist

---
* [Install](#install)
* [Config](#config)

---

## Install
```sh
go get -u github.com/jckimble/cachedns
```

## Configuration
Config File must exist as `/etc/cdns.yaml` or `./cdns.yaml`
```yaml
preload:
  - http://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext
resolvers:
  filtering:
    enabled: true
    filters:
      ".*google.*": 127.0.0.1
  forwarding:
    enabled: false
    resolvers:
      - 8.8.8.8
      - 8.8.4.4
  docker:
    enabled: true
port: ":5353"
```

## License

Copyright 2018 James Kimble

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
