/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2025 emschu[aet]mailbox.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 */
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ApiService } from '../api.service';
import { of, BehaviorSubject } from 'rxjs';

import { LogDashboardComponent } from './log-dashboard.component';

describe('LogdashboardComponent', () => {
  let component: LogDashboardComponent;
  let fixture: ComponentFixture<LogDashboardComponent>;
  let apiServiceMock: any;

  beforeEach(async () => {
    apiServiceMock = {
      updateStatus: jasmine.createSpy('updateStatus'),
      logEntries: jasmine.createSpy('logEntries').and.returnValue(of([])),
      statusSubject: new BehaviorSubject<any>(null),
      isWindowOpenedSubject: new BehaviorSubject<boolean>(true)
    };

    await TestBed.configureTestingModule({
      declarations: [ LogDashboardComponent ],
      providers: [
        { provide: ApiService, useValue: apiServiceMock }
      ]
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
