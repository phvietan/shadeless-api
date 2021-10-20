import { Db, MongoClient } from 'mongodb';
import * as Bluebird from 'bluebird';
import { ConfigService } from 'config/config.service';
import { Controller } from '@nestjs/common';
import AllDatabases from 'libs/databases/all.database';
import { sleep } from 'libs/helper';
import PathFuzzer from 'libs/pathFuzzer/pathFuzzer';

@Controller('background')
export class BackgroundController {
  private static instance?: BackgroundController;

  constructor() {
    console.log('Initialized BackgroundController');
  }

  private static async bootstrapDatabases(): Promise<Db> {
    const { databaseUrl } = ConfigService.getInstance().getConfig();
    return new Promise<Db>((resolve, reject) => {
      MongoClient.connect(databaseUrl, function (err, db) {
        if (err || !db) reject(err);
        const dbo = db.db('shadeless');
        AllDatabases.getInstance(dbo);
        resolve(dbo);
      });
    });
  }

  static async bootstrapBackground() {
    if (this.instance) return;
    this.instance = new BackgroundController();
    await this.bootstrapDatabases();

    const { parsedPathDb, botPathDb, projectDb } = AllDatabases.getInstance();
    await parsedPathDb.resetScanning();

    while (true) {
      const botPathRunning = await botPathDb.getRunningProject();
      await Bluebird.map(botPathRunning, async (bp) => {
        const project = await projectDb.getOneProjectByName(bp.project);
        const fuzzPaths = await parsedPathDb.getTodo(project);
        await Bluebird.map(fuzzPaths, async (path) => {
          console.log(`Fuzzing path: ${path.origin}${path.path}`);
          const pathFuzzer = new PathFuzzer(path, {
            ...bp,
          });
          await pathFuzzer.run();
          await sleep(3000);
        });
      });
      await sleep(5000);
    }
  }
}
