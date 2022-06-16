# signatories
A smart contract that allows signatories to sign a contract.

A Contract is simply a string. It has a collection of key aliases that are the parties to the contract.

Signatures are stored by signers. A Signature is created by a KA and bears a contract ID.

## API
### contract
Reads the contract, with Signatures.

Available to Contract creator and parties.

### contract_create
Creates a contract that can be signed by signatory parties.

Creates a channel for the Contract parties.

A contract contains data and a list of parties (identified by their KA).

When all of the parties have created signatures for the contract, the contract can be considered to be effective/active/binding.

### contract_add_party
Share the Contract channel with the KeyAlias of the party being added to the Contract.parties
list.

### contract_sign
Allows any of the Cotnract parties, by key alias, to sign a contract identified by contract ID.

Must not have already signed.