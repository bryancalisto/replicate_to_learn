"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const crypto = __importStar(require("crypto"));
class Transaction {
    constructor(amount, payer, payee) {
        this.amount = amount;
        this.payer = payer;
        this.payee = payee;
    }
    toString() {
        return JSON.stringify(this);
    }
}
class Block {
    constructor(prevHash, transaction, ts = Date.now()) {
        this.prevHash = prevHash;
        this.transaction = transaction;
        this.ts = ts;
    }
    get hash() {
        const str = JSON.stringify(this);
        const hash = crypto.createHash('SHA256');
        hash.update(str).end();
        return hash.digest('hex');
    }
}
Block.nonce = Math.round(Math.random() * 999999999);
class Chain {
    constructor() {
        this.chain = [new Block('', new Transaction(100, 'genesis', 'satoshi'))];
    }
    get lastBlock() {
        return this.chain[this.chain.length - 1];
    }
    addBlock(transaction, senderPublicKey, signature) {
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
    mineBlock(nonce) {
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
Chain.instance = new Chain();
class Wallet {
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
    sendMoney(amount, payeePublicKey) {
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
