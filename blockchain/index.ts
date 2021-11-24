import * as crypto from 'crypto';

class Transaction {
  constructor(
    public amount: number,
    public payer: string,
    public payee: string,
  ) { }

  toString() {
    return JSON.stringify(this);
  }
}

class Block {
  public static nonce = Math.round(Math.random() * 999999999);

  constructor(
    public prevHash: string,
    public transaction: Transaction,
    public ts = Date.now()
  ) { }

  get hash() {
    const str = JSON.stringify(this);
    const hash = crypto.createHash('SHA256');
    hash.update(str).end();
    return hash.digest('hex');
  }
}

class Chain {
  public static instance = new Chain();

  chain: Block[]

  constructor() {
    this.chain = [new Block('', new Transaction(100, 'genesis', 'satoshi'))]
  }

  get lastBlock() {
    return this.chain[this.chain.length - 1];
  }

  addBlock(transaction: Transaction, senderPublicKey: string, signature: Buffer) {
    // Verify the transaction origin's signature
    const verifier = crypto.createVerify('SHA256');

    verifier.update(transaction.toString()).end();
    const isValid = verifier.verify(senderPublicKey, signature);

    if (isValid) {
      // Create a new block by mining
      this.mineBlock(Block.nonce);
      this.chain.push(new Block(this.lastBlock.hash, transaction));
    }
  }

  mineBlock(nonce: number) {
    // The number that will be added to the nonce to adjust it so that the required hash be generated
    let solution = 0;

    console.log('MINING...');

    // Loop until the hash that starts with '0000' is found
    while (true) {
      const hash = crypto.createHash('SHA256');
      hash.update((nonce + solution).toString()).end();
      if (hash.digest('hex').substr(0, 4) === '0000') {
        console.log('BLOCK FOUND ON TRY # ' + solution);
        return solution;
      }

      solution++;
    }
  }
}

class Wallet {
  public publicKey: string;
  public privateKey: string;

  constructor() {
    const keyPair = crypto.generateKeyPairSync('rsa', {
      modulusLength: 2048,
      publicKeyEncoding: { type: 'spki', format: 'pem' },
      privateKeyEncoding: { type: 'pkcs8', format: 'pem' }
    });

    // Create public/private key for this wallet
    this.privateKey = keyPair.privateKey;
    this.publicKey = keyPair.publicKey;
  }

  sendMoney(amount: number, payeePublicKey: string) {
    // Create transaction
    const transaction = new Transaction(amount, this.publicKey, payeePublicKey);

    // Sign the transaction with this wallet private key
    const sign = crypto.createSign('SHA256');
    sign.update(transaction.toString()).end();
    const signature = sign.sign(this.privateKey);

    // Add transaction to blockchain
    Chain.instance.addBlock(transaction, this.publicKey, signature);
  }

}


/* SIMULATION */
const Goku = new Wallet();
const Sienna = new Wallet();
const Colmillo = new Wallet();

/* CHAIN AT THE BEGINNING */
// console.log(Chain.instance);

/* MAKE SOME TRANSACTIONS */
Goku.sendMoney(10, Sienna.publicKey);

/* CHAIN AFTER TRANSACTIONS*/
console.log(Chain.instance);