# INX-POI

INX-POI is an INX plugin which creates and validates Proof-of-Inclusions for blocks.

## Usage

### Create a Proof-of-Inclusion

To create a POI send a GET request to `http://{your-node-address}:14265/api/plugins/poi/v1/create/{blockID}`

### Validate a Proof-of-Inclusion

To validate a POI send the `application/json` payload of the created proof via POST request to `http://{your-node-address}:14265/api/plugins/poi/v1/validate`

## Example
```json
{
  "milestone": {
    "type": 7,
    "index": 13,
    "timestamp": 1653768570,
    "protocolVersion": 2,
    "previousMilestoneId": "0x17c0a6a711857ea46158ca46ed20daa09cf7b3fa9e7dbab67b4ba3b90ebba77a",
    "parents": [
      "0x417aab094d8e73b439f8cc68f8e7d83be2239bb34d20332f52e9cd7d6534ae6c",
      "0x4a0dc52628bd688cfd83028d13ad4ab3b8ef9f28a44a3064fa22309660e7dc43",
      "0x5b7b045b8b09980bcc8229eb3eb304a960b035c4737e33ea1a24d65b065df83c",
      "0x9b7d35e3e17f00e8bf221890a55ae14bbd0a52a4624defa6a88d5235e00c7d80",
      "0xc8e8ca9c3c9a5111520b41c37086f7e0249ed1a8d619976f011be8abeb8771a8",
      "0xf5d25ae03293dc54115b78b100c41ac540df00925c9d0ae95431f09e3f7be1d1"
    ],
    "inclusionMerkleRoot": "0xee3c9836ae52b79163cd9f645099edf7e9305d669123a396d73e30e2c3bafdd1",
    "appliedMerkleRoot": "0x0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8",
    "signatures": [
      {
        "type": 0,
        "publicKey": "0xed3c3f1a319ff4e909cf2771d79fece0ac9bd9fd2ee49ea6c0885c9cb3b1248c",
        "signature": "0x84373ad012aefc4966cd53331d40e94183ecfc81aeaf20c71ed1b98ce8a07b1cf4370ea00d97e165b7ee9e8656f351f6010dfa584ebdb66d8233c6c51e840600"
      },
      {
        "type": 0,
        "publicKey": "0xf6752f5f46a53364e2ee9c4d662d762a81efd51010282a75cd6bd03f28ef349c",
        "signature": "0xf444bc745a7d651012dc6b43d4fecc1ea2b17402beed7981395db0c56cc69e4ff1f585e7e52fe6317de9890a1bad2ba89c8e9c5258dba2316c01dccc8472b00b"
      }
    ]
  },
  "block": {
    "protocolVersion": 2,
    "parents": [
      "0x14eef4f3923ba0301621775e7e6f4d550006637bec639e9f9afdf2ab9d715cdb",
      "0x428079a3dbb95f8411f8831dc1bf1d3ba723327fd3ae1741eaafd22bff9eb468",
      "0x5b7b045b8b09980bcc8229eb3eb304a960b035c4737e33ea1a24d65b065df83c",
      "0x9003301a44cd04bf1911f82de72ad5a050359a880b6e02507f5d2b793b3b7ce3"
    ],
    "payload": {
      "type": 5,
      "tag": "0x484f524e4554205370616d6d6572",
      "data": "0x57652061726520616c6c206d616465206f662073746172647573742e0a436f756e743a203030303138380a54696d657374616d703a20323032322d30352d32385432303a30393a32375a0a54697073656c656374696f6e3a20323732c2b573"
    },
    "nonce": "299"
  },
  "proof": {
    "l": {
      "l": {
        "h": "0x6e463cb72c8639dbfc820e7a0349907e2353ac2afea3c7cf1492771d18a8e789"
      },
      "r": {
        "l": {
          "l": {
            "h": "0xf5e591867dea12da2e9777f393af0d7eb7055c9ddbe08a9e235781cfb1b5bab2"
          },
          "r": {
            "l": {
              "value": "0xb00ff4ee4cc5aeb94d7e901d2afe9b27ab568442e683aa2e8e9be0f8e894eb1f"
            },
            "r": {
              "h": "0xac7edca5fef53bce504e52448d06b5b1d7da9232cb6e6407a126a1262f393768"
            }
          }
        },
        "r": {
          "h": "0x3757577f93f26bbe0db47b1465752ad49d220ee7ee57aa8902029f361dab6afb"
        }
      }
    },
    "r": {
      "h": "0x9f9be742aab1eeeb033d39f2f55c421ad08bc0c7508e26c3fd116d78c1500abc"
    }
  }
}
```
