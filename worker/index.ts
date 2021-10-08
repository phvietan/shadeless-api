import Bluebird from 'bluebird';
import PathFuzzer from './libs/fuzzer/path';
import { MongoClient, Db } from 'mongodb';
import { initDatabase } from './libs/database/index';
import ParsedPathDb from './libs/database/parsed_path';

const url = "mongodb://localhost:27017/";

MongoClient.connect(url, function(err, db) {
  if (err || !db) throw err;
  const dbo = db.db("shadeless");
  main(dbo);
});

async function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function main(dbo: Db) {
  initDatabase(dbo);
  const db = ParsedPathDb.getInstance();
  await db.resetScanning();

  while (true) {
    const fuzzPaths = await db.getTodo();
    await Bluebird.map(fuzzPaths, async (path) => {
      console.log(`Fuzzing path: ${path.origin}${path.path}`);
      const pathFuzzer = new PathFuzzer(path);
      await pathFuzzer.run();
      await sleep(3000);
    });
  }

}

