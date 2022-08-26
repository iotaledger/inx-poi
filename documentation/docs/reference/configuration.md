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

INX-POI uses [JSON](https://www.json.org) for configuration files.

You can specify which config file INX-POI will use. Follow the `inx-poi` command with the `--config` or `-c` parameter and a path to file:

```bash
inx-poi -c config_defaults.json
```

To get the description of configuration parameters for the installed version of `inx-poi`, you can run the `inx-poi` command with `-h` and `--full` parameters:

```bash
inx-poi -h --full
```

## Application

The `app` object contains the application-related settings.

| Name            | Description                                                                                            | Type    | Default value |
| --------------- | ------------------------------------------------------------------------------------------------------ | ------- | ------------- |
| checkForUpdates | If true, the application will check for updates.                                                       | boolean | true          |
| stopGracePeriod | The maximum time for background processes to finish before the app terminates on system shutdown.      | string  | "5m"          |

### Example

```json
{
  "app": {
    "checkForUpdates": true,
    "stopGracePeriod": "5m"
  }
}
```

## INX

The `inx` object contains the INX-related settings.

| Name    | Description                                        | Type   | Default value    |
| ------- | -------------------------------------------------- | ------ | ---------------- |
| address | To which INX address the extension should connect. | string | "localhost:9029" |

### Example

```json
{
  "inx": {
    "address": "localhost:9029"
  }
}
```

## Proof of Inclusion

The `poi` object contains settings that are related to proofs of inclusion.

| Name                      | Description                                              | Type    | Default value    |
| ------------------------- | -------------------------------------------------------- | ------- | ---------------- |
| bindAddress               | The bind address to which the PoI HTTP server listens.   | string  | "localhost:9687" |
| debugRequestLoggerEnabled | If true, enables the debug logging.                      | boolean | false            |

### Example

```json
{
  "poi": {
    "bindAddress": "localhost:9687",
    "debugRequestLoggerEnabled": false
  }
}
```

## Profiling

The `profiling` object contains settings for profiling.

| Name        | Description                                       | Type    | Default value    |
| ----------- | ------------------------------------------------- | ------- | ---------------- |
| enabled     | If true, enables the the profiling plugin.        | boolean | false            |
| bindAddress | The bind address for the profiler to listen.      | string  | "localhost:6060" |

### Example

```json
{
  "profiling": {
    "enabled": false,
    "bindAddress": "localhost:6060"
  }
}
```