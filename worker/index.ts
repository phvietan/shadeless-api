import { MongoClient, Db } from 'mongodb';
import { initDatabase } from './libs/database/index';
import ParsedPathDb from './libs/database/parsed_path';

const url = "mongodb://localhost:27017/";

MongoClient.connect(url, function(err, db) {
  if (err || !db) throw err;
  const dbo = db.db("shadeless");
  main(dbo);
});

async function main(dbo: Db) {
  initDatabase(dbo);
  const db = ParsedPathDb.getInstance();
  const fuzzPaths = await db.getTodo();
  console.log(fuzzPaths);
}

