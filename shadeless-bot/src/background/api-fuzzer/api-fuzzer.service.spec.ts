import { Test, TestingModule } from '@nestjs/testing';
import { ApiFuzzerService } from './api-fuzzer.service';

describe('ApiFuzzerService', () => {
  let service: ApiFuzzerService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [ApiFuzzerService],
    }).compile();

    service = module.get<ApiFuzzerService>(ApiFuzzerService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
