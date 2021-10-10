import { Db } from 'mongodb';
import ParsedPacketDb from './parsed_packet';
import ParsedPathDb from './parsed_path';

export function initDatabase(dbo: Db) {
  ParsedPathDb.getInstance(dbo);
  ParsedPacketDb.getInstance(dbo);
}
