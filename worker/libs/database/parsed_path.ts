import { Collection, Db } from 'mongodb';

export enum PathStatus {
  TODO = 'todo',
  SCANNING = 'scanning',
  DONE = 'done',
}
export interface ParsedPath {
  requestPacketId: string;
	origin: string;
	path: string;
	status: PathStatus;
  project: string;
  created_at: Date;
  updated_at: Date;
}

class ParsedPathDb {
  FUZZ_PATH_NUM = 1;
  db: Collection<Document>;
  private static instance?: ParsedPathDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('parsed_paths');
  }

  static getInstance(dbo?: Db): ParsedPathDb {
    if (!this.instance) {
      if (!dbo) {
        throw new Error('WTF? DBO is undefined?')
      }
      this.instance = new ParsedPathDb(dbo);
    }
    return this.instance;
  }

  async getTodo() {
    const documents = await this.db.find({
      status: PathStatus.TODO,
    }).sort({ created_at: 1 }).limit(this.FUZZ_PATH_NUM).toArray();
    return (documents as any) as ParsedPath[];
  }
}

export default ParsedPathDb;
