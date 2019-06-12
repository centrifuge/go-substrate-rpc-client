package main



/**
  /**
   * @description Generate a payload and pplies the signature from a keypair

sign (method: Method, account: KeyringPair, { blockHash, era, nonce, version }: SignatureOptions): ExtrinsicSignature {
	const signer = new Address(account.publicKey());
	const signingPayload = new SignaturePayload({
	nonce,
	method,
	era: era || this.era || IMMORTAL_ERA,
	blockHash
	});
	const signature = new Signature(signingPayload.sign(account, version as RuntimeVersion));

return this.injectSignature(signature, signer, signingPayload.nonce, signingPayload.era);
}
 */
