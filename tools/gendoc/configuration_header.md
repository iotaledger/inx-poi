---
description: This section describes the configuration parameters and their types for INX-POI.
keywords:
- IOTA Node 
- Hornet Node
- POI
- Proof-Of-Inclusion
- Proof
- Inclusion
- Configuration
- JSON
- Customize
- Config
- reference
---


# Core Configuration

INX-POI uses a JSON standard format as a config file. If you are unsure about JSON syntax, you can find more information in the [official JSON specs](https://www.json.org).

You can change the path of the config file by using the `-c` or `--config` argument while executing `inx-poi` executable.

For example:
```bash
inx-poi -c config_defaults.json
```

You can always get the most up-to-date description of the config parameters by running:

```bash
inx-poi -h --full
```

