import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChannelsComponent } from './channels.component';

describe('Channels', () => {
  let component: ChannelsComponent;
  let fixture: ComponentFixture<ChannelsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChannelsComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ChannelsComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
