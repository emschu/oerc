import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LogDashboardComponent } from './log-dashboard.component';

describe('LogdashboardComponent', () => {
  let component: LogDashboardComponent;
  let fixture: ComponentFixture<LogDashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ LogDashboardComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(LogDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
