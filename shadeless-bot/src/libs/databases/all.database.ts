import { Collection, Db } from 'mongodb';
import BotPathDb from './botPath.database';
import PacketDb from './packet.database';
import ParsedPacketDb from './parsedPacket.database';
import ParsedPathDb from './parsedPath.database';
import ProjectDb from './project.database';

export interface BotPath {
  timeout: number;
  project: string;
  running: boolean;
  sleepRequest: number;
  asyncRequest: number;
}

class AllDatabases {
  dbo: Db;
  db: Collection<Document>;

  packetDb: PacketDb;
  botPathDb: BotPathDb;
  projectDb: ProjectDb;
  parsedPathDb: ParsedPathDb;
  parsedPacketDb: ParsedPacketDb;

  private static instance?: AllDatabases = undefined;

  constructor(dbo: Db) {
    this.dbo = dbo;
    this.packetDb = PacketDb.getInstance(this.dbo);
    this.botPathDb = BotPathDb.getInstance(this.dbo);
    this.projectDb = ProjectDb.getInstance(this.dbo);
    this.parsedPathDb = ParsedPathDb.getInstance(this.dbo);
    this.parsedPacketDb = ParsedPacketDb.getInstance(this.dbo);
  }

  static getInstance(dbo?: Db): AllDatabases {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?');
      }
      this.instance = new AllDatabases(dbo);
    }
    return this.instance;
  }
}

export default AllDatabases;
