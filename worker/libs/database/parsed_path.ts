import { Collection, Db } from 'mongodb';

export enum PathStatus {
  TODO = 'todo',
  SCANNING = 'scanning',
  DONE = 'done',
}
export interface ParsedPath {
  _id?: string;
  requestPacketId: string;
	origin: string;
	path: string;
	status: PathStatus;
  project: string;
  force: boolean;
  created_at?: Date;
  updated_at?: Date;
  error: string;
}

type DocumentType = Pick<Document, keyof Document>;

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
      requestPacketId: { $ne: "" },
    }).sort({ created_at: 1 }).limit(this.FUZZ_PATH_NUM).toArray();
    return (documents as any) as ParsedPath[];
  }

  async update(query: any, value: any) {
    await this.db.updateOne(query, {
      $set: value,
    })
  }

  async updateError(query: any, error: string) {
    return this.update(query, { error });
  }

  async insertResult(arrParsedPath: ParsedPath[]) {
    if (arrParsedPath.length === 0) return;
    const bypassCheck = (arrParsedPath as any) as DocumentType[];
    try {
      await this.db.insertMany(bypassCheck);
    } catch (err) {}
  }

  async resetScanning() {
    return this.db.updateMany({
      status: PathStatus.SCANNING,
    }, {
      $set: {
        status: PathStatus.TODO,
      }
    });
  }
}

export default ParsedPathDb;
