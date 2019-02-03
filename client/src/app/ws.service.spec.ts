import { TestBed } from '@angular/core/testing';

import { WSService } from './ws.service';

describe('WSService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: WSService = TestBed.get(WSService);
    expect(service).toBeTruthy();
  });
});
