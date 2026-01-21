/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
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
import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';
import { of, BehaviorSubject } from 'rxjs';

import { TimelineComponent } from './timeline.component';

describe('TimelineComponent', () => {
  let component: TimelineComponent;
  let fixture: ComponentFixture<TimelineComponent>;
  let apiServiceMock: any;

  beforeEach(waitForAsync(() => {
    apiServiceMock = {
      updateStatus: jasmine.createSpy('updateStatus'),
      fetchChannels: jasmine.createSpy('fetchChannels').and.returnValue(of([])),
      fetchProgramForDay: jasmine.createSpy('fetchProgramForDay'),
      statusSubject: new BehaviorSubject<any>(null),
      channelSubjectVar: new BehaviorSubject<any[]>([]),
      programSubject: new BehaviorSubject<any[]>([]),
      isLoadingSubject: new BehaviorSubject<boolean>(false),
      isInErrorsSubject: new BehaviorSubject<boolean>(false),
      isWindowOpenedSubject: new BehaviorSubject<boolean>(true)
    };

    TestBed.configureTestingModule({
      declarations: [ TimelineComponent ],
      imports: [ FormsModule ],
      providers: [
        { provide: ApiService, useValue: apiServiceMock }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
