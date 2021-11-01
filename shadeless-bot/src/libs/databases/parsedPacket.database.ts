import { Collection, Db } from 'mongodb';
import { Packet } from './packet.database';
import { FuzzStatus } from './parsedPath.database';
import { getFilterByProjectForBW, Project } from './project.database';

export interface ParsedPacket extends Packet {
  hash: string;
  result: string[];
  status: FuzzStatus;
  staticScore: number;
  logDir: string;
}

class ParsedPacketDb {
  static FUZZ_API_NUM = 1;

  db: Collection<Document>;
  private static instance?: ParsedPacketDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('parsed_packets');
  }

  static getInstance(dbo?: Db): ParsedPacketDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?');
      }
      this.instance = new ParsedPacketDb(dbo);
    }
    return this.instance;
  }

  async getOneByRequestId(requestPacketId: string) {
    const document = await this.db.findOne({ requestPacketId });
    if (!document) return undefined;
    return document as any as ParsedPacket;
  }

  async update(query: any, value: any) {
    await this.db.updateOne(query, {
      $set: value,
    });
  }

  async resetScanning() {
    return this.db.updateMany(
      { status: FuzzStatus.SCANNING },
      {
        $set: { status: FuzzStatus.TODO },
      },
    );
  }

  async getTodo(project: Project) {
    const filter = getFilterByProjectForBW(project);
    const documents = await this.db
      .find({
        status: FuzzStatus.TODO,
        ...filter,
      })
      .sort({ created_at: 1 })
      .limit(ParsedPacketDb.FUZZ_API_NUM)
      .toArray();
    return documents as any as ParsedPacket[];
  }
}

export default ParsedPacketDb;
