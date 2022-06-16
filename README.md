# signatories
A smart contract that allows signatories to sign a contract.

A Contract is simply a string. It has a collection of key aliases that are the parties to the contract.

Signatures are stored by signers. A Signature is created by a KA and bears a contract ID.

## API
### create_contract
Creates a contract that can be signed by signatory parties.
A contract contains data and a list of 

### add_signer
Adds a signer, by key alias, to a contract. This operation merely appends a KA to the list of
signers/parties to the contract.

### sign
Allows a signer, by key alias, to sign a contract identified by contract ID.