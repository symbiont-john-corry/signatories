# signatories
A smart contract that allows signatories to sign a contract.

A Contract is simply a string. It has a collection of key aliases that are the parties to the contract.

Signatures are stored by signers. A Signature is created by a KA and bears a contract ID.

## API
### create_contract
Creates a contract that can be signed by signatory parties.
A contract contains data and a list of parties (identified by their KA).

When all of the parties have created signatures for the contract, the contract can be considered to be effective/active/binding.


### sign
Allows a signer, by key alias, to sign a contract identified by contract ID.