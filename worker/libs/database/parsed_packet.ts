import { Collection, Db } from 'mongodb';

export interface ParsedPacketRequest {
  method: string;
  project: string;
  requestBodyHash: string;
  requestHeaders: string[];
  origin: string;
  path: string;
  querystring: string;
}
export interface ParsedPacket extends ParsedPacketRequest {
  _id?: string;
  requestPacketId: string;
  requestPacketPrefix: string;
  requestPacketIndex: number;

  toolName: string;
  requestLength: number;
  requestHttpVersion: string;
  requestContentType: string;
  referer: string;
  protocol: string;
  port: number;
  requestCookies: string;
  hasBodyParam: boolean;
  parameters: string[];

  responseStatus: number;
  responseContentType: string;
  responseStatusText: string;
  responseLength: number;
  responseMimeType: string;
  responseHttpVersion: string;
  responseInferredMimeType: string;
  responseCookies: string;
  responseBodyHash: string;
  responseHeaders: string[];

  rtt: number;
  reflectedParameters: Record<string, string>;
  codeName: string;

  created_at?: string;
}

class ParsedPacketDb {
  db: Collection<Document>;
  private static instance?: ParsedPacketDb = undefined;

  constructor(dbo: Db) {
    this.db = dbo.collection('parsed_packets');
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
