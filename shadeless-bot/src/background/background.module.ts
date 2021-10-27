import { Module } from '@nestjs/common';
import { PathFuzzerService } from './path-fuzzer/path-fuzzer.service';
import { BackgroundController } from './background.controller';
import { ConfigService } from 'config/config.service';
import { ApiFuzzerService } from './api-fuzzer/api-fuzzer.service';

@Module({
  imports: [],
  controllers: [BackgroundController],
  providers: [PathFuzzerService, ConfigService, ApiFuzzerService],
})
export class BackgroundModule {}
