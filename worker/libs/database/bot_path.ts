import { Collection, Db } from 'mongodb';
import { Packet } from './packet';

export interface ParsedPacket extends Packet {
  hash: string;
}

class ParsedPacketDb {
  db: Collection<Document>;
  private static instance?: ParsedPacketDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('bot_path');
  }

  static getInstance(dbo?: Db): ParsedPacketDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?')
      }
      this.instance = new ParsedPacketDb(dbo);
    }
    return this.instance;
  }

  async getOneByRequestId(requestPacketId: string) {
    const document = await this.db.findOne({ requestPacketId });
    if (!document) return undefined;
    return (document as any) as ParsedPacket;
  }
}

export default ParsedPacketDb;
