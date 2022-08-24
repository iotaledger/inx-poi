---
description: INX-POI enables you to generate and verify Proof-of-Inclusion of blocks in the Tangle.
image: /img/Banner/banner_hornet.png
keywords:
- IOTA Node
- Hornet Node
- INX
- POI
- Proof-Of-Inclusion
- Proof
- Inclusion
- IOTA
- Shimmer
- Node Software
- Welcome
- explanation
---

# Welcome to INX-POI

INX-POI generates and verifies proofs of inclusion of blocks. By generating a proof, you can always verify that the block used to be part of the Tangle. The Tangle nodes truncate old blocks by necessity, but with proofs of inclusion you can preserve important blocks on your own.

## Setup

We recommend you to use the [Docker images](https://hub.docker.com/r/iotaledger/inx-poi).
These images are also used in the [Docker setup](http://wiki.iota.org/hornet/develop/how_tos/using_docker) of Hornet.

## Configuration

The extension will connect to your local HORNET instance by default.

You can find all the configuration options in the [configuration section](reference/configuration.md).

## API

The extension exposes a custom set of REST APIs that you can use to generate or validate proofs.

You can find more information about the API in the [API reference section](reference/api_reference.md).

## Source Code

The source code of the project is available on [GitHub](https://github.com/iotaledger/inx-poi).