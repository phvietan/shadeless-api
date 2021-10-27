import { Db, MongoClient } from 'mongodb';
import { ConfigService } from 'config/config.service';
import { Controller } from '@nestjs/common';
import AllDatabases from 'libs/databases/all.database';
import ShadelessLogger from 'libs/logger/logger';
import { PathFuzzerService } from './path-fuzzer/path-fuzzer.service';

@Controller('background')
export class BackgroundController {
  private readonly logger = new ShadelessLogger();

  constructor(
    private configService: ConfigService,
    private pathFuzzerService: PathFuzzerService,
  ) {
    this.logger.log('Initialized BackgroundController');
    this.bootstrapBackground();
  }

  private async bootstrapDatabases(): Promise<Db> {
    const { databaseUrl } = this.configService.getConfig();
    return new Promise<Db>((resolve, reject) => {
      MongoClient.connect(databaseUrl, function (err, db) {
        if (err || !db) reject(err);
        const dbo = db.db('shadeless');
        AllDatabases.getInstance(dbo);
        resolve(dbo);
      });
    });
  }

  async bootstrapBackground() {
    await this.bootstrapDatabases();
    await Promise.all([this.pathFuzzerService.runForever()]);
  }
}
