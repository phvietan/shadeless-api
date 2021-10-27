import { Test, TestingModule } from '@nestjs/testing';
import { PathFuzzerService } from './path-fuzzer.service';

describe('PathFuzzerService', () => {
  let service: PathFuzzerService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [PathFuzzerService],
    }).compile();

    service = module.get<PathFuzzerService>(PathFuzzerService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
