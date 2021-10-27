import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { ConfigService } from 'config/config.service';
import { BackgroundModule } from './background/background.module';

@Module({
  imports: [BackgroundModule],
  controllers: [AppController],
  providers: [ConfigService],
})
export class AppModule {}
