import {BigInt, cosmos} from "@graphprotocol/graph-ts";
import {BankTx} from "../generated/schema";

export function handleTx(data: cosmos.TransactionData): void {
    const messages = data.tx.tx.body.messages;
    for (let i = 0; i < messages.length; i++) {
        if (messages[i].typeUrl == "/cosmos.bank.v1beta1.MsgSend") {
            const hash = data.tx.hash.toHexString();
            const msg = new BankTx(hash);
            msg.height = BigInt.fromU64(data.block.header.height);
            msg.save();
            break
        }
    }
}
